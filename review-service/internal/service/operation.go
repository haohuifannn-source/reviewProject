package service

import (
	"context"
	"review-service/internal/biz"
	"review-service/internal/data/model"

	pb "review-service/api/operation/v1"
)

type OperationService struct {
	pb.UnimplementedOperationServer

	uc *biz.OperationUsecase
}

func NewOperationService(uc *biz.OperationUsecase) *OperationService {
	return &OperationService{uc: uc}
}

func (s *OperationService) AuditReview(ctx context.Context, req *pb.AuditReviewRequest) (*pb.AuditReviewReply, error) {
	return &pb.AuditReviewReply{}, nil
}
func (s *OperationService) AuditAppeal(ctx context.Context, req *pb.AuditAppealRequest) (*pb.AuditAppealReply, error) {
	// 1. 参数的转换
	audit := &model.ReviewAppealInfo{
		AppealID:  req.GetAppealID(),
		ReviewID:  req.GetReviewID(),
		Status:    int32(req.GetStatus()),
		OpUser:    req.GetOpUser(),
		Reason:    req.GetOpReason(),
		OpRemarks: req.GetOpRemarks(),
	}
	// 2. 调用逻辑
	ret, err := s.uc.ReplyAppeal(ctx, audit)
	if err != nil {
		return nil, err
	}
	return &pb.AuditAppealReply{
		ReviewID: ret.ReviewID,
	}, nil
}
