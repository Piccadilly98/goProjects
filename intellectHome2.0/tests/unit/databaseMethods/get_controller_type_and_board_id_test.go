package data_base_methods

import (
	"context"
	"net/http"
	"testing"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/server"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/tests/utilits/init_data_base"
)

type TestCaseGetControllerTypeAndBoardID struct {
	Name          string
	ControllerID  string
	ExpectType    string
	ExpectBoardID string
	ExpectedCode  int
}

const (
	TypeSensor = "sensor"
	TypeBinary = "binary"
)

func TestGetControllerTypeAndBoardID(t *testing.T) {
	serv, err := server.NewServer(true, 30, 120, false)
	if err != nil {
		init_data_base.Cleanup(serv.Db)
		t.Fatalf("error in create db: %s", err.Error())
	}
	err = init_data_base.InitDataBase(serv.Db)
	if err != nil {
		init_data_base.Cleanup(serv.Db)
		t.Fatalf("error in init db: %s", err.Error())
	}
	defer func() {
		init_data_base.Cleanup(serv.Db)
		serv.Db.Close()
	}()

	testCases := []TestCaseGetControllerTypeAndBoardID{
		{
			Name:          "valid_test_1",
			ControllerID:  "led1",
			ExpectType:    TypeBinary,
			ExpectBoardID: "esp32_2_test",
			ExpectedCode:  http.StatusOK,
		},
		{
			Name:          "valid_test_2",
			ControllerID:  "led2",
			ExpectType:    TypeBinary,
			ExpectBoardID: "esp32_1_test",
			ExpectedCode:  http.StatusOK,
		},
		{
			Name:          "valid_test_3",
			ControllerID:  "led3",
			ExpectType:    TypeBinary,
			ExpectBoardID: "esp32_3_test",
			ExpectedCode:  http.StatusOK,
		},
		{
			Name:          "valid_test_4",
			ControllerID:  "led4",
			ExpectType:    TypeSensor,
			ExpectBoardID: "esp32_3_test",
			ExpectedCode:  http.StatusOK,
		},
		{
			Name:          "invalid_test_1",
			ControllerID:  "led12",
			ExpectType:    "",
			ExpectBoardID: "",
			ExpectedCode:  http.StatusBadRequest,
		},
		{
			Name:          "invalid_test_2",
			ControllerID:  "",
			ExpectType:    "",
			ExpectBoardID: "",
			ExpectedCode:  http.StatusBadRequest,
		},
		{
			Name:          "invalid_test_3",
			ControllerID:  "random",
			ExpectType:    "",
			ExpectBoardID: "",
			ExpectedCode:  http.StatusBadRequest,
		},
		{
			Name:          "invalid_test_4",
			ExpectType:    "",
			ExpectBoardID: "",
			ExpectedCode:  http.StatusBadRequest,
		},
		{
			Name:          "invalid_test_5",
			ControllerID:  "'select * from boards'",
			ExpectType:    "",
			ExpectBoardID: "",
			ExpectedCode:  http.StatusBadRequest,
		},
		{
			Name:          "invalid_test_6_attack_true",
			ControllerID:  "' OR '1'='1",
			ExpectType:    "",
			ExpectBoardID: "",
			ExpectedCode:  http.StatusBadRequest,
		},
		{
			Name:          "invalid_test_7_attack_update",
			ControllerID:  "'; UPDATE boards SET name='hacked'; -",
			ExpectType:    "",
			ExpectBoardID: "",
			ExpectedCode:  http.StatusBadRequest,
		},
		{
			Name:          "invalid_test_8_attack_drop_table",
			ControllerID:  "'; DROP TABLE boards; --",
			ExpectType:    "",
			ExpectBoardID: "",
			ExpectedCode:  http.StatusBadRequest,
		},
		{
			Name:          "invalid_test_9_attack_union",
			ControllerID:  "' UNION SELECT 1,1,1 FROM users --",
			ExpectType:    "",
			ExpectBoardID: "",
			ExpectedCode:  http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			controllerType, boardID, code, _ := serv.Db.GetControllerTypeAndBoardID(context.Background(), tc.ControllerID)
			if code != tc.ExpectedCode {
				t.Errorf("CODE: got: %d, expect: %d\n", code, tc.ExpectedCode)
			}

			if controllerType != tc.ExpectType {
				t.Errorf("TYPE: got: %s, expect: %s\n", controllerType, tc.ExpectType)
			}
			if boardID != tc.ExpectBoardID {
				t.Errorf("BOARD ID: got: %s, expected: %s\n", boardID, tc.ExpectBoardID)
			}
		})
	}
}
