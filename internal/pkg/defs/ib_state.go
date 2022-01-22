package defs

// IBPhysicalState IB端口物理状态
type IBPhysicalState uint

const (
	// PhyNOChange 端口物理状态
	PhyNOChange IBPhysicalState = iota
	PhySleep
	PhyPolling
	PhyDisable
	PhyShift
	PhyLinkUp
	PhyLinkErrorRecover
	PhyPhyTest
)

// IBPhysicalStateList IB端口物理状态列表
var IBPhysicalStateList [8]string = [8]string{"no change", "sleep", "polling", "disable", "shift", "link up", "link error recover", "phytest"}

func (i IBPhysicalState) String() string {
	return IBPhysicalStateList[uint(i)]
}

// IBState IB端口状态
type IBState uint

const (
	// NOChange 端口物理状态
	NOChange IBState = iota
	Down
	Init
	Armed
	Active
	ActDefer
)

// IBStateList IB端口物理状态列表
var IBStateList [6]string = [6]string{"no change", "down", "init", "armed", "active", "act defer"}

func (i IBState) String() string {
	return IBStateList[uint(i)]
}
