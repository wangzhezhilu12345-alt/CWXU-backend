package server

import (
	"context"
	"net/http"

	"edu-evaluation-backed/api/v1/auth"
	"edu-evaluation-backed/api/v1/base_info/course"
	"edu-evaluation-backed/api/v1/base_info/student"
	"edu-evaluation-backed/api/v1/base_info/teacher"
	eva_task2 "edu-evaluation-backed/api/v1/eva_task"
	"edu-evaluation-backed/internal/conf"
	authSvc "edu-evaluation-backed/internal/service/auth"
	"edu-evaluation-backed/internal/service/base_info"
	"edu-evaluation-backed/internal/service/eva_task"

	"github.com/go-kratos/kratos/v2/log"
	jwtMiddleware "github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

// jwtSecretKey JWT 签名密钥，需与 auth utility 保持一致
var jwtSecretKey = []byte("edu-evaluation-secret-key")

// whitelistOperations 不需要 JWT 校验的 operation 列表
var whitelistOperations = map[string]bool{
	"/api.v1.auth.Auth/AdminLogin":   true,
	"/api.v1.auth.Auth/StudentLogin": true,
	"/api.v1.eva_task.Task/List":     true,
}

// whitelistMatcher 白名单匹配：白名单 operation 返回 false（不应用 JWT 中间件），其余返回 true
func whitelistMatcher(ctx context.Context, operation string) bool {
	return !whitelistOperations[operation]
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	authService *authSvc.AuthService,
	studentService *base_info.StudentService,
	teacherService *base_info.TeacherService,
	courseService *base_info.CourseService,
	evaService *eva_task.EvaTaskService,
	logger log.Logger,
) *khttp.Server {
	var opts = []khttp.ServerOption{
		khttp.Middleware(
			recovery.Recovery(),
			selector.Server(
				jwtMiddleware.Server(func(token *jwt.Token) (interface{}, error) {
					return jwtSecretKey, nil
				}),
			).Match(whitelistMatcher).Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, khttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, khttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := khttp.NewServer(opts...)
	// 注册认证服务
	auth.RegisterAuthHTTPServer(srv, authService)
	// 上传路由
	b := srv.Route("/api/v1/base-info")
	b.POST("/student/import", studentService.Import)
	b.POST("/teacher/import", teacherService.Import)
	b.POST("/course/import", courseService.Import)
	student_i.RegisterStudentHTTPServer(srv, studentService)
	teacher_i.RegisterTeacherHTTPServer(srv, teacherService)
	course.RegisterCourseHTTPServer(srv, courseService)
	eva_task2.RegisterTaskHTTPServer(srv, evaService)

	// 静态文件下载 - 使用可执行文件所在目录
	router := mux.NewRouter()
	router.PathPrefix("/res").Handler(http.StripPrefix("/res", http.FileServer(http.Dir("./res"))))
	srv.HandlePrefix("/", router)
	return srv
}
