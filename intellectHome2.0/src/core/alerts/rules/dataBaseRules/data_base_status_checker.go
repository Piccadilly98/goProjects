package database_rules

import (
	"fmt"
	"strings"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
)

type DataBaseStatusChecker struct {
	StatusStartRecovery bool
	StatusNotRecovered  bool
	StatusFinishRecover bool
	StatusOK            bool
}

func NewDataBaseStatusChecker(StatusStartRecovery, StatusNotRecovered, StatusFinishRecover, StatusOK bool) (*DataBaseStatusChecker, error) {
	if !StatusFinishRecover && !StatusStartRecovery && !StatusNotRecovered && !StatusOK {
		return nil, fmt.Errorf("no rules, dataBaseStatusChecker no point")
	}
	return &DataBaseStatusChecker{
		StatusStartRecovery: StatusStartRecovery,
		StatusNotRecovered:  StatusNotRecovered,
		StatusFinishRecover: StatusFinishRecover,
		StatusOK:            StatusOK,
	}, nil
}

func (d *DataBaseStatusChecker) Check(payload any) (string, string) {
	res := ""
	status := ""
	str, ok := payload.(string)
	if !ok {
		return "", ""
	}

	if d.StatusStartRecovery && strings.Contains(str, DataBaseFail) {
		res = str
		status = rules.TypeAlertWarning
	}
	if d.StatusNotRecovered && strings.Contains(str, DataBaseNotRecover) {
		res = str
		status = rules.TypeAlertCritical
	}
	if d.StatusFinishRecover && strings.Contains(str, DataBaseFinishRecover) {
		res = str
		status = rules.TypeAlertNormal
	}
	if d.StatusOK && strings.Contains(str, DataBaseStatusOK) {
		res = str
		status = rules.TypeAlertNormal
	}
	return status, res
}
