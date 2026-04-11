package eva_question

import (
	"edu-evaluation-backed/internal/data/dal"
	"edu-evaluation-backed/internal/data/model"
)

// QuestionUseCase 评教问题业务用例
type QuestionUseCase struct {
	questionDal *dal.QuestionDal
}

// ListQuestions 获取评教问题列表
func (uc *QuestionUseCase) ListQuestions() ([]model.EvaluationQuestion, error) {
	return uc.questionDal.ListQuestions()
}

// UpdateQuestions 修改评教问题列表
func (uc *QuestionUseCase) UpdateQuestions(questions []model.EvaluationQuestion) error {
	return uc.questionDal.UpdateQuestions(questions)
}

// NewQuestionUseCase 创建评教问题业务用例实例
func NewQuestionUseCase(questionDal *dal.QuestionDal) *QuestionUseCase {
	return &QuestionUseCase{
		questionDal: questionDal,
	}
}
