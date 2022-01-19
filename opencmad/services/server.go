package services

import (
	"github.com/iskylite/opencm/opencmad/collector"
)

// Server Grpc服务
type Server struct {
	NodeCollector collector.NodeCollector
}
