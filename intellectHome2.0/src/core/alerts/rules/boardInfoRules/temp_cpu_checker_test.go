package board_info_rules_test

import (
	"testing"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	board_info_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardInfoRules"
)

type TestCaseTemperatureCpuCheck struct {
	Name           string
	InputTemp      float64
	ExpectedStatus string
	ExpectedText   string
}

func TestTemperatureCheck(t *testing.T) {
	check := board_info_rules.NewTemperatureCpuCheck()

	testCases := []TestCaseTemperatureCpuCheck{
		{
			Name:           "normal_test_1",
			InputTemp:      80,
			ExpectedStatus: rules.TypeAlertNormal,
			ExpectedText:   "",
		},
		{
			Name:           "normal_test_2",
			InputTemp:      89.9,
			ExpectedStatus: rules.TypeAlertNormal,
			ExpectedText:   "",
		},
		{
			Name:           "normal_test_3",
			InputTemp:      90.99999999999,
			ExpectedStatus: rules.TypeAlertNormal,
			ExpectedText:   "",
		},
		{
			Name:           "normal_test_4",
			InputTemp:      -10,
			ExpectedStatus: rules.TypeAlertNormal,
			ExpectedText:   "",
		},
		{
			Name:           "normal_test_5",
			InputTemp:      -19.9,
			ExpectedStatus: rules.TypeAlertNormal,
			ExpectedText:   "",
		},
		{
			Name:           "normal_test_6",
			InputTemp:      99.9,
			ExpectedStatus: rules.TypeAlertNormal,
			ExpectedText:   "",
		},
		{
			Name:           "normal_test_7",
			InputTemp:      99.9999,
			ExpectedStatus: rules.TypeAlertNormal,
			ExpectedText:   "",
		},
		{
			Name:           "warning_high_test_1",
			InputTemp:      100,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedText:   "warning high temp cpu",
		},
		{
			Name:           "warning_high_test_2",
			InputTemp:      100.1,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedText:   "warning high temp cpu",
		},
		{
			Name:           "warning_high_test_3",
			InputTemp:      119.9,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedText:   "warning high temp cpu",
		},
		{
			Name:           "warning_high_test_4",
			InputTemp:      115.98,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedText:   "warning high temp cpu",
		},
		{
			Name:           "critical_high_test_1",
			InputTemp:      120,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedText:   "critical high temp cpu",
		},
		{
			Name:           "critical_high_test_2",
			InputTemp:      120.0,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedText:   "critical high temp cpu",
		},
		{
			Name:           "critical_high_test_3",
			InputTemp:      140,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedText:   "critical high temp cpu",
		},
		{
			Name:           "critical_high_test_4",
			InputTemp:      150,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedText:   "critical high temp cpu",
		},
		{
			Name:           "critical_high_test_5",
			InputTemp:      120.000001,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedText:   "critical high temp cpu",
		},

		{
			Name:           "warning_low_test_1",
			InputTemp:      -20,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedText:   "warning low temp cpu",
		},
		{
			Name:           "warning_low_test_2",
			InputTemp:      -25,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedText:   "warning low temp cpu",
		},
		{
			Name:           "warning_low_test_3",
			InputTemp:      -25.5,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedText:   "warning low temp cpu",
		},
		{
			Name:           "warning_low_test_4",
			InputTemp:      -30,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedText:   "warning low temp cpu",
		},
		{
			Name:           "warning_low_test_5",
			InputTemp:      -34.5,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedText:   "warning low temp cpu",
		},
		{
			Name:           "warning_low_test_6",
			InputTemp:      -34.99999,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedText:   "warning low temp cpu",
		},

		{
			Name:           "critical_low_test_1",
			InputTemp:      -35,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedText:   "critical low temp cpu",
		},
		{
			Name:           "critical_low_test_2",
			InputTemp:      -35.00001,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedText:   "critical low temp cpu",
		},
		{
			Name:           "critical_low_test_3",
			InputTemp:      -50.5,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedText:   "critical low temp cpu",
		},
		{
			Name:           "critical_low_test_4",
			InputTemp:      -100,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedText:   "critical low temp cpu",
		},
	}
	//critical high - 120
	//warm high 100 - 119.9
	//normal -19.9 - 99.9
	//warn low - -20 - -34.9
	//critical low - -35

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			status, text := check.Check(tc.InputTemp)
			if status != tc.ExpectedStatus {
				t.Errorf("STATUS: got: %s, expect: %s\n", status, tc.ExpectedStatus)
			}
			if text != tc.ExpectedText {
				t.Errorf("TEXT: got: %s, expect: %s\n", text, tc.ExpectedText)
			}
		})
	}
}
