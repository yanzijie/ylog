# ylog
日志打印

### 引入本库
```
go get github.com/yanzijie/ylog
```
或者使用go mod模式
```
import (
	"github.com/yanzijie/ylog"
)
```

### 使用
```
func main(){
	//写入文件
	//初始化时候设置最低打印级别，如果想打印全部级别，就设置最低的debugLevel级别就可以
	//NewFileLogger()函数的第一个参数是打印级别，
	//第二个参数日志存放路径（相对于程序可执行文件的路径）,建议写绝对路径，写到配置文件里面
	logger := ylog.NewFileLogger(ylog.DebugLevel,"./")
	cb := "彩笔"
	defer logger.Close()
	logger.Debug("%s222是个好演员",cb)   //打印debug级别日志
	logger.Info("test2 Info")           //打印info级别日志......
	logger.Warn("test2 Warn")
	logger.Error("test2 error")
	logger.Fatal("test2 fatal")

	//输出到终端
	//初始化时候设置最低打印级别，如果想打印全部级别，就设置最低的debugLevel级别就可以
	//cFile := ylog.NewConsoleLogger(ylog.DebugLevel)
	//cFile.Debug("DebugConsole")
	//cFile.Info("InfoConsole")
	//cFile.Warn("WarnConsole")
	//cFile.Error("ErrorConsole")
	//cFile.Fatal("FatalConsole")
}
```
