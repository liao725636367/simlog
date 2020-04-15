# simlog
## 简单的日志库 实现了 多等级记录日志 日志自动按照配置大小切割，异步写日志，日志文件缓存写入

1. 支持向不同地方写日志
   1. 控制台输出日志 -适用于开发环境
   2. 文件输出日志 --适合线上环境
2. 日志分级别
   1. debug 调试
   2. info 普通信息
   3. warning 警告
   4. error 错误
   5. fatal 严重错误
3. 日志支持等级开关控制，根据配置日志等级 只输出 该等级r日志 或者更高等级的日志日志信息
4. 完整的日志记录信息有 时间、 行号、 文件名、 日志级别 、日志信息 
5. 日志可以设置按照文件大小切割
6. 异步日志写入 ，调用日志库只记录日志暂不写入日志，后台  goroutine 定期写入日志，不过多影响业务代码
7. 日志缓存 可以设置缓存日志数量，减少大量写入日志 对磁盘的 io



## 安装

```shell
go mod download github.com/liao725636367/simlog 
```



## 使用

1.命令行输出方式

```golang
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
}

```

​	效果：

​	![image-20200415174405629](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20200415174405629.png) 

2.记录到文件方式

```golang
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

```

​	效果：

​	![image-20200415174435625](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20200415174435625.png) 

![image-20200415174450333](C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20200415174450333.png) 