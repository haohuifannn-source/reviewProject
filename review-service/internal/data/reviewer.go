package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"review-service/internal/biz"
	"review-service/internal/data/model"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"golang.org/x/sync/singleflight"
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

func (r *reviewRepo) ListReviewByStoreID(ctx context.Context, storeID int64, offset, limit int32) ([]*biz.MyReviewInfo, error) {
	//return r.getData2(ctx, storeID, offset, limit) // 第一版，直接查ES
	return r.getData2(ctx, storeID, offset, limit) // 第一版，先查redis，再查ES
}

func (r *reviewRepo) getData1(ctx context.Context, storeID int64, offset, limit int32) ([]*biz.MyReviewInfo, error) {
	// 1. 构建查询请求
	res, err := r.data.elasticsearch.Search().
		Index(r.data.elasticsearch.index).
		// 分页设置：offset 对应 From，limit 对应 Size
		From(int(offset)).
		Size(int(limit)).
		Query(&types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{
						Term: map[string]types.TermQuery{
							// 注意：这里字段名要和 ES 索引中的 mapping 保持一致
							"store_id": {Value: storeID},
						},
					},
				},
			},
		}).
		Do(ctx) // 执行请求
	fmt.Printf("------> elasticsearch: %v %v", res, err)
	if err != nil {
		return nil, err
	}
	// 2. 解析结果
	reviews := make([]*biz.MyReviewInfo, 0, len(res.Hits.Hits))
	for idx, hit := range res.Hits.Hits {
		item := new(biz.MyReviewInfo)
		// hit.Source_ 是原始 JSON 字节流
		if err := json.Unmarshal(hit.Source_, item); err != nil {
			// 如果单条解析失败，记录日志并继续，防止整页崩溃
			r.log.Errorf("第 %v 信息 json.Unmarshal 失败", idx, err)
			continue
		}
		reviews = append(reviews, item)
	}
	fmt.Printf("------> elasticsearch reviews: %v\n", reviews)
	return reviews, nil
}

var g singleflight.Group

// key review:76901:1:10

// getData2带缓存版本的查询函数
func (r *reviewRepo) getData2(ctx context.Context, storeID int64, offset, limit int32) ([]*biz.MyReviewInfo, error) {
	// 1. 先查redis缓存
	// 2. 缓存没有则查询ES
	// 3. 通过singleflight合并短时间内的大量并发查询
	key := fmt.Sprintf("review:%d:%d:%d", storeID, offset, limit)
	b, err := r.getDataBySingleflighr(ctx, key)
	if err != nil {
		return nil, err
	}
	// 1.1 反序列化
	hm := new(types.HitsMetadata)
	if err := json.Unmarshal(b, hm); err != nil {
		return nil, err
	}
	list := make([]*biz.MyReviewInfo, 0, hm.Total.Value)
	for idx, hit := range hm.Hits {
		item := new(biz.MyReviewInfo)
		// hit.Source_ 是原始 JSON 字节流
		if err := json.Unmarshal(hit.Source_, item); err != nil {
			// 如果单条解析失败，记录日志并继续，防止整页崩溃
			r.log.Errorf("第 %v 信息 json.Unmarshal 失败", idx, err)
			continue
		}
		list = append(list, item)
	}
	return list, nil
}

// getData2带缓存版本的查询函数
func (r *reviewRepo) getDataBySingleflighr(ctx context.Context, key string) ([]byte, error) {
	v, err, shared := g.Do(key, func() (any, error) {
		// 查缓存
		data, err := r.getDataFromCache(ctx, key)
		r.log.Debugf("r.getDataFromCache data:%s, err :%v\n", data, err)
		if err == nil {
			return data, err
		}
		// 只有缓存查不到的时候才会去查ES
		if errors.Is(err, redis.Nil) {
			// 说明缓存失效，需要查ES
			data, err := r.getDataFromES(ctx, key)
			if err == nil {
				// 设置缓存
				return data, r.setCache(ctx, key, data)
			}
			return nil, err
		}
		// 查缓存失败了，直接返回错误，不继续向下传到压力，出错了一点不要查ES，熔断操作
		return nil, err
	})
	r.log.Debugf("Singleflighr ret : v : %v err : %v shared : %v\n", v, err, shared)
	if err != nil {
		return nil, err
	}
	return v.([]byte), nil
}

// getDataFromCache读缓存
func (r *reviewRepo) getDataFromCache(ctx context.Context, key string) ([]byte, error) {
	// 之所以用Bytes是因为ES查出来的也是字符切片，这里为了跟ES一致，也转化为String，这样后续的处理可以一样
	fmt.Printf("--------->getDataFromCache\n")
	ret, err := r.data.rdb.WithContext(ctx).Get(key).Bytes()
	if err != nil {
		fmt.Printf("--------->getDataFromCache failed % v\n", err)
		return nil, err
	}
	return ret, nil
}

// setCach设置缓存
func (r *reviewRepo) setCache(ctx context.Context, key string, data []byte) error {
	fmt.Printf("--------->setCache : key %s\n", key)
	err := r.data.rdb.WithContext(ctx).Set(key, data, time.Second*10).Err()
	if err != nil {
		fmt.Printf("--------->setCache failed % v\n", err)
		return err
	}
	return nil
}

// getDataFromES读ES
func (r *reviewRepo) getDataFromES(ctx context.Context, key string) ([]byte, error) {
	fmt.Printf("--------->getDataFromES\n")
	values := strings.Split(key, ":")
	if len(values) < 4 {
		return nil, errors.New("invalid key")
	}
	index, storeID, offsetStr, limiStr := values[0], values[1], values[2], values[3]
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return nil, err
	}
	limit, err := strconv.Atoi(limiStr)
	if err != nil {
		return nil, err
	}
	// 1. 构建查询请求
	res, err := r.data.elasticsearch.Search().
		Index(index).
		// 分页设置：offset 对应 From，limit 对应 Size
		From(int(offset)).
		Size(int(limit)).
		Query(&types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{
						Term: map[string]types.TermQuery{
							// 注意：这里字段名要和 ES 索引中的 mapping 保持一致
							"store_id": {Value: storeID},
						},
					},
				},
			},
		}).
		Do(ctx) // 执行请求
	fmt.Printf("------> elasticsearch: %v %v", res, err)
	if err != nil {
		return nil, err
	}
	return json.Marshal(res.Hits)
}
