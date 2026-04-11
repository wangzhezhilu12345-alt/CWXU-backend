package dal

import (
	"context"
	"time"

	"edu-evaluation-backed/internal/common/data/cache"
	"edu-evaluation-backed/internal/data"
	"edu-evaluation-backed/internal/data/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// QuestionDal 评教问题数据访问层
type QuestionDal struct {
	db  *gorm.DB
	rdb *redis.Client
	hc  *cache.HealthChecker
}

// ListQuestions 获取评教问题列表（带缓存）
func (d *QuestionDal) ListQuestions() ([]model.EvaluationQuestion, error) {
	ctx := context.Background()
	key := cache.QuestionListKey()

	result, err := cache.Get[[]model.EvaluationQuestion](ctx, d.rdb, d.hc, key, 10*time.Minute, func() (*[]model.EvaluationQuestion, error) {
		var questions []model.EvaluationQuestion
		if err := d.db.Order("sort ASC, id ASC").Find(&questions).Error; err != nil {
			return nil, err
		}
		return &questions, nil
	})
	if err != nil {
		return nil, err
	}
	return *result, nil
}

// UpdateQuestions 全量替换评教问题
func (d *QuestionDal) UpdateQuestions(questions []model.EvaluationQuestion) error {
	err := d.db.Transaction(func(tx *gorm.DB) error {
		// 删除所有旧问题
		if err := tx.Exec("DELETE FROM evaluation_questions").Error; err != nil {
			return err
		}
		// 批量插入新问题
		if len(questions) > 0 {
			if err := tx.Create(&questions).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 异步清除缓存
	go cache.Delete(context.Background(), d.rdb, d.hc, cache.QuestionListKey())
	return nil
}

// NewQuestionDal 创建评教问题数据访问层实例
func NewQuestionDal(data *data.Data) *QuestionDal {
	return &QuestionDal{
		db:  data.DB,
		rdb: data.RDB,
		hc:  data.HC,
	}
}
