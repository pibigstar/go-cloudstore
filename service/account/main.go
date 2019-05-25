package main

import (
	"log"
	"time"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/pibigstar/go-cloudstore/config"

	"github.com/micro/go-micro"
	"github.com/pibigstar/go-cloudstore/service/account/proto"
	userService "github.com/pibigstar/go-cloudstore/service/account/service"
)

// 启动rpc User服务，为gateway提供rpc远程调用
func main() {
	// 注册中心
	reg := consul.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.ConsulServerAddr,
		}
	})
	//创建一个服务
	service := micro.NewService(micro.Name(config.APIUserServiceName),
		micro.Registry(reg),
		micro.RegisterTTL(time.Second*10),     //10s检查等待时间
		micro.RegisterInterval(time.Second*5), // 服务每5s发一次心跳
	)
	// 将User相关服务注册进去
	proto.RegisterUserServiceHandler(service.Server(), new(userService.User))
	err := service.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
