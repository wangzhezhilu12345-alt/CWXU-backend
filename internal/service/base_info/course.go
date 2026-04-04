package base_info

import (
	"context"
	"errors"

	"edu-evaluation-backed/api/v1/base_info/course"
	student_i "edu-evaluation-backed/api/v1/base_info/student"
	"edu-evaluation-backed/internal/biz/base_info"
	"edu-evaluation-backed/internal/data/dal"
	"edu-evaluation-backed/internal/data/model"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// CourseService 课程信息服务
type CourseService struct {
	baseDal  *dal.BaseDal
	courseDal *dal.CourseDal
	courseUC  *base_info.CourseUseCase
}

// toCourseInfo 将 model.Course 转换为 CourseList
func toCourseInfo(c *model.Course) *course.CourseList {
	info := &course.CourseList{
		Id:          uint32(c.ID),
		CourseName:  c.CourseName,
		ClassName:   c.ClassName,
		Status:      int32(c.Status),
		TeacherList: ToTeacherInfoList(c.Teachers),
	}
	return info
}

// toCourseDetail 转换课程详情（含学生）
func toCourseDetail(c *model.Course) *course.CourseList {
	info := toCourseInfo(c)
	info.StudentList = toStudentListBrief(c.Students)
	return info
}

// toStudentListBrief 将 []model.Student 转换为 []*student_i.StudentInfo (简洁版)
func toStudentListBrief(students []model.Student) []*student_i.StudentInfo {
	result := make([]*student_i.StudentInfo, 0, len(students))
	for _, s := range students {
		result = append(result, &student_i.StudentInfo{
			Id:        uint32(s.ID),
			Name:      s.Name,
			StudentNo: s.StudentNo,
			Sex:       s.Sex,
		})
	}
	return result
}

// Detail 获取课程详情
func (c CourseService) Detail(ctx context.Context, req *course.GetCourseDetailReq) (*course.GetCourseDetailResp, error) {
	cs, err := c.courseDal.Detail(uint(req.CourseId))
	if err != nil {
		return nil, err
	}
	return &course.GetCourseDetailResp{
		Message: "success",
		Data:    toCourseDetail(cs),
	}, nil
}

// Edit 编辑课程信息
func (c CourseService) Edit(ctx context.Context, req *course.EditCourseReq) (*course.EditCourseResp, error) {
	if req.CourseId == 0 {
		return nil, errors.New("课程ID不能为空")
	}

	if req.CourseName != "" || req.ClassName != "" {
		if err := c.courseDal.UpdateCourse(uint(req.CourseId), req.CourseName, req.ClassName); err != nil {
			return nil, err
		}
	}

	if len(req.TeacherIds) > 0 {
		if err := c.courseDal.AddTeachers(uint(req.CourseId), req.TeacherIds); err != nil {
			return nil, err
		}
	}

	cs, err := c.courseDal.Detail(uint(req.CourseId))
	if err != nil {
		return nil, err
	}
	return &course.EditCourseResp{
		Message: "修改成功",
		Data:    toCourseInfo(cs),
	}, nil
}

// List 获取课程列表
func (c CourseService) List(ctx context.Context, req *course.GetCourseListReq) (*course.GetCourseListResp, error) {
	courses, tot, err := c.courseDal.List(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	data := make([]*course.CourseList, 0, len(*courses))
	for _, cItem := range *courses {
		data = append(data, toCourseInfo(&cItem))
	}
	return &course.GetCourseListResp{
		Message: "success",
		Data:    data,
		Total:   tot,
	}, nil
}

// Import 导入课程信息Excel文件
func (c CourseService) Import(ctx http.Context) error {
	req := ctx.Request()
	file, _, err := req.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	msg := c.courseUC.Import(file)
	if msg == "" {
		msg = "导入成功"
	}
	writeJSONResponse(ctx, msg)
	return nil
}

// Delete 删除课程
func (c CourseService) Delete(ctx context.Context, req *course.DeleteCourseReq) (*course.DeleteCourseResp, error) {
	if err := c.courseUC.DeleteCourse(uint(req.Id)); err != nil {
		return nil, err
	}
	return &course.DeleteCourseResp{Message: "删除成功"}, nil
}

// Reset 重置所有数据，只保留学生和教师
func (c CourseService) Reset(ctx context.Context, req *course.ResetReq) (*course.ResetResp, error) {
	if err := c.baseDal.ResetAll(); err != nil {
		return nil, err
	}
	return &course.ResetResp{Message: "重置成功"}, nil
}

// NewCourseService 创建课程信息服务实例
func NewCourseService(baseDal *dal.BaseDal, courseDal *dal.CourseDal, courseUC *base_info.CourseUseCase) *CourseService {
	return &CourseService{baseDal: baseDal, courseDal: courseDal, courseUC: courseUC}
}
