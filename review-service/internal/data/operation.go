package data

import (
	"context"
	"review-service/internal/biz"
	"review-service/internal/data/model"
	"review-service/internal/data/query"

	"github.com/go-kratos/kratos/v2/log"
)

type OperationRepo struct {
	data *Data
	log  *log.Helper
}

// 这里定义结构体作为key是防止其他第三方库或者你的同事命名一样的方法
type ctxTransactionKey struct{}

// GetDB 这是一个关键的辅助方法
func (o *OperationRepo) GetDB(ctx context.Context) *query.Query {
	// 从 context 中尝试获取事务对象 tx
	tx, ok := ctx.Value(ctxTransactionKey{}).(*query.Query)
	if ok {
		return tx
	}
	// 如果没有事务，返回默认的 query 对象
	return o.data.query
}

// NewGreeterRepo .
func NewoperationRepo(data *Data, logger log.Logger) biz.OperationRepo {
	return &OperationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (o *OperationRepo) UpdateAppealStatus(ctx context.Context, param *model.ReviewAppealInfo) (*model.ReviewAppealInfo, error) {
	db := o.GetDB(ctx)
	ret, err := db.ReviewAppealInfo.
		WithContext(ctx).
		Where(db.ReviewAppealInfo.AppealID.Eq(param.AppealID)).
		Updates(map[string]interface{}{
			"status":     param.Status,
			"op_user":    param.OpUser,
			"reason":     param.Reason,
			"op_remarks": param.OpRemarks,
		})
	log.Debugf("db.ReviewAppealInfo roweffect %#v", ret.RowsAffected)
	log.Debugf("db.ReviewAppealInfo roweffect %#v", param)
	return param, err

}
func (o *OperationRepo) UpdateReviewStatus(ctx context.Context, param *model.ReviewAppealInfo) error {
	db := o.GetDB(ctx)
	_, err := db.ReviewInfo.
		WithContext(ctx).
		Where(db.ReviewInfo.ReviewID.Eq(param.ReviewID)).
		Update(db.ReviewInfo.Status, 40)
	return err
}
