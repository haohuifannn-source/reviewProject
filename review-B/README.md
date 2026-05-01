# 模拟Review-B模块，实现GRPC调用Review-service

## 搭建一个ReviewReply通信服务
### 1. 前面的执行步骤和实现Review-service一样
### 2. 不一样的是需要通过RPC调用Review-service模块的功能
2.1 在data层跟以往不一样的是不需要嵌入gorm.DB,而是嵌入一个grpc的客户端，来调用review-service的服务
### 3. 使用服务发现
3.1 跟review-service一样，先定义conf和.yaml文件
3.2 与服务注册不一样的是，这里是调用，所以应该在data层链接grpc的地方做修改。
```go
func NewReviewServiceClient(d registry.Discovery) v1.ReviewClient {
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
	return v1.NewReviewClient(conn)
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
```
4. 修改wire.go函数的传参以及main.go中newapp的传参

## 搭建一个商家回复评论的服务
## 搭建一个商家申述的服务
