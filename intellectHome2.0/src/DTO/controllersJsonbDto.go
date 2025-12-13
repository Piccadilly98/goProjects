package dto

import (
	"time"
)

type ControllersJsonbDto struct {
	Binary   []BinaryController `json:"binary"`
	Sensor   []SensorController `json:"sensor"`
	MetaData MetaData           `json:"meta_data"`
}

type BinaryController struct {
	ControllerID   *string    `json:"controller_id"`
	Status         *bool      `json:"status"`
	ControllerType *string    `json:"type"`
	Name           *string    `json:"name"`
	CreatedDate    *time.Time `json:"created_date"`
	UpdatedDate    *time.Time `json:"updated_date"`
	PinNumber      *int       `json:"pin_number"`
}

func (bc *BinaryController) Validate() bool {
	if bc.ControllerID == nil || *bc.ControllerID == "" {
		return false
	}
	if bc.Status == nil {
		b := false
		bc.Status = &b
	}
	if bc.ControllerType == nil || *bc.ControllerType == "" {
		return false
	}
	if bc.CreatedDate == nil {
		time := time.Now()
		bc.CreatedDate = &time
	}
	return true
}

type SensorController struct {
	ControllerID   *string    `json:"controller_id"`
	ControllerType *string    `json:"type"`
	Name           *string    `json:"name"`
	Value          *int       `json:"value"`
	Unit           *string    `json:"unit"`
	CreatedDate    *time.Time `json:"created_date"`
	UpdatedDate    *time.Time `json:"updated_date"`
	PinNumber      *int       `json:"pin_number"`
}

func (sc *SensorController) Validate() bool {
	if sc.ControllerID == nil || *sc.ControllerID == "" {
		return false
	}
	if sc.ControllerType == nil || *sc.ControllerType == "" {
		return false
	}
	if sc.Value == nil {
		i := 0
		sc.Value = &i
	}
	if sc.Unit == nil || *sc.Unit == "" {
		return false
	}
	if sc.CreatedDate == nil {
		time := time.Now()
		sc.CreatedDate = &time
	}
	return true
}

type MetaData struct {
	ConfigVersion *string    `json:"config_version"`
	LastUpdate    *time.Time `json:"last_update"`
}
