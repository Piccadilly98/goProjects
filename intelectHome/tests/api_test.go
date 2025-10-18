package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Piccadilly98/goProjects/intelectHome/src/auth"
	"github.com/Piccadilly98/goProjects/intelectHome/src/handlers"
	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type testCases struct {
	name                string
	method              string
	login               string
	password            string
	headerKey           string
	headerValue         string
	expectedCode        int
	accessTokenContains bool
}

type logsTestCases struct {
	name         string
	method       string
	url          string
	role         string
	headerKey    string
	headerValue  string
	header       string
	expectedCode int
	body         string
}

type ServerSettings struct {
	r         *chi.Mux
	st        *storage.Storage
	sm        *auth.SessionManager
	tw        *auth.TokenWorker
	mid       any
	control   any
	boardsID  any
	boards    any
	devices   any
	devicesID any
	logs      any
	login     any
}

const (
	acessToken  = "accessToken"
	key         = "Authorization"
	method      = "Bearer"
	headerKey   = "format"
	headerValue = "txt"
)

func initTests() *ServerSettings {
	godotenv.Load("/Users/flowerma/Desktop/goProjects/intelectHome/.env")
	server := &ServerSettings{
		r:  chi.NewRouter(),
		st: storage.MakeStorage("ADMIN", "ESP32_1", "ESP32_2"),
		sm: auth.MakeSessionManager(),
		tw: &auth.TokenWorker{},
	}
	middleware := auth.MiddlewareAuth(server.st, server.sm)
	control := handlers.MakeHandlerControl(server.st)
	boardsID := handlers.MakeBoarsIDHandler(server.st)
	boards := handlers.MakeBoarsHandler(server.st)
	devices := handlers.MakeDevicesHandler(server.st)
	devicesID := handlers.MakeDevicesIDHandler(server.st)
	logs := handlers.MakeLogsHandler(server.st)
	login := auth.MakeLoginHandlers(server.st, server.sm, server.tw)
	server.mid = middleware
	server.control = control
	server.boards = boards
	server.boardsID = boardsID
	server.devices = devices
	server.devicesID = devicesID
	server.logs = logs
	server.login = login

	server.r.With(middleware).Route("/", func(r chi.Router) {
		r.Post("/control", control.Control)
		r.HandleFunc("/boards/{boardID}", boardsID.BoardsIDHandler)
		r.Get("/boards", boards.BoardsHandler)
		r.HandleFunc("/devices", devices.DevicesHandler)
		r.Get("/devices/{deviceID}", devicesID.DevicesIDHandler)
		r.Get("/logs", logs.LogsHandler)
		r.Post("/login", login.LoginHandler)
	})
	return server
}

func getEspHeader(server *ServerSettings) string {
	resp := &models.TokenResponseJSON{}
	reader := strings.NewReader(fmt.Sprintf(`{"login":"%s","password":"%s"}`, os.Getenv("ESP32_1_LOGIN"), os.Getenv("ESP32_1_PASSWORD")))
	req := httptest.NewRequest(http.MethodPost, "/login", reader)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.r.ServeHTTP(w, req)
	b, err := io.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, resp)
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(resp.Role, "ESP") {
		return fmt.Sprintf("%s %s", method, resp.AccessToken)
	}
	return ""
}

func getAdminHeader(server *ServerSettings) string {
	resp := &models.TokenResponseJSON{}
	reader := strings.NewReader(fmt.Sprintf(`{"login":"%s","password":"%s"}`, os.Getenv("ADMIN_LOGIN"), os.Getenv("ADMIN_PASSWORD")))
	req := httptest.NewRequest(http.MethodPost, "/login", reader)
	w := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")
	server.r.ServeHTTP(w, req)
	b, err := io.ReadAll(w.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, resp)
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasPrefix(resp.Role, "ADMIN") {

		return fmt.Sprintf("%s %s", method, resp.AccessToken)
	}
	return ""
}

func TestLoginHandler(t *testing.T) {
	server := initTests()
	testTable := []testCases{
		{
			name:                "validTest1",
			method:              http.MethodPost,
			login:               os.Getenv("ESP32_1_LOGIN"),
			password:            os.Getenv("ESP32_1_PASSWORD"),
			headerKey:           "Content-Type",
			headerValue:         "application/json",
			expectedCode:        200,
			accessTokenContains: true,
		},
		{
			name:                "inValidTest:testRepeatTest1",
			method:              http.MethodPost,
			login:               os.Getenv("ESP32_1_LOGIN"),
			password:            os.Getenv("ESP32_1_PASSWORD"),
			headerKey:           "Content-Type",
			headerValue:         "application/json",
			expectedCode:        http.StatusConflict,
			accessTokenContains: false,
		},
		{
			name:                "validTest2",
			method:              http.MethodPost,
			login:               os.Getenv("ADMIN_LOGIN"),
			password:            os.Getenv("ADMIN_PASSWORD"),
			headerKey:           "Content-Type",
			headerValue:         "application/json",
			expectedCode:        http.StatusOK,
			accessTokenContains: true,
		},
		{
			name:                "inValidTest:testRepeatTest2",
			method:              http.MethodPost,
			login:               os.Getenv("ADMIN_LOGIN"),
			password:            os.Getenv("ADMIN_PASSWORD"),
			headerKey:           "Content-Type",
			headerValue:         "application/json",
			expectedCode:        http.StatusConflict,
			accessTokenContains: false,
		},
		{
			name:                "50/50validTest",
			method:              http.MethodPost,
			login:               os.Getenv("ESP32_2_LOGIN"),
			password:            os.Getenv("ESP32_2_PASSWORD"),
			headerKey:           "Content-Type",
			headerValue:         "text",
			expectedCode:        http.StatusBadRequest,
			accessTokenContains: false,
		},
		{
			name:                "invalidTest1:Get",
			method:              http.MethodGet,
			expectedCode:        http.StatusMethodNotAllowed,
			accessTokenContains: false,
		},
		{
			name:                "invalidTest2:Delete",
			method:              http.MethodDelete,
			expectedCode:        http.StatusMethodNotAllowed,
			accessTokenContains: false,
		},
		{
			name:                "invalidTest3:Connect",
			method:              http.MethodConnect,
			expectedCode:        http.StatusMethodNotAllowed,
			accessTokenContains: false,
		},
		{
			name:                "invalidTest4:Put",
			method:              http.MethodPut,
			expectedCode:        http.StatusMethodNotAllowed,
			accessTokenContains: false,
		},
		{
			name:                "invalidTest5:Patch",
			method:              http.MethodPatch,
			expectedCode:        http.StatusMethodNotAllowed,
			accessTokenContains: false,
		},
		{
			name:                "invalidTest6:invalid Data",
			method:              http.MethodPost,
			login:               "123",
			password:            "sdsa",
			headerKey:           "Content-Type",
			headerValue:         "application/json",
			expectedCode:        http.StatusUnauthorized,
			accessTokenContains: false,
		},
		{
			name:                "invalidTest7:invalid Data",
			method:              http.MethodPost,
			login:               "153",
			password:            "sdsa",
			headerKey:           "Content-Type",
			headerValue:         "application/json",
			expectedCode:        http.StatusUnauthorized,
			accessTokenContains: false,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(fmt.Sprintf(`{"login":"%s","password":"%s"}`, tc.login, tc.password))
			req := httptest.NewRequest(tc.method, "/login", reader)
			req.Header.Set(tc.headerKey, tc.headerValue)
			w := httptest.NewRecorder()
			server.r.ServeHTTP(w, req)
			if tc.name == "invalidTest0:invalid header value" {
				t.Log(w.Body.String())
			}
			if w.Code != tc.expectedCode {
				t.Errorf("HTTPCODE:\ngot %d, want %d\n", w.Code, tc.expectedCode)
			}
			if strings.Contains(w.Body.String(), acessToken) != tc.accessTokenContains {
				t.Errorf("CONTAINS ACCESSTOKEN:\ngot %v, want %v\n", !tc.accessTokenContains, tc.accessTokenContains)
			}
		})
	}
}

func TestMiddleWare(t *testing.T) {
	server := initTests()
	adminHeader := getAdminHeader(server)
	espHeader := getEspHeader(server)
	oldToken, _, _ := server.tw.CreateToken("admin_login", "ADMIN", 1*time.Microsecond)
	notLoginToken, _, _ := server.tw.CreateToken("admin_login", "ADMIN", 24*time.Hour)
	invalidRoleToken, _, _ := server.tw.CreateToken("admin1_login", "GOD", 24*time.Hour)
	oldHeader := method + " " + oldToken
	notLoginHeader := method + " " + notLoginToken
	invalidRoleHeader := method + " " + invalidRoleToken

	testCases := []logsTestCases{
		{
			name:         "validTestLogs1: admin",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "validTestLogs2:admin|?logsID",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/logs?logsID=1",
			body:         "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "validTestLogs3:admin|?jwt",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/logs?jwtID=1",
			body:         "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "validTestLogs4:admin|format",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/logs",
			body:         "",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusOK,
		},
		{
			name:         "invalidTestLogs1:admin|oldjwt",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs2:admin|oldjwt",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/logs?jwtID=1",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs3:admin|oldjwt",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/logs?logsID=1",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs3:admin|oldjwt_format",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/logs?logsID=1",
			body:         "",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs4:admin|not login jwt",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs5:admin|not login jwt|jwtid",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/logs?jwtID=1",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs6:admin|not login jwt|logsid",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/logs?logsID=1",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs7:admin|not login jwt|format",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/logs",
			body:         "",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs8:esp",
			method:       http.MethodGet,
			role:         "ESP",
			header:       espHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "invalidTestLogs9:esp|jwtID",
			method:       http.MethodGet,
			role:         "ESP",
			header:       espHeader,
			url:          "/logs?jwtID=1",
			body:         "",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "invalidTestLogs10:esp|logsID",
			method:       http.MethodGet,
			role:         "ESP",
			header:       espHeader,
			url:          "/logs?logsID=1",
			body:         "",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "invalidTestLogs11:esp|format",
			method:       http.MethodGet,
			role:         "ESP",
			header:       espHeader,
			url:          "/logs",
			body:         "",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "validTestControl1:admin",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         `{"id":"led1","status":"on","boardsID":"esp32_1"}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "validTestControl2:admin",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         `[{"id":"led1","status":"on","boardsID":"esp32_1"},{"id":"led2","status":"on","boardsID":"esp32_1"}]`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "ivalidTestControl1:admin|method",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         `[{"id":"led1","status":"on","boardsID":"esp32_1"},{"id":"led2","status":"on","boardsID":"esp32_1"}]`,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "ivalidTestControl2:admin|invalid id device",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         `{"id":"led3","status":"on","boardsID":"esp32_1"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "ivalidTestControl3:admin|invalid id device in array begin",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         `[{"id":"led3","status":"on","boardsID":"esp32_1"},{"id":"led2","status":"on","boardsID":"esp32_1"}]`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "ivalidTestControl4:admin|invalid id device in array end",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         `[{"id":"led1","status":"on","boardsID":"esp32_1"},{"id":"led3","status":"on","boardsID":"esp32_1"}]`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "ivalidTestControl5:admin|invalid id boards",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         `{"id":"led3","status":"on","boardsID":"esp32_4"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "ivalidTestControl6:admin|invalid id boards in begin array",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         `[{"id":"led1","status":"on","boardsID":"esp32_3"},{"id":"led3","status":"on","boardsID":"esp32_1"}]`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "ivalidTestControl7:admin|invalid id boards in end array",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         `[{"id":"led1","status":"on","boardsID":"esp32_1"},{"id":"led3","status":"on","boardsID":"esp32_3"}]`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "ivalidTestControl8:admin|empty body", //поправить что бы пустой json не приниался
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "ivalidTestControl9:admin|oldtoken", //поправить на неплавильный токен возвращает 400
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl10:admin|oldtoken|body",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/control",
			body:         `[{"id":"led1","status":"on","boardsID":"esp32_1"},{"id":"led2","status":"on","boardsID":"esp32_3"}]`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl10:admin|oldtoken|invalid body",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/control",
			body:         `[{"id":"led1","status":"on","boardsID":"esp32_1"},{"id":"led3","status":"on","boardsID":"esp32_3"}]`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:admin|get",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "ivalidTestControl:admin|delete",
			method:       http.MethodDelete,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "ivalidTestControl:admin|put",
			method:       http.MethodPut,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "ivalidTestControl:admin|patch",
			method:       http.MethodPatch,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "ivalidTestControl:admin|head",
			method:       http.MethodHead,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "ivalidTestControl:admin|options",
			method:       http.MethodOptions,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "ivalidTestControl:not login jwt",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:esp",
			method:       http.MethodPost,
			role:         "ESP",
			header:       espHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "ivalidTestControl:esp|body",
			method:       http.MethodPost,
			role:         "GOD",
			header:       invalidRoleHeader,
			url:          "/control",
			body:         `[{"id":"led1","status":"on","boardsID":"esp32_1"},{"id":"led2","status":"on","boardsID":"esp32_3"}]`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:invalid role",
			method:       http.MethodPost,
			role:         "GOD",
			header:       invalidRoleHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},

		// {
		// 	name:         "invalidTestCase2",
		// 	method:       http.MethodPost,
		// 	role:         "ESP",
		// 	expectedCode: http.StatusForbidden,
		// },
		// {
		// 	name:         "validTestCase2",
		// 	method:       http.MethodGet,
		// 	role:         "ADMIN",
		// 	header:       getEspHeader(server),
		// 	url:          "/logs?logsID=1",
		// 	expectedCode: http.StatusOK,
		// },
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var reader io.Reader = nil
			if tc.body != "" {
				reader = strings.NewReader(tc.body)
			}
			req := httptest.NewRequest(tc.method, tc.url, reader)
			req.Header.Set(key, tc.header)
			if tc.headerKey != "" {
				req.Header.Set(tc.headerKey, tc.headerValue)
			}
			w := httptest.NewRecorder()

			server.r.ServeHTTP(w, req)
			if w.Code != tc.expectedCode {
				t.Errorf("got %d, expect: %d\nbody: %s", w.Code, tc.expectedCode, w.Body.String())
			}
		})

	}
}
