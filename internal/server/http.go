package server

import (
	"encoding/json"
	"iot-server/internal"
	"log"
	"net/http"
)

type MetricHandler struct {
	App *App
}

func CreateMetricHandler(app *App) MetricHandler {
	return MetricHandler{
		App: app,
	}
}

func (h MetricHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		var metric internal.Metric
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&metric)
		if err != nil {
			log.Println("Can't process POST")
			w.WriteHeader(400)
			return
		}
		if err = h.App.InfluxHandler.WriteMetric(metric); err != nil {
			log.Printf("Can't write metric: %s\n", err.Error())
			w.WriteHeader(500)
			return
		}
	}
}
