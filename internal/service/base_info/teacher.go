package base_info

import (
	"context"

	teacher_i "edu-evaluation-backed/api/v1/base_info/teacher"
	"edu-evaluation-backed/internal/biz/base_info"
	"edu-evaluation-backed/internal/data/dal"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// TeacherService 教师信息服务
type TeacherService struct {
	teacherUc *base_info.TeacherUseCase
	baseDal   *dal.BaseInfoDal
}

// List 获取教师列表
func (s TeacherService) List(ctx context.Context, req *teacher_i.GetTeacherListReq) (*teacher_i.GetTeacherListResp, error) {
	teachers, tot, err := s.baseDal.QueryTeacher(int(req.Page), int(req.PageSize), req.WorkNo, req.Name)
	if err != nil {
		return nil, err
	}
	return &teacher_i.GetTeacherListResp{
		Message: "success",
		Data:    ToTeacherInfoList(*teachers),
		Total:   tot,
	}, nil
}

// Import 导入教师信息Excel文件
func (s TeacherService) Import(ctx http.Context) error {
	req := ctx.Request()
	file, _, err := req.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	if err := s.teacherUc.Import(file); err != nil {
		return err
	}
	writeJSONResponse(ctx, "导入成功")
	return nil
}

// Update 更新教师信息
func (s TeacherService) Update(ctx context.Context, req *teacher_i.UpdateTeacherReq) (*teacher_i.UpdateTeacherResp, error) {
	teacher, err := s.teacherUc.UpdateTeacher(uint(req.Id), req.Name, req.Sex, req.WorkNo, req.Email)
	if err != nil {
		return nil, err
	}
	return &teacher_i.UpdateTeacherResp{
		Message: "修改成功",
		Data:    ToTeacherInfo(teacher),
	}, nil
}

// Delete 删除教师
func (s TeacherService) Delete(ctx context.Context, req *teacher_i.DeleteTeacherReq) (*teacher_i.DeleteTeacherResp, error) {
	if err := s.teacherUc.DeleteTeacher(uint(req.Id)); err != nil {
		return nil, err
	}
	return &teacher_i.DeleteTeacherResp{Message: "删除成功"}, nil
}

// NewTeacherService 创建教师信息服务实例
func NewTeacherService(teacherUc *base_info.TeacherUseCase, baseDal *dal.BaseInfoDal) *TeacherService {
	return &TeacherService{teacherUc: teacherUc, baseDal: baseDal}
}
