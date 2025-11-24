package dto

import "time"

type GetBoardInfoDTO struct {
	BoardId          *string    `json:"board_id"`
	CreatedDate      *time.Time `json:"created_date"`
	UpdatedDate      *time.Time `json:"updated_date"`
	CpuTemp          *float64   `json:"cpu_temp"`
	AvalibleRam      *int       `json:"avalible_ram"`
	RssiWifi         *int       `json:"rssi_wifi"`
	TotalRunTime     *int       `json:"runtime"`
	IpAddress        *string    `json:"ip"`
	Voltage          *float64   `json:"voltage"`
	TotalDeviceCount *int       `json:"device_count"`
	MacAddress       *string    `json:"mac_address"`
}
