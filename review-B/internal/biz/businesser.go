package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ReplyParam struct {
	ReviewID  int64
	StoreID   int64
	Content   string
	PicInfo   string
	VideoInfo string
}

type AppealParam struct {
	ReviewID  int64
	StoreID   int64
	Content   string
	Reason    string
	PicInfo   string
	VideoInfo string
}

// BusinesserRepo is a Greater repo.
type BusinesserRepo interface {
	ReplyReview(context.Context, *ReplyParam) (int64, error)

	CreateAppealReview(context.Context, *AppealParam) (int64, error)
}

// BusinesserUsecase is a Businesser usecase.
type BusinesserUsecase struct {
	repo BusinesserRepo
	log  *log.Helper
}

// NewBusinesserUsecase new a Businesser usecase.
func NewBusinesserUsecase(repo BusinesserRepo, logger log.Logger) *BusinesserUsecase {
	return &BusinesserUsecase{
		repo: repo,
		log:  log.NewHelper(logger)}
}

func (uc *BusinesserUsecase) CreateReply(ctx context.Context, param *ReplyParam) (int64, error) {
	return uc.repo.ReplyReview(ctx, param)
}

func (uc *BusinesserUsecase) AppealReview(ctx context.Context, param *AppealParam) (int64, error) {
	return uc.repo.CreateAppealReview(ctx, param)
}
