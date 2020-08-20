package gallop

import "github.com/spf13/pflag"

//Options 配置项
type Options struct {
	AddrPort   string
	ConfigPath string
}

// DefaultOptions 初始化命令行参数
func DefaultOptions() *Options {
	return &Options{
		AddrPort:   ":8080",
		ConfigPath: "./config.yaml",
	}
}

//AddFlags 增加命令行参数
func (s *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.AddrPort, "addr", s.AddrPort, "The ip address and port for the server on")
	fs.StringVarP(&s.ConfigPath, "config", "c", s.ConfigPath, "the config for config path")
}
