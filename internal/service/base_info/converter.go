package base_info

import (
	"encoding/json"

	student_i "edu-evaluation-backed/api/v1/base_info/student"
	teacher_i "edu-evaluation-backed/api/v1/base_info/teacher"
	"edu-evaluation-backed/internal/data/model"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// writeJSONResponse 写入 JSON 响应
func writeJSONResponse(ctx http.Context, message string) {
	ctx.Response().WriteHeader(200)
	resp, _ := json.Marshal(map[string]string{"message": message})
	_, _ = ctx.Response().Write(resp)
}

// ToStudentInfo 将 model.Student 转换为 StudentInfo
func ToStudentInfo(s *model.Student) *student_i.StudentInfo {
	return &student_i.StudentInfo{
		Id:        uint32(s.ID),
		Name:      s.Name,
		Sex:       s.Sex,
		StudentNo: s.StudentNo,
		IdCardNo:  s.IdCardNo,
	}
}

// ToTeacherInfo 将 model.Teacher 转换为 TeacherInfo
func ToTeacherInfo(t *model.Teacher) *teacher_i.TeacherInfo {
	return &teacher_i.TeacherInfo{
		Id:     uint32(t.ID),
		Name:   t.Name,
		Sex:    t.Sex,
		WorkNo: t.WorkNo,
		Email:  t.Email,
	}
}

// ToStudentInfoList 将 []model.Student 转换为 []*student_i.StudentInfo
func ToStudentInfoList(students []model.Student) []*student_i.StudentInfo {
	result := make([]*student_i.StudentInfo, 0, len(students))
	for _, s := range students {
		result = append(result, ToStudentInfo(&s))
	}
	return result
}

// ToTeacherInfoList 将 []model.Teacher 转换为 []*teacher_i.TeacherInfo
func ToTeacherInfoList(teachers []model.Teacher) []*teacher_i.TeacherInfo {
	result := make([]*teacher_i.TeacherInfo, 0, len(teachers))
	for _, t := range teachers {
		result = append(result, ToTeacherInfo(&t))
	}
	return result
}
