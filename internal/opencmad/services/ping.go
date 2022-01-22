package services

import (
	"context"

	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/pb"
)

// Ping 确认服务端是否在运行
func (s *Server) Ping(ctx context.Context, req *pb.GenericMsg) (*pb.GenericMsg, error) {
	level.Debug(s.logger).Log("msg", "recv ping.", "node", req.Node)
	return &pb.GenericMsg{Msg: "Ping OK", Node: s.host, Flag: serverName}, nil
}
