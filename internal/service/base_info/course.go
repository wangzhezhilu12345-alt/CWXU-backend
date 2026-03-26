package base_info

import (
	"edu-evaluation-backed/internal/biz/base_info"
	"edu-evaluation-backed/internal/data/dal"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/transport/http"
)

type CourseService struct {
	courseDal *dal.CourseDal
	courseUC  *base_info.CourseUseCase
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
