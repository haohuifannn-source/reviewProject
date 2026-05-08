package biz

import (
	"context"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewBusinessUsecase, NewReviewerUsecase, NewOperationUsecase)

// Transaction 定义事务接口
type Transaction interface {
	InTx(context.Context, func(context.Context) error) error
}
