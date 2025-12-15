package board_status_rules

import (
	"fmt"
	"strings"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

type BoardStatusChecker struct {
	active  bool
	lost    bool
	offline bool
}

func NewBoardStatusChecker(active, lost, offline bool) *BoardStatusChecker {
	return &BoardStatusChecker{
		active:  active,
		lost:    lost,
		offline: offline,
	}
}

func (b *BoardStatusChecker) Check(event events.Event) (*rules.Alert, error) {
	str, ok := event.Payload.(string)
	if !ok {
		return nil, fmt.Errorf("invalid type in payload")
	}
	al := &rules.Alert{
		Type:    rules.TypeAlertNormal,
		BoardID: &event.BoardID,
	}

	if b.active {
		if strings.Contains(str, "active") {
			al.Data = str
		}
	}
	if b.lost {
		if strings.Contains(str, "lost") {
			al.Data = str
			al.Type = rules.TypeAlertWarning
		}
	}
	if b.offline {
		if strings.Contains(str, "offline") {
			al.Data = str
			al.Type = rules.TypeAlertCritical
		}
	}
	if al.Data == "" {
		return nil, nil
	}
	return al, nil
}
