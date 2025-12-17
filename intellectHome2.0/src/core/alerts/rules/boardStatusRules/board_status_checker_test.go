package board_status_rules_test

import (
	"fmt"
	"testing"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	board_status_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardStatusRules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

type TestCaseBoardStatusChecker struct {
	name          string
	active        bool
	lost          bool
	offline       bool
	body          events.Event
	expectedError error
	expectedAlert bool
	expectedData  string
	expectedType  string
}

func TestBoardStatusChecker(t *testing.T) {
	testCases := []TestCaseBoardStatusChecker{
		{
			name:    "normal_all_check_input_active",
			active:  true,
			lost:    true,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "update status to active",
			},
			expectedAlert: true,
			expectedData:  "update status to active",
			expectedType:  rules.TypeAlertNormal,
		},
		{
			name:    "normal_all_check_input_lost",
			active:  true,
			lost:    true,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "updated status to lost",
			},
			expectedAlert: true,
			expectedData:  "updated status to lost",
			expectedType:  rules.TypeAlertWarning,
		},
		{
			name:    "normal_all_check_input_offline",
			active:  true,
			lost:    true,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "updated status to offline",
			},
			expectedAlert: true,
			expectedData:  "updated status to offline",
			expectedType:  rules.TypeAlertCritical,
		},
		{
			name:    "normal_no_active_input_active",
			active:  false,
			lost:    true,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "updated status to active",
			},
			expectedAlert: false,
		},
		{
			name:    "normal_no_active_input_lost",
			active:  false,
			lost:    true,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "updated status to lost",
			},
			expectedAlert: true,
			expectedData:  "updated status to lost",
			expectedType:  rules.TypeAlertWarning,
		},
		{
			name:    "normal_no_active_input_offline",
			active:  false,
			lost:    true,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "updated status to offline",
			},
			expectedAlert: true,
			expectedData:  "updated status to offline",
			expectedType:  rules.TypeAlertCritical,
		},
		{
			name:    "normal_no_lost_input_lost",
			active:  true,
			lost:    false,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "updated status to lost",
			},
			expectedAlert: false,
		},
		{
			name:    "normal_no_lost_input_active",
			active:  true,
			lost:    false,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "updated status to active",
			},
			expectedAlert: true,
			expectedData:  "updated status to active",
			expectedType:  rules.TypeAlertNormal,
		},
		{
			name:    "normal_no_lost_input_offline",
			active:  true,
			lost:    false,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "updated status to offline",
			},
			expectedAlert: true,
			expectedData:  "updated status to offline",
			expectedType:  rules.TypeAlertCritical,
		},
		{
			name:    "normal_no_offline_input_offline",
			active:  true,
			lost:    true,
			offline: false,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "update status to offline",
			},
			expectedAlert: false,
		},
		{
			name:    "normal_no_offline_input_active",
			active:  true,
			lost:    true,
			offline: false,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "update status to active",
			},
			expectedAlert: true,
			expectedType:  rules.TypeAlertNormal,
			expectedData:  "update status to active",
		},
		{
			name:    "normal_no_offline_input_lost",
			active:  true,
			lost:    true,
			offline: false,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "update status to lost",
			},
			expectedAlert: true,
			expectedType:  rules.TypeAlertWarning,
			expectedData:  "update status to lost",
		},
		{
			name:    "no_rules",
			active:  false,
			lost:    false,
			offline: false,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "updated status to lost",
			},
			expectedAlert: false,
			expectedError: fmt.Errorf("no rules, BoardStatusChecker no point"),
		},
		{
			name:    "no_board_id",
			active:  true,
			lost:    true,
			offline: true,
			body: events.Event{
				BoardID: "",
				Payload: "updated status to lost",
			},
			expectedAlert: false,
			expectedError: fmt.Errorf("invalid boardID"),
		},
		{
			name:    "empty_payload",
			active:  true,
			lost:    true,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "",
			},
			expectedAlert: false,
		},
		{
			name:    "no_type_payload",
			active:  true,
			lost:    true,
			offline: true,
			body: events.Event{
				BoardID: "",
				Payload: fmt.Errorf("lost"),
			},
			expectedAlert: false,
			expectedError: fmt.Errorf("invalid type in payload"),
		},
		{
			name:    "edge_case_active",
			active:  true,
			lost:    true,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "update status to inactive",
			},
			expectedAlert: false,
		},
		{
			name:    "edge_case_lost",
			active:  true,
			lost:    true,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "update status to no_lost",
			},
			expectedAlert: false,
		},
		{
			name:    "edge_case_offline",
			active:  true,
			lost:    true,
			offline: true,
			body: events.Event{
				BoardID: "esp32_2",
				Payload: "update status to no_offline",
			},
			expectedAlert: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			check, err := board_status_rules.NewBoardStatusChecker(tc.active, tc.lost, tc.offline)
			if err != nil {
				if tc.expectedError != nil {
					if err.Error() != tc.expectedError.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.expectedError.Error())
					}
				} else {
					t.Errorf("unexpected error: %v\n", err)
				}
				return
			}
			al, err := check.Check(tc.body)
			if err != nil {
				if tc.expectedError != nil {
					if err.Error() != tc.expectedError.Error() {
						t.Errorf("ERROR: got: %s, expect: %s\n", err.Error(), tc.expectedError.Error())
					}
				} else {
					t.Errorf("unexpected err: %s\n", err.Error())
				}
			}

			if tc.expectedAlert {
				if al.Type != tc.expectedType {
					t.Errorf("STATUS: got: %s, expect: %s\n", al.Type, tc.expectedType)
				}
				if al.Data != tc.expectedData {
					t.Errorf("DATA: got: %s, expect: %s\n", al.Data, tc.expectedData)
				}
				if *al.BoardID != tc.body.BoardID {
					t.Errorf("BOARD ID: got: %s, expect: %s\n", *al.BoardID, tc.body.BoardID)
				}
			} else {
				if al != nil {
					t.Errorf("unexpected alert: %s\n", al.Data)
				}
			}
		})
	}
}

func TestNewBoardStatusChecker_Validation(t *testing.T) {
	tests := []struct {
		name    string
		active  bool
		lost    bool
		offline bool
		wantErr bool
	}{
		{"no_rules", false, false, false, true},
		{"only_active", true, false, false, false},
		{"only_lost", false, true, false, false},
		{"only_offline", false, false, true, false},
		{"all", true, true, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := board_status_rules.NewBoardStatusChecker(
				tt.active, tt.lost, tt.offline)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
