package utils

import (
	"net"

	"github.com/iskylite/opencm/pb"
)

// GetInterfaces 获取所有可用得网络配置
func GetInterfaces() ([]*pb.Interface, error) {
	interfaces := make([]*pb.Interface, 0)
	allInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, val := range allInterfaces {
		if (val.Flags & net.FlagLoopback) == net.FlagLoopback {
			continue
		}
		addrs, err := val.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, err
			}
			if ip.To4() == nil {
				continue
			}
			ones, _ := ip.DefaultMask().Size()
			interfaces = append(interfaces, &pb.Interface{
				Dev:          val.Name,
				HardwareAddr: val.HardwareAddr.String(),
				Flags:        val.Flags.String(),
				IP:           ip.String(),
				Mask:         int32(ones),
			})
		}
	}
	return interfaces, nil
}
