package object

type SysConfig struct {
	Total     total     `toml:"total"`
	Generator generator `toml:"generator"`
}

type total struct {
	IsDebug bool `toml:"isDebug"`
}

type generator struct {
	TimeoutNS  uint32 `toml:"timeoutNS"`
	Lps        uint32 `toml:"lps"`
	DurationNS uint32 `toml:"durationNS"`
}
