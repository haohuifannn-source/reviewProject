package biz

import (
	"context"
	"review-service/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

// OperationRepo is a Greater repo.
type OperationRepo interface {
	UpdateAppealStatus(context.Context, *model.ReviewAppealInfo) (*model.ReviewAppealInfo, error)
	UpdateReviewStatus(context.Context, *model.ReviewAppealInfo) error
}

// OperationUsecase is a Operation usecase.
type OperationUsecase struct {
	repo OperationRepo
	tx   Transaction
	log  *log.Helper
}

// NewOperationUsecase new a Operation usecase.
func NewOperationUsecase(repo OperationRepo, tx Transaction, logger log.Logger) *OperationUsecase {
	return &OperationUsecase{repo: repo, tx: tx, log: log.NewHelper(logger)}
}

func (o *OperationUsecase) ReplyAppeal(ctx context.Context, param *model.ReviewAppealInfo) (*model.ReviewAppealInfo, error) {
	// 1. 逻辑处理
	// 设计到事务的操作
	err := o.tx.InTx(ctx, func(ctx context.Context) error {
		// 1. 首先要去更新申述的表的状态
		ret, err := o.repo.UpdateAppealStatus(ctx, param)
		if err != nil {
			o.log.Errorf("o.repo.UpdateAppealStatus err :", err)
			return err
		}
		// 2. 然后当申述通过的时候需要隐藏意见表的意见
		if ret.Status == 20 {
			if err := o.repo.UpdateReviewStatus(ctx, param); err != nil {
				o.log.Errorf("o.repo.UpdateReviewStatus err :", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil
	}
	return param, nil
}
