package database_rules

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

type DataBaseChecker struct {
	errorCheck    bool
	statusChecker *DataBaseStatusChecker
}

func NewDataBaseChecker(errorCheck bool, statusChecker *DataBaseStatusChecker) (*DataBaseChecker, error) {
	if !errorCheck && statusChecker == nil {
		return nil, fmt.Errorf("no rules, DataBaseChecker no point")
	}
	return &DataBaseChecker{
		errorCheck:    errorCheck,
		statusChecker: statusChecker,
	}, nil
}

func (d *DataBaseChecker) Check(event events.Event) (*rules.Alert, error) {
	al := &rules.Alert{
		Type: rules.TypeAlertNormal,
	}
	if d.errorCheck {
		if err, ok := event.Payload.(error); ok {
			al.Data = err.Error()
			if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "many clients") ||
				strings.Contains(err.Error(), "bad connection") || errors.Is(err, sql.ErrConnDone) || strings.Contains(err.Error(), "failed to connect") {
				al.Type = rules.TypeAlertCritical
			} else {
				al.Type = rules.TypeAlertWarning
			}
		}
	}
	if !d.errorCheck {
		if d.statusChecker != nil {
			if _, ok := event.Payload.(error); ok {
				return nil, nil
			}
			str, status, err := d.statusChecker.Check(event.Payload)
			if err != nil {
				return nil, err
			}
			al.Data = str
			al.Type = status
		}
	}
	if al.Data == "" {
		return nil, nil
	}
	return al, nil
}

func (d *DataBaseChecker) ErrorCheck() bool {
	return d.errorCheck
}

func (d *DataBaseChecker) StatusCheck() bool {
	return d.statusChecker != nil
}
