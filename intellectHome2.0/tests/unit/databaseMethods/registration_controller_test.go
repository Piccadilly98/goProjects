package data_base_methods

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/server"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/tests"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/tests/utilits/init_data_base"
)

type TestCaseRegistrationController struct {
	Name                  string
	BoardID               string
	Body                  *dto.RegistrationController
	IsValidBody           bool
	Binary                bool
	Sensor                bool
	ExpectedCode          int
	ExpectedExistInBoards bool
	ExpectedError         error
	Ctx                   context.Context
}

func TestRegistrationController(t *testing.T) {
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

	testCases := []TestCaseRegistrationController{
		{
			Name:    "valid_test_1_binary",
			BoardID: "esp32_1_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeBinary,
				ControllerID:   "test1",
				Name:           tests.GetPtrStr("new name"),
				PinNumber:      tests.GetPtrInt(1),
				Type:           "new_type",
				Status:         tests.GetPtrBool(true),
			},
			IsValidBody:           true,
			ExpectedCode:          http.StatusCreated,
			ExpectedExistInBoards: true,
			ExpectedError:         nil,
			Binary:                true,
			Sensor:                false,
		},
		{
			Name:    "valid_test_2_binary_no_name",
			BoardID: "esp32_1_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeBinary,
				ControllerID:   "test1_no_name",
				Name:           nil,
				PinNumber:      tests.GetPtrInt(1),
				Type:           "new_type",
				Status:         tests.GetPtrBool(true),
			},
			IsValidBody:           true,
			ExpectedCode:          http.StatusCreated,
			ExpectedExistInBoards: true,
			ExpectedError:         nil,
			Binary:                true,
			Sensor:                false,
		},
		{
			Name:    "valid_test_3_binary_no_name_pin",
			BoardID: "esp32_1_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeBinary,
				ControllerID:   "test1_no_name_pin",
				Name:           nil,
				PinNumber:      nil,
				Type:           "new_type",
				Status:         tests.GetPtrBool(true),
			},
			IsValidBody:           true,
			ExpectedCode:          http.StatusCreated,
			ExpectedExistInBoards: true,
			ExpectedError:         nil,
			Binary:                true,
			Sensor:                false,
		},
		{
			Name:    "valid_test_4_sensor",
			BoardID: "esp32_2_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeSensor,
				ControllerID:   "test2",
				Name:           tests.GetPtrStr("new name"),
				PinNumber:      tests.GetPtrInt(1),
				Type:           "new_type",
				Unit:           tests.GetPtrStr("%"),
				Value:          tests.GetPtrInt(12),
			},
			IsValidBody:           true,
			ExpectedCode:          http.StatusCreated,
			ExpectedExistInBoards: true,
			ExpectedError:         nil,
			Binary:                false,
			Sensor:                true,
		},
		{
			Name:    "valid_test_5_sensor_no_name",
			BoardID: "esp32_2_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeSensor,
				ControllerID:   "test2_no_name",
				Name:           nil,
				PinNumber:      tests.GetPtrInt(1),
				Type:           "new_type",
				Unit:           tests.GetPtrStr("%"),
				Value:          tests.GetPtrInt(12),
			},
			IsValidBody:           true,
			ExpectedCode:          http.StatusCreated,
			ExpectedExistInBoards: true,
			ExpectedError:         nil,
			Binary:                false,
			Sensor:                true,
		},
		{
			Name:    "valid_test_6_sensor_no_name_no_pin",
			BoardID: "esp32_2_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeSensor,
				ControllerID:   "test2_no_name_pin",
				Name:           nil,
				PinNumber:      nil,
				Type:           "new_type",
				Unit:           tests.GetPtrStr("%"),
				Value:          tests.GetPtrInt(12),
			},
			IsValidBody:           true,
			ExpectedCode:          http.StatusCreated,
			ExpectedExistInBoards: true,
			ExpectedError:         nil,
			Binary:                false,
			Sensor:                true,
		},
		{
			Name:    "valid_test_7_sensor",
			BoardID: "esp32_3_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeSensor,
				ControllerID:   "test2_board_esp_32_3_test",
				Name:           tests.GetPtrStr("new name"),
				PinNumber:      tests.GetPtrInt(1),
				Type:           "new_type",
				Unit:           tests.GetPtrStr("%"),
				Value:          tests.GetPtrInt(12),
			},
			IsValidBody:           true,
			ExpectedCode:          http.StatusCreated,
			ExpectedExistInBoards: true,
			ExpectedError:         nil,
			Binary:                false,
			Sensor:                true,
		},

		//      INVALID

		{
			Name:    "invalid_test_0_empty_body",
			BoardID: "esp32_1_test",
			Body: &dto.RegistrationController{
				ControllerType: "",
				ControllerID:   "",
				Name:           nil,
				PinNumber:      nil,
				Type:           "",
				Status:         nil,
			},
			IsValidBody:           false,
			ExpectedCode:          http.StatusBadRequest,
			ExpectedExistInBoards: false,
			ExpectedError:         fmt.Errorf("invalid boardID"),
			Binary:                true,
			Sensor:                false,
		},
		{
			Name:    "invalid_test_1_binary_no_valid_board_id",
			BoardID: "esp32_4_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeBinary,
				ControllerID:   "test2_invalid_board",
				Name:           tests.GetPtrStr("new name"),
				PinNumber:      tests.GetPtrInt(1),
				Type:           "new_type",
				Status:         tests.GetPtrBool(true),
			},
			IsValidBody:           true,
			ExpectedCode:          http.StatusBadRequest,
			ExpectedExistInBoards: false,
			ExpectedError:         fmt.Errorf("invalid boardID"),
			Binary:                true,
			Sensor:                false,
		},
		{
			Name:    "invalid_test_2_binary_repeat_controller_id",
			BoardID: "esp32_3_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeBinary,
				ControllerID:   "test1",
				Name:           tests.GetPtrStr("new name"),
				PinNumber:      tests.GetPtrInt(1),
				Type:           "new_type",
				Status:         tests.GetPtrBool(true),
			},
			IsValidBody:           true,
			ExpectedCode:          http.StatusBadRequest,
			ExpectedExistInBoards: false,
			ExpectedError:         fmt.Errorf("invalid controllerID"),
			Binary:                true,
			Sensor:                false,
		},
		{
			Name:    "invalid_test_3_binary_invalid_controller_type",
			BoardID: "esp32_3_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeSensor,
				ControllerID:   "test1",
				Name:           tests.GetPtrStr("new name"),
				PinNumber:      tests.GetPtrInt(1),
				Type:           "new_type",
				Status:         tests.GetPtrBool(true),
			},
			IsValidBody:           false,
			ExpectedCode:          http.StatusBadRequest,
			ExpectedExistInBoards: false,
			ExpectedError:         fmt.Errorf("invalid controllerID"),
			Binary:                true,
			Sensor:                false,
		},
		{
			Name:    "invalid_test_4_binary_invalid_body_sensor",
			BoardID: "esp32_3_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeBinary,
				ControllerID:   "test1",
				Name:           tests.GetPtrStr("new name"),
				PinNumber:      tests.GetPtrInt(1),
				Type:           "new_type",
				Unit:           tests.GetPtrStr("%"),
				Value:          tests.GetPtrInt(12),
			},
			IsValidBody:           false,
			ExpectedCode:          http.StatusBadRequest,
			ExpectedExistInBoards: false,
			ExpectedError:         fmt.Errorf("invalid controllerID"),
			Binary:                true,
			Sensor:                false,
		},
		{
			Name:    "invalid_test_5_sensor_no_valid_board_id",
			BoardID: "esp32_4_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeSensor,
				ControllerID:   "test2_invalid_board",
				Name:           tests.GetPtrStr("new name"),
				PinNumber:      tests.GetPtrInt(1),
				Type:           "new_type",
				Unit:           tests.GetPtrStr("%"),
				Value:          tests.GetPtrInt(1),
			},
			IsValidBody:           true,
			ExpectedCode:          http.StatusBadRequest,
			ExpectedExistInBoards: false,
			ExpectedError:         fmt.Errorf("invalid boardID"),
			Binary:                false,
			Sensor:                true,
		},
		{
			Name:    "invalid_test_6_sensor_invalid_type",
			BoardID: "esp32_1_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeBinary,
				ControllerID:   "test2_invalid_board",
				Name:           tests.GetPtrStr("new name"),
				PinNumber:      tests.GetPtrInt(1),
				Type:           "new_type",
				Unit:           tests.GetPtrStr("%"),
				Value:          tests.GetPtrInt(1),
			},
			IsValidBody:           false,
			ExpectedCode:          http.StatusBadRequest,
			ExpectedExistInBoards: false,
			ExpectedError:         fmt.Errorf("invalid boardID"),
			Binary:                false,
			Sensor:                true,
		},
		{
			Name:                  "invalid_test_7_sensor_invalid_body_empty",
			BoardID:               "esp32_1_test",
			Body:                  &dto.RegistrationController{},
			IsValidBody:           false,
			ExpectedCode:          http.StatusBadRequest,
			ExpectedExistInBoards: false,
			ExpectedError:         fmt.Errorf("invalid boardID"),
			Binary:                false,
			Sensor:                true,
		},
		{
			Name:    "invalid_test_8_sensor_invalid_body_binary",
			BoardID: "esp32_1_test",
			Body: &dto.RegistrationController{
				ControllerType: TypeSensor,
				ControllerID:   "test2_invalid_board",
				Name:           tests.GetPtrStr("new name"),
				PinNumber:      tests.GetPtrInt(1),
				Type:           "new_type",
				Status:         tests.GetPtrBool(true),
			},
			IsValidBody:           false,
			ExpectedCode:          http.StatusBadRequest,
			ExpectedExistInBoards: false,
			ExpectedError:         fmt.Errorf("invalid boardID"),
			Binary:                false,
			Sensor:                true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			validFirst := tc.Body.Validate()
			validSecond := false

			if validFirst {
				if tc.Body.ControllerType == TypeSensor {
					dto := tc.Body.ToSensorController()
					if dto == nil {
						validSecond = false
					} else {
						validSecond = dto.Validate()
					}
				} else if tc.Body.ControllerType == TypeBinary {
					dto := tc.Body.ToBinaryController()
					if dto == nil {
						validSecond = false
					} else {
						validSecond = dto.Validate()
					}
				} else {
					validFirst = false
					validSecond = false
				}
			}

			if tc.IsValidBody {
				if !validFirst || !validSecond {
					t.Errorf("ERROR IN VALIDATE CHECK: got: %v, expect: %v\n", (validFirst && validSecond), tc.IsValidBody)
				}
			} else {
				if validFirst && validSecond {
					t.Errorf("ERROR IN VALIDATE CHECK: got: %v, expect: %v\n", (validFirst && validSecond), tc.IsValidBody)
				}
			}
			if tc.IsValidBody && (validFirst && validSecond) {
				data, err := tc.Body.GetJson()
				if err != nil {
					t.Errorf("ERROR IM MARHSHAL: %v\n", err)
					return
				}
				code, err := serv.Db.RegistrationController(context.Background(), data, tc.BoardID, tc.Binary, tc.Sensor, tc.Body.ControllerID)
				if err != nil {
					if tc.ExpectedError != nil {
						if err.Error() != tc.ExpectedError.Error() {
							t.Errorf("ERROR IN EXPECTED ERROR: got: %s, expect: %s\n", err.Error(), tc.ExpectedError.Error())
						}
					} else {
						t.Error(err)
					}
				}
				if code != tc.ExpectedCode {
					t.Errorf("ERROR CODE: got: %d, expect: %d\n", code, tc.ExpectedCode)
				}
			}
			if tc.ExpectedExistInBoards {
				data, _, err := serv.Db.GetControllersInfoWithType(context.Background(), tc.BoardID, tc.Body.ControllerType, tc.Body.ControllerID)
				if err != nil {
					if tc.ExpectedError != nil {
						if err.Error() != tc.ExpectedError.Error() {
							t.Errorf("ERROR IN EXPECTED ERROR: got: %s, expect: %s\n", err.Error(), tc.ExpectedError.Error())
						}
					} else {
						t.Error(err)
					}
				}
				dto := &dto.RegistrationController{}
				err = json.Unmarshal(data, dto)
				if err != nil {
					t.Error(err)

				}
				if dto.ControllerID != tc.Body.ControllerID {
					t.Errorf("ERROR IN DB CONTROLLER TYPE: got: %s, expect: %s\n", dto.ControllerID, tc.Body.ControllerID)
				}
				if tc.Body.Name != nil {
					if dto.Name != nil {
						if *dto.Name != *tc.Body.Name {
							t.Errorf("ERROR IN DB CONTROLLER NAME: got: %s, expect: %s\n", *dto.Name, *tc.Body.Name)
						}
					} else {
						t.Errorf("ERROR IN DB CONTROLLER NAME: got: nil, expect: %s\n", *tc.Body.Name)
					}
				}
				if tc.Body.PinNumber != nil {
					if dto.PinNumber != nil {
						if *dto.PinNumber != *tc.Body.PinNumber {
							t.Errorf("ERROR IN DB CONTROLLER PIN NUMBER got: %d, expect: %d\n", *dto.PinNumber, *tc.Body.PinNumber)
						}
					} else {
						t.Errorf("ERROR IN DB CONTROLLER PIN NUMBER: got: nil, expect: %d\n", *tc.Body.PinNumber)
					}
				}
				if tc.Body.Type != dto.Type {
					t.Errorf("ERROR IN DB CONTROLLER TYPE: got: %s, expect: %s\n", dto.Type, tc.Body.Type)
				}

				if tc.Body.Status != nil {
					if dto.Status != nil {
						if *dto.Status != *tc.Body.Status {
							t.Errorf("ERROR IN DB CONTROLLER STATUS: got: %v, expect: %v\n", *dto.Status, tc.Body.Status)
						}
					} else {
						t.Errorf("ERROR IN DB CONTROLLER STATUS: got: nil, expect: %v\n", *tc.Body.Status)
					}
				}
				if tc.Body.Unit != nil {
					if dto.Unit != nil {
						if *dto.Unit != *tc.Body.Unit {
							t.Errorf("ERROR IN DB CONTROLLER UNIT: got: %s, expect: %s\n", *dto.Unit, *tc.Body.Unit)
						}
					} else {
						t.Errorf("ERROR IN DB CONTROLLER UNIT: got: nil, expect: %s\n", *tc.Body.Unit)
					}
				}
				if tc.Body.Value != nil {
					if dto.Value != nil {
						if *dto.Value != *tc.Body.Value {
							t.Errorf("ERROR IN DB CONTROLLER VALUE: got: %d, expect: %d\n", *dto.Value, *tc.Body.Value)
						}
					} else {
						t.Errorf("ERROR IN DB CONTROLLER VALUE: got: nil, expect: %d\n", *tc.Body.Value)
					}
				}
			} else {
				_, _, err := serv.Db.GetControllersInfoWithType(context.Background(), tc.BoardID, tc.Body.ControllerType, tc.Body.ControllerID)
				if err == nil {
					t.Errorf("ERROR: CONTROLLER REGISTRATION")
				}
			}
		})
	}
}
