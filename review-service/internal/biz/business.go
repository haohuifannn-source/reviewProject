package biz

import (
	"context"
	"errors"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// BusinessRepo is a Greater repo.
type BusinessRepo interface {
	GetReviewByReviewId(context.Context, int64) (*model.ReviewInfo, error)
	ReplyReview(context.Context, *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error)
	TaskToReply(context.Context, *model.ReviewReplyInfo) error

	CreateAppealReview(context.Context, *model.ReviewAppealInfo) error
	GetAppealReviewByReviewId(context.Context, int64) (*model.ReviewAppealInfo, error)
}

// BusinessUsecase is a Business usecase.
type BusinessUsecase struct {
	repo BusinessRepo
	log  *log.Helper
}

// NewBusinessUsecase new a Business usecase.
func NewBusinessUsecase(repo BusinessRepo, logger log.Logger) *BusinessUsecase {
	return &BusinessUsecase{repo: repo, log: log.NewHelper(logger)}
}

// ReplyReview 商家回复评价的接口
func (uc *BusinessUsecase) ReplyReview(ctx context.Context, reviewRep *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
	// 1. 数据校验
	// 1.1 数据合法性校验，不允许重复回复
	uc.log.Infof("开始回复评价")
	rreplys, err := uc.repo.GetReviewByReviewId(ctx, reviewRep.ReviewID)
	if err != nil {
		return nil, errors.New("数据库查询错误")
	}
	if rreplys.HasReply == 1 {
		return nil, errors.New("已经回复评价了")
	}
	// 1.2 水平越权校验（A商家不可以给B商家回复）
	// 举例子： 用户A删除订单，uerID + orderID 当条件去查询订单然后删除
	if reviewRep.StoreID != rreplys.StoreID {
		return nil, errors.New("水平越权")
	}
	// 2. 生成ReplyID
	rrID := snowflake.GenerateReviewID()
	reviewRep.ReplyID = rrID
	// 3. 更新数据(评价回复表和评价表都需要更新，因为评价表里面有一个是否回复的字段，涉及到一个事务操作)
	if err := uc.repo.TaskToReply(ctx, reviewRep); err != nil {
		return nil, err
	}
	// 4. 直接返回结果
	return reviewRep, nil
}

// ApealReview 商家申述用户评价的接口
func (uc *BusinessUsecase) ApealReview(ctx context.Context, reviewRep *model.ReviewAppealInfo) (*model.ReviewAppealInfo, error) {
	// 1. 查询是否有这个评价
	ret, err := uc.repo.GetReviewByReviewId(ctx, reviewRep.ReviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("评价不存在")
		}
		return nil, errors.New("数据库查询错误")
	}
	if ret.StoreID != reviewRep.StoreID {
		return nil, errors.New("越权操作：您无权申诉非本门店的评价")
	}
	// 2. 查询是否有申述的记录
	aret, err := uc.repo.GetAppealReviewByReviewId(ctx, reviewRep.ReviewID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("数据库查询错误")
	}
	if err == nil && aret.Status != 10 {
		return nil, errors.New("已经申述完成，不可以修改/重申")
	}
	// 4. 生成AppealID
	if aret == nil {
		appealID := snowflake.GenerateReviewID()
		reviewRep.AppealID = appealID
	} else {
		reviewRep.AppealID = aret.AppealID
	}
	reviewRep.Status = 10
	// 5. 保存结果
	if err = uc.repo.CreateAppealReview(ctx, reviewRep); err != nil {
		return nil, errors.New("保存申诉失败")
	}
	// 6. 直接返回结果
	return &model.ReviewAppealInfo{
		AppealID: reviewRep.AppealID,
	}, nil
}
