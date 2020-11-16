package ylog

import (
	"fmt"
	"os"
	"time"
)

//在终端打印日志
type ConsoleLog struct {
	logLv level
}

//日志文件结构构造函数
func NewConsoleLogger(logLv level)*ConsoleLog{
	c := &ConsoleLog{
		logLv:logLv,
	}
	return c
}

//把日志写入文件`	format日志内容，args-占位符的参数
func (c *ConsoleLog)consoleContent(lv level,format string, args ...interface{}){
	//设置日志打印门槛
	if c.logLv > lv{
		return
	}
	msg := fmt.Sprintf(format,args...)
	//日志格式 : [时间][文件名:行号][函数名][日志级别] 信息
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	fileName,line,funcName := getCaller(3)
	logLvStr := getLevelStr(lv)
	logMsg := fmt.Sprintf("[%s][%s:%d][%s][%s]%s",nowTime,fileName,line,funcName,logLvStr,msg)
	_,_ = fmt.Fprintln(os.Stdout,logMsg)
}

//各种日志级别的写文件方法
func (c *ConsoleLog)Debug(format string, args ...interface{}){
	c.consoleContent(DebugLevel,format,args...)
}

func (c *ConsoleLog)Info(format string, args ...interface{}){
	c.consoleContent(InfoLevel,format,args...)
}

func (c *ConsoleLog)Warn(format string, args ...interface{}){
	c.consoleContent(WarningLevel,format,args...)
}

func (c *ConsoleLog)Error(format string, args ...interface{}){
	c.consoleContent(ErrorLevel,format,args...)
}

func (c *ConsoleLog)Fatal(format string, args ...interface{}){
	c.consoleContent(FatalLevel,format,args...)
}

//单纯实现接口
func (c *ConsoleLog)Close(){

}

