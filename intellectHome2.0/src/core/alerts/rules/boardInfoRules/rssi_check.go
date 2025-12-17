package board_info_rules

import (
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
)

type RssiChecker struct {
	WarningRSSI  int
	CriticalRSSI int
}

func NewRSsiChecker(warning, critical int) *RssiChecker {
	if warning >= 0 || warning < -100 {
		warning = DefaultWarningRSSI
		// log.Printf("Rssi checker: set warning default: %f\n", DefaultWarningRSSI)
	}
	if critical >= 0 || critical < -100 {
		critical = DefaultCriticalRSSI
		// log.Printf("Rssi checker: set critical default: %f\n", DefaultCriticalRSSI)
	}

	return &RssiChecker{
		WarningRSSI:  warning,
		CriticalRSSI: critical,
	}
}

func (r *RssiChecker) Check(rssi int) (string, string) {
	status := rules.TypeAlertNormal
	text := ""
	if rssi <= r.WarningRSSI && rssi > r.CriticalRSSI {
		status = rules.TypeAlertWarning
		text = "warning low ethernet signal"
	}
	if rssi <= r.CriticalRSSI {
		status = rules.TypeAlertCritical
		text = "critical low ethernet signal"
	}
	return status, text
}
