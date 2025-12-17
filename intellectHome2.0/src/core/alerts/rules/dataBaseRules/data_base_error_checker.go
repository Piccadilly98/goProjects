package database_rules

import (
	"fmt"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
)

var DefaultWarnings = []string{
	"slow_query_detected",
	"low_disk_space",
	"high_cpu_usage",
	"high_connection_count",
	"long_running_transaction",
	"autovacuum_lag",
	"replication_lag",
	"low_cache_hit_ratio",
	"high_fragmentation",
	"non_optimal_config",
	"stale_statistics",
	"deadlock_detected",
	"no_recent_backup",
	"non_critical_log_errors",
	"high_memory_usage",
}

type ErrorDBChecker struct {
	criticalErr       bool
	warningErr        bool
	otherErrors       bool
	otherErrorsStatus string
	CriticalErrCfg    []string
	WarningErrCfg     []string
}

// default value otherErrorrsStatus - "WARNING"
// if transferred only otherErrorrsStatus, then alerts reaction on all error and set status - otherErrorrsStatus
// otherErrorrsStatus may be just: "normal", "WARNING", "CRITICAL"
func NewErrorDBChecker(criticalErr, warningErr, otherErrors bool, criticalErrCfg, warningErrCfg []string, otherErrorsStatus string) (*ErrorDBChecker, error) {
	if !criticalErr && !warningErr && !otherErrors {
		return nil, fmt.Errorf("no rules, ErrorDBChecker no point")
	}
	if criticalErr && len(criticalErrCfg) == 0 {
		return nil, fmt.Errorf("no critical config")
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

// func (e *ErrorDBChecker) Check(payload any) (string, string) {
// 	status := ""
// 	data := ""

// 	if err, ok := payload.(error); ok {

// 	}
// 	return status, data
// }
