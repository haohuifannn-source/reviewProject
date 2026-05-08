package data

import (
	"context"
	v1 "review-O/api/operation/v1"
	"review-O/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type operationRepo struct {
	data *Data
	log  *log.Helper
}

// NewOperationRepo .
func NewOperationRepo(data *Data, logger log.Logger) biz.OperationRepo {
	return &operationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (o *operationRepo) AuditAppeal(ctx context.Context, req *biz.AuditParam) (int64, error) {
	ret, err := o.data.rc.AuditAppeal(ctx, &v1.AuditAppealRequest{
		AppealID:  req.AppealID,
		ReviewID:  req.ReviewID,
		Status:    req.Status,
		OpUser:    req.OpUser,
		OpReason:  req.OpReason,
		OpRemarks: &req.OpRemark,
	})
	if err != nil {
		return 0, err
	}
	return ret.ReviewID, err

}
