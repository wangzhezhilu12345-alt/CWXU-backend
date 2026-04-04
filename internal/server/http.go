package server

import (
	"edu-evaluation-backed/api/v1/auth"
	"edu-evaluation-backed/api/v1/base_info/course"
	"edu-evaluation-backed/api/v1/base_info/student"
	"edu-evaluation-backed/api/v1/base_info/teacher"
	eva_task2 "edu-evaluation-backed/api/v1/eva_task"
	"edu-evaluation-backed/internal/conf"
	authSvc "edu-evaluation-backed/internal/service/auth"
	"edu-evaluation-backed/internal/service/base_info"
	"edu-evaluation-backed/internal/service/eva_task"
	http2 "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

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
	router.PathPrefix("/res").Handler(http2.StripPrefix("/res", http2.FileServer(http2.Dir("./res"))))
	srv.HandlePrefix("/", router)
	return srv
}
