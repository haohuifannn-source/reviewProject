package service

import (
	"context"
	"fmt"
	"review-service/internal/biz"
	"review-service/internal/data/model"

	pb "review-service/api/business/v1"
)

type BusinessService struct {
	pb.UnimplementedBusinessServer

	uc *biz.BusinessUsecase
}

func NewBusinessService(uc *biz.BusinessUsecase) *BusinessService {
	return &BusinessService{
		uc: uc,
	}
}

func (s *BusinessService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	// 参数转化
	// 调用biz逻辑代码
	fmt.Println("调用了回复函数")
	res, err := s.uc.ReplyReview(ctx, &model.ReviewReplyInfo{
		ReviewID:  req.GetReviewID(),
		StoreID:   req.GetStoreID(),
		Content:   req.GetContent(),
		PicInfo:   req.GetPicInfo(),
		VideoInfo: req.GetVideoInfo(),
	})
	if err != nil {
		return nil, err
	}
	// 构建返回的结构体
	return &pb.ReplyReviewReply{
		ReplyID: res.ReplyID,
	}, nil
}

// AppealReview 商家申述用户评价
func (s *BusinessService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error) {
	// 参数转化
	appealInfo := &model.ReviewAppealInfo{
		ReviewID:  req.GetReviewID(),
		StoreID:   req.GetStoreID(),
		Reason:    req.GetReason(),
		Content:   req.GetContent(),
		PicInfo:   req.GetPicInfo(),
		VideoInfo: req.GetVideoInfo(),
	}
	// 调用biz的逻辑函数
	res, err := s.uc.ApealReview(ctx, appealInfo)
	if err != nil {
		return nil, err
	}
	// 构建返回的结构体
	return &pb.AppealReviewReply{
		AppealID: res.AppealID,
	}, nil
}
