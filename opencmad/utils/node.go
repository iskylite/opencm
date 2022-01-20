package utils

import (
	"bytes"
	"net"
	"os"
	"sort"

	"github.com/iskylite/opencm/transport"
	"github.com/iskylite/procfs"
	"golang.org/x/sys/unix"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	tsFilePath = kingpin.Flag("node.ts_file_path",
		"diskless image .ts path").Default("/.ts").String()
)

// GetOS 获取系统基本配置
func GetOS() (*transport.OS, error) {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		return nil, err
	}
	memInfo, err := fs.Meminfo()
	if err != nil {
		return nil, err
	}

	stats, err := fs.Stat()
	if err != nil {
		return nil, err
	}

	var utsname unix.Utsname
	if err := unix.Uname(&utsname); err != nil {
		return nil, err
	}

	var ts string
	tsBytes, err := ReadAll(*tsFilePath)
	if err != nil {
		ts = ""
	}
	ts = BytesToString(tsBytes)

	return &transport.OS{
		BootTime:       stats.BootTime,
		CPUNum:         int32(len(stats.CPU)),
		MemTotal:       *memInfo.MemTotal / 1024 / 1024,
		Arch:           string(utsname.Machine[:bytes.IndexByte(utsname.Machine[:], 0)]),
		Kernel:         string(utsname.Release[:bytes.IndexByte(utsname.Release[:], 0)]),
		Version:        string(utsname.Version[:bytes.IndexByte(utsname.Version[:], 0)]),
		ImageBuildTime: ts,
	}, nil
}

// GetHostname 获取主机名
func GetHostname(allInterfaces []*transport.Interface) string {
	hostname, err := os.Hostname()
	if err == nil && hostname != "localhost" {
		return hostname
	}
	hostnames := make([]string, 0, len(allInterfaces))
	for _, ifcfg := range allInterfaces {
		hosts, err := net.LookupAddr(ifcfg.IP)
		if err != nil {
			continue
		}
		for _, host := range hosts {
			hostnames = append(hostnames, host)
		}
	}
	if len(hostnames) == 0 {
		return "localhost"
	}
	sort.Slice(hostnames, func(i, j int) bool {
		if len(hostnames[i]) < len(hostnames[j]) {
			return true
		}
		return false
	})
	return hostnames[0]
}
