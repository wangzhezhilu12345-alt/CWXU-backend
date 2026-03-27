package auth

import (
	"context"
	authApi "edu-evaluation-backed/api/v1/auth"
	authBiz "edu-evaluation-backed/internal/biz/auth"
)

// AuthService 认证服务
type AuthService struct {
	authUc *authBiz.AuthUseCase
}

// NewAuthService 创建认证服务实例
func NewAuthService(authUc *authBiz.AuthUseCase) *AuthService {
	return &AuthService{
		authUc: authUc,
	}
}

// AdminLogin 管理员登录
func (s *AuthService) AdminLogin(ctx context.Context, req *authApi.AdminLoginReq) (*authApi.AdminLoginResp, error) {
	_, err := s.authUc.AdminLogin(req.Username, req.Password)
	if err != nil {
		return &authApi.AdminLoginResp{Message: err.Error()}, nil
	}
	return &authApi.AdminLoginResp{Message: "登录成功"}, nil
}

// StudentLogin 学生登录
func (s *AuthService) StudentLogin(ctx context.Context, req *authApi.StudentLoginReq) (*authApi.StudentLoginResp, error) {
	student, err := s.authUc.StudentLogin(req.StuNo, req.CardNo, uint(req.TaskId))
	if err != nil {
		return &authApi.StudentLoginResp{Message: err.Error()}, nil
	}
	return &authApi.StudentLoginResp{
		Message: "登录成功",
		Data: &authApi.StudentData{
			StudentNo: student.StudentNo,
			Name:      student.Name,
		},
	}, nil
}

// StudentInfo 获取学生个人信息
func (s *AuthService) StudentInfo(ctx context.Context, req *authApi.StudentInfoReq) (*authApi.StudentInfoResp, error) {
	student, err := s.authUc.GetStudentInfo(req.StuNo)
	if err != nil {
		return &authApi.StudentInfoResp{Message: err.Error()}, nil
	}
	return &authApi.StudentInfoResp{
		Message: "success",
		Data: &authApi.StudentDetailInfo{
			Id:        uint32(student.ID),
			Name:      student.Name,
			Sex:       student.Sex,
			StudentNo: student.StudentNo,
			IdCardNo:  student.IdCardNo,
		},
	}, nil
}
