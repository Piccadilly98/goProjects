package models

import "fmt"

type Device_data struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	BoadrId string `json:"boardID"`
}

func (d *Device_data) String() string {
	return fmt.Sprintf("ID: %s  Status: %s\nIn Board: %s\n", d.ID, d.Status, d.BoadrId)
}
