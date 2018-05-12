package main

import (
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"time"
)

const (
	DB = "log_process"
	username = "bruce"
	password = "bruce"
)

func main() {
	//create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	//Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:DB,
		Precision:"s",
	})

	if err != nil {
		log.Fatal(err)
	}


	// Create a point and add to batch
	//172.0.0.12 - - [04/May/2018:17:56:59 +0000] http "GET /foo HTTP/1.0" 200 2427 "-" "KeepAliveClient" "-" - 2.164

	loc, _ := time.LoadLocation("Asia/Shanghai")

	t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", "04/May/2018:17:56:59 +0000", loc)
	tags := map[string]string{"Path": "/foo", "Method": "GET", "Scheme": "KeepAliveClient", "Status": string(200)}
	fields := map[string]interface{}{
		"UpstreamTime":   2.164,
		"RequestTime": 2.164,
		"BytesSent":   2427,
	}

	pt, err := client.NewPoint("log_info", tags, fields, t)
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

	// Close client resources
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
	log.Println("Write Success")

}
