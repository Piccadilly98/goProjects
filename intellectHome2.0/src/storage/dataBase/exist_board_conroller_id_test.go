package database_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/server"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/tests/utilits/init_data_base"
	_ "github.com/lib/pq"
)

type TestCaseExistBoardID struct {
	Name    string
	BoardID string
	Expect  bool
	Ctx     context.Context
}

type TestCaseExistControllerID struct {
	Name         string
	ControllerID string
	Expect       bool
}

func TestExistBoardIDControllerID(t *testing.T) {
	serv, err := server.NewServer(true, 30, 120, false, 0, 0)
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

	testCasesBoards := []TestCaseExistBoardID{
		{
			Name:    "BoardID: valid_test_1",
			BoardID: "esp32_1_test",
			Expect:  true,
		},
		{
			Name:    "BoardID: valid_test_2",
			BoardID: "esp32_2_test",
			Expect:  true,
		},
		{
			Name:    "BoardID: valid_test_3",
			BoardID: "esp32_3_test",
			Expect:  true,
		},
		{
			Name:    "BoardID: valid_test_4_repeat",
			BoardID: "esp32_2_test",
			Expect:  true,
		},
		{
			Name:    "BoardID: invalid_test_1_invalid_id",
			BoardID: "esp32_4_test",
			Expect:  false,
		},
		{
			Name:    "BoardID: invalid_test_2_empty_id",
			BoardID: "",
			Expect:  false,
		},
		{
			Name:    "BoardID: invalid_test_3_long_id",
			BoardID: strings.Repeat("a", 100000),
			Expect:  false,
		},
		{
			Name:    "BoardID: invalid_test_4_unicode symbols",
			BoardID: "esp32_тест_кириллица",
			Expect:  false,
		},
		{
			Name:    "BoardID: invalid_test_5_tab_id",
			BoardID: "\t\t\t\t\t\t\t\t",
			Expect:  false,
		},
		{
			Name:    "BoardID: invalid_test_6_tab_space_id",
			BoardID: "\t\t\t\t \t\t t\t\t",
			Expect:  false,
		},
		{
			Name:    "BoardID: invalid_test_7_attack_or",
			BoardID: "' OR '1'='1",
			Expect:  false,
		},
		{
			Name:    "BoardID: invalid_test_8_attack_comment",
			BoardID: "' OR 1=1 --",
			Expect:  false,
		},
		{
			Name:    "BoardID: invalid_test_9_attack_and",
			BoardID: "esp32_1' AND '1'='1",
			Expect:  false,
		},
		{
			Name:    "BoardID: invalid_test_10_attack_drop",
			BoardID: "; DROP TABLE boards; --",
			Expect:  false,
		},
		{
			Name:    "BoardID: invalid_test_11",
			BoardID: "`asd`",
			Expect:  false,
		},
		{
			Name:    "BoardID: invalid_test_12",
			BoardID: "`'DROP TABLE boards;'`",
			Expect:  false,
		},
		{
			Name:    "BoardID: invalid_test_13",
			BoardID: "''`a",
			Expect:  false,
		},
	}

	testCasesControllers := []TestCaseExistControllerID{
		{
			Name:         "ControllerID: valid_test_1",
			ControllerID: "led1",
			Expect:       true,
		},
		{
			Name:         "ControllerID: valid_test_2",
			ControllerID: "led2",
			Expect:       true,
		},
		{
			Name:         "ControllerID: valid_test_3",
			ControllerID: "led3",
			Expect:       true,
		},
		{
			Name:         "ControllerID: valid_test_4_senor",
			ControllerID: "led4",
			Expect:       true,
		},
		{
			Name:         "ControllerID: invalid_test_1",
			ControllerID: "led43",
			Expect:       false,
		},
		{
			Name:         "ControllerID: invalid_test_2",
			ControllerID: "423",
			Expect:       false,
		},
		{
			Name:         "ControllerID: invalid_test_3",
			ControllerID: "",
			Expect:       false,
		},
		{
			Name:         "ControllerID: invalid_test_4_long_id_edge",
			ControllerID: strings.Repeat("a", 50),
			Expect:       false,
		},
		{
			Name:         "ControllerID: invalid_test_4_long_id",
			ControllerID: strings.Repeat("a", 100000),
			Expect:       false,
		},
		{
			Name:         "ControllerID: invalid_test_5_unicode",
			ControllerID: "esp32_1_кухня",
			Expect:       false,
		},
		{
			Name:         "ControllerID: invalid_test_6_tab",
			ControllerID: "\t\t\t\t\t\t\t\t",
			Expect:       false,
		},
		{
			Name:         "ControllerID: invalid_test_7_tab_space",
			ControllerID: "\t \t\t \t\t \t\t\t",
			Expect:       false,
		},
		{
			Name:         "ControllerID: invalid_test_8_attack_or",
			ControllerID: "' OR '1'='1",
			Expect:       false,
		},
		{
			Name:         "ControllerID: invalid_test_9_attack_comment",
			ControllerID: "' OR 1=1 --",
			Expect:       false,
		},
		{
			Name:         "ControllerID: invalid_test_10_attack_drop",
			ControllerID: "; DROP TABLE boards; --",
			Expect:       false,
		},
		{
			Name:         "ControllerID: invalid_test_11",
			ControllerID: "`randomSymbol!`'s'`",
			Expect:       false,
		},
		{
			Name:         "ControllerID: invalid_test_12_long_id_unicode",
			ControllerID: strings.Repeat("у", 100000),
			Expect:       false,
		},
	}
	for _, tc := range testCasesBoards {
		t.Run(tc.Name, func(t *testing.T) {
			exist, code, err := serv.Db.GetExistWithBoardId(context.Background(), tc.BoardID)
			if err != nil {
				t.Error(err)
			}
			if code != http.StatusOK {
				t.Error("unexpected code!")
			}
			if exist != tc.Expect {
				t.Errorf("got: %v, expect: %v\n", exist, tc.Expect)
			}
		})
	}

	for _, tc := range testCasesControllers {
		t.Run(tc.Name, func(t *testing.T) {
			exist, code, err := serv.Db.GetExistWithControllerId(context.Background(), tc.ControllerID)
			if err != nil {
				t.Error(err)
			}
			if code != http.StatusOK {
				t.Error("unexpected code!")
			}
			if exist != tc.Expect {
				t.Errorf("got: %v, expect: %v\n", exist, tc.Expect)
			}
		})
	}
}
