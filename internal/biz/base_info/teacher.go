package base_info

import (
	"edu-evaluation-backed/internal/data/dal"
	"edu-evaluation-backed/internal/data/model"
	"mime/multipart"

	"github.com/xuri/excelize/v2"
)

// TeacherUseCase 教师信息业务用例
type TeacherUseCase struct {
	baseDal *dal.BaseInfoDal
}

// NewTeacherUseCase 创建教师信息业务用例实例
func NewTeacherUseCase(baseDal *dal.BaseInfoDal) *TeacherUseCase {
	return &TeacherUseCase{baseDal: baseDal}
}

// Import 从Excel文件导入教师数据
func (s TeacherUseCase) Import(f multipart.File) error {
	list, err := excelize.OpenReader(f)
	if err != nil {
		return err
	}
	defer list.Close()

	rows, err := list.GetRows("Sheet1")
	if err != nil {
		return err
	}

	var batch []*model.Teacher
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		batch = append(batch, &model.Teacher{
			WorkNo: row[1],
			Name:   row[2],
			Sex:    row[3],
			Email:  row[4],
		})
		if len(batch) >= 200 {
			s.baseDal.InsertTeacher(batch)
			batch = nil
		}
	}
	if len(batch) > 0 {
		s.baseDal.InsertTeacher(batch)
	}
	return nil
}

// UpdateTeacher 更新教师信息
func (s TeacherUseCase) UpdateTeacher(id uint, name, sex, workNo, email *string) (*model.Teacher, error) {
	return s.baseDal.UpdateTeacher(id, name, sex, workNo, email)
}

// DeleteTeacher 删除教师
func (s TeacherUseCase) DeleteTeacher(id uint) error {
	return s.baseDal.DeleteTeacher(id)
}

// GetTeacherByID 根据ID获取教师详情
func (s TeacherUseCase) GetTeacherByID(id uint) (*model.Teacher, error) {
	return s.baseDal.GetTeacherByID(id)
}
