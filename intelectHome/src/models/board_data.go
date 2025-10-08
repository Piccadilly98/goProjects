package models

import "fmt"

type DataBoard struct {
	BoardId        string  `json:"boardId"`
	TempCP         float32 `json:"tempCP"`
	FreeMemory     int     `json:"freeMemory"`
	WorkTimeSecond int64   `json:"workTime"`
	RSSI           int     `json:"rssi"`
	LocalIP        string  `json:"localIP"`
	NetworkIP      string  `json:"networkIP"`
	BatteryVoltage float32 `json:"voltage"`
}

func (d *DataBoard) String() string {
	return fmt.Sprintf(" Board ID: %s\n TempCP: %.2f\n  FreeMemory: %d\n TimeWork: %d\n RSSI: %d\nLocalIP: %s\nNetworkIP: %s\nVoltage: %.2f\n",
		d.BoardId, d.TempCP, d.FreeMemory, d.WorkTimeSecond, d.RSSI, d.LocalIP, d.NetworkIP, d.BatteryVoltage)
}
