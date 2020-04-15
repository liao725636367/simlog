package simlog

import (
	"fmt"
	"os"
	"time"
)


//终端写日志内容
type ConsoleLogger struct{
	Level logLevel
}
//生成 窗口日志输出实例
func NewConsoleLogger(level logLevel) *ConsoleLogger{
	return &ConsoleLogger{
		Level:level,
	}
}
// Debug formats according to a format specifier and returns the resulting string.
func (l *ConsoleLogger) Debug(format string, a ...interface{}){
	l.outPut( DEBUG,format,a...)
}
// Info formats according to a format specifier and returns the resulting string.
func (l *ConsoleLogger) Info(format string, a ...interface{}){
	l.outPut(  INFO,format,a...)
}
func (l *ConsoleLogger) Warning(format string, a ...interface{}){
	l.outPut( WARNING,format,a...)
}
func (l *ConsoleLogger) Error(format string, a ...interface{}){
	l.outPut( ERROR,format,a...)
}
func (l *ConsoleLogger) Fatal(format string, a ...interface{}){
	l.outPut(FATAL,format,a...)
}

//output log
func (l *ConsoleLogger) outPut( level logLevel,format string, a ...interface{}){
	now := time.Now().Format("2006-01-02 15:04:05")
	if l.Level > level { //如果记录级别大于 日志设置级别就不记录日志
		return
	}
	funcName,fileName,lineNo := getInfo(3) //跳过堆栈
	//格式化信息
	 format=fmt.Sprintf(format,a...)

	fmt.Fprintf(os.Stdout,"[%s] [%s] [%s:%s:%d]  %s",now,level.String(),fileName,funcName,lineNo,format )
}


