package dal

import (
	"context"
	"time"

	"edu-evaluation-backed/internal/common/data/cache"
	"edu-evaluation-backed/internal/common/utils"
	"edu-evaluation-backed/internal/data"
	"edu-evaluation-backed/internal/data/model"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// CourseDal 课程数据访问层
type CourseDal struct {
	db  *gorm.DB
	rdb *redis.Client
	hc  *cache.HealthChecker
}

// Detail 获取课程详情（带缓存）
func (c CourseDal) Detail(courseID uint) (*model.Course, error) {
	ctx := context.Background()
	key := cache.CourseDetailKey(courseID)

	result, err := cache.Get[model.Course](ctx, c.rdb, c.hc, key, 30*time.Minute, func() (*model.Course, error) {
		course := model.Course{}
		err := c.db.Where("id = ?", courseID).Preload("Teachers").Preload("Students").First(&course).Error
		if err != nil {
			return nil, err
		}
		return &course, nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateCourse 创建新课程
func (c CourseDal) CreateCourse(courseName, className string) (uint, error) {
	cs := model.Course{
		Status:     1,
		CourseName: courseName,
		ClassName:  className,
	}
	err := c.db.Create(&cs)
	return cs.ID, err.Error
}

// AddStudent 添加学生到课程
func (c CourseDal) AddStudent(courseID uint, studentNos []string) error {
	course := model.Course{}
	course.ID = courseID
	var students []model.Student
	log.Println("students: ", studentNos)
	if err := c.db.Where("student_no IN ?", studentNos).Find(&students).Error; err != nil {
		return err
	}
	if len(students) == 0 {
		return errors.New("未找到匹配的学生信息")
	}
	err := c.db.Model(&course).Association("Students").Append(&students)
	if err != nil {
		return err
	}
	return nil
}

// List 获取课程列表（page=-1 时返回全部）
func (c CourseDal) List(page, pageSize int) (*[]model.Course, int64, error) {
	var courses []model.Course
	var tot int64
	q := c.db.Model(&model.Course{}).Preload("Teachers")
	if page == -1 {
		err := q.Count(&tot).Find(&courses).Error
		return &courses, tot, err
	}
	page, pageSize = utils.PageNumHandle(page, pageSize)
	err := q.Count(&tot).Limit(pageSize).Offset(utils.CalculateOffset(page, pageSize)).Find(&courses).Error
	return &courses, tot, err
}

// QueryCourseByIds 批量获取课程信息根据ID列表
func (c CourseDal) QueryCourseByIds(ids []int32) (*[]model.Course, error) {
	var courses []model.Course
	err := c.db.Where("id IN ?", ids).Find(&courses).Error
	return &courses, err
}

// UpdateCourse 更新课程基本信息
func (c CourseDal) UpdateCourse(courseID uint, courseName, className string) error {
	var count int64
	err := c.db.Model(&model.Course{}).Where("class_name = ? AND id != ?", className, courseID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("班级名称已存在")
	}

	updates := make(map[string]interface{})
	if courseName != "" {
		updates["course_name"] = courseName
	}
	if className != "" {
		updates["class_name"] = className
	}

	if len(updates) == 0 {
		return nil
	}

	err = c.db.Model(&model.Course{}).Where("id = ?", courseID).Updates(updates).Error
	return err
}

// AddTeachers 绑定教师到课程
func (c CourseDal) AddTeachers(courseID uint, teacherWorkNos []int32) error {
	err := c.db.Model(&model.Course{Model: gorm.Model{ID: courseID}}).
		Association("Teachers").
		Clear()
	if err != nil {
		return err
	}
	course := model.Course{}
	course.ID = courseID
	var teachers []model.Teacher
	if err := c.db.Where("id IN ?", teacherWorkNos).Find(&teachers).Error; err != nil {
		return err
	}
	if len(teachers) == 0 {
		return errors.New("未找到匹配的教师信息")
	}
	err = c.db.Model(&course).Association("Teachers").Append(&teachers)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCourse 删除课程（硬删除），清理所有关联数据
func (c CourseDal) DeleteCourse(id uint) error {
	var course model.Course
	if err := c.db.First(&course, id).Error; err != nil {
		return err
	}

	// 检查课程是否在活跃评教任务中（status=1 进行中）
	var activeCount int64
	err := c.db.Table("evaluation_courses ec").
		Joins("JOIN evaluation_tasks et ON et.id = ec.evaluation_task_id").
		Where("ec.course_id = ? AND et.status = 1", id).
		Count(&activeCount).Error
	if err != nil {
		return err
	}
	if activeCount > 0 {
		return errors.New("该课程正在进行评教任务中，无法删除")
	}

	// 事务中按层级清理关联数据
	err = c.db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除 evaluation_details 中引用该课程的记录（硬删除）
		if err := tx.Unscoped().Where("course_id = ?", id).Delete(&model.EvaluationDetail{}).Error; err != nil {
			return err
		}
		// 2. 删除 evaluation_courses 中间表记录
		if err := tx.Exec("DELETE FROM evaluation_courses WHERE course_id = ?", id).Error; err != nil {
			return err
		}
		// 3. 清除 course_teachers 关联
		if err := tx.Model(&course).Association("Teachers").Clear(); err != nil {
			return err
		}
		// 4. 清除 course_students 关联
		if err := tx.Model(&course).Association("Students").Clear(); err != nil {
			return err
		}
		// 5. 硬删除课程本身
		if err := tx.Unscoped().Delete(&course).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 事务成功后异步失效缓存
	go func() {
		ctx := context.Background()
		cache.Delete(ctx, c.rdb, c.hc, cache.CourseDetailKey(id))
		cache.DeleteByPattern(ctx, c.rdb, c.hc, cache.TaskListPattern())
	}()

	return nil
}

// UpdateCourseStatus 更新课程状态
func (c CourseDal) UpdateCourseStatus(courseID uint, status int) error {
	return c.db.Model(&model.Course{}).Where("id = ?", courseID).Update("status", status).Error
}

// ResetEvaluationStats 重置课程评教统计
func (c CourseDal) ResetEvaluationStats(courseID uint) error {
	return c.db.Model(&model.Course{}).Where("id = ?", courseID).Updates(map[string]interface{}{
		"evaluation_score": 0,
		"evaluation_num":   0,
	}).Error
}

// NewCourseDal 创建课程数据访问层实例
func NewCourseDal(data *data.Data) *CourseDal {
	return &CourseDal{
		db:  data.DB,
		rdb: data.RDB,
		hc:  data.HC,
	}
}
