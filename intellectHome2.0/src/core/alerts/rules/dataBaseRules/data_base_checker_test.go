package database_rules_test

import (
	database_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/dataBaseRules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
// 	database_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/dataBaseRules"
// 	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
// )

// func TestNewDataBaseChecker(t *testing.T) {
// 	testCases := []struct {
// 		Name                       string
// 		ErrorCheck                 bool
// 		StatusStartRecovery        bool
// 		StatusNotRecovered         bool
// 		StatusFinishRecover        bool
// 		StatusOK                   bool
// 		ExpectedErrorStatusChecker error
// 		ExpectedErrorChecker       error
// 	}{
// 		{"normal", true, true, true, true, true, nil, nil},
// 		{"normal_no_error", false, true, true, true, true, nil, nil},
// 		{"normal_no_error_no_start", false, false, true, true, true, nil, nil},
// 		{"normal_no_err_no_start_no_not_recover", false, false, false, true, true, nil, nil},
// 		{"normal_no_err_no_start_no_not_recover_no_finish", false, false, false, false, true, nil, nil},
// 		{"normal_nil_status_checker", true, false, false, false, false, fmt.Errorf("no rules, dataBaseStatusChecker no point"), nil},
// 		{"invalid_nil_status_checker_no_error", false, false, false, false, false, fmt.Errorf("no rules, dataBaseStatusChecker no point"), fmt.Errorf("no rules, DataBaseChecker no point")},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			statusChecker, err := database_rules.NewDataBaseStatusChecker(tc.StatusStartRecovery, tc.StatusNotRecovered, tc.StatusFinishRecover, tc.StatusOK)
// 			if err != nil {
// 				if tc.ExpectedErrorStatusChecker != nil {
// 					if tc.ExpectedErrorStatusChecker.Error() != err.Error() {
// 						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.ExpectedErrorStatusChecker.Error())
// 						return
// 					}
// 				} else {
// 					t.Errorf("unexpected error status checker: %s\n", err.Error())
// 					return
// 				}
// 			}
// 			_, err = database_rules.NewDataBaseChecker(tc.ErrorCheck, statusChecker)
// 			if err != nil {
// 				if tc.ExpectedErrorChecker != nil {
// 					if tc.ExpectedErrorChecker.Error() != err.Error() {
// 						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.ExpectedErrorChecker.Error())
// 					}
// 				} else {
// 					t.Errorf("unexpected error status checker: %s\n", err.Error())
// 				}
// 				return
// 			}

// 		})
// 	}
// }

type TestCaseDataBaseCheckerCheck struct {
	Name          string
	ErrorCheck    *database_rules.ErrorDBChecker
	StatusChecker *database_rules.DataBaseStatusChecker
	Body          *events.Event
	ExpectedAl    bool
	ExpectData    string
	ExpectStatus  string
	ExpectedError error
}

// func TestDataBaseCheckerCheck(t *testing.T) {
// 	testCases := []TestCaseDataBaseCheckerCheck{
// 		{
// 			Name:       "normal_all_check_input_err_syntax",
// 			ErrorCheck: true,
// 			StatusChecker: func() *database_rules.DataBaseStatusChecker {
// 				s, _ := database_rules.NewDataBaseStatusChecker(true, true, true, true)
// 				return s
// 			}(),
// 			Body: &events.Event{
// 				Payload: fmt.Errorf("psql: syntax error"),
// 			},
// 			ExpectedAl:   true,
// 			ExpectData:   "psql: syntax error",
// 			ExpectStatus: rules.TypeAlertWarning,
// 		},
// 		{
// 			Name:       "normal_all_check_input_err_bad_connect",
// 			ErrorCheck: true,
// 			StatusChecker: func() *database_rules.DataBaseStatusChecker {
// 				s, _ := database_rules.NewDataBaseStatusChecker(true, true, true, true)
// 				return s
// 			}(),
// 			Body: &events.Event{
// 				Payload: fmt.Errorf("pq: dial tcp [::1]:5432: connect: connection refused"),
// 			},
// 			ExpectedAl:   true,
// 			ExpectData:   "pq: dial tcp [::1]:5432: connect: connection refused",
// 			ExpectStatus: rules.TypeAlertCritical,
// 		},
// 		{
// 			Name:       "normal_all_check_input_err_bad_connect_pgx",
// 			ErrorCheck: true,
// 			StatusChecker: func() *database_rules.DataBaseStatusChecker {
// 				s, _ := database_rules.NewDataBaseStatusChecker(true, true, true, true)
// 				return s
// 			}(),
// 			Body: &events.Event{
// 				Payload: fmt.Errorf("ping failed: dial tcp 127.0.0.1:5432: connect: connection refused sql: database is closed driver: bad connection"),
// 			},
// 			ExpectedAl:   true,
// 			ExpectData:   "ping failed: dial tcp 127.0.0.1:5432: connect: connection refused sql: database is closed driver: bad connection",
// 			ExpectStatus: rules.TypeAlertCritical,
// 		},
// 		{
// 			Name:       "normal_all_check_input_err_bad_connect_many_clients",
// 			ErrorCheck: true,
// 			StatusChecker: func() *database_rules.DataBaseStatusChecker {
// 				s, _ := database_rules.NewDataBaseStatusChecker(true, true, true, true)
// 				return s
// 			}(),
// 			Body: &events.Event{
// 				Payload: fmt.Errorf("pq: sorry, too many clients already"),
// 			},
// 			ExpectedAl:   true,
// 			ExpectData:   "pq: sorry, too many clients already",
// 			ExpectStatus: rules.TypeAlertCritical,
// 		},
// 		{
// 			Name:       "normal_all_check_input_err_timeout",
// 			ErrorCheck: true,
// 			StatusChecker: func() *database_rules.DataBaseStatusChecker {
// 				s, _ := database_rules.NewDataBaseStatusChecker(true, true, true, true)
// 				return s
// 			}(),
// 			Body: &events.Event{
// 				Payload: fmt.Errorf("NetworkError: failed to connect to `host=db.prod.internal`: dial tcp 10.0.1.5:5432: i/o timeout could not connect to server: Connection timed out"),
// 			},
// 			ExpectedAl:   true,
// 			ExpectData:   "NetworkError: failed to connect to `host=db.prod.internal`: dial tcp 10.0.1.5:5432: i/o timeout could not connect to server: Connection timed out",
// 			ExpectStatus: rules.TypeAlertCritical,
// 		},
// 		{
// 			Name:       "normal_no_err_input_err_bad_connect_many_clients",
// 			ErrorCheck: false,
// 			StatusChecker: func() *database_rules.DataBaseStatusChecker {
// 				s, _ := database_rules.NewDataBaseStatusChecker(true, true, true, true)
// 				return s
// 			}(),
// 			Body: &events.Event{
// 				Payload: fmt.Errorf("pq: sorry, too many clients already"),
// 			},
// 			ExpectedAl: false,
// 		},
// 		{
// 			Name:       "normal_no_err_input_err_timeout",
// 			ErrorCheck: false,
// 			StatusChecker: func() *database_rules.DataBaseStatusChecker {
// 				s, _ := database_rules.NewDataBaseStatusChecker(true, true, true, true)
// 				return s
// 			}(),
// 			Body: &events.Event{
// 				Payload: fmt.Errorf("NetworkError: failed to connect to `host=db.prod.internal`: dial tcp 10.0.1.5:5432: i/o timeout could not connect to server: Connection timed out"),
// 			},
// 			ExpectedAl: false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			checker, err := database_rules.NewDataBaseChecker(tc.ErrorCheck, tc.StatusChecker)
// 			if err != nil {
// 				t.Errorf("unexpected error: %s\n", err.Error())
// 				return
// 			}
// 			al, err := checker.Check(*tc.Body)
// 			if err != nil {
// 				t.Errorf("unexpected error: %s\n", err.Error())
// 				return
// 			}

// 			if tc.ExpectedAl {
// 				if al == nil {
// 					t.Error("EXPECTED AL: got: nil, expected: true\n")
// 					return
// 				}
// 				if al.Type != tc.ExpectStatus {
// 					t.Errorf("EXPECTED STATUS: got: %s, expect: %s\n", al.Type, tc.ExpectStatus)
// 				}
// 				if al.Data != tc.ExpectData {
// 					t.Errorf("EXPECTED DATA: got: %s, expected: %s\n", al.Data, tc.ExpectData)
// 				}

// 			} else {
// 				if al != nil {
// 					t.Errorf("unexpected body: %s\n", al.Data)
// 				}
// 			}
// 		})
// 	}
// }
