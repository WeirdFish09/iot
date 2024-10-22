package main

import (
	"iot-server/internal/server"
	"log"
	"net/http"
)

func main() {
	var app server.App
	app.InfluxHandler = server.CreateInfluxHandler()

	http.Handle("/metric", server.CreateMetricHandler(&app))

	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}
