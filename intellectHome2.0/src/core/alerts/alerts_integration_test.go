package alerts_test

import (
	"strings"
	"testing"
	"time"

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/notifiers"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	board_info_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardInfoRules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

func TestAlertManager_CheckTemperature(t *testing.T) {
	testCases := []struct {
		Name        string
		Temperature float64
		ExpectedMsg bool
		ExpectMsg   []string
	}{
		{"normal_1", 40, false, nil},
		{"normal_2", -10, false, nil},
		{"normal_3", 99.999, false, nil},
		{"normal_4", -19.999, false, nil},
		{"warning_high_temp", 100.0, true, []string{"WARNING", "esp32_2", "warning high temp cpu"}},
		{"critical_high_temp", 120.01, true, []string{"CRITICAL", "esp32_2", "critical high temp cpu"}},
		{"warning_low_temp", -20.5, true, []string{"WARNING", "esp32_2", "warning low temp cpu"}},
		{"critical_low_temp", -40.3, true, []string{"CRITICAL", "esp32_2", "critical low temp cpu"}},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			eb := events.NewEventBus(10, 5*time.Second)

			checker := board_info_rules.NewBoardInfoChecker(nil,
				board_info_rules.NewTemperatureCpuCheck(),
				nil)

			msgMock := notifiers.NewGetMessageMock(10)

			am, err := alerts.NewAlertsManager(eb, []rules.Rule{checker}, 0, 0, msgMock)
			if err != nil {
				t.Fatal(err)
			}

			am.Start()
			sub := eb.Subscribe(alerts.TopicforBoardInfoChecker, "test_check_temperature")
			err = eb.Publish(sub.Topic, events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					CpuTemp: getPtrFloat(tc.Temperature),
				},
			}, sub.ID)
			if err != nil {
				t.Fatal(err)
			}
			select {
			case message := <-msgMock.MsgCh:
				if tc.ExpectedMsg {
					if !checkMessage(message, tc.ExpectMsg) {
						t.Errorf("ERROR MESSAGE: got: %s, expect: %v\n", message, tc.ExpectMsg)
					}
				} else {
					t.Errorf("unexpected msg: %s\n", message)
				}
			case <-time.After(200 * time.Millisecond):
				if tc.ExpectedMsg {
					t.Errorf("got: block chan, expect: %v\n", tc.ExpectMsg)
				}
			}
		})
	}

}

func TestAlertManager_CheckRSSI(t *testing.T) {
	testCases := []struct {
		Name        string
		RSSI        int
		ExpectedMsg bool
		ExpectMsg   []string
	}{
		{"normal_1", -40, false, nil},
		{"normal_2", -20, false, nil},
		{"normal_3", -15, false, nil},
		{"normal_4", -59, false, nil},
		{"warning", -65, true, []string{"WARNING", "warning low ethernet signal"}},
		{"critical", -75, true, []string{"CRITICAL", "critical low ethernet signal"}},
		{"edge_1_rssi==0", 0, false, nil},
		{"edge_2_rssi>0", 10, false, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			eb := events.NewEventBus(10, 5*time.Second)

			checker := board_info_rules.NewBoardInfoChecker(board_info_rules.NewRSsiChecker(0, 0),
				nil,
				nil)

			msgMock := notifiers.NewGetMessageMock(10)

			am, err := alerts.NewAlertsManager(eb, []rules.Rule{checker}, 0, 0, msgMock)
			if err != nil {
				t.Fatal(err)
			}
			am.Start()
			sub := eb.Subscribe(alerts.TopicforBoardInfoChecker, "test_check_rssi")
			err = eb.Publish(sub.Topic, events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					RssiWifi: getPtrInt(tc.RSSI),
				},
			}, sub.ID)
			if err != nil {
				t.Fatal(err)
			}

			select {
			case message := <-msgMock.MsgCh:
				if tc.ExpectedMsg {
					if !checkMessage(message, tc.ExpectMsg) {
						t.Errorf("ERROR MESSAGE: got: %s, expect: %v\n", message, tc.ExpectMsg)
					}
				} else {
					t.Errorf("unexpected msg: %s\n", message)
				}
			case <-time.After(200 * time.Millisecond):
				if tc.ExpectedMsg {
					t.Errorf("got: block chan, expect: %v\n", tc.ExpectMsg)
				}
			}
		})
	}
}

func TestAlertManager_CheckVoltage(t *testing.T) {
	testCases := []struct {
		Name        string
		voltage     float64
		ExpectedMsg bool
		ExpectMsg   []string
	}{
		{"normal_1", 3.3, false, nil},
		{"normal_2", 3.2, false, nil},
		{"normal_3", 3.15, false, nil},
		{"normal_4", 3.4, false, nil},
		{"warning_high_voltage", 3.5, true, []string{"WARNING", "esp32_2", "warning high voltage"}},
		{"critical_high_voltage", 3.8, true, []string{"CRITICAL", "esp32_2", "critical high voltage"}},
		{"warning_low_voltage", 3.1, true, []string{"WARNING", "esp32_2", "warning low voltage"}},
		{"critical_low_voltage", 2.9, true, []string{"CRITICAL", "esp32_2", "critical low voltage"}},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			eb := events.NewEventBus(10, 5*time.Second)

			checker := board_info_rules.NewBoardInfoChecker(nil,
				nil,
				board_info_rules.NewVoltageChecker(0, 0, 0, 0))

			msgMock := notifiers.NewGetMessageMock(10)

			am, err := alerts.NewAlertsManager(eb, []rules.Rule{checker}, 0, 0, msgMock)
			if err != nil {
				t.Fatal(err)
			}

			am.Start()
			sub := eb.Subscribe(alerts.TopicforBoardInfoChecker, "test_check_voltage")
			err = eb.Publish(sub.Topic, events.Event{
				BoardID: "esp32_2",
				Payload: &dto.UpdateBoardInfo{
					Voltage: getPtrFloat(tc.voltage),
				},
			}, sub.ID)
			if err != nil {
				t.Fatal(err)
			}
			select {
			case message := <-msgMock.MsgCh:
				if tc.ExpectedMsg {
					if !checkMessage(message, tc.ExpectMsg) {
						t.Errorf("ERROR MESSAGE: got: %s, expect: %v\n", message, tc.ExpectMsg)
					}
				} else {
					t.Errorf("unexpected msg: %s\n", message)
				}
			case <-time.After(200 * time.Millisecond):
				if tc.ExpectedMsg {
					t.Errorf("got: block chan, expect: %v\n", tc.ExpectMsg)
				}
			}
		})
	}
}

func TestAlertManager_ComplexCsenarios(t *testing.T) {
	testCases := []struct {
		Name        string
		Voltage     float64
		RSSI        int
		Temperature float64
		ExpectedMsg bool
		ExpectMsg   []string
	}{
		{
			Name:        "normal_1",
			Voltage:     3.35,
			RSSI:        -40,
			Temperature: 0,
			ExpectedMsg: false,
			ExpectMsg:   nil,
		},
		{
			Name:        "normal_2",
			Voltage:     3.4,
			RSSI:        -50,
			Temperature: 90,
			ExpectedMsg: false,
			ExpectMsg:   nil,
		},
		{
			Name:        "normal_3",
			Voltage:     3.2,
			RSSI:        -60,
			Temperature: 99.9,
			ExpectedMsg: false,
			ExpectMsg:   nil,
		},
		{
			Name:        "normal_4",
			Voltage:     3.11111111,
			RSSI:        -64,
			Temperature: -19.99999,
			ExpectedMsg: false,
			ExpectMsg:   nil,
		},

		{
			Name:        "alert_warning_high_temp",
			Voltage:     3.2,
			RSSI:        -60,
			Temperature: 100,
			ExpectedMsg: true,
			ExpectMsg: []string{
				rules.TypeAlertWarning,
				"esp32_1",
				"warning high temp cpu",
			},
		},
		{
			Name:        "alert_warning_high_temp_and_voltage",
			Voltage:     3.51,
			RSSI:        -60,
			Temperature: 100,
			ExpectedMsg: true,
			ExpectMsg: []string{
				rules.TypeAlertWarning,
				"esp32_1",
				"warning high temp cpu",
				"warning high voltage",
			},
		},
		{
			Name:        "alert_warning_high_temp_and_voltage_and_rssi",
			Voltage:     3.51,
			RSSI:        -65,
			Temperature: 100,
			ExpectedMsg: true,
			ExpectMsg: []string{
				rules.TypeAlertWarning,
				"esp32_1",
				"warning high temp cpu",
				"warning high voltage",
				"warning low ethernet signal",
			},
		},

		{
			Name:        "alert_critical_high_temp",
			Voltage:     3.2,
			RSSI:        -60,
			Temperature: 120,
			ExpectedMsg: true,
			ExpectMsg: []string{
				rules.TypeAlertCritical,
				"esp32_1",
				"critical high temp cpu",
			},
		},
		{
			Name:        "alert_critical_high_temp_warning_rssi",
			Voltage:     3.2,
			RSSI:        -65,
			Temperature: 120,
			ExpectedMsg: true,
			ExpectMsg: []string{
				rules.TypeAlertCritical,
				"esp32_1",
				"critical high temp cpu",
				"warning low ethernet signal",
			},
		},
		{
			Name:        "alert_critical_high_temp_warning_low_voltage",
			Voltage:     3.1,
			RSSI:        -60,
			Temperature: 120,
			ExpectedMsg: true,
			ExpectMsg: []string{
				rules.TypeAlertCritical,
				"esp32_1",
				"critical high temp cpu",
				"warning low voltage",
			},
		},
		{
			Name:        "alert_critical_high_temp_critical_rssi",
			Voltage:     3.2,
			RSSI:        -100,
			Temperature: 120,
			ExpectedMsg: true,
			ExpectMsg: []string{
				rules.TypeAlertCritical,
				"esp32_1",
				"critical high temp cpu",
				"critical low ethernet signal",
			},
		},
		{
			Name:        "alert_critical_high_temp_critical_low_voltage",
			Voltage:     2,
			RSSI:        -60,
			Temperature: 120,
			ExpectedMsg: true,
			ExpectMsg: []string{
				rules.TypeAlertCritical,
				"esp32_1",
				"critical high temp cpu",
				"critical low voltage",
			},
		},

		{
			Name:        "alert_all_critical",
			Voltage:     2,
			RSSI:        -200,
			Temperature: 140,
			ExpectedMsg: true,
			ExpectMsg: []string{
				rules.TypeAlertCritical,
				"esp32_1",
				"critical high temp cpu",
				"critical low ethernet signal",
				"critical low voltage",
			},
		},
		{
			Name:        "alert_all_warning",
			Voltage:     3.1,
			RSSI:        -65,
			Temperature: 100,
			ExpectedMsg: true,
			ExpectMsg: []string{
				rules.TypeAlertWarning,
				"esp32_1",
				"warning high temp cpu",
				"warning low ethernet signal",
				"warning low voltage",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			eb := events.NewEventBus(10, 5*time.Second)

			checker := board_info_rules.NewBoardInfoChecker(board_info_rules.NewRSsiChecker(0, 0),
				board_info_rules.NewTemperatureCpuCheck(),
				board_info_rules.NewVoltageChecker(0, 0, 0, 0))

			msgMock := notifiers.NewGetMessageMock(10)

			am, err := alerts.NewAlertsManager(eb, []rules.Rule{checker}, 0, 0, msgMock)
			if err != nil {
				t.Fatal(err)
			}

			am.Start()
			sub := eb.Subscribe(alerts.TopicforBoardInfoChecker, "complex_tests")
			err = eb.Publish(sub.Topic, events.Event{
				BoardID: "esp32_1",
				Payload: &dto.UpdateBoardInfo{
					CpuTemp:  getPtrFloat(tc.Temperature),
					Voltage:  getPtrFloat(tc.Voltage),
					RssiWifi: getPtrInt(tc.RSSI),
				},
			}, sub.ID)
			if err != nil {
				t.Fatal(err)
			}
			select {
			case message := <-msgMock.MsgCh:
				if tc.ExpectedMsg {
					if !checkMessage(message, tc.ExpectMsg) {
						t.Errorf("ERROR MESSAGE: got: %s, expect: %v\n", message, tc.ExpectMsg)
					}
				} else {
					t.Errorf("unexpected msg: %s\n", message)
				}
			case <-time.After(200 * time.Millisecond):
				if tc.ExpectedMsg {
					t.Errorf("got: block chan, expect: %v\n", tc.ExpectMsg)
				}
			}
		})
	}
}
func getPtrFloat(f float64) *float64 {
	return &f
}

func getPtrInt(i int) *int {
	return &i
}

func checkMessage(msg string, expect []string) bool {
	match := 0

	for _, word := range expect {
		if strings.Contains(msg, word) {
			match++
		}
	}
	return match == len(expect)
}
