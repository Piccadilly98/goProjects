package board_info_rules_test

import (
	"testing"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	board_info_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardInfoRules"
)

type TestCaseNewVoltageChecker struct {
	Name         string
	warningHigh  float64
	criticalHigh float64
	warningLow   float64
	criticalLow  float64

	expectedWarningHigh  float64
	expectedCriticalHigh float64
	expectedWarningLow   float64
	expectedCriticalLow  float64
}

func TestNewVoltageChecker(t *testing.T) {
	testCases := []TestCaseNewVoltageChecker{
		{
			Name:                 "normal_test",
			warningHigh:          3.5,
			criticalHigh:         3.7,
			warningLow:           3.1,
			criticalLow:          2.9,
			expectedWarningHigh:  3.5,
			expectedCriticalHigh: 3.7,
			expectedWarningLow:   3.1,
			expectedCriticalLow:  2.9,
		},
		{
			Name:                 "normal_test_2",
			warningHigh:          3,
			criticalHigh:         3,
			warningLow:           3,
			criticalLow:          2,
			expectedWarningHigh:  3,
			expectedCriticalHigh: 3,
			expectedWarningLow:   3,
			expectedCriticalLow:  2,
		},
		{
			Name:                 "normal_test_3",
			warningHigh:          10,
			criticalHigh:         12,
			warningLow:           14,
			criticalLow:          16,
			expectedWarningHigh:  10,
			expectedCriticalHigh: 12,
			expectedWarningLow:   14,
			expectedCriticalLow:  16,
		},
		{
			Name:                 "default_test",
			warningHigh:          0,
			criticalHigh:         0,
			warningLow:           0,
			criticalLow:          0,
			expectedWarningHigh:  3.5,
			expectedCriticalHigh: 3.7,
			expectedWarningLow:   3.1,
			expectedCriticalLow:  2.9,
		},
		{
			Name:                 "default_test_2",
			warningHigh:          -1,
			criticalHigh:         -1,
			warningLow:           -1,
			criticalLow:          -1,
			expectedWarningHigh:  3.5,
			expectedCriticalHigh: 3.7,
			expectedWarningLow:   3.1,
			expectedCriticalLow:  2.9,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			check := board_info_rules.NewVoltageChecker(tc.warningHigh, tc.criticalHigh, tc.warningLow, tc.criticalLow)
			if check.CriticalHigh != tc.expectedCriticalHigh {
				t.Errorf("CRITICAL HIGH: got: %f, expect: %f\n", check.CriticalHigh, tc.expectedCriticalHigh)
			}
			if check.CriticalLow != tc.expectedCriticalLow {
				t.Errorf("CRITICAL LOW: got: %f, expect: %f\n", check.CriticalLow, tc.expectedCriticalLow)
			}
			if check.WarningHigh != tc.expectedWarningHigh {
				t.Errorf("WARNING HIGH: got: %f, expect: %f\n", check.WarningHigh, tc.expectedWarningHigh)
			}
			if check.WarningLow != tc.expectedWarningLow {
				t.Errorf("WARNING LOW: got: %f, expect: %f\n", check.WarningLow, tc.expectedWarningLow)
			}
		})
	}
}

type TestCaseVoltageChecker struct {
	Name           string
	inputVoltage   float64
	expectedStatus string
	expectedText   string
}

func TestVoltageChecker(t *testing.T) {
	check := board_info_rules.NewVoltageChecker(0, 0, 0, 0)
	testCases := []TestCaseVoltageChecker{
		{
			Name:           "test_normal_1",
			inputVoltage:   3.3,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},
		{
			Name:           "test_normal_2",
			inputVoltage:   3.2,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},
		{
			Name:           "test_normal_3",
			inputVoltage:   3.21,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},
		{
			Name:           "test_normal_4",
			inputVoltage:   3.4,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},
		{
			Name:           "test_normal_5",
			inputVoltage:   3.19,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},
		{
			Name:           "test_normal_6",
			inputVoltage:   3.111,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},
		{
			Name:           "test_normal_7",
			inputVoltage:   3.1000001,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},
		{
			Name:           "test_warning_low_1",
			inputVoltage:   3.1,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning low voltage",
		},
		{
			Name:           "test_warning_low_2",
			inputVoltage:   3.09,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning low voltage",
		},
		{
			Name:           "test_warning_low_3",
			inputVoltage:   3,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning low voltage",
		},
		{
			Name:           "test_warning_low_4",
			inputVoltage:   3.0000001,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning low voltage",
		},
		{
			Name:           "test_warning_low_5",
			inputVoltage:   2.95,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning low voltage",
		},
		{
			Name:           "test_warning_low_6",
			inputVoltage:   2.91,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning low voltage",
		},
		{
			Name:           "test_warning_low_7",
			inputVoltage:   2.9000000006,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning low voltage",
		},
		{
			Name:           "test_critical_low_1",
			inputVoltage:   2.9,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical low voltage",
		},
		{
			Name:           "test_critical_low_2",
			inputVoltage:   2.8,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical low voltage",
		},
		{
			Name:           "test_critical_low_3",
			inputVoltage:   2.899999999999999999999999,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical low voltage",
		},
		{
			Name:           "test_critical_low_4",
			inputVoltage:   -10,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical low voltage",
		},
		{
			Name:           "test_critical_low_5",
			inputVoltage:   2,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical low voltage",
		},
		{
			Name:           "test_critical_low_6",
			inputVoltage:   2.5,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical low voltage",
		},

		{
			Name:           "test_warning_high_1",
			inputVoltage:   3.5,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning high voltage",
		},
		{
			Name:           "test_warning_high_2",
			inputVoltage:   3.6,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning high voltage",
		},
		{
			Name:           "test_warning_high_3",
			inputVoltage:   3.599999999999999999999999999,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning high voltage",
		},
		{
			Name:           "test_warning_high_4",
			inputVoltage:   3.6,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning high voltage",
		},
		{
			Name:           "test_warning_high_5",
			inputVoltage:   3.68,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning high voltage",
		},
		{
			Name:           "test_warning_high_6",
			inputVoltage:   3.699999999,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning high voltage",
		},
		{
			Name:           "test_critical_high_1",
			inputVoltage:   3.7,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical high voltage",
		},
		{
			Name:           "test_critical_high_2",
			inputVoltage:   3.8,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical high voltage",
		},
		{
			Name:           "test_critical_high_3",
			inputVoltage:   3.9,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical high voltage",
		},
		{
			Name:           "test_critical_high_4",
			inputVoltage:   10,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical high voltage",
		},
		{
			Name:           "test_critical_high_5",
			inputVoltage:   100000000000000000000000,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical high voltage",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			status, text := check.Check(tc.inputVoltage)
			if status != tc.expectedStatus {
				t.Errorf("STATUS: got: %s, expect: %s\n", status, tc.expectedStatus)
			}
			if text != tc.expectedText {
				t.Errorf("TEXT: got: %s, expect: %s\n", text, tc.expectedText)
			}
		})
	}
}
