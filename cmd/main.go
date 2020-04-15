package main
 import (
 	"github.com/liao725636367/simlog"
	 "sync/atomic"
	 "time"
 )
var log simlog.Logger
var sub  int64 = 0
func Sub()int64{
	atomic.AddInt64(&sub,1)
	return sub
}
func main() {
	//控制台日志
	log =simlog.NewConsoleLogger(
		simlog.DEBUG,//记录的日志等级 等级排序为 DEBUG INFO WARNING ERROR FATAL 配置的当前日志一下等级的日志不会记录
		)
		log.Debug("调试一下 %v",Sub()) //日志参数 跟 fmt.Sprintf 参数一样
		log.Info("调试一下 %v",Sub())
		log.Warning("调试一下 %v",Sub())
		log.Error("调试一下 %v",Sub())
		log.Fatal("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Error("出现错误 %v",Sub())
		time.Sleep(3*time.Second)
	//文件日志
	log  = simlog.NewFileLogger( //实例化日志
		simlog.DEBUG,//记录的日志等级 等级排序为 DEBUG INFO WARNING ERROR FATAL 配置的当前日志一下等级的日志不会记录
		"./log/",//日志目录 不存在的目录会自动创建
		"log", //日志文件名称，不需要带文件后缀
		2*1024,//文件大小超过此字节数会将日志切割多个文件 0 表示不切割文件
		simlog.SplitDir, //日志是不同等级日志分割方式  SplitDir 按照目录分割 SplitFile 按照文件名分割 SplitNone 所有日志用一个文件名
		100, //缓存日志数量
		)
	for{
		log.Debug("调试一下 %v",Sub()) //日志参数 跟 fmt.Sprintf 参数一样
		log.Info("调试一下 %v",Sub())
		log.Warning("调试一下 %v",Sub())
		log.Error("调试一下 %v",Sub())
		log.Fatal("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Debug("调试一下 %v",Sub())
		log.Error("出现错误 %v",Sub())
		time.Sleep(3*time.Second)
	}
}
