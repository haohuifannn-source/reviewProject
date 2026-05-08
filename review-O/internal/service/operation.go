package service

import (
	"context"

	pb "review-O/api/operation/v1"
	"review-O/internal/biz"
)

type OperationService struct {
	pb.UnimplementedOperationServer

	uc *biz.OperationUsecase
}

func NewOperationService(uc *biz.OperationUsecase) *OperationService {
	return &OperationService{
		uc: uc,
	}
}

func (s *OperationService) AuditReview(ctx context.Context, req *pb.AuditReviewRequest) (*pb.AuditReviewReply, error) {
	return &pb.AuditReviewReply{}, nil
}
func (s *OperationService) AuditAppeal(ctx context.Context, req *pb.AuditAppealRequest) (*pb.AuditAppealReply, error) {
	// 1. 参数转化
	param := &biz.AuditParam{
		AppealID: req.GetAppealID(),
		ReviewID: req.GetReviewID(),
		Status:   req.GetStatus(),
		OpUser:   req.GetOpUser(),
		OpReason: req.GetOpReason(),
		OpRemark: req.GetOpRemarks(),
	}
	// 2. 调用biz层逻辑
	id, err := s.uc.AuditAppeal(ctx, param)
	return &pb.AuditAppealReply{
		ReviewID: id,
	}, err
}
