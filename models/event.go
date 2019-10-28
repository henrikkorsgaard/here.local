package models

import "time"

type DeviceEvent struct {
	Event       string
	DeviceMAC   string
	LocationMAC string
	Timestamp   time.Time
}

const (
	DEVICE_JOINED = "device joined"
	DEVICE_LEFT   = "device left"
)
