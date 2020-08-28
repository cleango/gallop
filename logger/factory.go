package logger

type LogFactory struct {
	Cfg LogConfig `value:"logger"`
}

func NewLogFactory() *LogFactory {
	return &LogFactory{}
}
func (l *LogFactory) InitLog() {
	InitLog(l.Cfg)
}
