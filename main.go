package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)
//var filePath string ="/home/ubuntu/xsys/url.txt"
var filePath string ="url.txt"

const (
	calldepth = 3
)

var logger Logger = &defaultLogger{log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
}

type Handler interface {
	Handle(v ...interface{})
}

func SetLogger(logger Logger) {
	logger = logger
}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}
func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	logger.Info(v...)
}
func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

func Warn(v ...interface{}) {
	logger.Warn(v...)
}
func Warnf(format string, v ...interface{}) {
	logger.Warnf(format, v...)
}

func Error(v ...interface{}) {
	logger.Error(v...)
}
func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}
func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

func Panic(v ...interface{}) {
	logger.Panic(v...)
}
func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v...)
}

func Handle(v ...interface{}) {
	if handle, ok := logger.(Handler); ok {
		handle.Handle(v...)
	}
}

type defaultLogger struct {
	*log.Logger
}

func (l *defaultLogger) Debug(v ...interface{}) {
	l.Output(calldepth, header("DEBUG", fmt.Sprint(v...)))
}

func (l *defaultLogger) Debugf(format string, v ...interface{}) {
	l.Output(calldepth, header("DEBUG", fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Info(v ...interface{}) {
	l.Output(calldepth, header(color.GreenString("INFO "), fmt.Sprint(v...)))
}

func (l *defaultLogger) Infof(format string, v ...interface{}) {
	l.Output(calldepth, header(color.GreenString("INFO "), fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Warn(v ...interface{}) {
	l.Output(calldepth, header(color.YellowString("WARN "), fmt.Sprint(v...)))
}

func (l *defaultLogger) Warnf(format string, v ...interface{}) {
	l.Output(calldepth, header(color.YellowString("WARN "), fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Error(v ...interface{}) {
	l.Output(calldepth, header(color.RedString("ERROR"), fmt.Sprint(v...)))
}

func (l *defaultLogger) Errorf(format string, v ...interface{}) {
	l.Output(calldepth, header(color.RedString("ERROR"), fmt.Sprintf(format, v...)))
}

func (l *defaultLogger) Fatal(v ...interface{}) {
	l.Output(calldepth, header(color.MagentaString("FATAL"), fmt.Sprint(v...)))
	os.Exit(1)
}

func (l *defaultLogger) Fatalf(format string, v ...interface{}) {
	l.Output(calldepth, header(color.MagentaString("FATAL"), fmt.Sprintf(format, v...)))
	os.Exit(1)
}

func (l *defaultLogger) Panic(v ...interface{}) {
	l.Logger.Panic(v...)
}

func (l *defaultLogger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(format, v...)
}

func (l *defaultLogger) Handle(v ...interface{}) {
	l.Error(v...)
}

func header(lvl, msg string) string {
	return fmt.Sprintf("%s: %s", lvl, msg)
}


func main() {
	//SetLogger(&defaultLogger{})
	logger.Info("fs")
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}
		line = strings.ReplaceAll(line,"\n","")
		name := strings.ReplaceAll(line,"http://ppp.downloadxx.com/assets/","")
		var vlen int64 = 0
		downloadFile(line,name, func(length, downLen int64) {
			//logger.Infof("file: %s size %d,ed %d ,percent %v \n",name,length/1024/1024, downLen, float32(downLen)/float32(length))
			//time.Sleep(time.Duration(100) * time.Millisecond)
			vlen = length
		})

		logger.Infof("%s is downloaded..  [%d MB]",name,vlen/1024/1024 )
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

}
func downloadFile(url,savename string, fb func(length, downLen int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	//创建一个http client
	client := new(http.Client)
	//get方法获取资源
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	//读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	//创建文件
	file, err := os.Create(savename)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		//没有错误了快使用 callback

		fb(fsize, written)

	}
	return err
}