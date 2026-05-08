package biz

import (
	"context"
	"fmt"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// ReviewerRepo is a Greater repo.
type ReviewerRepo interface {
	CreateReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderId(context.Context, int64) ([]*model.ReviewInfo, error)
	ListReviewByStoreID(context.Context, int64, int32, int32) ([]*MyReviewInfo, error)
	// ReplyReview(context.Context, *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error)
	// TaskToReply(context.Context, *model.ReviewReplyInfo) error
}

// ReviewerUsecase is a Reviewer usecase.
type ReviewerUsecase struct {
	repo ReviewerRepo
	log  *log.Helper
}

// NewReviewerUsecase new a Reviewer usecase.
func NewReviewerUsecase(repo ReviewerRepo, logger log.Logger) *ReviewerUsecase {
	return &ReviewerUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateReview 创建评价
// 实现业务逻辑的地方
func (uc *ReviewerUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReview, req:%v\n", review)
	// 1. 数据校验
	// 1.1 参数基础校验：正常来说不应该在这一层，在上一层或者框架层都应该能拦住（validator参数校验）
	// 1.2 参数业务校验：带业务逻辑的参数校验，比如已经评价过的不能再创建评价
	ret, err := uc.repo.GetReviewByOrderId(ctx, review.OrderID)
	if err != nil {
		fmt.Printf("GetReviewById fail, err:%v\n", err)
		return nil, v1.ErrorDbFiled("查询数据库失败")
	}
	if len(ret) > 0 {
		return nil, v1.ErrorOvderReviewed("订单%d已经评价", review.OrderID)
	}
	// 2. 生成reviewID
	// 这里使用雪花算法自己生成
	// 也可以接入公司内部的分布式雪花算法生成器
	urID := snowflake.GenerateReviewID()
	review.ReviewID = urID
	// 3. 查询订单和商品的快照信息
	// 实际业务场景下需要查询订单服务和商家服务（通过RPC调用订单服务和商家服务）
	// 4. 拼装数据入库
	return uc.repo.CreateReview(ctx, review)
}

func (uc *ReviewerUsecase) ListReviewByStoreID(ctx context.Context, storeID int64, page int32, size int32) ([]*MyReviewInfo, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}
	offset := (page - 1) * size
	limit := size
	return uc.repo.ListReviewByStoreID(ctx, storeID, offset, limit)
}

type MyReviewInfo struct {
	*model.ReviewInfo
	CreateAt MyTime `json:"create_at"`
	UpdateAt MyTime `json:"update_at"`
	// 因为ela得到的数据都是string类型的，需要把他们int的提出来，加上string的tag
	ID           int64 `json:"id,string"`
	Version      int32 `json:"version,string"`
	ReviewID     int64 `json:"review_id,string"` // id
	Score        int32 `json:"score,string"`
	ServiceScore int32 `json:"service_score,string"`
	ExpressScore int32 `json:"express_score,string"`
	HasMedia     int32 `json:"has_media,string"`
	OrderID      int64 `json:"order_id,string"` // id
	SkuID        int64 `json:"sku_id,string"`   // sku id
	SpuID        int64 `json:"spu_id,string"`   // spu id
	StoreID      int64 `json:"store_id,string"` // id
	UserID       int64 `json:"user_id,string"`  // id
	Anonymous    int32 `json:"anonymous,string"`
	Status       int32 `json:"status,string"`
	IsDefault    int32 `json:"is_default,string"`
	HasReply     int32 `json:"has_reply,string"` // :0;1
}

type MyTime time.Time

func (t *MyTime) UnmarshalJSON(data []byte) error {
	// data = `"2006-01-02 15:04:05"`
	// 因为本来就有引号，转化为string的时候又加了一个引号
	s := strings.Trim(string(data), `"`)
	tmp, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}
	*t = MyTime(tmp) // 给标准时间对象贴上自定义的标签，同时通过指针去修改对象
	return nil
}
