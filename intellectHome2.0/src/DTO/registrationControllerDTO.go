package dto

import "time"

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
