package data

import (
	"context"
	"review-service/internal/biz"
	"review-service/internal/data/model"
	"review-service/internal/data/query"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/clause"
)

type businessRepo struct {
	data *Data
	log  *log.Helper
}

// NewbusinessRepo .
func NewbusinessRepo(data *Data, logger log.Logger) biz.BusinessRepo {
	return &businessRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (b *businessRepo) GetReviewByReviewId(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return b.data.query.ReviewInfo.WithContext(ctx).Where(b.data.query.ReviewInfo.ReviewID.Eq(reviewID)).First()
}

func (b *businessRepo) ReplyReview(ctx context.Context, reviewReply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
	err := b.data.query.ReviewReplyInfo.WithContext(ctx).Save(reviewReply)
	return reviewReply, err
}

// TaskToReply商家回复后的事务操作，即修改评论表和保存商家回复
func (b *businessRepo) TaskToReply(ctx context.Context, revieReply *model.ReviewReplyInfo) error {
	// Basic transaction
	err := b.data.query.Transaction(func(tx *query.Query) error {
		// 保存评价回复
		if err := tx.ReviewReplyInfo.WithContext(ctx).Save(revieReply); err != nil {
			b.log.Errorf("Save revieReply reply failed, err :%v", err)
			return err
		}

		if _, err := tx.ReviewInfo.WithContext(ctx).
			Where(tx.ReviewInfo.ReviewID.Eq(revieReply.ReviewID)).
			Update(tx.ReviewInfo.HasReply, 1); err != nil {
			b.log.Errorf("Update ReviewInfo failed, err :%v", err)
			return err
		}
		// return nil will commit the whole transaction
		return nil
	})
	return err
}

// AppealReview 商家申述
func (b *businessRepo) CreateAppealReview(ctx context.Context, appealReview *model.ReviewAppealInfo) error {
	return b.data.query.ReviewAppealInfo.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "review_id"}},
		// 只有当 review_id 冲突时，才执行以下更新
		DoUpdates: clause.Assignments(map[string]interface{}{
			"status":     appealReview.Status,
			"content":    appealReview.Content,
			"reason":     appealReview.Reason,
			"pic_info":   appealReview.PicInfo,
			"video_info": appealReview.VideoInfo,
		}),
	}).Create(appealReview)
}

// GetAppealReviewByReviewId 通过ReviewID查询是否有申述的记录
func (b *businessRepo) GetAppealReviewByReviewId(ctx context.Context, reviewID int64) (*model.ReviewAppealInfo, error) {
	return b.data.query.ReviewAppealInfo.WithContext(ctx).Where(b.data.query.ReviewAppealInfo.ReviewID.Eq(reviewID)).First()
}
