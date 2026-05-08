package data

import (
	"context"
	"errors"
	"fmt"
	"review-service/internal/biz"
	"review-service/internal/conf"
	"review-service/internal/data/query"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/glebarez/sqlite"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDB, NewRedis, NewEla, NewData, NewbusinessRepo, NewreviewRepo, NewoperationRepo, NewTransaction)

// 这个方式是告诉 Wire，当你需要 biz.Transaction 时，就用 *Data 实例

// 判断是否实现了事务的接口
var _ biz.Transaction = (*Data)(nil)

// InTx 实现 biz 层的接口定义
func (d *Data) InTx(ctx context.Context, fn func(context.Context) error) error {
	// 这里的 d.query 是通过 gen 生成的 *query.Query
	return d.query.Transaction(func(tx *query.Query) error {
		// 使用我们之前讨论的 ctxTransactionKey 注入事务对象
		ctx = context.WithValue(ctx, ctxTransactionKey{}, tx)
		return fn(ctx)
	})
}
func NewTransaction(d *Data) biz.Transaction {
	return d // Data 实现了 biz.Transaction
}

// Data .
type Data struct {
	// TODO wrapped database client
	query         *query.Query
	log           *log.Helper
	elasticsearch *Elasticsearch
	rdb           *redis.Client
}

type Elasticsearch struct {
	*elasticsearch.TypedClient
	index string
}

func NewEla(conf *conf.Elasticsearch) (*Elasticsearch, error) {
	// ES 配置
	cfg := elasticsearch.Config{
		Addresses: conf.GetAddresses(),
	}

	// 创建客户端连接
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		fmt.Printf("elasticsearch.NewTypedClient failed, err:%v\n", err)
		return nil, err
	}
	return &Elasticsearch{
		TypedClient: client,
		index:       conf.GetIndex(),
	}, nil
}

// NewData .
func NewData(db *gorm.DB, rdb *redis.Client, elasticsearch *Elasticsearch, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 非常重要！为GEN生成的query代码设置数据库链接对象
	query.SetDefault(db)
	return &Data{
		log:           log.NewHelper(logger),
		query:         query.Q,
		elasticsearch: elasticsearch,
		rdb:           rdb,
	}, cleanup, nil
}

func NewDB(cfg *conf.Data) (*gorm.DB, error) {
	if cfg == nil {
		panic(errors.New("GEN:connectDB fail cfg is nil"))
	}
	switch strings.ToLower(cfg.Database.GetDriver()) {
	case "mysql":
		return gorm.Open(mysql.Open(cfg.Database.GetSource()))
	case "sqlite":
		return gorm.Open(sqlite.Open(cfg.Database.GetSource()))
	}
	return nil, errors.New("connectDB fail unsupported db driver")
}

func NewRedis(cfg *conf.Data) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		WriteTimeout: cfg.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  cfg.Redis.ReadTimeout.AsDuration(),
	})
}
