package ylog

//往文件里面写日志

import (
	"fmt"
	"os"
	"path"
	"strconv"
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
func NewFileLogger(logLv level,filePath string)*FileLogger{
	fileName := getFileName()
	//名字统一格式，表示创建这个日志的时间，例如: 2020-10-12-13 - 2020年10月12日13点的日志
	f := &FileLogger{
		logLv:logLv,
		fileName: fileName,
		filePath: filePath,
		maxSize: 10 * 1024 * 1024,	//10M
	}
	f.initFile()	//初始化文件句柄
	return f
}

func getFileName()(fileName string){
	year := time.Now().Year()
	month := time.Now().Month()
	day := time.Now().Day()
	hour := time.Now().Hour()

	strYear := strconv.Itoa(year)
	var strMonth string
	switch month {
	case time.January:
		strMonth = "01"
	case time.February:
		strMonth = "02"
	case time.March:
		strMonth = "03"
	case time.April:
		strMonth = "04"
	case time.May:
		strMonth = "05"
	case time.June:
		strMonth = "06"
	case time.July:
		strMonth = "07"
	case time.August:
		strMonth = "08"
	case time.September:
		strMonth = "09"
	case time.October:
		strMonth = "10"
	case time.November:
		strMonth = "11"
	case time.December:
		strMonth = "12"
	}
	strDay := strconv.Itoa(day)
	strHour := strconv.Itoa(hour)

	fileName = strYear + "-" + strMonth + "-" + strDay + "-" + strHour + ".log"
	return fileName
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
	//fmt.Println("Name ====",fileInfo.Name())  //这里得到的是文件的名字，不是路径+名字
	fileSize := fileInfo.Size()
	if fileSize >= f.maxSize{
		return true		//要切分文件
	}
	return false	//不用切分文件
}

//判断文件时间，看一下是否要切分
func (f *FileLogger)checkFileTime(file *os.File)bool{
	fileInfo,_ := file.Stat()
	newFileName := getFileName()
	if fileInfo.Name() != newFileName{
		return true //需要拆分
	}
	return false	//不用切分文件
}
func (f *FileLogger)checkFileTimeErr(file *os.File)bool{
	fileInfo,_ := file.Stat()
	newFileName := getFileName()
	newFileName = newFileName + ".err"
	if fileInfo.Name() != newFileName{
		return true //需要拆分
	}
	return false	//不用切分文件
}

//按照时间切分文件,正常输出的log文件
func (f *FileLogger)splitLogFileTime(file *os.File)*os.File{
	//关闭原来的文件
	_ = file.Close()

	//1.拼接日志文件路径名字,打开新文件得到文件句柄
	newFileName := getFileName()
	logName := path.Join(f.filePath,newFileName)
	fileObj,err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0664)
	if err!=nil{
		panic(fmt.Errorf("open file %s error : %v",fileObj,err))
	}

	//返回新的
	return fileObj
}
//按照时间切分文件,错误输出的log文件
func (f *FileLogger)splitLogFileTimeErr(file *os.File)*os.File{
	//关闭原来的文件
	_ = file.Close()
	//1.拼接日志文件路径名字,打开新文件得到文件句柄
	newFileName := getFileName()
	newFileName = newFileName + ".err"
	logName := path.Join(f.filePath,newFileName)
	fileObj,err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND,0664)
	if err!=nil{
		panic(fmt.Errorf("open file %s error : %v",fileObj,err))
	}

	//返回新的
	return fileObj
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

	var res bool
	//判断文件大小
	//res = f.checkFileSize(f.file)
	//if res == true{
	//	f.file = f.splitLogFile(f.file)
	//}
	////写入文件
	//_,_ = fmt.Fprintln(f.file,logMsg)

	//判断时间,看看需不需要拆分, 这里判断的是正常的log名字
	res = f.checkFileTime(f.file)
	if res == true{
		f.file = f.splitLogFileTime(f.file)
	}
	//写入文件
	_,_ = fmt.Fprintln(f.file,logMsg)

	//如果是错误或者崩溃的日志，还要多记录进错误日志文件
	if lv >= ErrorLevel{
		//res = f.checkFileSize(f.errFile)
		//if res == true{
		//	f.errFile = f.splitLogFile(f.errFile)
		//}
		//_,_ = fmt.Fprintln(f.errFile,logMsg)

		//这里判断的是err的log名字
		res = f.checkFileTimeErr(f.errFile)
		if res == true{
			f.errFile = f.splitLogFileTimeErr(f.errFile)
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


