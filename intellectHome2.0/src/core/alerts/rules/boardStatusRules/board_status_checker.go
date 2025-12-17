package board_status_rules

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

type BoardStatusChecker struct {
	active  bool
	lost    bool
	offline bool
}

func NewBoardStatusChecker(active, lost, offline bool) (*BoardStatusChecker, error) {
	if !active && !lost && !offline {
		return nil, fmt.Errorf("no rules, BoardStatusChecker no point")
	}
	return &BoardStatusChecker{
		active:  active,
		lost:    lost,
		offline: offline,
	}, nil
}

func (b *BoardStatusChecker) Check(event events.Event) (*rules.Alert, error) {
	str, ok := event.Payload.(string)
	if !ok {
		return nil, fmt.Errorf("invalid type in payload")
	}
	if event.BoardID == "" {
		return nil, fmt.Errorf("invalid boardID")
	}
	strSlice := strings.Fields(str)
	al := &rules.Alert{
		Type:    rules.TypeAlertNormal,
		BoardID: &event.BoardID,
	}

	if b.active {
		if slices.Contains(strSlice, "active") {
			al.Data = str
		}
	}
	if b.lost {
		if slices.Contains(strSlice, "lost") {
			al.Data = str
			al.Type = rules.TypeAlertWarning
		}
	}
	if b.offline {
		if slices.Contains(strSlice, "offline") {
			al.Data = str
			al.Type = rules.TypeAlertCritical
		}
	}
	if al.Data == "" {
		return nil, nil
	}
	return al, nil
}
