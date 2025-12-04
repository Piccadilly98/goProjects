package dto

import (
	"time"
)

const (
	typeBinary = "binary"
	typeSensor = "sensor"
)

type ControllerUpdateDTO struct {
	//
	Name        *string `json:"name"`
	PinNumber   *int    `json:"pin_number"`
	Type        *string `json:"type"`
	UpdatedDate time.Time

	//binary
	Status *bool `json:"status"`

	//sensor
	Value *int    `json:"value"`
	Unit  *string `json:"unit"`
}

func (cu *ControllerUpdateDTO) ValidateWithType(controllerType string) bool {

	if controllerType == typeBinary && (cu.Value != nil || cu.Unit != nil) {
		return false
	}
	if controllerType == typeSensor && cu.Status != nil {
		return false
	}

	return (cu.Status != nil) || (cu.Name != nil && *cu.Name != "") || (cu.PinNumber != nil && *cu.PinNumber > 0) || (cu.Type != nil && *cu.Type != "") ||
		(cu.Value != nil) || (cu.Unit != nil && *cu.Unit != "")
}
