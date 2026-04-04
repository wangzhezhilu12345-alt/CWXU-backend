package base_info

import (
	"edu-evaluation-backed/internal/data/dal"
	"edu-evaluation-backed/internal/data/model"
	"mime/multipart"

	"github.com/xuri/excelize/v2"
)

// StudentUseCase 学生信息业务用例
type StudentUseCase struct {
	studentDal *dal.BaseInfoDal
}

// NewStudentUseCase 创建学生信息业务用例实例
func NewStudentUseCase(studentDal *dal.BaseInfoDal) *StudentUseCase {
	return &StudentUseCase{studentDal: studentDal}
}

// ImportStudent 从Excel文件导入学生数据
func (s StudentUseCase) ImportStudent(f multipart.File) error {
	list, err := excelize.OpenReader(f)
	if err != nil {
		return err
	}
	defer list.Close()

	rows, err := list.GetRows("Sheet1")
	if err != nil {
		return err
	}

	var batch []*model.Student
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		batch = append(batch, &model.Student{
			StudentNo: row[1],
			Name:      row[2],
			Sex:       row[3],
			IdCardNo:  row[4],
		})
		if len(batch) >= 200 {
			s.studentDal.InsertStudent(batch)
			batch = nil
		}
	}
	if len(batch) > 0 {
		s.studentDal.InsertStudent(batch)
	}
	return nil
}

// UpdateStudent 更新学生信息
func (s StudentUseCase) UpdateStudent(id uint, name, sex, studentNo, idCardNo *string) (*model.Student, error) {
	return s.studentDal.UpdateStudent(id, name, sex, studentNo, idCardNo)
}

// DeleteStudent 删除学生
func (s StudentUseCase) DeleteStudent(id uint) error {
	return s.studentDal.DeleteStudent(id)
}

// GetStudentByID 根据ID获取学生详情
func (s StudentUseCase) GetStudentByID(id uint) (*model.Student, error) {
	return s.studentDal.GetStudentByID(id)
}
