package database_rules

import (
	"fmt"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

type DataBaseChecker struct {
	statusChecker *DataBaseStatusChecker
	errorChecker  *ErrorDBChecker
}

func NewDataBaseChecker(statusChecker *DataBaseStatusChecker, errorChecker *ErrorDBChecker) (*DataBaseChecker, error) {
	if statusChecker == nil && errorChecker == nil {
		return nil, fmt.Errorf("no rules, DataBaseChecker no point")
	}
	return &DataBaseChecker{
		statusChecker: statusChecker,
		errorChecker:  errorChecker,
	}, nil
}

func (d *DataBaseChecker) Check(event events.Event) (*rules.Alert, error) {
	al := &rules.Alert{
		BoardID: &event.BoardID,
	}

	if d.errorChecker != nil {
		al.Type, al.Data = d.errorChecker.Check(event.Payload)
		if al.Type != "" {
			return al, nil
		}
	}
	if d.statusChecker != nil {
		al.Type, al.Data = d.statusChecker.Check(event.Payload)
		if al.Type != "" {
			return al, nil
		}
	}
	return nil, nil
}
