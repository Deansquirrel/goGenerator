package object

type SysConfig struct {
	TotalConfig totalConfig			`toml:"total"`
	GeneratorConfig generatorConfig	`toml:"generator"`
}

type totalConfig struct {
	IsDebug bool					`toml:"isDebug"`
	T string	`toml:"t"`
}

type generatorConfig struct {
	TimeoutNS uint32		`toml:"timeoutNS"`
	Lps uint32						`toml:"lps"`
	DurationNS uint32		`toml:"durationNS"`
}