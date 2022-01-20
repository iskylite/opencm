package services

import (
	"context"

	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/transport"
)

func (s *Server) Ping(ctx context.Context, req *transport.GenericNode) (*transport.GenericMsg, error) {
	level.Debug(s.logger).Log("msg", "recv ping.", "node", req.Node)
	return &transport.GenericMsg{Msg: "OK"}, nil
}
