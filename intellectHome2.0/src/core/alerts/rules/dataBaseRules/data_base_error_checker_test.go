package database_rules_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	database_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/dataBaseRules"
)

// strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "many clients") ||
//
//	strings.Contains(err.Error(), "bad connection") || errors.Is(err, sql.ErrConnDone) || strings.Contains(err.Error(), "failed to connect") {
//	al.Type = rules.TypeAlertCritical
func TestNewErrorChecker(t *testing.T) {
	testCases := []struct {
		Name              string
		criticalErr       bool
		warningErr        bool
		otherErrors       bool
		otherErrorsStatus string
		criticalErrCfg    []string
		warningErrCfg     []string
		ExpectedErr       error
	}{
		{"normal_all_check_and_cfg", true, true, true, rules.TypeAlertWarning, []string{"connection refused", "many clients", "bad connection", "failed to connect"}, nil, nil},
		{"normal_no_warning", true, false, true, rules.TypeAlertWarning, []string{"connection refused", "many clients", "bad connection", "failed to connect"}, nil, nil},
		{"normal_no_other", true, true, false, rules.TypeAlertWarning, []string{"connection refused", "many clients", "bad connection", "failed to connect"}, nil, nil},
		{"invalid_other_status", true, true, true, "norm", []string{"connection refused", "many clients", "bad connection", "failed to connect"}, nil, fmt.Errorf("no valid other error status")},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			check, err := database_rules.NewErrorDBChecker(tc.criticalErr, tc.warningErr, tc.otherErrors, tc.criticalErrCfg, tc.warningErrCfg, tc.otherErrorsStatus)
			if err != nil {
				if tc.ExpectedErr != nil {
					if tc.ExpectedErr.Error() != err.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.ExpectedErr.Error())
					}
				} else {
					t.Errorf("unexpected err: %s\n", err.Error())
				}
				return
			}
			if tc.warningErrCfg == nil && tc.warningErr {
				if !slices.Equal(check.WarningErrCfg, database_rules.DefaultWarnings) {
					t.Errorf("invalid warning cfg: %v\n", check.WarningErrCfg)
				}
			}
			if tc.criticalErrCfg != nil && tc.criticalErr {
				if !slices.Equal(tc.criticalErrCfg, check.CriticalErrCfg) {
					t.Errorf("CRITICAL CFG: got: %v, expect: %v\n", check.CriticalErrCfg, tc.criticalErrCfg)
				}
			}
		})
	}
}
