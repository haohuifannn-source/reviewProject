package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type AuditParam struct {
	AppealID int64
	ReviewID int64
	Status   int64
	OpUser   string
	OpReason string
	OpRemark string
}

// OperationRepo is a Greater repo.
type OperationRepo interface {
	AuditAppeal(context.Context, *AuditParam) (int64, error)
}

// OperationUsecase is a Operation usecase.
type OperationUsecase struct {
	repo OperationRepo
	log  *log.Helper
}

// NewOperationUsecase new a Operation usecase.
func NewOperationUsecase(repo OperationRepo, logger log.Logger) *OperationUsecase {
	return &OperationUsecase{
		repo: repo,
		log:  log.NewHelper(logger)}
}

func (o *OperationUsecase) AuditAppeal(ctx context.Context, param *AuditParam) (int64, error) {
	return o.repo.AuditAppeal(ctx, param)
}
