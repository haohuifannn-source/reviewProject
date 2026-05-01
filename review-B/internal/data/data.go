package data

import (
	"context"
	v1 "review-B/api/business/v1"
	"review-B/internal/conf"

	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDiscover, NewReviewServiceClient, NewData, NewBusinesserRepo)

// Data .
type Data struct {
	//应该嵌入一个GRPC的客户端，通过Client去调用review-service的服务
	rc  v1.BusinessClient
	log *log.Helper
}

// NewData .
func NewData(rc v1.BusinessClient, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{
		rc:  rc,
		log: log.NewHelper(logger),
	}, cleanup, nil
}

func NewReviewServiceClient(d registry.Discovery) v1.BusinessClient {
	// import "github.com/go-kratos/kratos/v2/transport/grpc"
	conn, err := grpc.DialInsecure(
		context.Background(),
		//grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithEndpoint("discovery:///review-service"), //这里对应于review-service中main函数定义的服务名称
		grpc.WithDiscovery(d),
		grpc.WithMiddleware(
			recovery.Recovery(),
			validate.ProtoValidate(),
		),
	)
	if err != nil {
		panic(err)
	}
	return v1.NewBusinessClient(conn)
}

// NewDiscover服务发现的构造函数
func NewDiscover(conf *conf.Registry) registry.Discovery {
	// new consul client
	c := api.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	client, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	// new dis with consul client
	dis := consul.New(client)
	return dis
}
