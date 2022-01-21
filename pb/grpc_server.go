package pb

import (
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	Max_Message_Size = 1 << 30 // 1 GB
)

func NewGrpcServer(opts ...grpc.ServerOption) *grpc.Server {
	var options []grpc.ServerOption
	options = append(options,
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:             10 * time.Second, // wait time before ping if no activity
			Timeout:          20 * time.Second, // ping timeout
			MaxConnectionAge: 10 * time.Hour,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             60 * time.Second, // min time a client should wait before sending a ping
			PermitWithoutStream: true,
		}),
		grpc.MaxRecvMsgSize(Max_Message_Size),
		grpc.MaxSendMsgSize(Max_Message_Size),
	)
	for _, opt := range opts {
		if opt != nil {
			options = append(options, opt)
		}
	}
	return grpc.NewServer(options...)
}
