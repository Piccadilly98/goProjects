package board_info_rules

import "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"

type TemperatureCpuCheck struct {
	WarningHighTemp  float64
	CriticalHighTemp float64
	WarningLowTemp   float64
	CriticalLowTemp  float64
}

func NewTemperatureCpuCheck() *TemperatureCpuCheck {
	return &TemperatureCpuCheck{
		WarningHighTemp:  DefaultWarningHighCpuTemp,
		CriticalHighTemp: DefaultCriticalHighCpuTemp,
		WarningLowTemp:   DefaultWarningLowCpuTemp,
		CriticalLowTemp:  DefaultCriticalLowCpuTemp,
	}
}

func (tc *TemperatureCpuCheck) Check(temp float64) (string, string) {
	status := rules.TypeAlertNormal
	text := ""
	if temp <= tc.WarningLowTemp && temp > tc.CriticalLowTemp {
		status = rules.TypeAlertWarning
		text = "warning low temp cpu"
	}
	if temp <= tc.CriticalLowTemp {
		status = rules.TypeAlertCritical
		text = "critical low temp cpu"
	}
	if temp >= tc.WarningHighTemp && temp < tc.CriticalHighTemp {
		status = rules.TypeAlertWarning
		text = "warning high temp cpu"
	}
	if temp >= tc.CriticalHighTemp {
		status = rules.TypeAlertCritical
		text = "critical high temp cpu"
	}
	return status, text
}
