package auth

import (
	"context"
	authApi "edu-evaluation-backed/api/v1/auth"
	authBiz "edu-evaluation-backed/internal/biz/auth"
	authUtil "edu-evaluation-backed/internal/common/utils/auth"
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
	admin, err := s.authUc.AdminLogin(req.Username, req.Password)
	if err != nil {
		return &authApi.AdminLoginResp{Message: err.Error()}, nil
	}
	token, err := authUtil.GenerateToken(map[string]interface{}{
		"userId":   admin.ID,
		"username": admin.Username,
		"role":     "admin",
	})
	if err != nil {
		return &authApi.AdminLoginResp{Message: "token生成失败"}, nil
	}
	return &authApi.AdminLoginResp{Message: "登录成功", Token: token}, nil
}

// StudentLogin 学生登录
func (s *AuthService) StudentLogin(ctx context.Context, req *authApi.StudentLoginReq) (*authApi.StudentLoginResp, error) {
	student, err := s.authUc.StudentLogin(req.StuNo, req.CardNo, uint(req.TaskId))
	if err != nil {
		return &authApi.StudentLoginResp{Message: err.Error()}, nil
	}
	token, err := authUtil.GenerateToken(map[string]interface{}{
		"userId":    student.ID,
		"studentNo": student.StudentNo,
		"role":      "student",
	})
	if err != nil {
		return &authApi.StudentLoginResp{Message: "token生成失败"}, nil
	}
	return &authApi.StudentLoginResp{
		Message: "登录成功",
		Data: &authApi.StudentData{
			StudentNo: student.StudentNo,
			Name:      student.Name,
			Token:     token,
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

// AdminChangePassword 管理员修改密码
func (s *AuthService) AdminChangePassword(ctx context.Context, req *authApi.AdminChangePasswordReq) (*authApi.AdminChangePasswordResp, error) {
	err := s.authUc.AdminChangePassword(req.Username, req.OldPassword, req.NewPassword)
	if err != nil {
		return &authApi.AdminChangePasswordResp{Message: err.Error()}, nil
	}
	return &authApi.AdminChangePasswordResp{Message: "密码修改成功"}, nil
}
