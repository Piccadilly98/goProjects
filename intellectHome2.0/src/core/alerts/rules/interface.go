package rules

import "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"

const (
	TypeAlertNormal   = "normal"
	TypeAlertWarning  = "WARNING"
	TypeAlertCritical = "CRITICAL"
)

type Rule interface {
	Check(event events.Event) (*Alert, error)
}

type Alert struct {
	Type    string
	BoardID *string
	Data    string
}
