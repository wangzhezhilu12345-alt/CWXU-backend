package eva_question

import (
	"context"

	eva_question2 "edu-evaluation-backed/api/v1/eva_question"
	"edu-evaluation-backed/internal/biz/eva_question"
	"edu-evaluation-backed/internal/data/model"
)

// QuestionService 评教问题服务
type QuestionService struct {
	uc *eva_question.QuestionUseCase
}

// ListQuestions 获取评教问题列表
func (s *QuestionService) ListQuestions(ctx context.Context, req *eva_question2.ListQuestionsReq) (*eva_question2.ListQuestionsResp, error) {
	questions, err := s.uc.ListQuestions()
	if err != nil {
		return nil, err
	}

	items := make([]*eva_question2.QuestionItem, 0, len(questions))
	for _, q := range questions {
		items = append(items, &eva_question2.QuestionItem{
			Id:      uint32(q.ID),
			Content: q.Content,
			Score:   int32(q.Score),
			Sort:    int32(q.Sort),
		})
	}

	return &eva_question2.ListQuestionsResp{
		Message: "获取成功",
		Data:    items,
	}, nil
}

// UpdateQuestions 修改评教问题列表
func (s *QuestionService) UpdateQuestions(ctx context.Context, req *eva_question2.UpdateQuestionsReq) (*eva_question2.UpdateQuestionsResp, error) {
	questions := make([]model.EvaluationQuestion, 0, len(req.Questions))
	for _, q := range req.Questions {
		questions = append(questions, model.EvaluationQuestion{
			Content: q.Content,
			Score:   int(q.Score),
			Sort:    int(q.Sort),
		})
	}

	if err := s.uc.UpdateQuestions(questions); err != nil {
		return nil, err
	}

	return &eva_question2.UpdateQuestionsResp{
		Message: "更新成功",
	}, nil
}

// NewQuestionService 创建评教问题服务实例
func NewQuestionService(uc *eva_question.QuestionUseCase) *QuestionService {
	return &QuestionService{
		uc: uc,
	}
}
