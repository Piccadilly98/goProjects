package e_2_e_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/server"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/tests"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/tests/utilits/init_data_base"
	_ "github.com/lib/pq"
)

func TestBoardLifeCycle(t *testing.T) {
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
	dtoRegistrationBoard := &dto.RegistrationBoardDTO{
		BoardId:    "esp32_2_new",
		Name:       nil,
		BoardState: nil,
		BoardType:  nil,
	}
	t.Run("Board_registration_and_update_name", func(t *testing.T) {
		t.Run("Stage_1_Registration_board", func(t *testing.T) {
			b, err := json.Marshal(dtoRegistrationBoard)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/boards", bytes.NewBuffer(b))
			serv.R.ServeHTTP(w, req)

			if w.Code != http.StatusCreated {
				t.Errorf("ERROR CODE: got: %d, expect: 201\n", w.Code)
				t.Fatal()
			}
			dtoAfterRegistration := &dto.GetBoardDataDto{}
			err = json.NewDecoder(w.Body).Decode(dtoAfterRegistration)
			if err != nil {
				t.Fatal(err)
			}
			if *dtoAfterRegistration.BoardId != dtoRegistrationBoard.BoardId {
				t.Errorf("ERROR IN RESPONSE ID: got: %s, expect: %s\n", *dtoAfterRegistration.BoardId, dtoRegistrationBoard.BoardId)
				t.Fatal()
			}
			if dtoAfterRegistration.Name != nil {
				t.Errorf("ERROR IN RESPONSE: got: %s, expect: nil\n", *dtoAfterRegistration.Name)
				t.Fatal()
			}
			if dtoAfterRegistration.BoardState != nil {
				if *dtoAfterRegistration.BoardState != "registred" {
					t.Errorf("ERROR IN RESPONSE DEFAULT STATE: got: %s, expect: registred\n", *dtoAfterRegistration.BoardState)
					t.Fatal()
				}
			} else {
				t.Errorf("ERROR IN RESPONSE DEFAULT STATE: got: nil, expect: registred\n")
				t.Fatal()
			}
			if dtoAfterRegistration.BoardType != nil {
				if *dtoAfterRegistration.BoardType != "esp32_all_task" {
					t.Errorf("ERROR IN RESPONSE DEFAULT TYPE: got: %s, expect: esp32_all_task\n", *dtoAfterRegistration.BoardType)
					t.Fatal()
				}
			} else {
				t.Errorf("ERROR IN RESPONSE DEFAULT TYPE: got: nil, expect: esp32_all_task\n")
				t.Fatal()
			}
			if dtoAfterRegistration.CreatedDate == nil {
				t.Errorf("ERROR IN RESPONSE: time create == nil\n")
				t.Fatal()
			}
			if dtoAfterRegistration.UpdatedDate != nil {
				t.Errorf("ERROR IN RESPONSE: time update != nil\n")
				t.Fatal()
			}
			exist, _, err := serv.Db.GetExistWithBoardId(context.Background(), dtoRegistrationBoard.BoardId)
			if err != nil {
				t.Error(err)
				return
			}
			if !exist {
				t.Errorf("ERROR IN DB: board not created\n")
				t.Fatal()
			}

			var existCheck bool
			err = serv.Db.GetPointer().QueryRow(`
			SELECT EXISTS(SELECT 1 FROM boards
			WHERE board_id = 'esp32_2_new');
			`).Scan(&existCheck)
			if err != nil {
				t.Fatal(err)
			}
			if existCheck != exist {
				t.Errorf("ERROR IN DB METHOD check exist != result getExistBoardID")
				t.Fatal()
			}

			t.Log("\nREGISTRATION_COMPLETE!\n")
		})

		t.Run("Stage_2_update_boards_data_name_and_type_and_get", func(t *testing.T) {
			update := &dto.UpdateBoardDataDto{
				Name: tests.GetPtrStr("new_name"),
				Type: tests.GetPtrStr("new_type"),
			}
			b, err := json.Marshal(update)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPatch, "/boards/esp32_2_new", bytes.NewBuffer(b))

			serv.R.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("ERROR IN RESPONSE CODE: got: %d, expect: 200\n", w.Code)
				t.Fatal()
			}

			dtoDB, _, err := serv.Db.GetDtoWithId(context.Background(), dtoRegistrationBoard.BoardId)
			if err != nil {
				t.Fatal(err)
			}

			if *dtoDB.BoardId != dtoRegistrationBoard.BoardId {
				t.Errorf("ERROR IN DB ID: got: %s, expect: %s\n", *dtoDB.BoardId, dtoRegistrationBoard.BoardId)
				t.Fatal()
			}

			if dtoDB.UpdatedDate == nil {
				t.Errorf("ERROR IN DB: update date == nil\n")
				t.Fatal()
			}

			if dtoDB.CreatedDate == nil {
				t.Errorf("ERROR IN DB: created date == nil\n")
				t.Fatal()
			}

			if dtoDB.Name != nil {
				if *dtoDB.Name != *update.Name {
					t.Errorf("ERROR IN DB NAME: got: %s, expect: %s\n", *dtoDB.Name, *update.Name)
					t.Fatal()
				}
			} else {
				t.Errorf("ERROR IN DB: name == nil\n")
				t.Fatal()
			}
			if dtoDB.BoardType != nil {
				if *dtoDB.BoardType != *update.Type {
					t.Errorf("ERROR IN DB TYPE: got: %s, expect: %s\n", *dtoDB.BoardType, *update.Type)
					t.Fatal()
				}
			} else {
				t.Errorf("ERROR IN DB: type == nil\n")
				t.Fatal()
			}

			if dtoDB.BoardState != nil {
				if *dtoDB.BoardState != "registred" {
					t.Errorf("ERROR IN DB STATE: got: %s, expect: registred\n", *dtoDB.BoardState)
					t.Fatal()
				}
			} else {
				t.Errorf("ERROR IN DB: state == nil\n")
				t.Fatal()
			}
			w = httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/boards/esp32_2_new", nil)

			serv.R.ServeHTTP(w, req)

			dtoGet := &dto.GetBoardDataDto{}
			err = json.NewDecoder(w.Body).Decode(dtoGet)
			if err != nil {
				t.Fatal(err)
			}
		})
	})
}
