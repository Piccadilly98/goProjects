package board_info_rules

import (
	"fmt"

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

type BoardInfoChecker struct {
	rssiChecker    *RssiChecker
	tempChecker    *TemperatureCpuCheck
	voltageChecker *VoltageChecker
}

func NewBoardInfoChecker(rssiChecker *RssiChecker, tempChecker *TemperatureCpuCheck, voltageChecker *VoltageChecker) *BoardInfoChecker {
	// if rssiChecker == nil && tempChecker == nil && voltageChecker == nil {
	// 	return nil
	// }
	return &BoardInfoChecker{
		rssiChecker:    rssiChecker,
		tempChecker:    tempChecker,
		voltageChecker: voltageChecker,
	}
}

func (b *BoardInfoChecker) Check(event events.Event) (*rules.Alert, error) {
	dto, ok := event.Payload.(*dto.UpdateBoardInfo)
	if !ok {
		return nil, fmt.Errorf("invalid type in event")
	}
	alert := &rules.Alert{
		Type:    rules.TypeAlertNormal,
		BoardID: &event.BoardID,
	}
	res := ""
	if b.rssiChecker != nil {
		if dto.RssiWifi == nil {
			return nil, fmt.Errorf("invalid dto")
		}
		status, text := b.rssiChecker.Check(*dto.RssiWifi)
		if status != rules.TypeAlertNormal {
			if len(res) != 0 {
				res += fmt.Sprintf("\n%s", text)
			} else {
				res += text
			}
			if alert.Type != rules.TypeAlertCritical {
				alert.Type = status
			}
		}
	}
	if b.tempChecker != nil {
		if dto.CpuTemp == nil {
			return nil, fmt.Errorf("invalid dto")
		}
		status, text := b.tempChecker.Check(*dto.CpuTemp)
		if status != rules.TypeAlertNormal {
			if len(res) != 0 {
				res += fmt.Sprintf("\n%s", text)
			} else {
				res += text
			}
			if alert.Type != rules.TypeAlertCritical {
				alert.Type = status
			}
		}
	}
	if b.voltageChecker != nil {
		if dto.Voltage == nil {
			return nil, fmt.Errorf("invalid dto")
		}
		status, text := b.voltageChecker.Check(*dto.Voltage)
		if status != rules.TypeAlertNormal {
			if len(res) != 0 {
				res += fmt.Sprintf("\n%s", text)
			} else {
				res += text
			}
			if alert.Type != rules.TypeAlertCritical {
				alert.Type = status
			}
		}
	}
	if alert.Type == rules.TypeAlertNormal {
		return nil, nil
	}
	alert.Data = res
	return alert, nil
}
