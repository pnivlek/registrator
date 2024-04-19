package influxdb

import (
	"context"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

type Metrics struct {
	ServiceName   string
	ContainerID   string
	HostName      string
	ServicePort   int
	ServiceIP     string
	ServiceStatus string
	ServiceTags   []string
	ServiceID     string
}

type InfluxDBClient struct {
	BucketName  string
	InfluxToken string
	InfluxdbURL string
	OrgName     string
	client      influxdb2.Client
}

func New() InfluxDBClient {
	return InfluxDBClient{
		BucketName: os.Getenv("bucket"),
		OrgName:    os.Getenv("org_name"),
		client:     influxdb2.NewClient(os.Getenv("influx_url"), os.Getenv("influx_token")),
	}
}

func (c *InfluxDBClient) WriteData(metrics *Metrics) {
	client := c.client
	defer client.Close()
	// Use blocking write client for writes to desired bucket

	writeAPI := client.WriteAPIBlocking(c.OrgName, c.BucketName)
	// write some points

	p := influxdb2.NewPointWithMeasurement("stat").
		AddTag("service_name", metrics.ServiceName).
		AddTag("container_id", metrics.ContainerID).
		AddField("host", metrics.HostName).
		AddField("port", metrics.ServicePort).
		AddField("ip", metrics.ServiceIP).
		AddField("status", metrics.ServiceStatus).
		AddField("tags", metrics.ServiceTags).
		SetTime(time.Now())

	// write synchronously
	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		log.Println("Influx: Error writing to influxdb. Error is: ", err)
	}
}
