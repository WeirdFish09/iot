package server

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"iot-server/internal"
	"os"
	"time"
)

const (
	ORG    = "iot"
	BUCKET = "sensors"
)

type InfluxHandler struct {
	InfluxClient influxdb2.Client
}

func CreateInfluxHandler() InfluxHandler {
	url := os.Getenv("INFLUXDB_URL")
	token := os.Getenv("INFLUXDB_TOKEN")
	return InfluxHandler{
		InfluxClient: influxdb2.NewClient(url, token),
	}
}

func (ih *InfluxHandler) WriteMetric(metric internal.Metric) error {
	writeApi := ih.InfluxClient.WriteAPIBlocking(ORG, BUCKET)
	p := influxdb2.NewPointWithMeasurement("sensor").
		AddTag("name", metric.Name).
		AddField("value", metric.Value).
		SetTime(time.Unix(metric.Time, 0))
	return writeApi.WritePoint(context.Background(), p)
}
