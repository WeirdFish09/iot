package rpi_sensor

import (
	"encoding/json"
	"fmt"
	"github.com/stianeikeland/go-rpio/v4"
	"iot-server/internal"
	"net/http"
	"strconv"
	"time"
)

var (
	initialPins = []rpio.Pin{18, 23, 24, 25}
)

func InitPins() Pins {
	var pins Pins = make(Pins)
	for _, pin := range initialPins {
		pinStatus := &PinStatus{}
		pinStatus.Pin = pin
		pinStatus.Mode = 0  // Read
		pinStatus.State = 0 // Low
		pinStatus.Pin.Input()
		pins[pin] = pinStatus
	}
	return pins
}

type PinHandler struct {
	App App
}

func (handler *PinHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pinRaw := r.PathValue("num")
	pinNum, err := strconv.ParseInt(pinRaw, 10, 8)
	if err != nil {
		handleError(w, 404)
		return
	}
	val, exists := handler.App.CurrentPins[rpio.Pin(pinNum)]
	if !exists {
		handleError(w, 404)
		return
	}
	if r.Method == "GET" {
		if val.Mode != rpio.Input {
			handleError(w, 400)
			return
		}
		res := val.Pin.Read()
		w.WriteHeader(200)
		fmt.Fprint(w, res)
		return
	}
	if r.Method == "POST" {
		if val.Mode != rpio.Output {
			handleError(w, 400)
			return
		}
		val.Pin.Toggle()
		w.WriteHeader(200)
		return
	}
}

type PinStateHandler struct {
	App App
}

func (handler *PinStateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pinRaw := r.PathValue("num")
	pinNum, err := strconv.ParseInt(pinRaw, 10, 8)
	if err != nil {
		handleError(w, 404)
		return
	}
	val, exists := handler.App.CurrentPins[rpio.Pin(pinNum)]

	if r.Method == "GET" {
		if !exists {
			handleError(w, 404)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(val)
		return
	}
	if r.Method == "POST" {
		if !exists {
			handler.App.CurrentPins[rpio.Pin(pinNum)] = &PinStatus{}
			val = handler.App.CurrentPins[rpio.Pin(pinNum)]
			val.Pin = rpio.Pin(pinNum)
		}
		decoder := json.NewDecoder(r.Body)
		var req PinChangeRequest
		err := decoder.Decode(&req)
		if err != nil {
			handleError(w, 400)
			return
		}
		val.State = req.State
		val.Mode = req.Mode
		val.Pin.Mode(val.Mode)
		val.Pin.Write(val.State)
	}
}

type InfoHandler struct {
	App App
}

func (handler *InfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		handleError(w, http.StatusMethodNotAllowed)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(handler.App.CurrentPins)
}

func handleError(w http.ResponseWriter, code int) {
	http.Error(w, "error", code)
}

type TempHandler struct {
	App App
}

func (handler *TempHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		handleError(w, http.StatusMethodNotAllowed)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	val, err := handler.App.SPIHandler.ReadTemp()
	if err != nil {
		handleError(w, 500)
		return
	}
	var metric internal.Metric
	metric.Value = float32(val)
	metric.Time = time.Now().Unix()
	metric.Name = "temperature"
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(metric)
	return
}
