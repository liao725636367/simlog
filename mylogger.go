package simLog

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

type logLevel uint16
//定义日志接口
type Logger interface {
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warning(format string, a ...interface{})
	Error(format string, a ...interface{})
	Fatal(format string, a ...interface{})
}
const (
	//定义日志级别
	DEBUG logLevel = iota //0
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
	NONE //不记录日志
)
func (l logLevel ) String () string{
	switch l {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "DEBUG"
	}
}
func getInfo(skip int )(funcName , fileName string,lineNo int){
	pc,file,lineNo,ok:=runtime.Caller(skip)
	if !ok{
		fmt.Printf("runtime.Caller() failed \n")
		return
	}
	funcName =runtime.FuncForPC(pc).Name() //获取函数名称
	//去掉 funcName 的文件名
	funcName = strings.Split(funcName,".")[1]
	fileName = path.Base(file) //去除路径前缀
	return
}
