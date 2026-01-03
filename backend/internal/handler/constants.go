package handler

// 网关处理器共用常量
const (
	// DefaultMaxAccountSwitches 默认账号切换次数上限
	// 用于处理上游错误时的账号切换重试
	DefaultMaxAccountSwitches = 3

	// MaxAccountSwitchesHigh 高切换次数上限
	// 用于 Claude API 等需要更多重试的场景
	MaxAccountSwitchesHigh = 10
)
