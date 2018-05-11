package main

import (
	"bytes"
	"strings"
	"fmt"
	"time"
)

type LogProcess struct {
	rc chan []byte
	wc chan string
	path string //文件路径
	influxDns string //influx data source
}

func (l *LogProcess)ReadFromFile() {
	l.rc <- bytes.NewBufferString("Message").Bytes()
}

func (l *LogProcess)Process() {
	message := <- l.rc
	l.wc <- strings.ToUpper(string(message))
}

func (l *LogProcess)WriteIntoInfluxDB() {
	fmt.Println(<-l.wc)
}

func main() {
	lp := &LogProcess{
		rc: make(chan []byte),
		wc: make(chan string),
		path: "./access.log",
		influxDns: "bruce@bruce...",

	}

	go lp.ReadFromFile()
	go lp.Process()
	go lp.WriteIntoInfluxDB()
	time.Sleep(time.Second * 2)
	fmt.Println("H")
}
