package rpi_sensor

import "github.com/stianeikeland/go-rpio/v4"

type PinStatus struct {
	Pin   rpio.Pin   `json:"pin"`
	Mode  rpio.Mode  `json:"mode"`
	State rpio.State `json:"state"`
}

type PinChangeRequest struct {
	Mode  rpio.Mode  `json:"mode"`
	State rpio.State `json:"state"`
}

type Pins map[rpio.Pin]*PinStatus
