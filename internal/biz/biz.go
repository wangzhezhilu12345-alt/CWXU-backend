package biz

import (
	"edu-evaluation-backed/internal/biz/base_info"
	"edu-evaluation-backed/internal/biz/eva_task"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(base_info.NewStudentUseCase, base_info.NewTeacherUseCase, base_info.NewCourseUseCase, eva_task.NewEvaTaskUseCase)
