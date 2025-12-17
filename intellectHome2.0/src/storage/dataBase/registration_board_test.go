package database_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/server"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/tests"
	data_base_methods "github.com/Piccadilly98/goProjects/intellectHome2.0/tests/unit/databaseMethods"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/tests/utilits/init_data_base"
)

type TestCaseRegistrationBoard struct {
	Name                   string
	BoardID                *string
	BoardName              *string
	BoardType              *string
	State                  *string
	ExpectedExistBoards    bool
	ExpectedExistBoardInfo bool
	ExpectedCode           int
	ExpectedBody           *dto.GetBoardDataDto
	Ctx                    context.Context
	ExpectedError          error
}

func TestRegistrationBoard(t *testing.T) {
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

	testCases := []TestCaseRegistrationBoard{
		{
			Name:                   "valid_test_1",
			BoardID:                tests.GetPtrStr("esp32_1_new"),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_1_new"),
				Name:       nil,
				BoardType:  nil,
				BoardState: tests.GetPtrStr("registred"),
			},
		},
		{
			Name:                   "valid_test_2",
			BoardID:                tests.GetPtrStr("esp32_2_new"),
			BoardName:              tests.GetPtrStr("new boards"),
			BoardType:              tests.GetPtrStr("all_type"),
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_2_new"),
				Name:       tests.GetPtrStr("new boards"),
				BoardType:  tests.GetPtrStr("all_type"),
				BoardState: tests.GetPtrStr("registred"),
			},
		},
		{
			Name:                   "valid_test_3",
			BoardID:                tests.GetPtrStr("esp32_3_new"),
			BoardName:              nil,
			BoardType:              tests.GetPtrStr("all_type"),
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_3_new"),
				Name:       nil,
				BoardType:  tests.GetPtrStr("all_type"),
				BoardState: tests.GetPtrStr("new"),
			},
		},
		{
			Name:                   "valid_test_4",
			BoardID:                tests.GetPtrStr("esp32_4_new"),
			BoardName:              tests.GetPtrStr("new name"),
			BoardType:              nil,
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_4_new"),
				Name:       tests.GetPtrStr("new name"),
				BoardType:  nil,
				BoardState: tests.GetPtrStr("new"),
			},
		},
		{
			Name:                   "valid_test_5_symbol_in_id",
			BoardID:                tests.GetPtrStr("esp32_5_new'"),
			BoardName:              tests.GetPtrStr("new name"),
			BoardType:              nil,
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_5_new'"),
				Name:       tests.GetPtrStr("new name"),
				BoardType:  nil,
				BoardState: tests.GetPtrStr("new"),
			},
		},
		{
			Name:                   "valid_test_6_tab_in_id",
			BoardID:                tests.GetPtrStr("esp32_6_new\t"),
			BoardName:              tests.GetPtrStr("new name"),
			BoardType:              nil,
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_6_new\t"),
				Name:       tests.GetPtrStr("new name"),
				BoardType:  nil,
				BoardState: tests.GetPtrStr("new"),
			},
		},
		{
			Name:                   "valid_test_7_space_in_id",
			BoardID:                tests.GetPtrStr("esp32_7 new"),
			BoardName:              tests.GetPtrStr("new name"),
			BoardType:              nil,
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_7 new"),
				Name:       tests.GetPtrStr("new name"),
				BoardType:  nil,
				BoardState: tests.GetPtrStr("new"),
			},
		},
		{
			Name:                   "valid_test_8_space_and_tab_in_id",
			BoardID:                tests.GetPtrStr("esp32_7 new\t"),
			BoardName:              tests.GetPtrStr("new name"),
			BoardType:              nil,
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_7 new\t"),
				Name:       tests.GetPtrStr("new name"),
				BoardType:  nil,
				BoardState: tests.GetPtrStr("new"),
			},
		},
		{
			Name:                   "valid_test_9_tab_in_name",
			BoardID:                tests.GetPtrStr("esp32_8_new"),
			BoardName:              tests.GetPtrStr("new name\t"),
			BoardType:              nil,
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_8_new"),
				Name:       tests.GetPtrStr("new name\t"),
				BoardType:  nil,
				BoardState: tests.GetPtrStr("new"),
			},
		},
		{
			Name:                   "valid_test_10_double_tab_in_name",
			BoardID:                tests.GetPtrStr("esp32_9_new"),
			BoardName:              tests.GetPtrStr("new name\t\t\t\t\t"),
			BoardType:              nil,
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_9_new"),
				Name:       tests.GetPtrStr("new name\t\t\t\t\t"),
				BoardType:  nil,
				BoardState: tests.GetPtrStr("new"),
			},
		},
		{
			Name:                   "valid_test_11_tab_in_type",
			BoardID:                tests.GetPtrStr("esp32_10_new"),
			BoardName:              nil,
			BoardType:              tests.GetPtrStr("\t\t\t"),
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_10_new"),
				Name:       nil,
				BoardType:  tests.GetPtrStr("\t\t\t"),
				BoardState: tests.GetPtrStr("new"),
			},
		},
		{
			Name:                   "valid_test_12_tab_and_space_in_type",
			BoardID:                tests.GetPtrStr("esp32_11_new"),
			BoardName:              nil,
			BoardType:              tests.GetPtrStr("\t\t   \t"),
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_11_new"),
				Name:       nil,
				BoardType:  tests.GetPtrStr("\t\t   \t"),
				BoardState: tests.GetPtrStr("new"),
			},
		},
		//          EDGE

		{
			Name:                   "valid_test_13_edge_len_id",
			BoardID:                tests.GetPtrStr(strings.Repeat("a", 50)),
			BoardName:              nil,
			BoardType:              tests.GetPtrStr("\t\t   \t"),
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr(strings.Repeat("a", 50)),
				Name:       nil,
				BoardType:  tests.GetPtrStr("\t\t   \t"),
				BoardState: tests.GetPtrStr("new"),
			},
		},
		{
			Name:                   "valid_test_14_edge_len_name",
			BoardID:                tests.GetPtrStr("esp32_"),
			BoardName:              tests.GetPtrStr(strings.Repeat("a", 50)),
			BoardType:              tests.GetPtrStr("\t\t   \t"),
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_"),
				Name:       tests.GetPtrStr(strings.Repeat("a", 50)),
				BoardType:  tests.GetPtrStr("\t\t   \t"),
				BoardState: tests.GetPtrStr("new"),
			},
		},
		{
			Name:                   "valid_test_15_edge_len_type",
			BoardID:                tests.GetPtrStr("esp32"),
			BoardName:              nil,
			BoardType:              tests.GetPtrStr(strings.Repeat("A", 50)),
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32"),
				Name:       nil,
				BoardType:  tests.GetPtrStr(strings.Repeat("A", 50)),
				BoardState: tests.GetPtrStr("new"),
			},
		},
		{
			Name:                   "valid_test_16_edge_len_state",
			BoardID:                tests.GetPtrStr("esp32_32"),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr(strings.Repeat("'", 50)),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_32"),
				Name:       nil,
				BoardType:  nil,
				BoardState: tests.GetPtrStr(strings.Repeat("'", 50)),
			},
		},
		{
			Name:                   "valid_test_17_edge_len_all_column",
			BoardID:                tests.GetPtrStr(strings.Repeat("'", 50)),
			BoardName:              tests.GetPtrStr(strings.Repeat("'", 50)),
			BoardType:              tests.GetPtrStr(strings.Repeat("'", 50)),
			State:                  tests.GetPtrStr(strings.Repeat("'", 50)),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr(strings.Repeat("'", 50)),
				Name:       tests.GetPtrStr(strings.Repeat("'", 50)),
				BoardType:  tests.GetPtrStr(strings.Repeat("'", 50)),
				BoardState: tests.GetPtrStr(strings.Repeat("'", 50)),
			},
		},
		{
			Name:                   "valid_test_18_edge_len_id_utf8",
			BoardID:                tests.GetPtrStr(strings.Repeat("щ", 50)),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr(strings.Repeat("щ", 50)),
				Name:       nil,
				BoardType:  nil,
				BoardState: tests.GetPtrStr("registred"),
			},
		},
		{
			Name:                   "valid_test_19_edge_len_name_utf8",
			BoardID:                tests.GetPtrStr("esp32_32_32"),
			BoardName:              tests.GetPtrStr(strings.Repeat("п", 50)),
			BoardType:              nil,
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_32_32"),
				Name:       tests.GetPtrStr(strings.Repeat("п", 50)),
				BoardType:  nil,
				BoardState: tests.GetPtrStr("registred"),
			},
		},
		{
			Name:                   "valid_test_20_edge_len_type_utf8",
			BoardID:                tests.GetPtrStr("esp32_32_32_2"),
			BoardName:              nil,
			BoardType:              tests.GetPtrStr(strings.Repeat("ц", 50)),
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_32_32_2"),
				Name:       nil,
				BoardType:  tests.GetPtrStr(strings.Repeat("ц", 50)),
				BoardState: tests.GetPtrStr("registred"),
			},
		},
		{
			Name:                   "valid_test_21_edge_len_state_utf8",
			BoardID:                tests.GetPtrStr("esp31"),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr(strings.Repeat("г", 50)),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp31"),
				Name:       nil,
				BoardType:  nil,
				BoardState: tests.GetPtrStr(strings.Repeat("г", 50)),
			},
		},
		{
			Name:                   "valid_test_22_edge_len_all_column_utf8",
			BoardID:                tests.GetPtrStr(strings.Repeat("б", 50)),
			BoardName:              tests.GetPtrStr(strings.Repeat("б", 50)),
			BoardType:              tests.GetPtrStr(strings.Repeat("б", 50)),
			State:                  tests.GetPtrStr(strings.Repeat("б", 50)),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr(strings.Repeat("б", 50)),
				Name:       tests.GetPtrStr(strings.Repeat("б", 50)),
				BoardType:  tests.GetPtrStr(strings.Repeat("б", 50)),
				BoardState: tests.GetPtrStr(strings.Repeat("б", 50)),
			},
		},
		{
			Name:                   "invalid_test_1_repeat_id",
			BoardID:                tests.GetPtrStr("esp32_1_new"),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr("new"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           http.StatusBadRequest,
		},
		{
			Name:                   "invalid_test_2_no_state",
			BoardID:                tests.GetPtrStr("esp32_5_new"),
			BoardName:              nil,
			BoardType:              nil,
			State:                  nil,
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_3_no_name",
			BoardID:                nil,
			BoardName:              nil,
			BoardType:              nil,
			State:                  nil,
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_4_no_name_with_body",
			BoardID:                nil,
			BoardName:              tests.GetPtrStr("new"),
			BoardType:              tests.GetPtrStr("new"),
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_5_no_name_with_body",
			BoardID:                nil,
			BoardName:              nil,
			BoardType:              tests.GetPtrStr("new"),
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_6_no_name_with_body",
			BoardID:                nil,
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_7_long_id",
			BoardID:                tests.GetPtrStr(strings.Repeat("a", 51)),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_8_very_long_id",
			BoardID:                tests.GetPtrStr(strings.Repeat("a", 1000)),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_9_very_long_id_utf",
			BoardID:                tests.GetPtrStr(strings.Repeat("ф", 1000)),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_10_long_name",
			BoardID:                tests.GetPtrStr("esp30"),
			BoardName:              tests.GetPtrStr(strings.Repeat("a", 51)),
			BoardType:              nil,
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_11_very_long_name",
			BoardID:                tests.GetPtrStr("esp30"),
			BoardName:              tests.GetPtrStr(strings.Repeat("a", 1000)),
			BoardType:              nil,
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_12_very_long_name_utf",
			BoardID:                tests.GetPtrStr("esp30"),
			BoardName:              tests.GetPtrStr(strings.Repeat("ф", 1000)),
			BoardType:              nil,
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_13_long_type",
			BoardID:                tests.GetPtrStr("esp30"),
			BoardName:              nil,
			BoardType:              tests.GetPtrStr(strings.Repeat("a", 51)),
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_14_very_long_type",
			BoardID:                tests.GetPtrStr("esp30"),
			BoardName:              nil,
			BoardType:              tests.GetPtrStr(strings.Repeat("a", 1000)),
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_15_very_long_type_utf",
			BoardID:                tests.GetPtrStr("esp30"),
			BoardName:              nil,
			BoardType:              tests.GetPtrStr(strings.Repeat("и", 1000)),
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_16_long_state",
			BoardID:                tests.GetPtrStr("esp30"),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr(strings.Repeat("g", 51)),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_16_very_long_state",
			BoardID:                tests.GetPtrStr("esp30"),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr(strings.Repeat("g", 1000)),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_17_very_long_state_utf",
			BoardID:                tests.GetPtrStr("esp30"),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr(strings.Repeat("ф", 1000)),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_18_long_all",
			BoardID:                tests.GetPtrStr(strings.Repeat("g", 51)),
			BoardName:              tests.GetPtrStr(strings.Repeat("g", 51)),
			BoardType:              tests.GetPtrStr(strings.Repeat("g", 51)),
			State:                  tests.GetPtrStr(strings.Repeat("g", 51)),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_19_very_long_all",
			BoardID:                tests.GetPtrStr(strings.Repeat("g", 1000)),
			BoardName:              tests.GetPtrStr(strings.Repeat("g", 1000)),
			BoardType:              tests.GetPtrStr(strings.Repeat("g", 1000)),
			State:                  tests.GetPtrStr(strings.Repeat("g", 1000)),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},
		{
			Name:                   "invalid_test_20_very_long_all_utf",
			BoardID:                tests.GetPtrStr(strings.Repeat("☺️", 1000)),
			BoardName:              tests.GetPtrStr(strings.Repeat("☺️", 1000)),
			BoardType:              tests.GetPtrStr(strings.Repeat("☺️", 1000)),
			State:                  tests.GetPtrStr(strings.Repeat("☺️", 1000)),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           http.StatusInternalServerError,
		},

		// 			CONTEXT

		{
			Name:                   "test_with_normal_context_valid_body",
			BoardID:                tests.GetPtrStr("esp32_1_ctx"),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    true,
			ExpectedExistBoardInfo: true,
			ExpectedCode:           201,
			Ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				return ctx
			}(),
			ExpectedBody: &dto.GetBoardDataDto{
				BoardId:    tests.GetPtrStr("esp32_1_ctx"),
				Name:       nil,
				BoardType:  nil,
				BoardState: tests.GetPtrStr("registred"),
			},
		},
		{
			Name:                   "test_with_canceled_context_valid_body",
			BoardID:                tests.GetPtrStr("esp32_2_ctx"),
			BoardName:              nil,
			BoardType:              nil,
			State:                  tests.GetPtrStr("registred"),
			ExpectedExistBoards:    false,
			ExpectedExistBoardInfo: false,
			ExpectedCode:           0,
			Ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Microsecond)
				return ctx
			}(),
			ExpectedError: context.Canceled,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var code int
			var err error
			if tc.Ctx == nil {
				code, _ = serv.Db.RegistrationBoard(context.Background(), tc.BoardID, tc.BoardName, tc.BoardType, tc.State)
			} else {
				code, err = serv.Db.RegistrationBoard(tc.Ctx, tc.BoardID, tc.BoardName, tc.BoardType, tc.State)
				if err != nil {
					if !errors.Is(err, tc.ExpectedError) {
						t.Errorf("ERROR IN EXPECTED ERR: got: %s, expect: %s\n", err.Error(), tc.ExpectedError.Error())
					}
				}
			}
			if code != tc.ExpectedCode {
				t.Errorf("CODE: got: %d, expect: %d\n", code, tc.ExpectedCode)
			}

			if tc.ExpectedExistBoards {
				exist, _, err := serv.Db.GetExistWithBoardId(context.Background(), *tc.BoardID)
				if err != nil {
					t.Errorf("ERROR: %v\n", err.Error())
				}
				if exist != tc.ExpectedExistBoards {
					t.Errorf("EXIST BOARDS: got: %v, expect: %v\n", exist, tc.ExpectedExistBoards)
				}
			}
			if tc.ExpectedExistBoardInfo {
				exist, err := data_base_methods.CheckExistBoardInfo(serv.Db, *tc.BoardID)
				if err != nil {
					t.Errorf("ERROR: %v\n", err.Error())
				}
				if exist != tc.ExpectedExistBoardInfo {
					t.Errorf("EXIST BOARDINFO: got: %v, expect: %v\n", exist, tc.ExpectedExistBoardInfo)
				}
			}
			if tc.ExpectedBody != nil {
				dto, code, err := serv.Db.GetDtoWithId(context.Background(), *tc.BoardID)
				if code != http.StatusOK {
					t.Errorf("CODE IN DTO: got: %d, expect 200\n", code)
				}
				if err != nil {
					t.Error(err)
				}
				if *dto.BoardId != *tc.ExpectedBody.BoardId {
					t.Errorf("ERROR IN DB BOARD ID: got: %s, expect: %s\n", *dto.BoardId, *tc.ExpectedBody.BoardId)
				}
				if *dto.BoardState != *tc.ExpectedBody.BoardState {
					t.Errorf("ERROR IN DB BOARD STATE: got: %s, expect: %s\n", *dto.BoardState, *tc.ExpectedBody.BoardState)
				}
				if tc.ExpectedBody.BoardType != nil {
					if dto.BoardType != nil {
						if *dto.BoardType != *tc.ExpectedBody.BoardType {
							t.Errorf("ERROR IN DB BOARD TYPE: got: %s, expect: %s\n", *dto.BoardType, *tc.ExpectedBody.BoardType)
						}
					} else {
						t.Errorf("ERROR IN DB BOARD TYPE: got: nil, expect: %s\n", *tc.ExpectedBody.BoardType)
					}
				}
				if tc.ExpectedBody.Name != nil {
					if dto.Name != nil {
						if *dto.Name != *tc.ExpectedBody.Name {
							t.Errorf("ERROR IN DB BOARD NAME: got: %s, expect: %s\n", *dto.Name, *tc.ExpectedBody.Name)
						}
					} else {
						t.Errorf("ERROR IN DB BOARD NAME: got: nil, expect: %s\n", *tc.ExpectedBody.Name)
					}
				}
			}
		})
	}
}
