package database_rules_test

import (
	"fmt"
	"testing"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	database_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/dataBaseRules"
)

func TestNewDataBaseStatusChecker(t *testing.T) {
	testCases := []struct {
		Name                string
		StatusStartRecovery bool
		StatusNotRecovered  bool
		StatusFinishRecover bool
		StatusOK            bool
		ExpectedError       error
	}{
		{"normal_all", true, true, true, true, nil},
		{"normal_no_start_recovery", false, true, true, true, nil},
		{"normal_no_no_recovered", true, false, true, true, nil},
		{"normal_no_finish_recovery", true, true, false, true, nil},
		{"normal_no_ok", true, true, true, false, nil},
		{"normal_only_start_recovery", true, false, false, false, nil},
		{"normal_only_not_recovered", false, true, false, false, nil},
		{"normal_only_finish_recovery", false, false, true, false, nil},
		{"normal_only_ok", false, false, false, true, nil},
		{"no_rules", false, false, false, false, fmt.Errorf("no rules, dataBaseStatusChecker no point")},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			check, err := database_rules.NewDataBaseStatusChecker(tc.StatusStartRecovery, tc.StatusNotRecovered, tc.StatusFinishRecover, tc.StatusOK)
			if err != nil {
				if tc.ExpectedError != nil {
					if tc.ExpectedError.Error() != err.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.ExpectedError.Error())
					}
				} else {
					t.Errorf("unexpected err: %s\n", err.Error())
				}
				return
			}
			if check.StatusFinishRecover != tc.StatusFinishRecover {
				t.Errorf("STATUS FINISH RECOVER: got: %v, expect: %v\n", check.StatusFinishRecover, tc.StatusFinishRecover)
			}
			if check.StatusNotRecovered != tc.StatusNotRecovered {
				t.Errorf("STATUS FINISH RECOVER: got: %v, expect: %v\n", check.StatusNotRecovered, tc.StatusNotRecovered)
			}
			if check.StatusOK != tc.StatusOK {
				t.Errorf("STATUS FINISH RECOVER: got: %v, expect: %v\n", check.StatusOK, tc.StatusOK)
			}
			if check.StatusStartRecovery != tc.StatusStartRecovery {
				t.Errorf("STATUS FINISH RECOVER: got: %v, expect: %v\n", check.StatusStartRecovery, tc.StatusStartRecovery)
			}
		})
	}
}

type TestCaseDataBaseStatusChecker struct {
	Name                string
	StatusStartRecovery bool
	StatusNotRecovered  bool
	StatusFinishRecover bool
	StatusOK            bool
	ExpectedError       error
	InputData           any
	ExpectedData        bool
	ExpectData          string
	ExpectStatus        string
}

// переписать что бы проверялось тело на "", ""
func TestDataBaseStatusChecker(t *testing.T) {
	testCases := []TestCaseDataBaseStatusChecker{
		{
			Name:                "normal_start_recover",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           fmt.Sprintf("message received by: %s, message: %v,  DataBase fail, start Recover\n", "sd", "v"),
			ExpectedData:        true,
			ExpectData:          fmt.Sprintf("message received by: %s, message: %v,  DataBase fail, start Recover\n", "sd", "v"),
			ExpectStatus:        rules.TypeAlertWarning,
		},
		{
			Name:                "normal_finish_recover",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase recovered successfully\n",
			ExpectedData:        true,
			ExpectData:          "DataBase recovered successfully\n",
			ExpectStatus:        rules.TypeAlertNormal,
		},
		{
			Name:                "normal_not_recovered",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase not recover, server off\n",
			ExpectedData:        true,
			ExpectData:          "DataBase not recover, server off\n",
			ExpectStatus:        rules.TypeAlertCritical,
		},
		{
			Name:                "normal_ok",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase ping status: ok, time: 21-09-22.02212",
			ExpectedData:        true,
			ExpectData:          "DataBase ping status: ok, time: 21-09-22.02212",
			ExpectStatus:        rules.TypeAlertNormal,
		},
		{
			Name:                "normal_start_recover_status_false",
			StatusStartRecovery: false,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           fmt.Sprintf("message received by: %s, message: %v,  DataBase fail, start Recover\n", "sd", "v"),
			ExpectedData:        false,
		},
		{
			Name:                "normal_no_start_input_not_recovered",
			StatusStartRecovery: false,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase not recover, server off\n",
			ExpectedData:        true,
			ExpectStatus:        rules.TypeAlertCritical,
			ExpectData:          "DataBase not recover, server off\n",
		},
		{
			Name:                "normal_no_start_input_recover_finish",
			StatusStartRecovery: false,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase recovered successfully\n",
			ExpectedData:        true,
			ExpectStatus:        rules.TypeAlertNormal,
			ExpectData:          "DataBase recovered successfully\n",
		},
		{
			Name:                "normal_no_start_input_ok",
			StatusStartRecovery: false,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase ping status: ok",
			ExpectedData:        true,
			ExpectStatus:        rules.TypeAlertNormal,
			ExpectData:          "DataBase ping status: ok",
		},
		{
			Name:                "normal_no_not_recovered_input_start",
			StatusStartRecovery: true,
			StatusNotRecovered:  false,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           fmt.Sprintf("message received by: %s, message: %v,  DataBase fail, start Recover\n", "sd", "v"),
			ExpectedData:        true,
			ExpectStatus:        rules.TypeAlertWarning,
			ExpectData:          fmt.Sprintf("message received by: %s, message: %v,  DataBase fail, start Recover\n", "sd", "v"),
		},
		{
			Name:                "normal_no_not_recovered_input_no_recovered",
			StatusStartRecovery: true,
			StatusNotRecovered:  false,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase not recover, server off\n",
			ExpectedData:        false,
		},
		{
			Name:                "normal_no_not_recovered_input_finish",
			StatusStartRecovery: true,
			StatusNotRecovered:  false,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase recovered successfully\n",
			ExpectedData:        true,
			ExpectStatus:        rules.TypeAlertNormal,
			ExpectData:          "DataBase recovered successfully\n",
		},
		{
			Name:                "normal_no_not_recovered_input_ok",
			StatusStartRecovery: true,
			StatusNotRecovered:  false,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase ping status: ok",
			ExpectedData:        true,
			ExpectStatus:        rules.TypeAlertNormal,
			ExpectData:          "DataBase ping status: ok",
		},
		{
			Name:                "normal_no_finish_input_start",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: false,
			StatusOK:            true,
			InputData:           fmt.Sprintf("message received by: %s, message: %v,  DataBase fail, start Recover\n", "sd", "v"),
			ExpectedData:        true,
			ExpectStatus:        rules.TypeAlertWarning,
			ExpectData:          fmt.Sprintf("message received by: %s, message: %v,  DataBase fail, start Recover\n", "sd", "v"),
		},
		{
			Name:                "normal_no_finish_input_no_recovered",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: false,
			StatusOK:            true,
			InputData:           "DataBase not recover, server off\n",
			ExpectedData:        true,
			ExpectData:          "DataBase not recover, server off\n",
			ExpectStatus:        rules.TypeAlertCritical,
		},
		{
			Name:                "normal_no_finish_input_finish",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: false,
			StatusOK:            true,
			InputData:           "DataBase recovered successfully\n",
			ExpectedData:        false,
		},
		{
			Name:                "normal_no_finish_input_ok",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: false,
			StatusOK:            true,
			InputData:           "DataBase ping status: ok",
			ExpectedData:        true,
			ExpectStatus:        rules.TypeAlertNormal,
			ExpectData:          "DataBase ping status: ok",
		},
		{
			Name:                "normal_no_ok_input_start",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            false,
			InputData:           fmt.Sprintf("message received by: %s, message: %v,  DataBase fail, start Recover\n", "sd", "v"),
			ExpectedData:        true,
			ExpectStatus:        rules.TypeAlertWarning,
			ExpectData:          fmt.Sprintf("message received by: %s, message: %v,  DataBase fail, start Recover\n", "sd", "v"),
		},
		{
			Name:                "normal_no_ok_input_no_recovered",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            false,
			InputData:           "DataBase not recover, server off\n",
			ExpectedData:        true,
			ExpectData:          "DataBase not recover, server off\n",
			ExpectStatus:        rules.TypeAlertCritical,
		},
		{
			Name:                "normal_no_ok_input_finish",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            false,
			InputData:           "DataBase recovered successfully\n",
			ExpectedData:        true,
			ExpectStatus:        rules.TypeAlertNormal,
			ExpectData:          "DataBase recovered successfully\n",
		},
		{
			Name:                "normal_no_ok_input_ok",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            false,
			InputData:           "DataBase ping status: ok",
			ExpectedData:        false,
		},

		{
			Name:                "edge_no_recovered_in_input",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase not_recovered successfully\n",
			ExpectedData:        false,
		},
		{
			Name:                "edge_no_recovered_2_in_input",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase not recovered successfully\n",
			ExpectedData:        false,
		},
		{
			Name:                "edge_no_recovered_3_in_input",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase not-recovered successfully\n",
			ExpectedData:        false,
		},
		{
			Name:                "edge_no_recovered_4_in_input",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase no recovered successfully\n",
			ExpectedData:        false,
		},
		{
			Name:                "edge_no_ok_in_input",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase ping status: not ok",
			ExpectedData:        false,
		},
		{
			Name:                "edge_no_ok_2_in_input",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase ping status: no ok",
			ExpectedData:        false,
		},
		{
			Name:                "edge_no_ok_3_in_input",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase ping status: no-ok",
			ExpectedData:        false,
		},
		{
			Name:                "edge_not_ok_4_in_input",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase ping status: not-ok",
			ExpectedData:        false,
		},
		{
			Name:                "edge_no_ok_5_in_input",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            true,
			InputData:           "DataBase ping status: no_ok",
			ExpectedData:        false,
		},
		{
			Name:                "invalid_type_payload_1",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            false,
			InputData:           []string{"DataBase ping status: ok"},
			ExpectedData:        false,
		},
		{
			Name:                "invalid_type_payload_2",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            false,
			InputData:           struct{ status string }{status: "DataBase ping status: ok"},
			ExpectedData:        false,
		},
		{
			Name:                "invalid_type_payload_3",
			StatusStartRecovery: true,
			StatusNotRecovered:  true,
			StatusFinishRecover: true,
			StatusOK:            false,
			InputData:           map[string]string{"status": "DataBase ping status: ok"},
			ExpectedData:        false,
		},
		{
			Name:                "no_rules",
			StatusStartRecovery: false,
			StatusNotRecovered:  false,
			StatusFinishRecover: false,
			StatusOK:            false,
			ExpectedError:       fmt.Errorf("no rules, dataBaseStatusChecker no point"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			check, err := database_rules.NewDataBaseStatusChecker(tc.StatusStartRecovery, tc.StatusNotRecovered, tc.StatusFinishRecover, tc.StatusOK)
			if err != nil {
				if tc.ExpectedError != nil {
					if err.Error() != tc.ExpectedError.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.ExpectedError.Error())
					}
				} else {
					t.Errorf("unexpected error: %s\n", err.Error())
				}
				return
			}
			str, status := check.Check(tc.InputData)

			if tc.ExpectedData {
				if str != tc.ExpectData {
					t.Errorf("DATA: got: %s, expect: %s\n", str, tc.ExpectData)
				}
				if status != tc.ExpectStatus {
					t.Errorf("STATUS: got: %s, expect: %s\n", status, tc.ExpectStatus)
				}
			} else {
				if str != "" {
					t.Errorf("unexpected str: %s\n", str)
				}
				if status != "" {
					t.Errorf("unexpected status: %s\n", status)
				}
			}
		})
	}
}
