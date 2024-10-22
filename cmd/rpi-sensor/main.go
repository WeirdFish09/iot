package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
	"iot-server/internal"
	rpi "iot-server/internal/rpi-sensor"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var app rpi.App

	SPIHandler, err := rpi.HWInit()
	if err != nil {
		log.Fatal(err)
	}
	app.SPIHandler = SPIHandler
	if err = rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()
	app.CurrentPins = rpi.InitPins()
	http.Handle("/info", &rpi.InfoHandler{App: app})
	http.Handle("/pinstate/{num}", &rpi.PinStateHandler{App: app})
	http.Handle("/pin/{num}", &rpi.PinHandler{App: app})
	http.Handle("/temp", &rpi.TempHandler{App: app})

	fmt.Println("Initialized.")

	go func(app rpi.App) {
		for {
			val, err := app.SPIHandler.ReadTemp()
			if err != nil {
				fmt.Println(err)
			} else {
				var metric internal.Metric
				metric.Value = float32(val)
				metric.Time = time.Now().Unix()
				metric.Name = "temperature"
				json_data, err := json.Marshal(metric)

				if err != nil {
					fmt.Println("Error processing json: %v", err)
				}
				resp, err := http.Post(os.Getenv("SERVER_URL")+"/metric", "application/json",
					bytes.NewBuffer(json_data))

				if err != nil || resp.StatusCode != http.StatusOK {
					fmt.Println("Error sending req: %v", err)
					fmt.Println("Status code: %v", resp.StatusCode)
				}
			}
			time.Sleep(5 * time.Second)
		}
	}(app)

	log.Fatal(http.ListenAndServe(":7000", nil))
}
