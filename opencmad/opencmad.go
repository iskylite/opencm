package main

import (
	_ "net/http/pprof"

	"github.com/iskylite/opencm/opencmad/services"
)

func main() {
	s := services.DefaultServer()
	// 日志
	s.Use(services.DefaultLoggerHandler)
	// 注册
	s.Use(services.RegisterOpenCMADHandler)
	// 注销
	defer s.Use(services.UnRegisterOpenCMADHandler)
	// 采集模块初始化
	s.Use(services.NodeCollectorHandler)
	// 开启grpc server
	if err := s.Serve(); err != nil {
		panic(err)
	}
}
