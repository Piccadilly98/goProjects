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

type RegistrationController struct {
	ControllerType string `json:"controller_type"`

	ControllerID string  `json:"controller_id"`
	Name         *string `json:"name"`
	PinNumber    *int    `json:"pin_number"`
	Type         string  `json:"type"`

	//binary
	Status *bool `json:"status"`

	//sensor
	Value *int    `json:"value"`
	Unit  *string `json:"unit"`
}

func (rc *RegistrationController) Validate() bool {
	if rc.ControllerType != "sensor" && rc.ControllerType != "binary" {
		return false
	}
	if rc.ControllerID == "" {
		return false
	}
	if rc.Type == "" {
		return false
	}
	return true
}

func (rc *RegistrationController) ToSensorController() *SensorController {
	if rc.ControllerType != "sensor" {
		return nil
	}

	sc := &SensorController{}
	sc.ControllerID = &rc.ControllerID
	sc.ControllerType = &rc.Type
	sc.Name = rc.Name
	sc.Value = rc.Value
	sc.Unit = rc.Unit
	sc.PinNumber = rc.PinNumber
	time := time.Now()
	sc.CreatedDate = &time
	return sc
}

func (rc *RegistrationController) ToBinaryController() *BinaryController {
	if rc.ControllerType != "binary" {
		return nil
	}

	bc := &BinaryController{}
	bc.ControllerID = &rc.ControllerID
	bc.ControllerType = &rc.Type
	bc.Name = rc.Name
	bc.PinNumber = rc.PinNumber
	bc.Status = rc.Status
	time := time.Now()
	bc.CreatedDate = &time
	return bc
}
