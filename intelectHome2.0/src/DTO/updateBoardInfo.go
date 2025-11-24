package dto

import (
	"time"
)

type UpdateBoardInfo struct {
	CpuTemp          *float64 `json:"cpu_temp"`
	AvalibleRam      *int     `json:"avalible_ram"`
	RssiWifi         *int     `json:"rssi_wifi"`
	TotalRunTime     *int     `json:"runtime"`
	IpAddress        *string  `json:"ip"`
	Voltage          *float64 `json:"voltage"`
	TotalDeviceCount *int     `json:"device_count"`
	MacAddress       *string  `json:"mac_address"`
	TimeUpload       time.Time
}

func (u *UpdateBoardInfo) Validate() bool {
	if u.CpuTemp == nil || *u.CpuTemp == 0 {
		return false
	}
	if u.AvalibleRam == nil || *u.AvalibleRam <= 0 {
		return false
	}
	if u.RssiWifi == nil || *u.RssiWifi == 0 {
		return false
	}
	if u.TotalRunTime == nil || *u.TotalRunTime <= 0 {
		return false
	}
	if u.IpAddress == nil || *u.IpAddress == "" {
		return false
	}
	if u.Voltage == nil || *u.Voltage <= 0 {
		return false
	}
	if u.TotalDeviceCount == nil {
		return false
	}
	if u.MacAddress == nil || *u.MacAddress == "" {
		return false
	}
	u.TimeUpload = time.Now()
	return true
}
