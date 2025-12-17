package board_info_rules_test

import (
	"testing"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	board_info_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardInfoRules"
)

type TestCaseNewRSSICheck struct {
	Name     string
	warning  int
	critical int

	expectedCritical int
	expectedWarning  int
}

func TestNewRSSIChecker(t *testing.T) {
	testCases := []TestCaseNewRSSICheck{
		{
			Name:     "default_test",
			warning:  0,
			critical: 0,

			expectedCritical: board_info_rules.DefaultCriticalRSSI,
			expectedWarning:  board_info_rules.DefaultWarningRSSI,
		},
		{
			Name:     "default_test_2",
			warning:  10,
			critical: 100,

			expectedCritical: board_info_rules.DefaultCriticalRSSI,
			expectedWarning:  board_info_rules.DefaultWarningRSSI,
		},
		{
			Name:     "normal_1",
			warning:  -20,
			critical: -10,

			expectedCritical: -10,
			expectedWarning:  -20,
		},
		{
			Name:     "normal_2",
			warning:  -100,
			critical: -100,

			expectedCritical: -100,
			expectedWarning:  -100,
		},
		{
			Name:     "normal_3",
			warning:  -101,
			critical: -100,

			expectedCritical: -100,
			expectedWarning:  board_info_rules.DefaultWarningRSSI,
		},
		{
			Name:     "normal_4",
			warning:  -101,
			critical: -101,

			expectedCritical: board_info_rules.DefaultCriticalRSSI,
			expectedWarning:  board_info_rules.DefaultWarningRSSI,
		},
		{
			Name:     "normal_5",
			warning:  -200,
			critical: -200,

			expectedCritical: board_info_rules.DefaultCriticalRSSI,
			expectedWarning:  board_info_rules.DefaultWarningRSSI,
		},
		{
			Name:     "normal_6",
			warning:  -50,
			critical: -70,

			expectedCritical: -70,
			expectedWarning:  -50,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			check := board_info_rules.NewRSsiChecker(tc.warning, tc.critical)
			if check.CriticalRSSI != tc.expectedCritical {
				t.Errorf("CRITICAL: got: %d, expect: %d\n", check.CriticalRSSI, tc.expectedCritical)
			}
			if check.WarningRSSI != tc.expectedWarning {
				t.Errorf("WARNING: got: %d, expect: %d\n", check.WarningRSSI, tc.expectedWarning)
			}
		})
	}
}

type TestCaseRSSICheck struct {
	Name           string
	inputRSSI      int
	expectedStatus string
	expectedText   string
}

func TestRSSICheck(t *testing.T) {
	check := board_info_rules.NewRSsiChecker(0, 0)

	testCases := []TestCaseRSSICheck{
		{
			Name:           "normal_1",
			inputRSSI:      -10,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},
		{
			Name:           "normal_2",
			inputRSSI:      -20,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},
		{
			Name:           "normal_3",
			inputRSSI:      -25,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},
		{
			Name:           "normal_4",
			inputRSSI:      -30,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},
		{
			Name:           "norma_5",
			inputRSSI:      -40,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},
		{
			Name:           "normal_6",
			inputRSSI:      -64,
			expectedStatus: rules.TypeAlertNormal,
			expectedText:   "",
		},

		{
			Name:           "warning_1",
			inputRSSI:      -65,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning low ethernet signal",
		},
		{
			Name:           "warning_2",
			inputRSSI:      -66,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning low ethernet signal",
		},
		{
			Name:           "warning_3",
			inputRSSI:      -74,
			expectedStatus: rules.TypeAlertWarning,
			expectedText:   "warning low ethernet signal",
		},

		{
			Name:           "critical_1",
			inputRSSI:      -75,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical low ethernet signal",
		},

		{
			Name:           "critical_2",
			inputRSSI:      -99,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical low ethernet signal",
		},
		{
			Name:           "critical_3",
			inputRSSI:      -100,
			expectedStatus: rules.TypeAlertCritical,
			expectedText:   "critical low ethernet signal",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			status, text := check.Check(tc.inputRSSI)
			if status != tc.expectedStatus {
				t.Errorf("STATUS: got: %s, expect: %s\n", status, tc.expectedStatus)
			}
			if text != tc.expectedText {
				t.Errorf("STATUS: got: %s, expect: %s\n", text, tc.expectedText)
			}
		})
	}
}
