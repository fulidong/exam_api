package biz

import (
	"context"
	"encoding/json"
	_const "exam_api/internal/const"
	"exam_api/internal/data/entity"
	"exam_api/internal/pkg/icontext"
	innErr "exam_api/internal/pkg/ierrors"
	"exam_api/internal/pkg/isnowflake"
	"github.com/go-kratos/kratos/v2/log"
)

type ExamEventRepo interface {
	ExamEvent(ctx context.Context, examEvent *entity.ExamEvent) error
}

type ExamEventUseCase struct {
	repo ExamEventRepo
	log  *log.Helper
}

func NewExamEventUseCase(repo ExamEventRepo, logger log.Logger) *ExamEventUseCase {
	return &ExamEventUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *ExamEventUseCase) ExamEvent(ctx context.Context, examineeAnswerId string, eventType _const.ExamEventType, meta map[string]interface{}) (err error) {
	l := uc.log.WithContext(ctx)
	id, _ := isnowflake.SnowFlake.NextID(_const.ExamiEventPrefix)
	sessionId, _ := icontext.SessionIdFrom(ctx)
	examToken, _ := icontext.ExamTokenFrom(ctx)
	client, _ := icontext.UserClientFrom(ctx)
	userId, _ := icontext.UserIdFrom(ctx)
	userAgent, _ := icontext.UserAgentFrom(ctx)
	metaStr, _ := json.Marshal(meta)
	examEvent := &entity.ExamEvent{
		ID:               id,
		SessionID:        sessionId,
		ExamToken:        examToken,
		ExamineeAnswerID: examineeAnswerId,
		EventType:        string(eventType),
		IP:               client.IP,
		UserAgent:        userAgent,
		Meta:             string(metaStr),
		CreatedBy:        userId,
	}
	err = uc.repo.ExamEvent(ctx, examEvent)
	if err != nil {
		data, _ := json.Marshal(examEvent)
		l.Errorf("ExamEvent.repo.ExamEvent Failed, examEvent:%v, err:%v", string(data), err.Error())
		err = innErr.ErrInternalServer
		return
	}
	return
}
