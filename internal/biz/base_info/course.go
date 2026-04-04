package base_info

import (
	"edu-evaluation-backed/internal/data/dal"
	"fmt"
	"mime/multipart"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/xuri/excelize/v2"
)

// CourseUseCase 课程信息业务用例
type CourseUseCase struct {
	courseDal *dal.CourseDal
}

// courseItem 课程分组键
type courseItem struct {
	courseName string
	className  string
}

// Import 从Excel文件导入课程数据
func (c CourseUseCase) Import(f multipart.File) string {
	list, err := excelize.OpenReader(f)
	if err != nil {
		return err.Error()
	}
	defer list.Close()

	rows, err := list.GetRows("Sheet1")
	if err != nil {
		return err.Error()
	}

	// 按课程名称+班级名称分组
	courseMap := make(map[courseItem][]string)
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		item := courseItem{courseName: row[1], className: row[2]}
		courseMap[item] = append(courseMap[item], row[3])
	}

	var logMsg string
	for k, v := range courseMap {
		id, err := c.courseDal.CreateCourse(k.courseName, k.className)
		if err != nil {
			logMsg += fmt.Sprintf("课程:%s,班级:%s,错误:班级已存在\n", k.courseName, k.className)
			continue
		}
		if err := c.courseDal.AddStudent(id, v); err != nil {
			log.Info(err)
			logMsg += fmt.Sprintf("课程:%s,班级:%s,添加学生错误:%s\n", k.courseName, k.className, err.Error())
		}
	}
	return logMsg
}

// DeleteCourse 删除课程
func (c CourseUseCase) DeleteCourse(id uint) error {
	return c.courseDal.DeleteCourse(id)
}

// NewCourseUseCase 创建课程信息业务用例实例
func NewCourseUseCase(courseDal *dal.CourseDal) *CourseUseCase {
	return &CourseUseCase{courseDal: courseDal}
}
