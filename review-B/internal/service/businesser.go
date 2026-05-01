package service

import (
	"context"
	pb "review-B/api/business/v1"
	"review-B/internal/biz"
)

// BusinesserService is a Businesser service.
type BusinesserService struct {
	pb.UnimplementedBusinessServer

	uc *biz.BusinesserUsecase
}

// NewBusinesserService new a Businesser service.
func NewBusinesserService(uc *biz.BusinesserUsecase) *BusinesserService {
	return &BusinesserService{uc: uc}
}

func (s *BusinesserService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	replyID, err := s.uc.CreateReply(ctx, &biz.ReplyParam{
		ReviewID:  req.GetReviewID(),
		StoreID:   req.GetStoreID(),
		Content:   req.GetContent(),
		PicInfo:   req.GetPicInfo(),
		VideoInfo: req.GetVideoInfo(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReplyReviewReply{
		ReplyID: replyID,
	}, nil
}

func (s *BusinesserService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error) {
	param := &biz.AppealParam{
		ReviewID:  req.GetReviewID(),
		StoreID:   req.GetStoreID(),
		Reason:    req.GetReason(),
		Content:   req.GetContent(),
		VideoInfo: req.GetVideoInfo(),
		PicInfo:   req.GetPicInfo(),
	}
	ret, err := s.uc.AppealReview(ctx, param)
	if err != nil {
		return nil, err
	}
	return &pb.AppealReviewReply{
		AppealID: ret,
	}, nil
}
