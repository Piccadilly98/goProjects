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
		Name                string
		criticalErr         bool
		warningErr          bool
		otherErrors         bool
		otherErrorsStatus   string
		criticalErrCfg      []string
		warningErrCfg       []string
		ExpectedErr         error
		expectedOtherStatus string
	}{
		{"normal_all_check_and_cfg", true, true, true, rules.TypeAlertWarning, []string{"connection refused", "many clients", "bad connection", "failed to connect"}, nil, nil, rules.TypeAlertWarning},
		{"normal_all_check_and_cfg_status_normal", true, true, true, rules.TypeAlertNormal, []string{"connection refused", "many clients", "bad connection", "failed to connect"}, nil, nil, rules.TypeAlertNormal},
		{"normal_all_check_and_cfg_status_critical", true, true, true, rules.TypeAlertCritical, []string{"connection refused", "many clients", "bad connection", "failed to connect"}, nil, nil, rules.TypeAlertCritical},
		{"normal_all_check_and_cfg_empty_status", true, true, true, "", []string{"connection refused", "many clients", "bad connection", "failed to connect"}, nil, nil, rules.TypeAlertWarning},
		{"normal_no_critical", false, true, true, rules.TypeAlertWarning, []string{"connection refused", "many clients", "bad connection", "failed to connect"}, nil, nil, rules.TypeAlertWarning},
		{"normal_no_warning", true, false, true, rules.TypeAlertWarning, []string{"connection refused", "many clients", "bad connection", "failed to connect"}, nil, nil, rules.TypeAlertWarning},
		{"normal_no_other", true, true, false, rules.TypeAlertWarning, []string{"connection refused", "many clients", "bad connection", "failed to connect"}, nil, nil, rules.TypeAlertWarning},
		{"valid_no_critical_cfg", true, true, false, rules.TypeAlertWarning, nil, nil, nil, rules.TypeAlertWarning},
		{"invalid_no_param", false, false, false, rules.TypeAlertWarning, nil, nil, fmt.Errorf("no rules, ErrorDBChecker no point"), rules.TypeAlertWarning},
		{"invalid_other_status", true, true, true, "norm", []string{"connection refused", "many clients", "bad connection", "failed to connect"}, nil, fmt.Errorf("no valid other error status"), rules.TypeAlertWarning},
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
			if tc.expectedOtherStatus != check.OtherErrorrsStatus() {
				t.Errorf("OTHER STATUS: got: %s, epect: %s\n", check.OtherErrorrsStatus(), tc.expectedOtherStatus)
			}
		})
	}
}

type TestCases struct {
	Name           string
	Payload        any
	ExpectedStatus string
	ExpectedData   string
}

func TestDefaultVars(t *testing.T) {
	checker, err := database_rules.NewErrorDBChecker(true, true, true, nil, nil, rules.TypeAlertWarning)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []TestCases{
		{
			Name:           "valid_test_all_check_input_critical",
			Payload:        fmt.Errorf("dial tcp [::1]:5432: connect: connection refused"),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   fmt.Errorf("dial tcp [::1]:5432: connect: connection refused").Error(),
		},
		{
			Name:           "valid_test_all_check_input_critical_2",
			Payload:        fmt.Errorf("pq: connection to server at \"192.168.1.100\" (192.168.1.100), port 5432 failed: Network is unreachable"),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "pq: connection to server at \"192.168.1.100\" (192.168.1.100), port 5432 failed: Network is unreachable",
		},
		{
			Name:           "valid_test_all_check_input_critical_3",
			Payload:        fmt.Errorf("pq: connection reset by peer"),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "pq: connection reset by peer",
		},
		{
			Name:           "valid_test_all_check_input_critical_4",
			Payload:        fmt.Errorf("dial tcp 192.168.1.100:5432: i/o timeout"),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "dial tcp 192.168.1.100:5432: i/o timeout",
		},
		{
			Name:           "valid_test_all_check_input_critical_5",
			Payload:        fmt.Errorf("pq: password authentication failed for user \"postgres\""),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "pq: password authentication failed for user \"postgres\"",
		},
		{
			Name:           "valid_test_all_check_input_critical_6",
			Payload:        fmt.Errorf("pq: out of memory"),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "pq: out of memory",
		},
		{
			Name:           "valid_test_all_check_input_critical_7",
			Payload:        fmt.Errorf("pq: could not write to file \"base/16384/12345\": No space left on device"),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "pq: could not write to file \"base/16384/12345\": No space left on device",
		},
		{
			Name:           "valid_test_all_check_input_critical_8",
			Payload:        fmt.Errorf("pq: sorry, too many clients already"),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "pq: sorry, too many clients already",
		},

		{
			Name:           "valid_test_all_check_input_warning_1",
			Payload:        fmt.Errorf("syntax error"),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "syntax error",
		},
		{
			Name:           "valid_test_all_check_input_warning_2",
			Payload:        fmt.Errorf("pq: duplicate key value violates unique constraint \"users_email_key\""),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: duplicate key value violates unique constraint \"users_email_key\"",
		},
		{
			Name:           "valid_test_all_check_input_warning_3",
			Payload:        fmt.Errorf("pq: DETAIL: Key (id)=(123) already exists."),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: DETAIL: Key (id)=(123) already exists.",
		},
		{
			Name:           "valid_test_all_check_input_warning_4",
			Payload:        fmt.Errorf("pq: insert or update on table \"boards\" violates foreign key constraint \"board_id_fkey\""),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: insert or update on table \"boards\" violates foreign key constraint \"board_id_fkey\"",
		},
		{
			Name:           "valid_test_all_check_input_warning_5",
			Payload:        fmt.Errorf("pq: null value in column \"temp_cpu\" violates not-null constraint"),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: null value in column \"temp_cpu\" violates not-null constraint",
		},
		{
			Name:           "valid_test_all_check_input_warning_6",
			Payload:        fmt.Errorf("pq: DETAIL: Failing row contains (123, null, John, Doe)."),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: DETAIL: Failing row contains (123, null, John, Doe).",
		},
		{
			Name:           "valid_test_all_check_input_warning_7",
			Payload:        fmt.Errorf("pq: invalid input syntax for type integer: \"abc\""),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: invalid input syntax for type integer: \"abc\"",
		},
		{
			Name:           "valid_test_all_check_input_warning_8",
			Payload:        fmt.Errorf("pq: cannot cast type text to integer"),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: cannot cast type text to integer",
		},
		{
			Name:           "valid_test_all_check_input_warning_9",
			Payload:        fmt.Errorf("pq: permission denied for table users"),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: permission denied for table users",
		},
		{
			Name:           "valid_test_all_check_input_warning_10",
			Payload:        fmt.Errorf("pq: relation \"nonexistent_table\" does not exist"),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: relation \"nonexistent_table\" does not exist",
		},

		{
			Name:           "only detail line - should be warning",
			Payload:        fmt.Errorf("pq: DETAIL: Key (id)=(123) already exists."),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: DETAIL: Key (id)=(123) already exists.",
		},
		{
			Name:           "failing row detail",
			Payload:        fmt.Errorf("pq: DETAIL: Failing row contains (123, null, John, Doe)."),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: DETAIL: Failing row contains (123, null, John, Doe).",
		},

		{
			Name:           "random_error_other_errors",
			Payload:        fmt.Errorf("random"),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "random",
		},
		{
			Name:           "nil_error",
			Payload:        nil,
			ExpectedStatus: "",
			ExpectedData:   "",
		},

		{
			Name:           "edge_case_register",
			Payload:        fmt.Errorf("CoNnEcTiOn FAIL"),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "CoNnEcTiOn FAIL",
		},
		{
			Name:           "edge_case_invalid_error",
			Payload:        fmt.Errorf("CoNnEcTiOn not FAIL"),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "CoNnEcTiOn not FAIL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			status, data := checker.Check(tc.Payload)
			if status != tc.ExpectedStatus {
				t.Errorf("STATUS: got: %s, expect: %s\n", status, tc.ExpectedStatus)
			}
			if data != tc.ExpectedData {
				t.Errorf("DATA: got: %s, expect; %s\n", data, tc.ExpectedData)
			}
		})
	}
}

type TestCasesSCenarios struct {
	Name           string
	Payload        any
	Checker        *database_rules.ErrorDBChecker
	ExpectedStatus string
	ExpectedData   string
}

func TestScenarios(t *testing.T) {
	testCases := []TestCasesSCenarios{
		//		 =====ONLY DEFAULT CRITICAL NOT OTHER======
		{
			Name:    "check_only_default_critical_other_false_input_critical",
			Payload: fmt.Errorf("dial tcp [::1]:5432: connect: connection refused"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(true, false, false, nil, nil, "")
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "dial tcp [::1]:5432: connect: connection refused",
		},
		{
			Name:    "check_only_default_critical_other_false_input_warning",
			Payload: fmt.Errorf("pq: permission denied for table users"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(true, false, false, nil, nil, "")
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: "",
			ExpectedData:   "",
		},
		{
			Name:    "check_only_default_critical_other_false_input_nil",
			Payload: nil,
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(true, false, false, nil, nil, "")
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: "",
			ExpectedData:   "",
		},
		{
			Name:    "check_only_default_critical_other_false_input_random",
			Payload: fmt.Errorf("random"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(true, false, false, nil, nil, "")
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: "",
			ExpectedData:   "",
		},

		// =====ONLY DEFAULT WARNING NOT OTHER=====

		{
			Name:    "check_only_default_warning_other_false_input_warning",
			Payload: fmt.Errorf("pq: invalid input syntax for type integer: \"abc\""),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(false, true, false, nil, nil, "")
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: invalid input syntax for type integer: \"abc\"",
		},
		{
			Name:    "check_only_default_warning_other_false_input_critical",
			Payload: fmt.Errorf("pq: password authentication failed for user \"postgres\""),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(false, true, false, nil, nil, "")
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: "",
			ExpectedData:   "",
		},
		{
			Name:    "check_only_default_warning_other_false_input_nil",
			Payload: nil,
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(false, true, false, nil, nil, "")
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: "",
			ExpectedData:   "",
		},
		{
			Name:    "check_only_default_warning_other_false_input_random",
			Payload: fmt.Errorf("random"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(true, false, false, nil, nil, "")
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: "",
			ExpectedData:   "",
		},

		// =====ONLY OTHER WARNING=====

		{
			Name:    "check_other_warning_input_nil",
			Payload: nil,
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(false, false, true, nil, nil, rules.TypeAlertWarning)
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: "",
			ExpectedData:   "",
		},
		{
			Name:    "check_other_warning_input_random",
			Payload: fmt.Errorf("random"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(false, false, true, nil, nil, rules.TypeAlertWarning)
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "random",
		},
		{
			Name:    "check_other_warning_input_critical",
			Payload: fmt.Errorf("pq: connection to server at \"192.168.1.100\" (192.168.1.100), port 5432 failed: Network is unreachable"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(false, false, true, nil, nil, rules.TypeAlertWarning)
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "pq: connection to server at \"192.168.1.100\" (192.168.1.100), port 5432 failed: Network is unreachable",
		},
		{
			Name:    "check_other_warning_input_warning",
			Payload: fmt.Errorf("syntax error"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(false, false, true, nil, nil, rules.TypeAlertWarning)
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "syntax error",
		},

		// =====ONLY OTHER CRITICAL=====

		{
			Name:    "check_other_critical_input_nil",
			Payload: nil,
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(false, false, true, nil, nil, rules.TypeAlertCritical)
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: "",
			ExpectedData:   "",
		},
		{
			Name:    "check_other_critical_input_random",
			Payload: fmt.Errorf("random"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(false, false, true, nil, nil, rules.TypeAlertCritical)
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "random",
		},
		{
			Name:    "check_other_critical_input_critical",
			Payload: fmt.Errorf("pq: connection to server at \"192.168.1.100\" (192.168.1.100), port 5432 failed: Network is unreachable"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(false, false, true, nil, nil, rules.TypeAlertCritical)
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "pq: connection to server at \"192.168.1.100\" (192.168.1.100), port 5432 failed: Network is unreachable",
		},
		{
			Name:    "check_other_critical_input_warning",
			Payload: fmt.Errorf("syntax error"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(false, false, true, nil, nil, rules.TypeAlertCritical)
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "syntax error",
		},

		// ====CUSTOM CFG =====

		{
			Name:    "custom_cfg_critical_1_input_invalid",
			Payload: fmt.Errorf("random"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(true, false, false, []string{"connection"}, nil, "")
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: "",
			ExpectedData:   "",
		},
		{
			Name:    "custom_cfg_critical_1_input_warning",
			Payload: fmt.Errorf("syntax error"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(true, false, false, []string{"connection"}, nil, "")
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: "",
			ExpectedData:   "",
		},
		{
			Name:    "custom_cfg_critical_1_input_critical_invalid",
			Payload: fmt.Errorf("pq: password authentication failed for user \"postgres\""),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(true, false, false, []string{"connection"}, nil, "")
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: "",
			ExpectedData:   "",
		},
		{
			Name:    "custom_cfg_critical_1_input_critical_valid",
			Payload: fmt.Errorf("pq: connection fail"),
			Checker: func() *database_rules.ErrorDBChecker {
				check, err := database_rules.NewErrorDBChecker(true, false, false, []string{"connection"}, nil, "")
				if err != nil {
					t.Fatal(err.Error())
				}
				return check
			}(),
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "pq: connection fail",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			status, data := tc.Checker.Check(tc.Payload)
			if status != tc.ExpectedStatus {
				t.Errorf("STATUS: got: %s, expect: %s\n", status, tc.ExpectedStatus)
			}
			if data != tc.ExpectedData {
				t.Errorf("DATA: got: %s, expect; %s\n", data, tc.ExpectedData)
			}
		})
	}
}

func TestCaseInsensitivity(t *testing.T) {
	checker, _ := database_rules.NewErrorDBChecker(
		true, true, false,
		[]string{"CONNECTION"},
		[]string{"DUPLICATE"},
		"",
	)

	cases := []struct {
		input          error
		expectedStatus string
	}{
		{fmt.Errorf("Connection refused"), rules.TypeAlertCritical},
		{fmt.Errorf("CONNECTION FAILED"), rules.TypeAlertCritical},
		{fmt.Errorf("CoNnEcTiOn error"), rules.TypeAlertCritical},
		{fmt.Errorf("Duplicate key"), rules.TypeAlertWarning},
		{fmt.Errorf("DUPLICATE VALUE"), rules.TypeAlertWarning},
		{fmt.Errorf("DuPlIcAtE entry"), rules.TypeAlertWarning},
	}

	for _, tc := range cases {
		t.Run(tc.input.Error(), func(t *testing.T) {
			status, _ := checker.Check(tc.input)
			if status != tc.expectedStatus {
				t.Errorf("got %s, expected %s", status, tc.expectedStatus)
			}
		})
	}
}
