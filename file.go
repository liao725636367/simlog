package simLog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)
type SplitType uint8
const (
	SplitNone SplitType = iota //不分割日志文件
	SplitFile //按照文件名切割日志
	SplitDir //按照目录名称切割日志

)
type LogMsg struct {
	level logLevel
	msg string //写入内容
}
//文件句柄和buffer结构体
type fileHandler struct {
	File  *os.File
	buffer []string //文件buffer 实现缓冲写入日志
}
//使用互斥锁避免文件操作冲突
var mutex sync.Mutex
//终端写日志内容
type FileLogger struct{
	Level logLevel
	filePath string
	fileName string //日志文件名
	maxFileSize int64
	splitType SplitType //日志是否分开存储
	bufferNum int //日志缓存数量
	bufferCount int //日志已经缓存数量
	fileHandlers map[string]*fileHandler //文件句柄 列表
	logChan chan *LogMsg //日志chan 实现异步日志
	lastFlush  time.Time //最后刷新写入时间 60秒至少写一次日志
}

func NewFileLogger(level logLevel,filePath ,fileName string,maxFileSize int64,splitType  SplitType,bufferNum int ) *FileLogger{
	fl := &FileLogger{
		level,
		filePath,
		fileName,
		maxFileSize,
		splitType,
		bufferNum,
		0,
		make(map[string]*fileHandler),
		make(chan *LogMsg,1000),
		time.Now(),
	}
	 err:=fl.initFile() //初始化文件句柄
	if err!=nil{
		panic(err)
		return nil
	}
	return fl
}
//initFile 初始化文件句柄
func (l *FileLogger) initFile() error {
	//os.OpenFile()
	splitType:=l.splitType

	logPath,err :=filepath.Abs(l.filePath)
	if err!=nil{
		return err
	}
	err=AutoCreatePath(logPath,0644)
	if err!=nil{
		return err
	}
	//fmt.Println(l.filePath,logPath)
	//return nil
	if splitType == SplitNone{ //默认不分割日志
		logFile  :=  filepath.Join(logPath,l.fileName+".log")
		file,err:=os.OpenFile(logFile,os.O_CREATE|os.O_APPEND|os.O_WRONLY,0644)
		if err!=nil{
			return err
		}
		l.fileHandlers["default"]=&fileHandler{
			File:   file,
			buffer: nil,
		}
	}else if splitType == SplitFile{ //按照文件名称分割日志
		levels:=[]logLevel{DEBUG,TRACE,INFO,WARNING,ERROR,FATAL}
		for _,level:=range levels{
			levelStr:=strings.ToLower(level.String())
			logFile:=filepath.Join(logPath,l.fileName+"."+levelStr+".log")
			file,err:=os.OpenFile(logFile,os.O_CREATE|os.O_APPEND|os.O_WRONLY,0644)
			if err!=nil{
				return err
			}
			l.fileHandlers[level.String()]=&fileHandler{
				File:   file,
				buffer: nil,
			}
		}
	}else if splitType == SplitDir{ //按照文件名称分割日志
		levels:=[]logLevel{DEBUG,TRACE,INFO,WARNING,ERROR,FATAL}
		for _,level:=range levels{
			levelStr:=strings.ToLower(level.String())
			sonLogPath := filepath.Join(logPath,levelStr)
			err:=AutoCreatePath(sonLogPath,0644)
			if err!=nil{
				return err
			}
			//fmt.Println(sonLogPath)

			logFile:=filepath.Join(sonLogPath,l.fileName+"."+levelStr+".log")
			file,err:=os.OpenFile(logFile,os.O_CREATE|os.O_APPEND|os.O_WRONLY,0644)
			if err!=nil{
				return err
			}
			l.fileHandlers[level.String()]= &fileHandler{
				File:   file,
				buffer: nil,
			}
		}
	}
	//开启一个后台goroutine写日志
	//for i:=0;i<6;i++{
		go l.backgroundLog()
	//}
	return nil
}
//Fatal 记录调试日志
func (l *FileLogger) Debug(format string, a ...interface{}){
	l.outPut( DEBUG,format,a...)
}
//Fatal 记录严重信息日志
func (l *FileLogger) Info(format string, a ...interface{}){
	l.outPut(  INFO,format,a...)
}
//Fatal 记录报警日志
func (l *FileLogger) Warning(format string, a ...interface{}){
	l.outPut( WARNING,format,a...)
}
//Fatal 记录错误日志
func (l *FileLogger) Error(format string, a ...interface{}){
	l.outPut( ERROR,format,a...)
}
//Fatal 记录严重错误日志
func (l *FileLogger) Fatal(format string, a ...interface{}){
	l.outPut(FATAL,format,a...)
}
//CheckFileSize 检测文件大小是否超出最大限制
func (l *FileLogger) CheckFileSize(f *os.File) bool {
	fileInfo,err:=f.Stat()
	if err!=nil{
		fmt.Printf("get file info failed err:%v\n",err)
		return false
	}
	//fmt.Println(fileInfo.Name(),fileInfo.Size(), l.maxFileSize)
	//panic("结束")
	return fileInfo.Size() > l.maxFileSize
}
//outPut 记录日志
func (l *FileLogger) outPut( level logLevel,format string, a ...interface{}){
	now := time.Now().Format("2006-01-02 15:04:05")
	if l.Level > level { //如果记录级别大于 日志设置级别就不记录日志
		return
	}
	funcName,fileName,lineNo := getInfo(3) //跳过堆栈
	//funcName:=""
	//fileName:=""
	//lineNo:=0
	//格式化信息
	format=fmt.Sprintf(format,a...)

	//fmt.Println(l.fileHandlers[fileIndexStr].Name());
	//return
	str:=fmt.Sprintf("[%s] [%s] [%s:%s:%d]  %s \n",now,level.String(),fileName,funcName,lineNo,format)
	select {
		case l.logChan <- &LogMsg{level:level,msg:str}:
	default:
		//避免channel满了阻塞导致业务代码无法顺畅执行-此时会丢弃多的日志
	}
	//写日志放到后台运行
}
//backgroundLog 后台记录日志
func (l *FileLogger) backgroundLog(){
	var fileIndexStr string
	go func(l *FileLogger) {
		//每60秒写入一下文件
		d:=time.Duration(time.Second*60)
		ticker:=time.NewTicker(d)
		defer  ticker.Stop()
		for   {
			<-ticker.C
			l.flush()
		}
	}(l)
	for {
		select {
			case logMsg := <- l.logChan: //获得channel里面的结构体

				level := logMsg.level
				msg:=logMsg.msg
				if _,ok:=l.fileHandlers[level.String()];ok{

					fileIndexStr = level.String()
				}else{
					fileIndexStr = "default"
				}
				l.fileHandlers[fileIndexStr].buffer = append(l.fileHandlers[fileIndexStr].buffer,msg)
				l.bufferCount++ //缓存每增加日志都计数+1
				if len(l.fileHandlers[fileIndexStr].buffer) >    l.bufferNum  {

					l.flush() //信息写入数据中
				}
				//mutex.Unlock()
		default:
			time.Sleep(time.Millisecond*200) //娶不到日志就休息 200毫秒
		}
	}

}
//flush 缓存中日志写入到磁盘 根据每个级别缓存数据
//fileIndexStr string
func (l *FileLogger) flush()error{
	mutex.Lock() //加锁避免文件操作冲突
	defer mutex.Unlock()
	var fileObj *os.File
	var msg string
	for levelStr,handler:=range l.fileHandlers{
		if len(handler.buffer) > 0{

				fileObj =handler.File
			//mutex.Lock()
			if ok:=l.CheckFileSize(fileObj);ok{

				//4.赋值给当前文件对应 句柄

				fileObj,err:=l.SplitFile(fileObj)
				if err!=nil{
					return err
				}
				handler.File=fileObj

			}
			//组合信息
			msg = strings.Join(handler.buffer,"")

			handler.buffer = nil
			l.lastFlush=time.Now()
			fmt.Println(l.fileHandlers[levelStr])
			fmt.Fprintf(fileObj,  msg)

		}

	}

	return nil

}
//SplitFile 切割日志文件
func (l *FileLogger) SplitFile(file *os.File)(*os.File,error){
	//切割文件
	//1.关闭当前日志文件
	err:=file.Close()
	if err!=nil{
		fmt.Printf("close file failed err:%v\n",err)
		return nil,err
	}
	//2.rename 备份
	nowStr:=time.Now().Format("20060102150405")
	oldName := file.Name()
	// 插入字符串
	//找到末尾索引
	index := len(oldName)-4
	newName := oldName[:index]+"-"+nowStr+oldName[index:]
	//fmt.Println(newName)
	//panic("结束")
	err=os.Rename(oldName,newName)//移动文件
	if err!=nil{
		fmt.Printf("rename file failed err:%v\n",err)
		return nil,err
	}
	//3.打开新文件
	fileObj,err:=os.OpenFile(oldName,os.O_CREATE|os.O_APPEND|os.O_WRONLY,0644)
	if err!=nil{
	fmt.Printf("open file failed err:%v\n",err)
		return nil,err
	}
	return fileObj,nil
}
//Close 关闭日志文件
func (l *FileLogger) Close(){
	for _,file:=range l.fileHandlers{
		_=file.File.Close()
	}
}
