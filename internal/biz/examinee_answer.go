package biz

import (
	"context"
	"encoding/json"
	"errors"
	v1 "exam_api/api/exam_api/v1"
	_const "exam_api/internal/const"
	"exam_api/internal/data/entity"
	"exam_api/internal/middleware"
	"exam_api/internal/pkg/icontext"
	innErr "exam_api/internal/pkg/ierrors"
	"exam_api/internal/pkg/isnowflake"
	"exam_api/internal/pkg/itask"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"math"
	"time"
)

type ExamineeAnswerRepo interface {
	GetByAssociationId(ctx context.Context, associationId string) (resEntity *entity.ExamineeAnswer, err error)
	GetByIDs(ctx context.Context, examineeId string) (list []*entity.ExamineeAnswer, err error)
	Create(ctx context.Context, examineeAnswer *entity.ExamineeAnswer) error
	UpdateAction(ctx context.Context, examineeAnswerId string, lastActionTime, lastActionTime2 time.Time, remaining int32, completeQuestionNum int32) (int64, error)
	UpdateResult(ctx context.Context, examineeAnswerId string, score float64, comparability, usability int32) error
	SubmitResult(ctx context.Context, examineeAnswerId string, score float64, comparability, usability int32) error
}

type ExamineeAnswerUseCase struct {
	repo                     ExamineeAnswerRepo
	associationUc            *ExamineeSalesPaperAssociationUseCase
	salesPaperUc             *SalesPaperUseCase
	examineeQuestionAnswerUC *ExamineeQuestionAnswerUseCase
	examEvent                *ExamEventUseCase
	redisRepo                RedisRepository
	log                      *log.Helper
}

func NewExamineeAnswerUseCase(repo ExamineeAnswerRepo,
	associationUc *ExamineeSalesPaperAssociationUseCase,
	salesPaperUc *SalesPaperUseCase,
	examineeQuestionAnswerUC *ExamineeQuestionAnswerUseCase,
	examEvent *ExamEventUseCase,
	redisRepo RedisRepository,
	logger log.Logger) *ExamineeAnswerUseCase {
	return &ExamineeAnswerUseCase{
		repo:                     repo,
		associationUc:            associationUc,
		salesPaperUc:             salesPaperUc,
		examineeQuestionAnswerUC: examineeQuestionAnswerUC,
		examEvent:                examEvent,
		redisRepo:                redisRepo,
		log:                      log.NewHelper(logger)}
}

func (uc *ExamineeAnswerUseCase) StartExam(ctx context.Context, req *v1.StartExamRequest) (resp *v1.StartExamResponse, err error) {
	l := uc.log.WithContext(ctx)
	resp = &v1.StartExamResponse{}
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
		return
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
		curTime := time.Now()
		examineeAnswer = &entity.ExamineeAnswer{
			ID:                              id,
			SalesPaperID:                    association.SalesPaperID,
			ExamineeID:                      association.ExamineeID,
			ExamineeSalesPaperAssociationID: association.ID,
			Score:                           0,
			BeginTestTime:                   curTime,
			LastActionTime:                  curTime,
			SubmitTime:                      nil,
			CompleteQuestionNum:             0,
			Comparability:                   0,
			Deadline:                        curTime.AddDate(0, 0, 3),
			Usability:                       0,
			RemainingTimelimit:              salesPaper.RecommendTimeLim * 60,
			CreatedBy:                       userId,
		}
		e = uc.repo.Create(ctx, examineeAnswer)
		if e != nil {
			l.Errorf("StartExam.repo.Create Failed, req:%v, err:%v", req, e.Error())
			err = innErr.ErrInternalServer
			return
		}
		// 将状态更新成进行中
		e = uc.associationUc.UpdateStageNumber(ctx, req.ExamineeAssociationId, v1.StageNumber_InProgress)
		if e != nil {
			l.Errorf("StartExam.associationUc.UpdateStageNumber Failed, req:%v, stage:%v, err:%v", req, v1.StageNumber_InProgress, e.Error())
			err = innErr.ErrInternalServer
			return
		}
	}
	if examineeAnswer.Deadline.Sub(time.Now()) < 0 {
		//已经过期， 设置状态
		err = uc.associationUc.UpdateStageNumber(ctx, req.ExamineeAssociationId, v1.StageNumber_Expire)
		if err != nil {
			l.Errorf("StartExam.associationUc.UpdateStageNumber Failed, req:%v, stage:%v, err:%v", req, v1.StageNumber_Expire, err.Error())
		}
		err = errors.New("该考试已过截止时间")
		return
	}
	clientInfo, _ := icontext.UserClientFrom(ctx)
	// 生成jwt
	examJWT, _, err := middleware.JWT.GenerateExamToken(accessToken, association.ID, time.Duration(examineeAnswer.RemainingTimelimit*60)*time.Second, clientInfo)
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
	resp = &v1.HeartbeatAndSaveResponse{}
	var (
		l                = uc.log.WithContext(ctx)
		userId, _        = icontext.UserIdFrom(ctx)
		associationId, _ = icontext.AssociationIdFrom(ctx)
	)
	// 1. 获取答题信息
	examineeAnswer, err := uc.repo.GetByAssociationId(ctx, associationId)
	if err != nil {
		l.Errorf("HeartbeatAndSave.repo.GetByAssociationId.NextID Failed, req:%v, err:%v", req, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	// 2. 防止频繁心跳
	thisDuration := 0.0
	activeTime := time.Now()
	if !examineeAnswer.LastActionTime.IsZero() {
		gap := activeTime.Sub(examineeAnswer.LastActionTime)
		if gap < 5*time.Second { // 小于 5 秒
			err = innErr.ErrHeartbeat
			return
		}
	}
	// 3. 判断是否“重新进入考试”（上次心跳超过 5 分钟），记录事件
	if !examineeAnswer.LastActionTime.IsZero() {
		elapsed := activeTime.Sub(examineeAnswer.LastActionTime).Seconds()
		if elapsed > 300 { // 超过 5 分钟
			//记录事件
			meta := map[string]interface{}{
				"gap_seconds": int(elapsed),
				"last_active": examineeAnswer.LastActionTime.Unix(),
				"current":     activeTime,
			}

			go itask.TaskWithContext(ctx, func() {
				if e := uc.examEvent.ExamEvent(ctx, examineeAnswer.ID, _const.ExamEventLongInactive, meta); e != nil {
					l.Errorf("HeartbeatAndSave.examEvent.ExamEvent Failed, req:%v, err:%v", req, e.Error())
				}
			}, l)
		}
		// 计算本次耗时（最多算 30 秒）
		thisDuration = math.Min(elapsed, 30)
	}
	// 4. 检查考试时间是否已经用完
	limit := examineeAnswer.RemainingTimelimit - int32(thisDuration)
	if limit <= 0 {
		// todo 自动提交
		resp.Remaining = 0
		return
	}
	// 5. 更新最后活跃时间和剩余时间(乐观锁)
	rowsAffected, err := uc.repo.UpdateAction(ctx, examineeAnswer.ID, activeTime, examineeAnswer.LastActionTime, limit, int32(len(req.AnswerData)))
	if err != nil {
		l.Errorf("HeartbeatAndSave.repo.UpdateAction Failed, req:%v, err:%v", req, err.Error())
		err = nil
	}
	if rowsAffected == 0 {
		// 6. 说明有其他请求抢先更新了 LastActionTime
		// 降级：查最新状态，但不再累加时间（防重复计时）
		examineeAnswer, e := uc.repo.GetByAssociationId(ctx, associationId)
		if e != nil {
			l.Errorf("HeartbeatAndSave.repo.GetByAssociationId.NextID Failed, req:%v, err:%v", req, e.Error())
			err = innErr.ErrInternalServer
			return
		}
		limit = examineeAnswer.RemainingTimelimit
	}
	// 7. 保存答案
	if len(req.AnswerData) > 0 {
		answers := make([]*entity.ExamineeAnswerQuestionAnswer, 0, len(req.AnswerData))
		for _, questionAnswerData := range req.AnswerData {
			id, _ := isnowflake.SnowFlake.NextID(_const.ExamineeAnswerQuestionAnswerPrefix)
			sign, _ := json.Marshal(questionAnswerData.OptionsSerialNumberData)
			answers = append(answers, &entity.ExamineeAnswerQuestionAnswer{
				ID:               id,
				ExamineeAnswerID: examineeAnswer.ID,
				QuestionID:       questionAnswerData.QuestionId,
				Score:            0,
				OptionSign:       string(sign),
				CreatedBy:        userId,
				UpdatedBy:        userId,
			})
		}
		err = uc.examineeQuestionAnswerUC.SaveAnswer(ctx, answers)
		if err != nil {
			l.Errorf("HeartbeatAndSave.examineeQuestionAnswerUC.SaveAnswer Failed, req:%v, err:%v", req, err.Error())
			err = innErr.ErrInternalServer
			return
		}
	}
	// 8. 记录心跳事件
	go itask.TaskWithContext(ctx, func() {
		if e := uc.examEvent.ExamEvent(ctx, examineeAnswer.ID, _const.ExamEventHeartbeat, make(map[string]interface{})); e != nil {
			l.Errorf("HeartbeatAndSave.examEvent.ExamEvent Failed, req:%v, err:%v", req, e.Error())
		}
	}, l)
	resp.Remaining = limit
	return
}

func (uc *ExamineeAnswerUseCase) SubmitExam(ctx context.Context, req *v1.SubmitExamRequest) (resp *v1.SubmitExamResponse, err error) {
	resp = &v1.SubmitExamResponse{}
	var (
		l                = uc.log.WithContext(ctx)
		userId, _        = icontext.UserIdFrom(ctx)
		associationId, _ = icontext.AssociationIdFrom(ctx)
		thisDuration     = 0.0
		activeTime       = time.Now()
		lockExpire       = 5 * time.Second // 锁过期时间（防死锁）
		lockValue, _     = isnowflake.SnowFlake.NextID("lock")
		submitKeyExpire  = 4 * time.Hour // 提交标记过期时间
	)
	lockKey := fmt.Sprintf(_const.RedisLockKey, associationId)
	submitKey := fmt.Sprintf(_const.RedisSubmitKey, associationId)

	// 1. 第一重校验：检查是否已提交（Redis 快速失败）
	isSubmitted, _ := uc.redisRepo.Exists(ctx, submitKey)
	if isSubmitted {
		err = errors.New("试卷已提交，请勿重复操作")
		return
	}
	// 3. 获取分布式锁（5s 过期）
	ok, err := uc.redisRepo.SetNX(ctx, lockKey, lockValue, lockExpire)
	if err != nil {
		l.Errorf("SubmitExam.redisRepo.SetNX Failed, req:%v, err:%v", req, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	if !ok {
		err = errors.New("系统繁忙，请稍后重试")
		return
	}

	// 4. 延迟释放锁（用 defer）
	defer func() {
		// 直接使用 Eval
		result, err := uc.redisRepo.Eval(ctx, _const.UnlockScript, []string{lockKey}, lockValue)
		if err != nil {
			l.Errorf("释放锁失败: %v", err)
		} else if n, ok := result.(int64); !ok || n == 0 {
			l.Errorf("未释放锁（可能已过期或被其他协程持有）: %s", lockKey)
		}
	}()
	examineeAnswer, err := uc.repo.GetByAssociationId(ctx, associationId)
	if err != nil {
		l.Errorf("SubmitExam.repo.GetByAssociationId Failed, associationId:%v, err:%v ", associationId, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	// 5. 第二重校验：查询数据库最新状态（防并发）
	association, err := uc.associationUc.GetById(ctx, associationId)
	if err != nil {
		l.Errorf("SubmitExam.associationUc.GetById Failed, associationId:%v, err:%v ", associationId, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	if association == nil || association.StageNumber != 1 {
		err = errors.New("考试记录不存在")
		return
	}
	if association.StageNumber == int32(v1.StageNumber_Submit) {
		err = errors.New("试卷已提交，请勿重复操作")
		return
	}
	if association.StageNumber != int32(v1.StageNumber_InProgress) {
		err = errors.New("考试状态异常")
		return
	}
	// 6. 更新最后活跃时间和作答题目数量
	if !examineeAnswer.LastActionTime.IsZero() {
		elapsed := activeTime.Sub(examineeAnswer.LastActionTime).Seconds()
		// 计算本次耗时（最多算 30 秒）
		thisDuration = math.Min(elapsed, 30)
	}
	// 7. 检查考试时间是否已经用完
	limit := examineeAnswer.RemainingTimelimit - int32(thisDuration)
	if limit <= 0 {
		limit = 0
	}
	// 8. 更新最后活跃时间和剩余时间
	_, err = uc.repo.UpdateAction(ctx, examineeAnswer.ID, activeTime, examineeAnswer.LastActionTime, limit, int32(len(req.AnswerData)))
	if err != nil {
		l.Errorf("HeartbeatAndSave.repo.UpdateAction Failed, req:%v, err:%v", req, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	// 9. 保存答案，更新状态
	if len(req.AnswerData) > 0 {
		answers := make([]*entity.ExamineeAnswerQuestionAnswer, 0, len(req.AnswerData))
		for _, questionAnswerData := range req.AnswerData {
			id, _ := isnowflake.SnowFlake.NextID(_const.ExamineeAnswerQuestionAnswerPrefix)
			sign, _ := json.Marshal(questionAnswerData.OptionsSerialNumberData)
			answers = append(answers, &entity.ExamineeAnswerQuestionAnswer{
				ID:               id,
				ExamineeAnswerID: examineeAnswer.ID,
				QuestionID:       questionAnswerData.QuestionId,
				Score:            0,
				OptionSign:       string(sign),
				CreatedBy:        userId,
				UpdatedBy:        userId,
			})
		}
		err = uc.examineeQuestionAnswerUC.SaveAnswer(ctx, answers)
		if err != nil {
			l.Errorf("HeartbeatAndSave.examineeQuestionAnswerUC.SaveAnswer Failed, req:%v, err:%v", req, err.Error())
			err = innErr.ErrInternalServer
			return
		}
	}
	// 10. 更新状态
	err = uc.associationUc.UpdateStageNumber(ctx, associationId, v1.StageNumber_Submit)
	if err != nil {
		l.Errorf("HeartbeatAndSave.examineeQuestionAnswerUC.SaveAnswer Failed, req:%v, err:%v", req, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	// 11. 添加提交成功key，防止重放
	_ = uc.redisRepo.Set(ctx, submitKey, "", submitKeyExpire)
	// 12. 记录心跳事件
	go itask.TaskWithContext(ctx, func() {
		if e := uc.examEvent.ExamEvent(ctx, examineeAnswer.ID, _const.ExamEventSubmit, make(map[string]interface{})); e != nil {
			l.Errorf("HeartbeatAndSave.examEvent.ExamEvent Failed, req:%v, err:%v", req, e.Error())
		}
	}, l)
	return
}

func (uc *ExamineeAnswerUseCase) ExamQuestionRecord(ctx context.Context, req *v1.ExamQuestionRecordRequest) (resp *v1.ExamQuestionRecordResponse, err error) {

	resp = &v1.ExamQuestionRecordResponse{AnswerData: make([]*v1.QuestionAnswerData, 0)}
	var (
		l                = uc.log.WithContext(ctx)
		associationId, _ = icontext.AssociationIdFrom(ctx)
	)
	examineeAnswer, err := uc.repo.GetByAssociationId(ctx, associationId)
	if err != nil {
		l.Errorf("ExamQuestionRecord.repo.GetByAssociationId Failed, associationId:%v, err:%v", associationId, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	if examineeAnswer == nil {
		err = errors.New("考试不存在")
		return
	}
	//查看试卷数据
	examineeAnswers, err := uc.examineeQuestionAnswerUC.GetByExamineeAnswerId(ctx, examineeAnswer.ID)
	if err != nil {
		l.Errorf("ExamQuestionRecord.examineeQuestionAnswerUC.GetByExamineeAnswerId Failed, examineeAnswer.ID:%v, err:%v", examineeAnswer.ID, err.Error())
		err = innErr.ErrInternalServer
		return
	}
	if len(examineeAnswers) == 0 {
		return
	}
	resp.AnswerData = make([]*v1.QuestionAnswerData, 0, len(examineeAnswers))
	for _, answer := range examineeAnswers {
		options := make([]string, 0)
		if answer.OptionSign != "" {
			json.Unmarshal([]byte(answer.OptionSign), &options)
		}
		resp.AnswerData = append(resp.AnswerData, &v1.QuestionAnswerData{
			QuestionId:              answer.QuestionID,
			OptionsSerialNumberData: options,
		})
	}
	return
}
