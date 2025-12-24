package alerts_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/notifiers"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	board_info_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardInfoRules"
	board_status_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardStatusRules"
	database_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/dataBaseRules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

type TestCaseNewAlertsManager struct {
	Name               string
	EventBus           *events.EventBus
	Ruls               []rules.Rule
	BufferMessage      int
	MaxMessageParallel int
	Recipients         []notifiers.Notifier
	ExpectedError      error
}

func TestNewAlertsManager(t *testing.T) {
	testCases := []TestCaseNewAlertsManager{
		{
			Name:     "real_test_all_rules",
			EventBus: events.NewEventBus(10, 10*time.Second),
			Ruls: func() []rules.Rule {
				ruls := []rules.Rule{}
				ruls = append(ruls, board_info_rules.NewBoardInfoChecker(
					board_info_rules.NewRSsiChecker(0, 0),
					board_info_rules.NewTemperatureCpuCheck(),
					board_info_rules.NewVoltageChecker(0, 0, 0, 0)),
				)
				st, _ := board_status_rules.NewBoardStatusChecker(true, true, true)
				ruls = append(ruls, st)
				dbStat, _ := database_rules.NewDataBaseStatusChecker(true, true, true, true)
				dbErr, _ := database_rules.NewErrorDBChecker(true, true, true, nil, nil, rules.TypeAlertWarning)
				dbCheck, err := database_rules.NewDataBaseChecker(dbStat, dbErr)
				if err != nil {
					t.Fatal(err)
				}
				ruls = append(ruls, dbCheck)
				return ruls
			}(),
			BufferMessage:      10,
			MaxMessageParallel: 10,
			Recipients:         []notifiers.Notifier{&notifiers.LogNotifier{}},
			ExpectedError:      nil,
		},
		{
			Name:     "real_test_no_board_info_checker",
			EventBus: events.NewEventBus(10, 10*time.Second),
			Ruls: func() []rules.Rule {
				ruls := []rules.Rule{}
				st, _ := board_status_rules.NewBoardStatusChecker(true, true, true)
				ruls = append(ruls, st)
				dbStat, _ := database_rules.NewDataBaseStatusChecker(true, true, true, true)
				dbErr, _ := database_rules.NewErrorDBChecker(true, true, true, nil, nil, rules.TypeAlertWarning)
				dbCheck, err := database_rules.NewDataBaseChecker(dbStat, dbErr)
				if err != nil {
					t.Fatal(err)
				}
				ruls = append(ruls, dbCheck)
				return ruls
			}(),
			BufferMessage:      10,
			MaxMessageParallel: 10,
			Recipients:         []notifiers.Notifier{&notifiers.LogNotifier{}},
			ExpectedError:      nil,
		},
		{
			Name:     "real_test_no_board_status",
			EventBus: events.NewEventBus(10, 10*time.Second),
			Ruls: func() []rules.Rule {
				ruls := []rules.Rule{}
				ruls = append(ruls, board_info_rules.NewBoardInfoChecker(
					board_info_rules.NewRSsiChecker(0, 0),
					board_info_rules.NewTemperatureCpuCheck(),
					board_info_rules.NewVoltageChecker(0, 0, 0, 0)),
				)
				dbStat, _ := database_rules.NewDataBaseStatusChecker(true, true, true, true)
				dbErr, _ := database_rules.NewErrorDBChecker(true, true, true, nil, nil, rules.TypeAlertWarning)
				dbCheck, err := database_rules.NewDataBaseChecker(dbStat, dbErr)
				if err != nil {
					t.Fatal(err)
				}
				ruls = append(ruls, dbCheck)
				return ruls
			}(),
			BufferMessage:      10,
			MaxMessageParallel: 10,
			Recipients:         []notifiers.Notifier{&notifiers.LogNotifier{}},
			ExpectedError:      nil,
		},
		{
			Name:     "real_test_no_dbChecker",
			EventBus: events.NewEventBus(10, 10*time.Second),
			Ruls: func() []rules.Rule {
				ruls := []rules.Rule{}
				ruls = append(ruls, board_info_rules.NewBoardInfoChecker(
					board_info_rules.NewRSsiChecker(0, 0),
					board_info_rules.NewTemperatureCpuCheck(),
					board_info_rules.NewVoltageChecker(0, 0, 0, 0)),
				)
				st, _ := board_status_rules.NewBoardStatusChecker(true, true, true)
				ruls = append(ruls, st)
				return ruls
			}(),
			BufferMessage:      10,
			MaxMessageParallel: 10,
			Recipients:         []notifiers.Notifier{&notifiers.LogNotifier{}},
			ExpectedError:      nil,
		},
		{
			Name:     "real_test_no_rulles",
			EventBus: events.NewEventBus(10, 10*time.Second),
			Ruls: func() []rules.Rule {
				ruls := []rules.Rule{}
				ruls = append(ruls, board_info_rules.NewBoardInfoChecker(
					board_info_rules.NewRSsiChecker(0, 0),
					board_info_rules.NewTemperatureCpuCheck(),
					board_info_rules.NewVoltageChecker(0, 0, 0, 0)),
				)
				st, _ := board_status_rules.NewBoardStatusChecker(true, true, true)
				ruls = append(ruls, st)
				return ruls
			}(),
			BufferMessage:      10,
			MaxMessageParallel: 10,
			Recipients:         []notifiers.Notifier{&notifiers.LogNotifier{}},
			ExpectedError:      fmt.Errorf("no rules, alertsManager no point"),
		},
		{
			Name:     "real_test_default_buffer_size",
			EventBus: events.NewEventBus(10, 10*time.Second),
			Ruls: func() []rules.Rule {
				ruls := []rules.Rule{}
				ruls = append(ruls, board_info_rules.NewBoardInfoChecker(
					board_info_rules.NewRSsiChecker(0, 0),
					board_info_rules.NewTemperatureCpuCheck(),
					board_info_rules.NewVoltageChecker(0, 0, 0, 0)),
				)
				st, _ := board_status_rules.NewBoardStatusChecker(true, true, true)
				ruls = append(ruls, st)
				return ruls
			}(),
			BufferMessage:      0,
			MaxMessageParallel: 10,
			Recipients:         []notifiers.Notifier{&notifiers.LogNotifier{}},
			ExpectedError:      nil,
		},
		{
			Name:     "real_test_default_max_message_parallel",
			EventBus: events.NewEventBus(10, 10*time.Second),
			Ruls: func() []rules.Rule {
				ruls := []rules.Rule{}
				ruls = append(ruls, board_info_rules.NewBoardInfoChecker(
					board_info_rules.NewRSsiChecker(0, 0),
					board_info_rules.NewTemperatureCpuCheck(),
					board_info_rules.NewVoltageChecker(0, 0, 0, 0)),
				)
				st, _ := board_status_rules.NewBoardStatusChecker(true, true, true)
				ruls = append(ruls, st)
				return ruls
			}(),
			BufferMessage:      10,
			MaxMessageParallel: 0,
			Recipients:         []notifiers.Notifier{&notifiers.LogNotifier{}},
			ExpectedError:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			checker, err := alerts.NewAlertsManager(tc.EventBus, tc.Ruls, tc.BufferMessage, tc.MaxMessageParallel, tc.Recipients...)
			if err != nil {
				if tc.ExpectedError != nil {
					if tc.ExpectedError.Error() != err.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", tc.ExpectedError.Error(), err.Error())
					}
				} else {
					t.Errorf("ERROR: unexpected error: %s\n", err.Error())
				}
				return
			}
			if checker == nil {
				t.Errorf("checker == nil!\n")
			}
		})
	}
}
