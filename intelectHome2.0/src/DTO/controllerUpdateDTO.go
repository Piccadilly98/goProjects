package dto

import (
	"time"
)

type ControllerUpdateDTO struct {
	//
	ControllerID string  `json:"controller_id"`
	Name         *string `json:"name"`
	PinNumber    *int    `json:"pin_number"`
	Type         *string `json:"type"`
	UpdatedDate  time.Time

	//binary
	Status *bool `json:"status"`

	//sensor
	Value *int    `json:"value"`
	Unit  *string `json:"unit"`
}

func (cu *ControllerUpdateDTO) Validate() bool {
	if cu.ControllerID == "" {
		return false
	}

	if cu.Status != nil && (cu.Value != nil || cu.Unit != nil) {
		return false
	}

	return (cu.Status != nil) || (cu.Name != nil && *cu.Name != "") || (cu.PinNumber != nil && *cu.PinNumber > 0) || (cu.Type != nil && *cu.Type != "") ||
		(cu.Value != nil) || (cu.Unit != nil && *cu.Unit != "")
}
