package ylog

import (
	"path"
	"runtime"
)

//得到文件名，行号,函数名  skip-需要跳几层去找行号
func getCaller(skip int)(fileName string,line int,funcName string){
	//file得到的是文件的全路径, pc可以拿到函数名
	pc,file,line,ok := runtime.Caller(skip)
	if !ok {
		return
	}
	//从file中得到文件名
	fileName = path.Base(file)		//得到最后一个 / 后面的字符串
	//拿到函数名
	funcName = runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName)
	return fileName,line,funcName
}


