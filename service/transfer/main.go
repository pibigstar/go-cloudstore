package main

import (
	"log"

	"github.com/pibigstar/go-cloudstore/config"
	"github.com/pibigstar/go-cloudstore/mq"
	"github.com/pibigstar/go-cloudstore/service/transfer/process"
)

func main() {

	log.Println("开启监听任务队列....")
	mq.StartConsume(config.TransOSSQueueName, "transfer_oss", process.ProcessTransfer)
}
