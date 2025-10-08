package models

import "fmt"

type Device_data struct {
	ID     string `json:"deviceID"`
	Status string `json:"deviceStatus"`
}

func (d *Device_data) String() string {
	return fmt.Sprintf("ID: %s  Status: %s\n", d.ID, d.Status)
}
