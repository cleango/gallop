package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//LogConfig 日志配置文件
type LogConfig struct {
	Type  string `value:"type"`
	Path  string `value:"path"`
	Level string `value:"level"`
	Stack bool `value:"stack"`
}

var (
	//  logger
	logger *zap.SugaredLogger = nil
)

//InitLog 注册日志
func InitLog(c LogConfig) {
	if c.Level == "" {
		c.Level = "debug"
	}
	if c.Path == "" {
		c.Path = "./logs"
	}
	zLogger := initLog(c)
	logger = zLogger.Sugar()
}

func getEnableLevelPag(min zapcore.Level, max zapcore.Level) zap.LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		return level >= min && level <= max
	}
}
func getEnableLevelMin(min zapcore.Level) zap.LevelEnablerFunc {
	return func(level zapcore.Level) bool {
		return level >= min
	}
}

// initLog create logger
func initLog(c LogConfig) *zap.Logger {
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		StacktraceKey:"stack",
		CallerKey:    "func",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
	if c.Stack{
		return zap.New(getCore(c, encoder), zap.AddCaller(),zap.AddCallerSkip(1), zap.AddStacktrace(zap.WarnLevel))
	}else{
		return zap.New(getCore(c, encoder), zap.AddCaller(),zap.AddCallerSkip(1))
	}

}
func getCore(c LogConfig, encoder zapcore.Encoder) zapcore.Core {
	l := zap.NewAtomicLevel().Level()
	_ = l.Set(c.Level)

	var cores []zapcore.Core
	switch c.Type {
	case "file":
		if zapcore.ErrorLevel > l {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(getWriter(c, "info")), getEnableLevelPag(zap.DebugLevel, zap.WarnLevel)))
		}
		if zapcore.ErrorLevel >= l {
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(getWriter(c, "error")), getEnableLevelMin(zap.ErrorLevel)))
		}
	default:
		cores = append(cores, zapcore.NewCore(encoder, os.Stdout, l))
	}
	return zapcore.NewTee(cores...)
}
func getWriter(c LogConfig, level string) io.Writer {

	return &lumberjack.Logger{
		Filename:   c.Path + "/" + level + ".log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     20,
		Compress:   true,
	}
}

//LogField 隔离引用
type LogField map[string]interface{}

//Put 属性赋值
func (l LogField) Put(key string, value interface{}) LogField {
	l[key] = value
	return l
}

func with(fields ...LogField) *zap.SugaredLogger {
	if len(fields) == 0 {
		return logger
	}
	var params []interface{}
	for _, f := range fields {
		for k, v := range f {
			params = append(params, k, v)
		}
	}
	return logger.With(params...)
}

//Debug 日志
func Debug(msg string, fields ...LogField) {
	with(fields...).Debug(msg)
}

//Info 日志
func Info(msg string, fields ...LogField) {

	with(fields...).Info(msg)
}

//Warn 日志
func Warn(msg string, fields ...LogField) {

	with(fields...).Warn(msg)
}

//Error 日志
func Error(err error, fields ...LogField) {
	with(fields...).Error(err.Error())
}

//Errorf 日志
func Errorf(msg string, err interface{}, fields ...LogField) {
	with(fields...).Errorf("%s %v", msg, err)
}

//Panic 日志
func Panic(msg error, fields ...LogField) {
	with(fields...).Panic(msg.Error())
}

//Fatal 日志
func Fatal(msg string, fields ...LogField) {
	with(fields...).Fatal(msg)
}
