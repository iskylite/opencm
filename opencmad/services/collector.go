package services

import (
	"context"

	"github.com/iskylite/opencm/opencmad/collector"
	"github.com/iskylite/opencm/pb"
)

// NodeCollectorHandler 初始化nc
func NodeCollectorHandler(s *Server) error {
	nc, err := collector.NewNodeCollector(s.logger, uint(s.collectorFlags))
	if err != nil {
		return err
	}
	s.nodeCollector = nc
	return nil
}

// DefaultNCHandler 默认nc接口测试运行
func DefaultNCHandler(s *Server) error {
	nc, err := collector.NewNodeCollector(s.logger, uint(1<<len(collector.CollectorFlag)-1))
	if err != nil {
		return err
	}
	s.nodeCollector = nc
	datas := nc.Gather()
	for _, data := range datas {
		collector.FormatCollectorData(data, s.logger)
	}
	return nil
}

// Collect 数据采集接口
func (s *Server) Collect(ctx context.Context, req *pb.CollectRequest) (*pb.CollectResponse, error) {
	return &pb.CollectResponse{
		NodeName:     s.host,
		NodeType:     s.NodeType,
		PullNodes:    req.PullNodes,
		CollectDatas: s.nodeCollector.Gather(),
	}, nil
}
