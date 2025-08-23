package biz

import (
	"context"
	"errors"
	v1 "exam_api/api/exam_api/v1"
	_const "exam_api/internal/const"
	"exam_api/internal/data/entity"
	"exam_api/internal/middleware"
	"exam_api/internal/pkg/icontext"
	innErr "exam_api/internal/pkg/ierrors"
	"exam_api/internal/pkg/isnowflake"
	"github.com/go-kratos/kratos/v2/log"
	"math"
	"time"
)

type ExamineeAnswerRepo interface {
	GetByAssociationId(ctx context.Context, associationId string) (resEntity *entity.ExamineeAnswer, err error)
	GetByIDs(ctx context.Context, examineeId string) (list []*entity.ExamineeAnswer, err error)
	Create(ctx context.Context, examineeAnswer *entity.ExamineeAnswer) error
	UpdateAction(ctx context.Context, associationId string, lastActionTime, lastActionTime2 time.Time, remaining int32, completeQuestionNum int32) error
	UpdateResult(ctx context.Context, associationId string, score float64, comparability, usability int32) error
	SubmitResult(ctx context.Context, associationId string, score float64, comparability, usability int32) error
}

type ExamineeAnswerUseCase struct {
	repo          ExamineeAnswerRepo
	associationUc *ExamineeSalesPaperAssociationUseCase
	salesPaperUc  *SalesPaperUseCase
	log           *log.Helper
}

func NewExamineeAnswerUseCase(repo ExamineeAnswerRepo,
	associationUc *ExamineeSalesPaperAssociationUseCase,
	salesPaperUc *SalesPaperUseCase,
	logger log.Logger) *ExamineeAnswerUseCase {
	return &ExamineeAnswerUseCase{
		repo:          repo,
		associationUc: associationUc,
		salesPaperUc:  salesPaperUc,
		log:           log.NewHelper(logger)}
}

func (uc *ExamineeAnswerUseCase) StartExam(ctx context.Context, req *v1.StartExamRequest) (resp *v1.StartExamResponse, err error) {
	l := uc.log.WithContext(ctx)
	userId, _ := icontext.UserIdFrom(ctx)
	accessToken, _ := icontext.UserTokenFrom(ctx)
	association, err := uc.associationUc.GetById(ctx, req.ExamineeAssociationId)
	if err != nil {
		l.Errorf("StartExam.associationUc.GetById Failed, req:%v, err:%v", req, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	if association == nil {
		err = errors.New("该考试不存在")
		return
	}
	//查看试卷数据
	salesPaper, err := uc.salesPaperUc.GetSalesPaperDetail(ctx, association.SalesPaperID)
	if err != nil {
		l.Errorf("StartExam.salesPaperUc.GetSalesPaperDetail Failed, association.SalesPaperID:%v, err:%v", association.SalesPaperID, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	if salesPaper == nil || !salesPaper.IsEnabled {
		err = errors.New("试卷不存在")
	}
	// 判断是否已经有answer数据，也就是说是否开始考试了。
	examineeAnswer, err := uc.repo.GetByAssociationId(ctx, req.ExamineeAssociationId)
	if err != nil {
		l.Errorf("StartExam.repo.GetByAssociationId Failed, req:%v, err:%v", req, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	if examineeAnswer == nil {
		id, e := isnowflake.SnowFlake.NextID(_const.ExamineeAnswerPrefix)
		if e != nil {
			l.Errorf("StartExam.isnowflake.SnowFlake.NextID Failed, req:%v, err:%v", req, e.Error())
			err = innErr.ErrInternalServer
			return
		}
		//第一次进入考试
		examineeAnswer = &entity.ExamineeAnswer{
			ID:                              id,
			SalesPaperID:                    association.SalesPaperID,
			ExamineeID:                      association.ExamineeID,
			ExamineeSalesPaperAssociationID: association.ID,
			Score:                           0,
			BeginTestTime:                   time.Now(),
			SubmitTime:                      nil,
			CompleteQuestionNum:             0,
			LastActionTime:                  nil,
			Comparability:                   0,
			Deadline:                        time.Now().AddDate(0, 0, 3),
			Usability:                       0,
			RemainingTimelimit:              salesPaper.RecommendTimeLim * 60,
			CreatedBy:                       userId,
		}
	}
	//判断是不是过期
	if examineeAnswer.Deadline.Sub(time.Now()) < 0 {
		err = errors.New("考试已过期")
		//设置association表状态为已过期
		_ = uc.associationUc.UpdateStageNumber(ctx, req.ExamineeAssociationId, v1.StageNumber_Expired)
		return
	}
	clientInfo, _ := icontext.UserClientFrom(ctx)
	// 生成jwt
	examJWT, _, err := middleware.JWT.GenerateExamToken(accessToken, req.ExamineeAssociationId, 0, clientInfo)
	if err != nil {
		// 处理错误
		err = errors.New("进入考试失败，请重试")
		return
	}
	resp.ExamToken = examJWT
	resp.TotalDuration = salesPaper.RecommendTimeLim * 60
	resp.Remaining = examineeAnswer.RemainingTimelimit
	resp.UsedDuration = resp.TotalDuration - resp.Remaining
	return
}

func (uc *ExamineeAnswerUseCase) HeartbeatAndSave(ctx context.Context, req *v1.HeartbeatAndSaveRequest) (resp *v1.HeartbeatAndSaveResponse, err error) {
	l := uc.log.WithContext(ctx)
	resp = &v1.HeartbeatAndSaveResponse{}
	// 判断上次登录时间间隔，计算本次耗时
	// 保存数据到examinee_answer中的信息
	// 保存作答信息
	// 心跳结束后可以自动提交

	associationId, _ := icontext.AssociationIdFrom(ctx)
	examineeAnswer, err := uc.repo.GetByAssociationId(ctx, associationId)
	if err != nil {
		l.Errorf("HeartbeatAndSave.repo.GetByAssociationId.NextID Failed, req:%v, err:%v", req, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	// 计算本次耗时（最多算 60 秒）
	thisDuration := 0.0
	activeTime := time.Now()
	if examineeAnswer.LastActionTime != nil && !examineeAnswer.LastActionTime.IsZero() {
		elapsed := activeTime.Sub(*examineeAnswer.LastActionTime).Seconds()
		if elapsed > 300 { // 超过 5 分钟
			//记录事件
			//go logEvent(sess.ID, "long_inactive", c.ClientIP(), map[string]interface{}{
			//	"gap_seconds": int(elapsed),
			//	"last_active": sess.LastActive.Unix(),
			//	"current":     now.Unix(),
			//})
		}
		thisDuration = math.Min(elapsed, 60)
	}
	// 计算本次耗时（最多算 60 秒）
	if examineeAnswer.LastActionTime != nil && !examineeAnswer.LastActionTime.IsZero() {
		elapsed := activeTime.Sub(*examineeAnswer.LastActionTime).Seconds()
		thisDuration = math.Min(elapsed, 60)
	}
	limit := examineeAnswer.RemainingTimelimit - int32(thisDuration)
	if limit <= 0 {
		// todo 自动提交
		resp.Remaining = 0
		return
	}
	// 更新最后活跃时间和剩余时间
	err = uc.repo.UpdateAction(ctx, req.ExamineeAssociationId, activeTime, *examineeAnswer.LastActionTime, limit, int32(len(req.AnswerData)))
	if err != nil {
		l.Errorf("HeartbeatAndSave.repo.GetByAssociationId.NextID Failed, req:%v, err:%v", req, err.Error())
		err = nil
	}
	// todo 保存答案
	resp.Remaining = limit
	return
}
