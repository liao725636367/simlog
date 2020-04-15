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
3. 日志支持等级开关控制，根据配置日志等级 只输出 该等级或者更高等级的日志日志信息
4. 完整的日志记录信息有 时间 行号 文件名 日志级别 日志信息 
5. 日志可以设置按照文件大小切割
6. 异步日志写入 ，调用日志库只记录日志暂不写入日志，后台  goroutine 定期写入日志，不过多影响业务代码
7. 日志缓存 可以设置缓存日志数量，减少大量写入日志 对磁盘的 io
