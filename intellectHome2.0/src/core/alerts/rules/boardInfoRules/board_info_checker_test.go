package board_info_rules_test

import (
	"fmt"
	"strings"
	"testing"

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	board_info_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardInfoRules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

func TestBoardInfoChecker_NoCheckers(t *testing.T) {
	checker := board_info_rules.NewBoardInfoChecker(nil, nil, nil)

	event := events.Event{
		BoardID: "board-1",
		Payload: &dto.UpdateBoardInfo{
			RssiWifi: intPtr(-50),
			CpuTemp:  float64Ptr(65.0),
			Voltage:  float64Ptr(12.5),
		},
	}

	alert, err := checker.Check(event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if alert != nil {
		t.Error("expected no alert when all checkers are nil")
	}
}

type TestCaseBoardInfoChecker struct {
	Name           string
	RssiChecker    *board_info_rules.RssiChecker
	TempChecker    *board_info_rules.TemperatureCpuCheck
	VoltageChecker *board_info_rules.VoltageChecker
	event          events.Event
	ExpectedAlert  bool
	ExpectedStatus string
	ExpectedData   string
	ExpectedError  error
}

func TestBoardInfoCheck(t *testing.T) {
	testCases := []TestCaseBoardInfoChecker{
		{
			Name:           "normal_test_all_checks_1_default",
			RssiChecker:    board_info_rules.NewRSsiChecker(0, 0),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(0, 0, 0, 0),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-50),
					CpuTemp:  float64Ptr(60),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert: false,
		},
		{
			Name:           "normal_test_all_checks_2",
			RssiChecker:    board_info_rules.NewRSsiChecker(-51, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-50),
					CpuTemp:  float64Ptr(60),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert: false,
		},
		{
			Name:           "normal_test_all_checks_3_warning_rssi",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-50),
					CpuTemp:  float64Ptr(60),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "warning low ethernet signal",
		},
		{
			Name:           "normal_test_all_checks_4_critical_rssi",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-70),
					CpuTemp:  float64Ptr(60),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "critical low ethernet signal",
		},
		{
			Name:           "normal_test_all_checks_5_warning_low_temp",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-40),
					CpuTemp:  float64Ptr(-20),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "warning low temp cpu",
		},
		{
			Name:           "normal_test_all_checks_6_critical_low_temp",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-40),
					CpuTemp:  float64Ptr(-35),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "critical low temp cpu",
		},
		{
			Name:           "normal_test_all_checks_7_warning_high_temp",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-40),
					CpuTemp:  float64Ptr(100),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "warning high temp cpu",
		},
		{
			Name:           "normal_test_all_checks_8_critical_high_temp",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-40),
					CpuTemp:  float64Ptr(120),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "critical high temp cpu",
		},
		{
			Name:           "normal_test_all_checks_9_warning_low_voltage",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-40),
					CpuTemp:  float64Ptr(80),
					Voltage:  float64Ptr(2.9),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "warning low voltage",
		},
		{
			Name:           "normal_test_all_checks_10_critical_low_voltage",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-40),
					CpuTemp:  float64Ptr(80),
					Voltage:  float64Ptr(2.7),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "critical low voltage",
		},
		{
			Name:           "normal_test_all_checks_11_warning_high_voltage",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-40),
					CpuTemp:  float64Ptr(80),
					Voltage:  float64Ptr(3.5),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertWarning,
			ExpectedData:   "warning high voltage",
		},
		{
			Name:           "normal_test_all_checks_12_critical_high_voltage",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-40),
					CpuTemp:  float64Ptr(80),
					Voltage:  float64Ptr(3.8),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "critical high voltage",
		},
		{
			Name:           "normal_test_all_checks_low_voltage_critical_low_rssi",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-90),
					CpuTemp:  float64Ptr(80),
					Voltage:  float64Ptr(3.5),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "critical low ethernet signal\nwarning high voltage",
		},
		{
			Name:           "normal_test_all_checks_critical_rssi_warning_temp",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-90),
					CpuTemp:  float64Ptr(100),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "critical low ethernet signal\nwarning high temp cpu",
		},
		{
			Name:           "normal_test_all_checks_critical_rssi_critical_temp",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-90),
					CpuTemp:  float64Ptr(200),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "critical low ethernet signal\ncritical high temp cpu",
		},
		{
			Name:           "normal_test_all_checks_warning_rssi_critical_temp",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-50),
					CpuTemp:  float64Ptr(200),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert:  true,
			ExpectedStatus: rules.TypeAlertCritical,
			ExpectedData:   "warning low ethernet signal\ncritical high temp cpu",
		},
		{
			Name:           "normal_test_no_check_rssi",
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(100),
					CpuTemp:  float64Ptr(80),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert: false,
		},
		{
			Name:           "normal_test_no_check_temp",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			VoltageChecker: board_info_rules.NewVoltageChecker(3.5, 3.7, 2.9, 2.7),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-40),
					CpuTemp:  float64Ptr(200),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert: false,
		},
		{
			Name:        "normal_test_no_check_voltage",
			RssiChecker: board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker: board_info_rules.NewTemperatureCpuCheck(),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-40),
					CpuTemp:  float64Ptr(80),
					Voltage:  float64Ptr(3.8),
				},
			},
			ExpectedAlert: false,
		},
		{
			Name:           "invalid_test_1_no_dto",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(0, 0, 0, 0),
			event: events.Event{
				BoardID: "esp32_2",
			},
			ExpectedAlert: false,
			ExpectedError: fmt.Errorf("invalid type in event"),
		},
		{
			Name:           "invalid_test_2_nil_voltage",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(0, 0, 0, 0),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-40),
					CpuTemp:  float64Ptr(80),
					Voltage:  nil,
				},
			},
			ExpectedAlert: false,
			ExpectedError: fmt.Errorf("invalid dto"),
		},
		{
			Name:           "invalid_test_3_nil_temp",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(0, 0, 0, 0),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: intPtr(-40),
					CpuTemp:  nil,
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert: false,
			ExpectedError: fmt.Errorf("invalid dto"),
		},
		{
			Name:           "invalid_test_4_nil_rssi",
			RssiChecker:    board_info_rules.NewRSsiChecker(-50, -70),
			TempChecker:    board_info_rules.NewTemperatureCpuCheck(),
			VoltageChecker: board_info_rules.NewVoltageChecker(0, 0, 0, 0),
			event: events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: nil,
					CpuTemp:  float64Ptr(40),
					Voltage:  float64Ptr(3.3),
				},
			},
			ExpectedAlert: false,
			ExpectedError: fmt.Errorf("invalid dto"),
		},
	}

	//rssi -> temp -> voltage
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			check := board_info_rules.NewBoardInfoChecker(tc.RssiChecker, tc.TempChecker, tc.VoltageChecker)
			al, err := check.Check(tc.event)
			if err != nil {
				if tc.ExpectedError != nil {
					if err.Error() != tc.ExpectedError.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.ExpectedError.Error())
					}
				} else {
					t.Errorf("unknown err: %v\n", err.Error())
				}
			}
			if tc.ExpectedAlert {
				if al.Type != tc.ExpectedStatus {
					t.Errorf("STATUS: got: %s, expect: %s\n", al.Type, tc.ExpectedStatus)
				}
				if al.Data != tc.ExpectedData {
					t.Errorf("DATA: got: %s, expect: %s\n", al.Data, tc.ExpectedData)
				}
			} else {
				if al != nil {
					t.Errorf("unexpected alert!\n%s\n", al.Data)
				}
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
