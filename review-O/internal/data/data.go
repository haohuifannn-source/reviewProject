package data

import (
	"context"
	v1 "review-O/api/operation/v1"
	"review-O/internal/conf"

	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/hashicorp/consul/api"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewOperationRepo, NewDiscover)

// Data .
type Data struct {
	// TODO wrapped database client
	rc  v1.OperationClient
	log *log.Helper
}

// NewData .
func NewData(rc v1.OperationClient, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{
		rc:  rc,
		log: log.NewHelper(logger),
	}, cleanup, nil
}

func NewDiscover(conf *conf.Registry) v1.OperationClient {
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

	endpoint := "discovery:///review-service"
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(dis),
		grpc.WithMiddleware(
			recovery.Recovery(),
			validate.ProtoValidate(),
		),
	)
	if err != nil {
		panic(err)
	}
	return v1.NewOperationClient(conn)
}
