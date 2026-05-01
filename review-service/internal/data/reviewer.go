package data

import (
	"context"
	"review-service/internal/biz"
	"review-service/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type reviewRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewreviewRepo(data *Data, logger log.Logger) biz.ReviewerRepo {
	return &reviewRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *reviewRepo) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Save(review)
	return review, err
}

func (r *reviewRepo) GetReviewByOrderId(ctx context.Context, oId int64) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.OrderID.Eq(oId)).Find()
}

func (r *reviewRepo) GetReviewByReviewId(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(reviewID)).First()
}

// func (r *reviewRepo) ReplyReview(ctx context.Context, reviewReply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
// 	err := r.data.query.ReviewReplyInfo.WithContext(ctx).Save(reviewReply)
// 	return reviewReply, err
// }

// func (r *reviewRepo) GetReviewReplyByReviewID(ctx context.Context, ReviewID int64) (*model.ReviewReplyInfo, error) {
// 	return r.data.query.ReviewReplyInfo.WithContext(ctx).Where(r.data.query.ReviewReplyInfo.ReplyID.Eq(ReviewID)).First()
// }

// // TaskToReply商家回复后的事务操作，即修改评论表和保存商家回复
// func (r *reviewRepo) TaskToReply(ctx context.Context, revieReply *model.ReviewReplyInfo) error {
// 	// Basic transaction
// 	err := r.data.query.Transaction(func(tx *query.Query) error {
// 		// 保存评价回复
// 		if err := tx.ReviewReplyInfo.WithContext(ctx).Save(revieReply); err != nil {
// 			r.log.Errorf("Save revieReply reply failed, err :%v", err)
// 			return err
// 		}

// 		if _, err := tx.ReviewInfo.WithContext(ctx).
// 			Where(tx.ReviewInfo.ReviewID.Eq(revieReply.ReviewID)).
// 			Update(tx.ReviewInfo.HasReply, 1); err != nil {
// 			r.log.Errorf("Update ReviewInfo failed, err :%v", err)
// 			return err
// 		}
// 		// return nil will commit the whole transaction
// 		return nil
// 	})
// 	return err
// }
