package services

import (
	"context"
	"time"

	"github.com/go-kit/log/level"
	"github.com/iskylite/opencm/opencmad/utils"
	"github.com/iskylite/opencm/pb"
	"google.golang.org/grpc"
)

func getRegisterData() (*pb.OpenCMADRegistry, error) {
	interfaces, err := utils.GetInterfaces()
	if err != nil {
		return nil, err
	}
	hostname := utils.GetHostname(interfaces)
	osConfig, err := utils.GetOS()
	if err != nil {
		return nil, err
	}
	return &pb.OpenCMADRegistry{
		Host:       hostname,
		OS:         osConfig,
		Interfaces: interfaces,
	}, nil
}

// RegisterOpenCMADHandler 向opencmd注册当前服务
func RegisterOpenCMADHandler(s *Server) error {
	registry, err := getRegisterData()
	if err != nil {
		return err
	}
	timeoutOption := grpc.WithTimeout(time.Second * time.Duration(*clientDialTimeout))
	cli, err := initOpenCMDClient(s, timeoutOption, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(*clientExecTimeout))
	defer cancel()
	config, err := cli.RegisterOpenCMAD(ctx, registry)
	if err != nil {
		return err
	}
	s.collectorFlags = int(config.CollectorFlags)
	s.NodeType = config.NodeType
	s.host = registry.Host
	return nil
}

// UnRegisterOpenCMADHandler 向opencmd注销当前服务
func UnRegisterOpenCMADHandler(s *Server) error {
	timeoutOption := grpc.WithTimeout(time.Second * time.Duration(*clientDialTimeout))
	cli, err := initOpenCMDClient(s, timeoutOption, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*time.Duration(*clientExecTimeout))
	defer cancel()
	msg, err := cli.UnRegisterOpenCMAD(ctx, &pb.GenericMsg{Node: s.host, Flag: s.serverName, Msg: "UnRegisterOpenCMAD"})
	if err != nil {
		return err
	}
	level.Info(s.logger).Log("msg", msg.Msg)
	return nil
}

func initOpenCMDClient(s *Server, opts ...grpc.DialOption) (pb.OpenCMDServiceClient, error) {
	conn, err := grpc.Dial(s.opencmdServer, opts...)
	if err != nil {
		return nil, err
	}
	return pb.NewOpenCMDServiceClient(conn), nil
}
