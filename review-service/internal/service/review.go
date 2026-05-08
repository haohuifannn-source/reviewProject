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
	fmt.Println("[service] CreateReview, req:%#v", req)
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

func (s *ReviewService) ListReviewByStoreID(ctx context.Context, req *pb.ListReviewByStoreIDRequest) (*pb.ListReviewByStoreIDReply, error) {
	fmt.Printf("[service] ListReviewByStoreID, req:%#v\n", req)
	rets, err := s.uc.ListReviewByStoreID(ctx, req.StoreID, req.Page, req.Size)
	if err != nil {
		return nil, nil
	}
	result := make([]*pb.ReviewInfo, 0, len(rets))
	for _, v := range rets {
		result = append(result, &pb.ReviewInfo{
			ReviewID:     v.ReviewID,
			UserID:       v.UserID,
			OrderID:      v.OrderID,
			Score:        int64(v.Score),
			ServiceScore: int64(v.ServiceScore),
			ExpressScore: int64(v.ExpressScore),
			Content:      v.Content,
			PicInfo:      v.PicInfo,
			VideoInfo:    v.VideoInfo,
		})
	}
	return &pb.ListReviewByStoreIDReply{
		List: result,
	}, nil
}
