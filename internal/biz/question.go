package biz

import (
	"context"
	"encoding/json"
	"errors"
	v1 "exam_api/api/exam_api/v1"
	_const "exam_api/internal/const"
	"exam_api/internal/data/entity"
	innErr "exam_api/internal/pkg/ierrors"
	"exam_api/internal/pkg/iutils"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"time"
)

type QuestionRepo interface {
	GetList(ctx context.Context, salesPaperId string) (res []*entity.Question, err error)
	GetPageListBySalesPaperId(ctx context.Context, salesPaperId string) (res []*entity.Question, err error)
	GetListBySalesPaperId(ctx context.Context, salesPaperId string) (res []*entity.Question, err error)
	GetOptionList(ctx context.Context, questionId string) (res []*entity.QuestionOption, err error)
	GetOptionListByQuestionIds(ctx context.Context, questionIds []string) (res map[string][]*entity.QuestionOption, err error)
	GetById(ctx context.Context, questionId string) (qEntity *entity.Question, qOptionsEntities []*entity.QuestionOption, err error)
}

type QuestionUseCase struct {
	repo      QuestionRepo
	redisRepo RedisRepository
	log       *log.Helper
}

func NewQuestionUseCase(repo QuestionRepo, redisRepo RedisRepository, logger log.Logger) *QuestionUseCase {
	return &QuestionUseCase{repo: repo, redisRepo: redisRepo, log: log.NewHelper(logger)}
}

func (uc *QuestionUseCase) ExamQuestion(ctx context.Context, salesPaperId string) (resp *v1.ExamQuestionResponse, err error) {

	resp = &v1.ExamQuestionResponse{QuestionData: make([]*v1.QuestionData, 0)}
	l := uc.log.WithContext(ctx)
	// 先查询缓存
	key := fmt.Sprintf(_const.GetQuestionsBySalesPaperIdRedisKey, salesPaperId)
	data, e := uc.redisRepo.Get(ctx, key)
	if e != nil && !errors.Is(err, redis.Nil) {
		l.Errorf("GetQuestionBySalesPaperId.redisRepo.Get Failed, key:%v, err:%v", key, e.Error())
	}
	if data != "" {
		if e = json.Unmarshal([]byte(data), &resp.QuestionData); e == nil {
			return
		}
	}
	// 缓存没有则查询数据库
	res, err := uc.repo.GetPageListBySalesPaperId(ctx, salesPaperId)
	if err != nil {
		l.Errorf("GetQuestionBySalesPaperId.repo.GetPageListBySalesPaperId Failed, salesPaperId:%v, err:%v", salesPaperId, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	questionIds := make([]string, 0, len(res))
	for _, re := range res {
		questionIds = append(questionIds, re.ID)
	}
	mQuestionOptions, err := uc.repo.GetOptionListByQuestionIds(ctx, questionIds)
	if err != nil {
		l.Errorf("GetQuestionBySalesPaperId.repo.GetOptionListByQuestionIds Failed, questionIds:%v, err:%v", questionIds, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	for _, re := range res {
		cur := &v1.QuestionData{
			QuestionId:     re.ID,
			Title:          re.Title,
			QuestionTypeId: v1.QuestionType(re.QuestionTypeID),
			Order:          re.Order_,
		}
		if v, ok := mQuestionOptions[cur.QuestionId]; ok {
			for _, option := range v {
				cur.QuestionOptionsData = append(cur.QuestionOptionsData, &v1.QuestionOptionData{
					QuestionOptionId: option.ID,
					Description:      option.Description,
					SerialNumber:     iutils.OrderToLetter(option.Order_),
				})
			}
		}
		resp.QuestionData = append(resp.QuestionData, cur)
	}
	value, _ := json.Marshal(resp.QuestionData)
	e = uc.redisRepo.Set(ctx, key, string(value), time.Duration(3)*time.Minute)
	if e != nil {
		l.Errorf("GetQuestionBySalesPaperId.repo.redisRepo Set Failed, key:%v, err:%v", key, e.Error())
	}
	return
}
