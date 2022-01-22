package defs

// NodeType 节点类型
type NodeType string

const (
	// MN 管理节点
	MN NodeType = "mn"
	// LN 登录节点
	LN NodeType = "ln"
	// MDS 元数据服务器
	MDS NodeType = "mds"
	// OSS 对象存储服务器
	OSS NodeType = "oss"
	// ION 数据转发服务器
	ION NodeType = "ion"
	// CN 计算节点
	CN NodeType = "cn"
	// DEV 测试开发节点
	DEV NodeType = "dev"
	// NAS NAS存储节点
	NAS NodeType = "nas"
)
