package service

import (
	"context"
	"fmt"

	pb "review-service/api/review/v1"
	"review-service/internal/biz"
	"review-service/internal/data/model"
)

type ReviewService struct {
	pb.UnimplementedReviewServer

	uc *biz.ReviewerUsecase
}

func NewReviewService(uc *biz.ReviewerUsecase) *ReviewService {
	return &ReviewService{
		uc: uc,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewReply, error) {
	fmt.Println("[service] CreateReview, req:%#v\n", req)
	// 参数转化
	// 调用biz层
	var anonymous int32
	if req.Anonymous {
		anonymous = 1
	}
	review, err := s.uc.CreateReview(ctx, &model.ReviewInfo{
		UserID:       req.UserID,
		OrderID:      req.OrderID,
		Score:        req.Score,
		ServiceScore: req.ServiceScore,
		ExpressScore: req.ExpressScore,
		Content:      req.Content,
		PicInfo:      req.PicInfo,
		VideoInfo:    req.VideoInfo,
		Anonymous:    anonymous,
		Status:       0,
	})
	if err != nil {
		return nil, err
	}
	// 拼接返回结果
	return &pb.CreateReviewReply{ReviewID: review.ReviewID}, nil
}
func (s *ReviewService) ReplyReview(ctx context.Context, revieReply *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	// 参数转化
	// 调用biz逻辑代码
	res, err := s.uc.ReplyReview(ctx, &model.ReviewReplyInfo{
		ReviewID:  revieReply.GetReviewID(),
		StoreID:   revieReply.GetStoreID(),
		Content:   revieReply.GetContent(),
		PicInfo:   revieReply.GetPicInfo(),
		VideoInfo: revieReply.GetVideoInfo(),
	})
	if err != nil {
		return nil, err
	}
	// 构建返回的结构体
	return &pb.ReplyReviewReply{
		ReplyID: res.ReplyID,
	}, nil
}
func (s *ReviewService) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewReply, error) {
	return &pb.DeleteReviewReply{}, nil
}
func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {
	return &pb.GetReviewReply{}, nil
}
func (s *ReviewService) ListReview(ctx context.Context, req *pb.ListReviewRequest) (*pb.ListReviewReply, error) {
	return &pb.ListReviewReply{}, nil
}
