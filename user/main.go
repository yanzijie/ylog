package main

import "ylog"

var logger ylog.Logger

func main(){
	//logger := ylog.NewFileLogger(ylog.DebugLevel,"xxx.log","./")
	//cb := "彩笔"
	//defer logger.Close()
	//logger.Debug("%s222是个好演员",cb)
	//logger.Info("test2 Info")
	//logger.Warn("test2 Warn")
	//logger.Error("test2 error")
	//logger.Fatal("test2 fatal")

	//初始化时候设置最低打印级别，如果想打印全部级别，就设置最低的debugLevel级别就可以
	cFile := ylog.NewConsoleLogger(ylog.DebugLevel)
	cFile.Debug("DebugConsole")
	cFile.Info("InfoConsole")
	cFile.Warn("WarnConsole")
	cFile.Error("ErrorConsole")
	cFile.Fatal("FatalConsole")
}
