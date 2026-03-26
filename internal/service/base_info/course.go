package base_info

import (
	context "context"
	"edu-evaluation-backed/api/v1/base_info/course"
	student_i "edu-evaluation-backed/api/v1/base_info/student"
	"edu-evaluation-backed/api/v1/base_info/teacher"
	"edu-evaluation-backed/internal/biz/base_info"
	"edu-evaluation-backed/internal/data/dal"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/transport/http"
)

type CourseService struct {
	courseDal *dal.CourseDal
	courseUC  *base_info.CourseUseCase
}

func (c CourseService) Detail(ctx context.Context, req *course.GetCourseDetailReq) (*course.GetCourseDetailResp, error) {
	cs, err := c.courseDal.Detail(uint(req.CourseId))
	if err != nil {
		return nil, err
	}
	resp := &course.GetCourseDetailResp{
		Message: "success",
		Data: &course.CourseList{
			Id:          uint32(cs.ID),
			CourseName:  cs.CourseName,
			ClassName:   cs.ClassName,
			TeacherList: make([]*teacher_i.TeacherInfo, 0),
			StudentList: make([]*student_i.StudentInfo, 0),
			Status:      int32(cs.Status),
		},
	}
	for _, t := range cs.Teachers {
		resp.Data.TeacherList = append(resp.Data.TeacherList, &teacher_i.TeacherInfo{
			Id:     uint32(t.ID),
			Name:   t.Name,
			WorkNo: t.WorkNo,
			Email:  t.Email,
		})
	}
	for _, s := range cs.Students {
		resp.Data.StudentList = append(resp.Data.StudentList, &student_i.StudentInfo{
			Id:        uint32(s.ID),
			Name:      s.Name,
			StudentNo: s.StudentNo,
			Sex:       s.Sex,
		})
	}
	return resp, nil
}

func (c CourseService) Edit(ctx context.Context, req *course.GetCourseListReq) (*course.GetCourseListResp, error) {
	//TODO implement me
	panic("implement me")
}

func (c CourseService) List(ctx context.Context, req *course.GetCourseListReq) (*course.GetCourseListResp, error) {
	cs, tot, err := c.courseDal.List(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	rsp := &course.GetCourseListResp{
		Message: "success",
		Data:    make([]*course.CourseList, 0),
		Total:   tot,
	}
	for _, c := range *cs {
		rsp.Data = append(rsp.Data, &course.CourseList{
			Id:          uint32(c.ID),
			CourseName:  c.CourseName,
			ClassName:   c.ClassName,
			TeacherList: make([]*teacher_i.TeacherInfo, 0),
			Status:      int32(c.Status),
		})
		for _, t := range c.Teachers {
			rsp.Data[len(rsp.Data)-1].TeacherList = append(rsp.Data[len(rsp.Data)-1].TeacherList, &teacher_i.TeacherInfo{
				Id:   uint32(t.ID),
				Name: t.Name,
			})
		}
	}
	return rsp, nil
}

func (c CourseService) Import(ctx http.Context) error {
	req := ctx.Request()
	file, _, err := req.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	iLog := c.courseUC.Import(file)
	if iLog == "" {
		iLog = "导入成功"
	}
	ctx.Response().WriteHeader(200)
	resp, _ := json.Marshal(map[string]string{
		"message": iLog,
	})
	_, _ = ctx.Response().Write(resp)
	return nil
}

func NewCourseService(courseDal *dal.CourseDal, courseUC *base_info.CourseUseCase) *CourseService {
	return &CourseService{
		courseDal: courseDal,
		courseUC:  courseUC,
	}
}
