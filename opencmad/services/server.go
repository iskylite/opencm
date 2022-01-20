package services

import (
	"net"

	"github.com/go-kit/log"
	"github.com/iskylite/opencm/opencmad/collector"
	"github.com/iskylite/opencm/transport"
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
	listenAddress  string
	opencmdServer  string
	host           string
	NodeType       string
	logger         log.Logger
	collectorFlags int
	nodeCollector  *collector.NodeCollector
	gs             *grpc.Server
}

func DefaultServer() *Server {
	return &Server{
		listenAddress: *listenAddress,
		opencmdServer: net.JoinHostPort(*opencmdHost, *opencmdPort),
	}
}

func (s *Server) Use(handler func(*Server) error) {
	if err := handler(s); err != nil {
		panic(err)
	}
}

func (s *Server) Serve() error {
	conn, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		return err
	}
	s.gs = grpc.NewServer()
	transport.RegisterOpenCMADServiceServer(s.gs, s)
	return s.gs.Serve(conn)
}
