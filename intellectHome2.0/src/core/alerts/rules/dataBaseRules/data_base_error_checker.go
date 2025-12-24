package database_rules

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
)

var DefaultWarnings = []string{
	"duplicate key",
	"violates unique constraint",
	"violates foreign key",
	"violates not-null",
	"violates check constraint",
	"invalid input syntax",
	"invalid byte sequence",
	"malformed array",
	"cannot cast",
	"numeric value out of range",
	"permission denied",
	"must be owner",
	"does not exist",
	"already exists",
	"transaction is aborted",
	"in failed sql transaction",
	"syntax error",
	"ambiguous",
	"operator does not exist",
	"function does not exist",
	"wrong number of parameters",
	"parameter",
	"division by zero",
	"array subscript out of range",
	"string data right truncation",
	"datetime field overflow",
	"cannot be changed",
	"unrecognized configuration",
	"failing row",
	"key (",
}

var DefaultCritial = []string{
	// === ОШИБКИ СОЕДИНЕНИЯ ===
	"connection", "connect", "conn",
	"reset", "refused", "broken", "closed",
	"EOF", "network", "socket", "host", "port",
	"dial", "timeout", "i/o timeout",

	// === СЕТЕВЫЕ ОШИБКИ ===
	"no route to host",
	"network is unreachable",
	"connection refused",
	"connection timed out",

	// === ОШИБКИ СЕРВЕРА ===
	"server closed",
	"server terminated",
	"database system",
	"shutting down",
	"starting up",
	"recovery",
	"crash",

	// === АВТОРИЗАЦИЯ ===
	"password authentication failed",
	"pg_hba.conf",
	"authentication failed",
	"login failed",

	// === РЕСУРСЫ И ЛИМИТЫ ===
	"too many connections",
	"too many clients",
	"out of memory",
	"disk full",
	"cannot allocate",

	// === ФАТАЛЬНЫЕ ОШИБКИ ===
	"fatal",
	"panic",
	"terminating connection",
	"canceling statement",

	// === ОШИБКИ ПРОТОКОЛА ===
	"protocol violation",
	"unexpected message",
	"invalid message",

	// === ОШИБКИ ДОСТУПА К ДАННЫМ ===
	"could not open file",
	"could not read file",
	"could not write",
	"read-only",
	"read only transaction",
}

type ErrorDBChecker struct {
	criticalErr       bool
	warningErr        bool
	otherErrors       bool
	otherErrorsStatus string
	CriticalErrCfg    []string
	WarningErrCfg     []string
	mu                sync.RWMutex
}

// default value otherErrorrsStatus - "WARNING"
// if transferred only otherErrorrsStatus, then alerts reaction on all error and set status - otherErrorrsStatus
// otherErrorrsStatus may be just: "normal", "WARNING", "CRITICAL"
func NewErrorDBChecker(criticalErr, warningErr, otherErrors bool, criticalErrCfg, warningErrCfg []string, otherErrorsStatus string) (*ErrorDBChecker, error) {
	if !criticalErr && !warningErr && !otherErrors {
		return nil, fmt.Errorf("no rules, ErrorDBChecker no point")
	}
	if criticalErr && len(criticalErrCfg) == 0 {
		criticalErrCfg = DefaultCritial
	}
	if warningErr && len(warningErrCfg) == 0 {
		warningErrCfg = DefaultWarnings
	}
	if otherErrorsStatus != rules.TypeAlertNormal && otherErrorsStatus != rules.TypeAlertCritical &&
		otherErrorsStatus != rules.TypeAlertWarning && otherErrorsStatus != "" {
		return nil, fmt.Errorf("no valid other error status")
	}
	if otherErrorsStatus == "" {
		otherErrorsStatus = rules.TypeAlertWarning
	}

	return &ErrorDBChecker{
		criticalErr:       criticalErr,
		warningErr:        warningErr,
		otherErrors:       otherErrors,
		otherErrorsStatus: otherErrorsStatus,
		CriticalErrCfg:    criticalErrCfg,
		WarningErrCfg:     warningErrCfg,
	}, nil
}

func (e *ErrorDBChecker) OtherErrorrsStatus() string {
	return e.otherErrorsStatus
}

func (e *ErrorDBChecker) Check(payload any) (string, string) {

	err, ok := payload.(error)
	if !ok || err == nil {
		return "", ""
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.criticalErr {
		if e.checkCriticalErr(err) {
			return rules.TypeAlertCritical, err.Error()
		}
	}
	if e.warningErr {
		if e.checkWarningErr(err) {
			return rules.TypeAlertWarning, err.Error()
		}
	}
	if e.otherErrors && e.otherErrorsStatus != rules.TypeAlertNormal {
		return e.otherErrorsStatus, err.Error()
	}

	return "", ""
}

func (e *ErrorDBChecker) checkCriticalErr(err error) bool {
	for _, val := range e.CriticalErrCfg {
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower(val)) {
			return true
		}
	}
	return false
}

func (e *ErrorDBChecker) checkWarningErr(err error) bool {
	for _, val := range e.WarningErrCfg {
		if strings.Contains(strings.ToLower(err.Error()), strings.ToLower(val)) {
			return true
		}
	}
	return false
}
