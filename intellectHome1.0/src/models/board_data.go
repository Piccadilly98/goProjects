package models

import (
	"fmt"
	"time"
)

type DataBoard struct {
	BoardId        string  `json:"boardId"`
	TempCP         float32 `json:"tempCP"`
	FreeMemory     int     `json:"freeMemory"`
	WorkTimeSecond int64   `json:"workTime"`
	RSSI           int     `json:"rssi"`
	LocalIP        string  `json:"localIP"`
	NetworkIP      string  `json:"networkIP"`
	BatteryVoltage float32 `json:"voltage"`
	QuantityDevice int     `json:"quantityDevice"`
	TimeUpload     time.Time
	TimeAdded      time.Time
}

func (d *DataBoard) String() string {
	return fmt.Sprintf("Board ID: %s\nTempCP: %.2f\nFreeMemory: %d\nTimeWork: %d\nRSSI: %d\nLocalIP: %s\nNetworkIP: %s\nVoltage: %.2f\nQuantityDevice: %d\nTime upload: %v\nTime added: %v",
		d.BoardId, d.TempCP, d.FreeMemory, d.WorkTimeSecond, d.RSSI, d.LocalIP, d.NetworkIP, d.BatteryVoltage, d.QuantityDevice, d.TimeUpload.String(), d.TimeAdded.String())
}
