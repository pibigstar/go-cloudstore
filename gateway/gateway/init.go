package gateway

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/pibigstar/go-cloudstore/config"
	"github.com/pibigstar/go-cloudstore/service/account/proto"
)

var userClient proto.UserService

func init() {
	//注册中心设为consul
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.ConsulServerAddr,
		}
	})
	service := micro.NewService(micro.Registry(reg))
	//解析命令行参数
	service.Init()
	cli := service.Client()
	//初始化一个rpcClient
	userClient = proto.NewUserService(config.APIUserServiceName, cli)
}
