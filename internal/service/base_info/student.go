package base_info

import (
	"context"

	baseinfo2 "edu-evaluation-backed/api/v1/base_info/student"
	"edu-evaluation-backed/internal/biz/base_info"
	"edu-evaluation-backed/internal/data/dal"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// StudentService 学生信息服务
type StudentService struct {
	studentUc *base_info.StudentUseCase
	baseDal   *dal.BaseInfoDal
}

// List 获取学生列表
func (s StudentService) List(ctx context.Context, req *baseinfo2.GetStudentListReq) (*baseinfo2.GetStudentListResp, error) {
	students, tot, err := s.baseDal.QueryStudent(int(req.Page), int(req.PageSize), req.StudentNo, req.Name)
	if err != nil {
		return nil, err
	}
	return &baseinfo2.GetStudentListResp{
		Message: "success",
		Data:    ToStudentInfoList(*students),
		Total:   tot,
	}, nil
}

// Import 导入学生信息Excel文件
func (s StudentService) Import(ctx http.Context) error {
	req := ctx.Request()
	file, _, err := req.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	if err := s.studentUc.ImportStudent(file); err != nil {
		return err
	}
	writeJSONResponse(ctx, "导入成功")
	return nil
}

// Update 更新学生信息
func (s StudentService) Update(ctx context.Context, req *baseinfo2.UpdateStudentReq) (*baseinfo2.UpdateStudentResp, error) {
	student, err := s.studentUc.UpdateStudent(uint(req.Id), req.Name, req.Sex, req.StudentNo, req.IdCardNo)
	if err != nil {
		return nil, err
	}
	return &baseinfo2.UpdateStudentResp{
		Message: "修改成功",
		Data:    ToStudentInfo(student),
	}, nil
}

// Delete 删除学生
func (s StudentService) Delete(ctx context.Context, req *baseinfo2.DeleteStudentReq) (*baseinfo2.DeleteStudentResp, error) {
	if err := s.studentUc.DeleteStudent(uint(req.Id)); err != nil {
		return nil, err
	}
	return &baseinfo2.DeleteStudentResp{Message: "删除成功"}, nil
}

// NewStudentService 创建学生信息服务实例
func NewStudentService(studentUc *base_info.StudentUseCase, baseDal *dal.BaseInfoDal) *StudentService {
	return &StudentService{studentUc: studentUc, baseDal: baseDal}
}
