package data

import (
	"context"
	v1 "review-B/api/business/v1"
	"review-B/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type BusinesserRepo struct {
	data *Data
	log  *log.Helper
}

// NewBusinesserRepo .
func NewBusinesserRepo(data *Data, logger log.Logger) biz.BusinesserRepo {
	return &BusinesserRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *BusinesserRepo) ReplyReview(ctx context.Context, ReplyReview *biz.ReplyParam) (int64, error) {
	r.log.WithContext(ctx).Infof("[data] Reply, param:%v", ReplyReview)
	// 不是查询数据库，而是发起一个RPC调用调用其他的服务
	ret, err := r.data.rc.ReplyReview(ctx, &v1.ReplyReviewRequest{
		ReviewID:  ReplyReview.ReviewID,
		StoreID:   ReplyReview.StoreID,
		Content:   ReplyReview.Content,
		PicInfo:   ReplyReview.PicInfo,
		VideoInfo: ReplyReview.VideoInfo,
	})
	r.log.WithContext(ctx).Debugf("ReplyReview return, ret: %v err : %v", ret, err)
	if err != nil {
		return 0, err
	}
	return ret.GetReplyID(), nil
}

func (r *BusinesserRepo) CreateAppealReview(ctx context.Context, AppealReview *biz.AppealParam) (int64, error) {
	ret, err := r.data.rc.AppealReview(ctx, &v1.AppealReviewRequest{
		ReviewID:  AppealReview.ReviewID,
		StoreID:   AppealReview.StoreID,
		Reason:    AppealReview.Reason,
		Content:   AppealReview.Content,
		PicInfo:   AppealReview.PicInfo,
		VideoInfo: AppealReview.VideoInfo,
	})
	r.log.WithContext(ctx).Debugf("AppealReview return, ret: %v err : %v", ret, err)
	if err != nil {
		return 0, err
	}
	return ret.GetAppealID(), nil
}
