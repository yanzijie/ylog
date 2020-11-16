package ylog

//日志级别
type level uint16
const(
	DebugLevel level = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel		//崩溃
)

type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	Close()
}

//得到日志级别字符串
func getLevelStr(lv level)string{
	switch lv {
	case 0:
		return "DebugLevel"
	case 1:
		return "InfoLevel"
	case 2:
		return "WarningLevel"
	case 3:
		return "ErrorLevel"
	case 4:
		return "FatalLevel"
	default:
		return "DebugLevel"
	}
}

