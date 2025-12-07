package tests

import (
	"testing"

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
)

const (
	TypeBinary = "binary"
	TypeSensor = "sensor"
)

type testDTO struct {
	dto    *dto.ControllerUpdateDTO
	Type   string
	Result bool
}

func TestLoginHandler(t *testing.T) {
	tcase := []testDTO{

		{

			dto: &dto.ControllerUpdateDTO{
				Name:      getPtr("stas"),
				PinNumber: getPtrInt(1),
				Type:      getPtr("non_type"),
				Status:    getPtrBool(true),
			},
			Type:   TypeBinary,
			Result: true,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      nil,
				PinNumber: nil,
				Type:      nil,
				Status:    getPtrBool(true),
			},
			Type:   TypeBinary,
			Result: true,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      getPtr("stas"),
				PinNumber: nil,
				Type:      nil,
				Status:    nil,
			},
			Type:   TypeBinary,
			Result: true,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      nil,
				PinNumber: getPtrInt(1),
				Type:      nil,
				Status:    nil,
			},
			Type:   TypeBinary,
			Result: true,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      nil,
				PinNumber: nil,
				Type:      getPtr("s"),
				Status:    nil,
			},
			Type:   TypeBinary,
			Result: true,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      getPtr("stas"),
				PinNumber: nil,
				Type:      nil,
				Status:    nil,
			},
			Type:   TypeSensor,
			Result: true,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      nil,
				PinNumber: getPtrInt(1),
				Type:      nil,
				Status:    nil,
			},
			Type:   TypeSensor,
			Result: true,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      nil,
				PinNumber: nil,
				Type:      getPtr("stas"),
				Status:    nil,
			},
			Type:   TypeSensor,
			Result: true,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      getPtr("stas"),
				PinNumber: nil,
				Type:      nil,
				Status:    nil,
				Value:     getPtrInt(1),
			},
			Type:   TypeSensor,
			Result: true,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      getPtr("stas"),
				PinNumber: nil,
				Type:      nil,
				Status:    nil,
				Value:     getPtrInt(1),
				Unit:      getPtr("%"),
			},
			Type:   TypeSensor,
			Result: true,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      getPtr("stas"),
				PinNumber: getPtrInt(1),
				Type:      getPtr("s"),
				Status:    nil,
				Value:     getPtrInt(1),
				Unit:      getPtr("%"),
			},
			Type:   TypeSensor,
			Result: true,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      getPtr("stas"),
				PinNumber: getPtrInt(1),
				Type:      getPtr("s"),
				Status:    nil,
				Value:     getPtrInt(1),
				Unit:      nil,
			},
			Type:   TypeBinary,
			Result: false,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      getPtr("stas"),
				PinNumber: getPtrInt(1),
				Type:      getPtr("s"),
				Status:    nil,
				Value:     nil,
				Unit:      getPtr("1"),
			},
			Type:   TypeBinary,
			Result: false,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      getPtr("stas"),
				PinNumber: getPtrInt(1),
				Type:      getPtr("s"),
				Status:    getPtrBool(false),
				Value:     getPtrInt(1),
				Unit:      nil,
			},
			Type:   TypeSensor,
			Result: false,
		},
		{
			dto: &dto.ControllerUpdateDTO{
				Name:      getPtr("stas"),
				PinNumber: getPtrInt(1),
				Type:      getPtr("s"),
				Status:    getPtrBool(false),
				Value:     nil,
				Unit:      getPtr("!"),
			},
			Type:   TypeSensor,
			Result: false,
		},
	}

	for i, r := range tcase {
		ok := r.dto.ValidateWithType(r.Type)
		if ok != r.Result {
			t.Errorf("case : %d\ngot %v, expect %v\n", i, ok, r.Result)
		}
	}
}

func getPtr(str string) *string {
	st := &str

	return st
}

func getPtrInt(i int) *int {
	ip := &i
	return ip
}

func getPtrBool(b bool) *bool {
	bp := &b
	return bp
}
