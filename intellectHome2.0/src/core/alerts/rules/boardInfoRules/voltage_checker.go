package board_info_rules

import "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"

type VoltageChecker struct {
	WarningHigh  float64
	CriticalHigh float64
	WarningLow   float64
	CriticalLow  float64
}

func NewVoltageChecker(warningHigh, criticalHigh, warningLow, criticalLow float64) *VoltageChecker {
	if warningHigh == 0 {
		warningHigh = DefaultWarningHighVoltage
	}
	if criticalHigh == 0 {
		criticalHigh = DefaultCriticalHighVoltage
	}
	if warningLow == 0 {
		warningLow = DefaultWarningLowVoltage
	}
	if criticalLow == 0 {
		criticalLow = DefaultCriticalLowVoltage
	}

	return &VoltageChecker{
		WarningHigh:  warningHigh,
		WarningLow:   warningLow,
		CriticalHigh: criticalHigh,
		CriticalLow:  criticalLow,
	}
}

func (v *VoltageChecker) Check(voltage float64) (string, string) {
	status := rules.TypeAlertNormal
	text := ""
	if voltage <= v.WarningLow && voltage > v.CriticalLow {
		status = rules.TypeAlertWarning
		text = "warning low voltage"
	}
	if voltage < v.CriticalLow {
		status = rules.TypeAlertCritical
		text = "critical low voltage"
	}

	if voltage >= v.WarningHigh && voltage < v.CriticalHigh {
		status = rules.TypeAlertWarning
		text = "warning high voltage"
	}
	if voltage >= v.CriticalHigh {
		status = rules.TypeAlertCritical
		text = "critical high voltage"
	}
	return status, text
}
