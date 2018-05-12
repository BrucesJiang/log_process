package main

import (
	"strings"
	"fmt"
	"time"
)

/**
 * 定义读取器接口， 便于扩展
 */
type Reader interface {
	Read(rc chan interface{})
}

/**
 * 定义写入器接口，便于扩展
 */
type Writer interface {
	Write(wc chan interface{})
}


/**
 * 定义一个读取器
 */
type ReadFromFile struct {
	path string //文件路径
}


/**
 * 定义一个写入器
 */
type WriteIntoInfluxDB struct {
	influxDBDsn string //influx data source
}


type LogProcess struct {
	rc chan interface{}
	wc chan interface{}
	reader Reader //读取器
	writer Writer //写入器
}

func (r *ReadFromFile)Read(rc chan interface{}) {
	rc <- "Message"
}

func (l *LogProcess)Process() {
	message := <- l.rc
	//convert interface{} to string
	msg := fmt.Sprintf("%v", message)
	l.wc <- strings.ToUpper(msg)
}

func (w *WriteIntoInfluxDB)Write(wc chan interface{}) {
	fmt.Println(<-wc)
}

func main() {
	r := &ReadFromFile{
		path: "./access.log",
	}

	w := &WriteIntoInfluxDB{
		influxDBDsn: "bruce@bruce...",
	}


	lp := &LogProcess{
		make(chan interface{}),
		make(chan interface{}),
		r,
		w,
	}

	go lp.reader.Read(lp.rc)
	go lp.Process()
	go lp.writer.Write(lp.wc)
	time.Sleep(time.Second * 2)
	fmt.Println("H")
}
