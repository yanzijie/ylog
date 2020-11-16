package ylog

//往文件里面写日志

import (
	"fmt"
	"os"
	"path"
	"time"
)

//日志文件结构
type FileLogger struct {
	logLv	level		//当前日志的级别
	fileName string
	filePath string
	file *os.File		//非错误信息的文件句柄
	errFile *os.File	//错误日志信息的文件句柄
	maxSize int64		//日志文件最大空间
}

//日志文件结构构造函数
func NewFileLogger(logLv level,fileName,filePath string)*FileLogger{
	f := &FileLogger{
		logLv:logLv,
		fileName: fileName,
		filePath: filePath,
		maxSize: 10 * 1024 * 1024,	//10M
	}
	f.initFile()	//初始化文件句柄
	return f
}

//打开指定日志文件，赋值给结构体字段
func (f *FileLogger)initFile(){
	//1.拼接日志文件路径名字, 2.打开文件得到文件句柄，3.赋值给结构体字段
	logName := path.Join(f.filePath,f.fileName)
	fileObj,err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0664)
	if err!=nil{
		panic(fmt.Errorf("open file %s error : %v",fileObj,err))
	}
	f.file = fileObj

	errLogName := fmt.Sprintf("%s.err",logName)
	errFileObj,err := os.OpenFile(errLogName, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0664)
	if err!=nil{
		panic(fmt.Errorf("open file %s error : %v",errFileObj,err))
	}
	f.errFile = errFileObj
}

//按照大小切分文件
func (f *FileLogger)splitLogFile(file *os.File)*os.File{
	//切分文件
	//1.得到当前文件的完整路径
	fileName := file.Name()
	backUpName := fmt.Sprintf("%s_%v.back",fileName,time.Now().Unix())	//备份文件完整路径
	//2.关闭原来的文件
	_ = file.Close()
	//3.备份原来的文件,把原来的文件 fileName 改名为 backUpName
	_ = os.Rename(fileName,backUpName)
	//4.新建一个, rename把原来的移走了，还是打开 fileName, 然后再赋值给结构体字段
	fileObj,err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0664)
	if err!=nil{
		panic(fmt.Errorf("open file %s error : %v",fileObj,err))
	}
	return fileObj

}

//判断文件大小，看一下是否要切分
func (f *FileLogger)checkFileSize(file *os.File)bool{
	fileInfo,_ := file.Stat()
	fileSize := fileInfo.Size()
	if fileSize >= f.maxSize{
		return true		//要切分文件
	}
	return false	//不用切分文件
}

//把日志写入文件`	format日志内容，args-占位符的参数
func (f *FileLogger)writeContent(lv level,format string, args ...interface{}){
	//设置日志打印门槛
	if f.logLv > lv{
		return
	}
	msg := fmt.Sprintf(format,args...)
	//日志格式 : [时间][文件名:行号][函数名][日志级别] 信息
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	fileName,line,funcName := getCaller(3)
	logLvStr := getLevelStr(lv)
	logMsg := fmt.Sprintf("[%s][%s:%d][%s][%s]%s",nowTime,fileName,line,funcName,logLvStr,msg)

	//判断文件大小
	var res bool
	res = f.checkFileSize(f.file)
	if res == true{
		f.file = f.splitLogFile(f.file)
	}

	//写入文件
	_,_ = fmt.Fprintln(f.file,logMsg)
	//如果是错误或者崩溃的日志，还要多记录进错误日志文件
	if lv >= ErrorLevel{
		res = f.checkFileSize(f.errFile)
		if res == true{
			f.errFile = f.splitLogFile(f.errFile)
		}
		_,_ = fmt.Fprintln(f.errFile,logMsg)
	}
}

//各种日志级别的写文件方法
func (f *FileLogger)Debug(format string, args ...interface{}){
	f.writeContent(DebugLevel,format,args...)
}

func (f *FileLogger)Info(format string, args ...interface{}){
	f.writeContent(InfoLevel,format,args...)
}

func (f *FileLogger)Warn(format string, args ...interface{}){
	f.writeContent(WarningLevel,format,args...)
}

func (f *FileLogger)Error(format string, args ...interface{}){
	f.writeContent(ErrorLevel,format,args...)
}

func (f *FileLogger)Fatal(format string, args ...interface{}){
	f.writeContent(FatalLevel,format,args...)
}

func (f *FileLogger)Close(){
	_ = f.file.Close()
	_ = f.errFile.Close()
}


