package services

import (
	"net"

	"github.com/go-kit/log"
	"github.com/iskylite/opencm/internal/opencmad/collector"
	"github.com/iskylite/opencm/pb"
	"google.golang.org/grpc"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	serverName = "opencmad"
)

var (
	listenAddress     = kingpin.Flag("grpc.listen-address", "Address on which to grpc server.").Default(":1995").String()
	opencmdHost       = kingpin.Flag("grpc.opencmd-host", "opencmd server host.").Default("mn0").String()
	opencmdPort       = kingpin.Flag("grpc.opencmd-port", "opencmd server port.").Default("1994").String()
	clientDialTimeout = kingpin.Flag("grpc.opencmd-dial-timeout", "opencmd connect timeout seconds.").Default("2").Int()
	clientExecTimeout = kingpin.Flag("grpc.opencmd-exec-timeout", "opencmd client execute timeout seconds after connected.").Default("2").Int()
)

// Server Grpc服务
type Server struct {
	serverName     string
	listenAddress  string
	opencmdServer  string
	host           string
	NodeType       string
	logger         log.Logger
	collectorFlags int
	nodeCollector  *collector.NodeCollector
	opts           []grpc.ServerOption
	gs             *grpc.Server
}

// DefaultServer 默认Server
func DefaultServer() *Server {
	return &Server{
		serverName:    serverName,
		listenAddress: *listenAddress,
		opencmdServer: net.JoinHostPort(*opencmdHost, *opencmdPort),
	}
}

// Use 配置Server，包括opencmad服务注册
func (s *Server) Use(handler func(*Server) error) {
	if err := handler(s); err != nil {
		panic(err)
	}
}

// Serve 启动grpc.Server
func (s *Server) Serve() error {
	conn, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		return err
	}
	s.gs = grpc.NewServer(s.opts...)
	pb.RegisterOpenCMADServiceServer(s.gs, s)
	return s.gs.Serve(conn)
}
