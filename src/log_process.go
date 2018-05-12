package main

import (
	"fmt"
	"time"
	"os"
	"log"
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
	"net/url"
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
	rc     chan interface{}
	wc     chan interface{}
	reader Reader //读取器
	writer Writer //写入器
}

//解析信息字段
type Message struct {
	// 本地时间
	TimeLocal                    time.Time
	//上传数据l量
	BytesSent                    int
	//路径、方法、模式、状态
	Path, Method, Schema, Status string
	//上传时间，响应时间
	UpstreamTime, RequestTime    float64
}

/**
 * 系统分为三个模块，利用接收者与struct进行关联
 */

//读取模块
func (r *ReadFromFile) Read(rc chan interface{}) {
	//rc <- "Message"
	//1.打开文件
	file, err := os.Open(r.path)

	if err != nil {
		log.Panic(fmt.Printf("Fail to Open file: %s\n", err.Error()))
	}

	//从文件末尾开始读取数据, 为了读取最新数据
	file.Seek(0, 2) //将文件指针定义到文件末尾

	br := bufio.NewReader(file)

	//按行读取
	for {
		line, err := br.ReadBytes('\n') //行标识符

		if err == io.EOF { //如果是文件末尾
			// 如果读取到文件末尾，则休眠500ms,继续读取文件
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil { //文件读取失败
			log.Panic(fmt.Sprintf("File Read Error: %s\n", err.Error()))
		}

		//去掉换行符
		//注意传输数据类型的统一
		rc <- string(line[:len(line)-1])
	}
}

//数据解析模块
func (l *LogProcess) Process() {

	//编译正则表达式
	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	//获取本地时间
	loc, _ := time.LoadLocation("Asia/Shanghai")

	//for {
	//	message := <- l.rc
	//	//convert interface{} to string
	//	msg := fmt.Sprintf("%v", message)
	//	l.wc <- strings.ToUpper(msg)
	//}

	//每行的数据
	//172.0.0.12 - - [04/May/2018:17:56:59 +0000] http "GET /foo HTTP/1.0" 200 2427 "-" "KeepAliveClient" "-" - 2.164
	for line := range l.rc {
		//用正则表达式解析没一行内容
		lm := fmt.Sprintf("%v", line)
		ret := r.FindStringSubmatch(lm)

		//没一行可以解析出14子块
		if len(ret) != 14 {
			log.Println("FindStringSubmatch fail:", ret)
			continue
		}

		message := &Message{}
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc)
		if err != nil {
			log.Println("ParseInLocation fail:", err.Error(), ret[4])
			continue
		}
		message.TimeLocal = t

		bytesSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = bytesSent

		//GET /foo?query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			log.Println("strings.Split fail :", ret[6])
			continue
		}
		message.Method = reqSli[0]

		u, err := url.Parse(reqSli[1])
		if err != nil {
			log.Println("url parse fail : ", err.Error())
			continue
		}
		message.Path = u.Path

		message.Schema = ret[5]
		message.Status = ret[7]

		upstreamTime, _ := strconv.ParseFloat(ret[12], 64)
		requestTime, _ := strconv.ParseFloat(ret[13], 64)
		message.UpstreamTime = upstreamTime
		message.RequestTime = requestTime

		l.wc <- message
		//fmt.Println("Process = ", message)
	}

}

//写入模块
func (w *WriteIntoInfluxDB) Write(wc chan interface{}) {
	for line := range wc {
		fmt.Println(line)
	}

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
	time.Sleep(time.Second * 10)
	fmt.Println("H")
}
