package server

import (
	v1Business "review-service/api/business/v1"
	v1Review "review-service/api/review/v1"
	"review-service/internal/conf"
	"review-service/internal/service"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, reviewer *service.ReviewService, business *service.BusinessService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			validate.Validator(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1Review.RegisterReviewServer(srv, reviewer)
	v1Business.RegisterBusinessServer(srv, business)
	return srv
}

func NewRegister(conf *conf.Registry) registry.Registrar {
	// new consul client
	c := api.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	// new reg with consul client
	reg := consul.New(client, consul.WithHealthCheck(true))
	return reg
}
