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
	"github.com/Piccadilly98/goProjects/intelectHome/src/middleware"
	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
	"github.com/Piccadilly98/goProjects/intelectHome/src/rate_limit"
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

type GlobalRateLimiterTestCase struct {
	name                       string
	endPoints                  string
	header                     string
	methods                    string
	quantityRequestsInSeconds  int
	codeInLastRequests         int
	firstRejectedIndex         int
	quantityRejected           int
	expectedCode               int
	expectedFirstRejectedIndex int
	expectedRejectedQuantity   int
}

type IpRateLimiterTestCase struct {
	name                   string
	ips                    []IpBehavior
	header                 string
	method                 string
	quantityOkRequests     int
	quantityRejected       int
	expectOkRequests       int
	expectQuantityRejected int
}

type IpBehavior struct {
	ip       string
	behavior map[string]int //endpoints->requests
}

type headerForTests struct {
	validAdminHeader    string
	validEspHeader      string
	oldAdminHeader      string
	oldEspHeader        string
	notLoginAdminHeader string
	notLoginEspHeader   string
	invalidRoleHeader   string
	invalidTokenHeader  string
}

func makeHeadersForTests(server *ServerSettings) *headerForTests {
	ht := &headerForTests{}

	ht.validAdminHeader = getAdminHeader(server)
	ht.validEspHeader = getEspHeader(server)
	oldToken, _, _ := server.tw.CreateToken("admin_login", "ADMIN", 1*time.Microsecond)
	oldTokenEsp, _, _ := server.tw.CreateToken("esp32_1_login", "ESP", 1*time.Microsecond)
	notLoginEspToken, _, _ := server.tw.CreateToken("esp_32_1_login", "ESP", 24*time.Hour)
	notLoginToken, _, _ := server.tw.CreateToken("admin_login", "ADMIN", 24*time.Hour)
	invalidRoleToken, _, _ := server.tw.CreateToken("admin1_login", "GOD", 24*time.Hour)
	ht.oldAdminHeader = method + " " + oldToken
	ht.oldEspHeader = method + " " + oldTokenEsp
	ht.notLoginEspHeader = method + " " + notLoginEspToken
	ht.notLoginAdminHeader = method + " " + notLoginToken
	ht.invalidRoleHeader = method + " " + invalidRoleToken
	ht.invalidTokenHeader = method + " " + "sdasdqweqdas.random.testsinvalid"
	return ht
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
	ipRl      *rate_limit.IpRateLimiter
	globalRl  *rate_limit.GlobalRateLimiter
}

const (
	acessToken  = "accessToken"
	key         = "Authorization"
	method      = "Bearer"
	headerKey   = "format"
	headerValue = "text"
)

func initTests(globalRL bool, ipRl bool) *ServerSettings {
	godotenv.Load("/Users/flowerma/Desktop/goProjects/intelectHome/.env")
	server := &ServerSettings{
		r:  chi.NewRouter(),
		st: storage.MakeStorage("ADMIN", "ESP32_1", "ESP32_2"),
		sm: auth.MakeSessionManager(),
		tw: &auth.TokenWorker{},
	}
	middlewareAuth := auth.MiddlewareAuth(server.st, server.sm)
	control := handlers.MakeHandlerControl(server.st)
	boardsID := handlers.MakeBoarsIDHandler(server.st)
	boards := handlers.MakeBoarsHandler(server.st)
	devices := handlers.MakeDevicesHandler(server.st)
	devicesID := handlers.MakeDevicesIDHandler(server.st)
	logs := handlers.MakeLogsHandler(server.st)
	login := auth.MakeLoginHandlers(server.st, server.sm, server.tw)
	server.mid = middlewareAuth
	server.control = control
	server.boards = boards
	server.boardsID = boardsID
	server.devices = devices
	server.devicesID = devicesID
	server.logs = logs
	server.login = login

	server.ipRl = rate_limit.MakeIpRateLimiter(10, 10)
	server.globalRl = rate_limit.MakeGlobalRateLimiter(50, 50)
	if globalRL {
		server.r.Use(middleware.GlobalRateLimiterToMiddleware(server.globalRl, server.st))
	}
	if ipRl {
		server.r.Use(middleware.IpRateLimiter(server.ipRl, server.st))
	}
	server.r.With(middlewareAuth).Route("/", func(r chi.Router) {
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
	server := initTests(false, false)
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
	server := initTests(false, false)
	adminHeader := getAdminHeader(server)
	espHeader := getEspHeader(server)
	oldToken, _, _ := server.tw.CreateToken("admin_login", "ADMIN", 1*time.Microsecond)
	oldTokenEsp, _, _ := server.tw.CreateToken("esp32_1_login", "ESP", 1*time.Microsecond)
	notLoginEspToken, _, _ := server.tw.CreateToken("esp_32_1_login", "ESP", 24*time.Hour)
	notLoginToken, _, _ := server.tw.CreateToken("admin_login", "ADMIN", 24*time.Hour)
	invalidRoleToken, _, _ := server.tw.CreateToken("admin1_login", "GOD", 24*time.Hour)
	oldHeader := method + " " + oldToken
	oldEspHeader := method + " " + oldTokenEsp
	notLoginHeaderEsp := method + " " + notLoginEspToken
	notLoginHeader := method + " " + notLoginToken
	invalidRoleHeader := method + " " + invalidRoleToken
	invalidTokenHeader := method + " " + "sdasdqweqdas.random.testsinvalid"

	testCases := []logsTestCases{
		{

			//				LOGS

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

		//				ADMIN VALID TOKEN, INVALID METHODS

		{
			name:         "ivalidTestLogs:admin|put",
			method:       http.MethodPut,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "ivalidTestLogs:admin|delete",
			method:       http.MethodDelete,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "ivalidTestLogs:admin|patch",
			method:       http.MethodPatch,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "ivalidTestLogs:admin|head",
			method:       http.MethodHead,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "ivalidTestLogs:admin|options",
			method:       http.MethodOptions,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusMethodNotAllowed,
		},

		//				ADMIN INVALID TOKEN
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
			name:         "invalidTestLogs2:admin|oldjwt|logsID=1",
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
			name:         "invalidTestLogs4:admin|oldjwt_format",
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
			name:         "invalidTestLogs5:admin|not login jwt",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs6:admin|not login jwt|jwtid",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/logs?jwtID=1",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs7:admin|not login jwt|logsid",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/logs?logsID=1",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs8:admin|not login jwt|format",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/logs",
			body:         "",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},

		//			ESP VALID TOKEN

		{
			name:         "invalidTestLogs9:esp",
			method:       http.MethodGet,
			role:         "ESP",
			header:       espHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "invalidTestLogs10:esp|jwtID",
			method:       http.MethodGet,
			role:         "ESP",
			header:       espHeader,
			url:          "/logs?jwtID=1",
			body:         "",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "invalidTestLogs11:esp|logsID",
			method:       http.MethodGet,
			role:         "ESP",
			header:       espHeader,
			url:          "/logs?logsID=1",
			body:         "",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "invalidTestLogs12:esp|format",
			method:       http.MethodGet,
			role:         "ESP",
			header:       espHeader,
			url:          "/logs",
			body:         "",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusForbidden,
		},

		// 				OTHER INVALID TOKEN
		{
			name:         "invalidTestLogs13:randomToken",
			method:       http.MethodGet,
			role:         "?",
			header:       invalidTokenHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs14:randomToken|format",
			method:       http.MethodGet,
			role:         "?",
			header:       invalidTokenHeader,
			url:          "/logs",
			body:         "",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs:invalid role token",
			method:       http.MethodGet,
			role:         "GOD",
			header:       invalidRoleHeader,
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs:invalid role token|format",
			method:       http.MethodGet,
			role:         "GOD",
			header:       invalidRoleHeader,
			url:          "/logs",
			body:         "",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs15:noheader",
			method:       http.MethodGet,
			role:         "ESP",
			url:          "/logs",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestLogs:noheader|format",
			method:       http.MethodGet,
			role:         "ESP",
			url:          "/logs",
			body:         "",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},

		//				CONTROL

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

		//					ADMIN VALID TOKEN, INVALID METHODS

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
			name:         "ivalidTestControl:admin|put",
			method:       http.MethodPut,
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

		//				ADMIN VALID TOKEN, INVALID BODY

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
			name:         "ivalidTestControl8:admin|empty body",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusBadRequest,
		},

		//				ADMIN INVALID TOKEN

		{
			name:         "ivalidTestControl9:admin|oldtoken",
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
			name:         "ivalidTestControl11:admin|oldtoken|invalid body",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/control",
			body:         `[{"id":"led1","status":"on","boardsID":"esp32_1"},{"id":"led3","status":"on","boardsID":"esp32_3"}]`,
			expectedCode: http.StatusUnauthorized,
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

		// 					OTHER INVALID TOKENS

		{
			name:         "ivalidTestControl:random token",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:invalidRoleToken",
			method:       http.MethodPost,
			role:         "GOD",
			header:       invalidRoleHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalidTestControl:noheader",
			method:       http.MethodGet,
			role:         "?",
			url:          "/control",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},

		//				ESP VALID TOKEN, NO ACCESS

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
			name:         "ivalidTestControl:esp|invalid id device",
			method:       http.MethodPost,
			role:         "ESP",
			header:       espHeader,
			url:          "/control",
			body:         `{"id":"led2","status":"on","boardsID":"esp32_1"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "ivalidTestControl:esp|body",
			method:       http.MethodPost,
			role:         "ESP",
			header:       espHeader,
			url:          "/control",
			body:         `{"id":"led1","status":"on","boardsID":"esp32_1"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "ivalidTestControl:esp|body_array",
			method:       http.MethodPost,
			role:         "ESP",
			header:       espHeader,
			url:          "/control",
			body:         `[{"id":"led1","status":"on","boardsID":"esp32_1"},{"id":"led2","status":"on","boardsID":"esp32_3"}]`,
			expectedCode: http.StatusForbidden,
		},

		//			ESP INVALID TOKEN

		{
			name:         "ivalidTestControl:espOldToken",
			method:       http.MethodPost,
			role:         "ESP",
			header:       oldEspHeader,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:espOldToken|body",
			method:       http.MethodPost,
			role:         "ESP",
			header:       oldEspHeader,
			url:          "/control",
			body:         `{"id":"led1","status":"on","boardsID":"esp32_1"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:espOldToken|body array",
			method:       http.MethodPost,
			role:         "ESP",
			header:       oldEspHeader,
			url:          "/control",
			body:         `[{"id":"led1","status":"on","boardsID":"esp32_1"},{"id":"led2","status":"on","boardsID":"esp32_1"}]`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:espOldToken|invalid body",
			method:       http.MethodPost,
			role:         "ESP",
			header:       oldEspHeader,
			url:          "/control",
			body:         `{"id":"led2","status":"on","boardsID":"esp32_1"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:espOldToken|invalid body array",
			method:       http.MethodPost,
			role:         "ESP",
			header:       oldEspHeader,
			url:          "/control",
			body:         `[{"id":"led3","status":"on","boardsID":"esp32_1"},{"id":"led2","status":"on","boardsID":"esp32_1"}]`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:espNotLogin|",
			method:       http.MethodPost,
			role:         "ESP",
			header:       notLoginHeaderEsp,
			url:          "/control",
			body:         "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:espNotLogin|body",
			method:       http.MethodPost,
			role:         "ESP",
			header:       notLoginHeaderEsp,
			url:          "/control",
			body:         `{"id":"led1","status":"on","boardsID":"esp32_1"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:espNotLogin|body array",
			method:       http.MethodPost,
			role:         "ESP",
			header:       notLoginHeaderEsp,
			url:          "/control",
			body:         `[{"id":"led1","status":"on","boardsID":"esp32_1"},{"id":"led2","status":"on","boardsID":"esp32_1"}]`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:espNotLogin|invalid body",
			method:       http.MethodPost,
			role:         "ESP",
			header:       notLoginHeaderEsp,
			url:          "/control",
			body:         `{"id":"led2","status":"on","boardsID":"esp32_1"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "ivalidTestControl:espNotLogin|invalid body array",
			method:       http.MethodPost,
			role:         "ESP",
			header:       notLoginHeaderEsp,
			url:          "/control",
			body:         `[{"id":"led3","status":"on","boardsID":"esp32_1"},{"id":"led2","status":"on","boardsID":"esp32_1"}]`,
			expectedCode: http.StatusUnauthorized,
		},

		//			BOARDS

		{
			name:         "validTestBoards1:admin",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards",
			expectedCode: http.StatusOK,
		},
		{
			name:         "validTestBoards2:admin|format",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusOK,
		},

		//			VALID ADMIN TOKEN, INVALID MATHODS

		{
			name:         "inValidTestBoards:admin|post",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoards:admin|post|format",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoards:admin|put",
			method:       http.MethodPut,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoards:admin|put|format",
			method:       http.MethodPut,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoards:admin|delete",
			method:       http.MethodDelete,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoards:admin|delete|format",
			method:       http.MethodDelete,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoards:admin|patch",
			method:       http.MethodPatch,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoards:admin|patch|format",
			method:       http.MethodPatch,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoards:admin|head",
			method:       http.MethodHead,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoards:admin|head|format",
			method:       http.MethodHead,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusMethodNotAllowed,
		},

		//		INVALID ADMIN TOKEN

		{
			name:         "inValidTestBoards:admin|old",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/boards",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoards:admin|old|format",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoards:admin|not login",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/boards",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoards:admin|not login|format",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},

		//			OTHER INVALID TOKEN

		{
			name:         "inValidTestBoards:random|",
			method:       http.MethodGet,
			role:         "?",
			header:       invalidTokenHeader,
			url:          "/boards",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoards:random|format",
			method:       http.MethodGet,
			role:         "?",
			header:       invalidTokenHeader,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoards:invalidRole",
			method:       http.MethodGet,
			role:         "GOD",
			header:       invalidRoleHeader,
			url:          "/boards",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoards:invalid role|format",
			method:       http.MethodGet,
			role:         "GOD",
			header:       invalidRoleHeader,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoards:no header",
			method:       http.MethodGet,
			role:         "?",
			header:       "",
			url:          "/boards",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoards:no header|format",
			method:       http.MethodGet,
			role:         "?",
			header:       "",
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},

		//		VALID ESP TOKEN, NO ACCESS

		{
			name:         "inValidTestBoards:esp",
			method:       http.MethodGet,
			role:         "ESP",
			header:       espHeader,
			url:          "/boards",
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "inValidTestBoards:esp|format",
			method:       http.MethodGet,
			role:         "ESP",
			header:       espHeader,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusForbidden,
		},

		//		INVALID ESP TOKEN

		{
			name:         "inValidTestBoards:esp|old",
			method:       http.MethodGet,
			role:         "ESP",
			header:       oldEspHeader,
			url:          "/boards",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoards:esp|old|format",
			method:       http.MethodGet,
			role:         "ESP",
			header:       oldEspHeader,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoards:esp|not login",
			method:       http.MethodGet,
			role:         "ESP",
			header:       notLoginHeaderEsp,
			url:          "/boards",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoards:esp|not login|format",
			method:       http.MethodGet,
			role:         "ESP",
			header:       notLoginHeaderEsp,
			url:          "/boards",
			headerKey:    headerKey,
			headerValue:  headerValue,
			expectedCode: http.StatusUnauthorized,
		},

		//			BOARDS{ID}
		//		   	   GET

		{
			name:         "inValidTestBoardsID:admin",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "inValidTestBoardsID:esp",
			method:       http.MethodGet,
			role:         "ESP32_1",
			header:       espHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusOK,
		},

		//			 POST

		{
			name:         "inValidTestBoardsID:admin",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards/esp32_1",
			body:         `{"boardId": "esp32_1","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "inValidTestBoardsID:esp",
			method:       http.MethodGet,
			role:         "ESP32_1",
			header:       espHeader,
			url:          "/boards/esp32_1",
			body:         `{"boardId": "esp32_1","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			expectedCode: http.StatusOK,
		},

		//			IVALID METHODS
		//			 ADMIN TOKEN

		{
			name:         "inValidTestBoardsID1:admin|put",
			method:       http.MethodPut,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoardsID:admin|patch",
			method:       http.MethodPatch,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoardsID:admin|delete",
			method:       http.MethodDelete,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoardsID:admin|options",
			method:       http.MethodOptions,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoardsID:admin|head",
			method:       http.MethodHead,
			role:         "ADMIN",
			header:       adminHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusMethodNotAllowed,
		},

		//			ESP TOKEN

		{
			name:         "inValidTestBoardsID:esp|put",
			method:       http.MethodPut,
			role:         "ESP",
			header:       espHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoardsID:esp|patch",
			method:       http.MethodPatch,
			role:         "ESP",
			header:       espHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoardsID:esp|delete",
			method:       http.MethodDelete,
			role:         "ESP",
			header:       espHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoardsID:esp|options",
			method:       http.MethodOptions,
			role:         "ESP",
			header:       espHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusMethodNotAllowed,
		},
		{
			name:         "inValidTestBoardsID:esp|head",
			method:       http.MethodHead,
			role:         "ESP",
			header:       espHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusMethodNotAllowed,
		},

		//			INVALID ADMIN TOKEN
		//				     GET

		{
			name:         "inValidTestBoardsID:admin|old",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoardsID:admin|not login",
			method:       http.MethodGet,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},

		//				POST: VALID BODY

		{
			name:         "inValidTestBoardsID:admin|post|old",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/boards/esp32_1",
			body:         `{"boardId": "esp32_1","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoardsID:admin|post|not login",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/boards/esp32_1",
			body:         `{"boardId": "esp32_1","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			expectedCode: http.StatusUnauthorized,
		},

		//				POST: INVALID BODY

		{
			name:         "inValidTestBoardsID:admin|post|old|invalid body",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       oldHeader,
			url:          "/boards/esp32_1",
			body:         `{"boardId": "esp32_2","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoardsID:admin|post|not login|invalid body",
			method:       http.MethodPost,
			role:         "ADMIN",
			header:       notLoginHeader,
			url:          "/boards/esp32_1",
			body:         `{"boardId": "esp32_2","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			expectedCode: http.StatusUnauthorized,
		},

		//				IVALID ESP TOKEN
		//					GET

		{
			name:         "inValidTestBoardsID:esp|old",
			method:       http.MethodGet,
			role:         "ESP",
			header:       oldEspHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoardsID:esp|not login",
			method:       http.MethodGet,
			role:         "ESP",
			header:       notLoginHeaderEsp,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},

		//					POST: VALID BODY

		{
			name:         "inValidTestBoardsID:esp|old|",
			method:       http.MethodGet,
			role:         "ESP",
			header:       oldEspHeader,
			body:         `{"boardId": "esp32_1","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoardsID:esp|not login",
			method:       http.MethodGet,
			role:         "ESP",
			header:       notLoginHeaderEsp,
			body:         `{"boardId": "esp32_1","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},

		//			POST: INVALID BODY

		{
			name:         "inValidTestBoardsID:esp|old|invalid body",
			method:       http.MethodGet,
			role:         "ESP",
			header:       oldEspHeader,
			body:         `{"boardId": "esp32_2","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoardsID:esp|not login|invalid body",
			method:       http.MethodGet,
			role:         "ESP",
			header:       notLoginHeaderEsp,
			body:         `{"boardId": "esp32_2","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},

		//			OTHER INVALID TOKEN
		//					GET

		{
			name:         "inValidTestBoardsID:ivalid role",
			method:       http.MethodGet,
			role:         "GOD",
			header:       invalidRoleHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoardsID:random token",
			method:       http.MethodGet,
			role:         "?",
			header:       invalidTokenHeader,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoardsID:no headers",
			method:       http.MethodGet,
			role:         "?",
			header:       "",
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},

		//			POST:VALID BODY

		{
			name:         "inValidTestBoardsID:ivalid role",
			method:       http.MethodPost,
			role:         "GOD",
			header:       invalidRoleHeader,
			url:          "/boards/esp32_1",
			body:         `{"boardId": "esp32_1","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoardsID:random token",
			method:       http.MethodPost,
			role:         "?",
			header:       invalidTokenHeader,
			body:         `{"boardId": "esp32_1","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoardsID:no headers",
			method:       http.MethodPost,
			role:         "?",
			header:       "",
			body:         `{"boardId": "esp32_1","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},

		// POST:INVALID BODY

		{
			name:         "inValidTestBoardsID:ivalid role|invalid body",
			method:       http.MethodPost,
			role:         "GOD",
			header:       invalidRoleHeader,
			url:          "/boards/esp32_1",
			body:         `{"boardId": "esp32_3","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoardsID:random token|invalid body",
			method:       http.MethodPost,
			role:         "?",
			header:       invalidTokenHeader,
			body:         `{"boardId": "esp32_3","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "inValidTestBoardsID:no headers|invalid body",
			method:       http.MethodPost,
			role:         "?",
			header:       "",
			body:         `{"boardId": "esp32_3","tempCP": 0,"freeMemory": 0,"workTime": 0,"rssi": 0,"localIP": "","networkIP": "","voltage": 0,"quantityDevice": 0,"TimeUpload": "0001-01-01T00:00:00Z","TimeAdded": "2025-10-18T19:40:05.125446+04:00"}`,
			url:          "/boards/esp32_1",
			expectedCode: http.StatusUnauthorized,
		},
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

func TestGlobalRateLimiter(t *testing.T) {
	server := initTests(true, false)
	ht := makeHeadersForTests(server)
	time.Sleep(1 * time.Second)

	testCases := []GlobalRateLimiterTestCase{
		{
			name:                       "validTestNoHeader",
			endPoints:                  "/login",
			methods:                    http.MethodPost,
			quantityRequestsInSeconds:  40,
			expectedCode:               http.StatusBadRequest,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "validTestBoundary",
			endPoints:                  "/login",
			methods:                    http.MethodPost,
			quantityRequestsInSeconds:  50,
			expectedCode:               http.StatusBadRequest,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "invalidTestExceeded",
			endPoints:                  "/login",
			methods:                    http.MethodPost,
			quantityRequestsInSeconds:  51,
			expectedCode:               http.StatusTooManyRequests,
			expectedFirstRejectedIndex: 50,
			expectedRejectedQuantity:   1,
		},
		{
			name:                       "invalidTestOverExceeded",
			endPoints:                  "/login",
			methods:                    http.MethodPost,
			quantityRequestsInSeconds:  100,
			expectedCode:               http.StatusTooManyRequests,
			expectedFirstRejectedIndex: 50,
			expectedRejectedQuantity:   50,
		},
		{
			name:                       "invalidTestSuperOverExceeded",
			endPoints:                  "/login",
			methods:                    http.MethodPost,
			quantityRequestsInSeconds:  1000,
			expectedCode:               http.StatusTooManyRequests,
			expectedFirstRejectedIndex: 50,
			expectedRejectedQuantity:   950,
		},
		{
			name:                       "validTestBoundaryBoardsWithHeader",
			endPoints:                  "/boards",
			methods:                    http.MethodGet,
			header:                     ht.validAdminHeader,
			quantityRequestsInSeconds:  50,
			expectedCode:               http.StatusOK,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "validTestBoundaryBoards",
			endPoints:                  "/boards",
			methods:                    http.MethodPost,
			quantityRequestsInSeconds:  50,
			expectedCode:               http.StatusUnauthorized,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "invalidTestExceededBoards",
			endPoints:                  "/boards",
			methods:                    http.MethodPost,
			quantityRequestsInSeconds:  51,
			expectedCode:               http.StatusTooManyRequests,
			expectedFirstRejectedIndex: 50,
			expectedRejectedQuantity:   1,
		},
		{
			name:                       "validTestBoundaryDevicesWithHeaders",
			endPoints:                  "/devices",
			methods:                    http.MethodGet,
			header:                     ht.validAdminHeader,
			quantityRequestsInSeconds:  50,
			expectedCode:               http.StatusOK,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "validTestBoundaryDevices",
			endPoints:                  "/devices",
			methods:                    http.MethodGet,
			quantityRequestsInSeconds:  50,
			expectedCode:               http.StatusUnauthorized,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "invalidTestExceededDevices",
			endPoints:                  "/devices",
			methods:                    http.MethodGet,
			quantityRequestsInSeconds:  51,
			expectedCode:               http.StatusTooManyRequests,
			expectedFirstRejectedIndex: 50,
			expectedRejectedQuantity:   1,
		},
		{
			name:                       "validTestBoundaryLogsWithHeaders",
			endPoints:                  "/logs",
			methods:                    http.MethodGet,
			quantityRequestsInSeconds:  50,
			header:                     ht.validAdminHeader,
			expectedCode:               http.StatusOK,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "validTestBoundaryLogs",
			endPoints:                  "/logs",
			methods:                    http.MethodGet,
			quantityRequestsInSeconds:  50,
			expectedCode:               http.StatusUnauthorized,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "invalidTestExceededLogs",
			endPoints:                  "/logs",
			methods:                    http.MethodGet,
			quantityRequestsInSeconds:  51,
			expectedCode:               http.StatusTooManyRequests,
			expectedFirstRejectedIndex: 50,
			expectedRejectedQuantity:   1,
		},
		{
			name:                       "validTestBoundaryControl",
			endPoints:                  "/control",
			methods:                    http.MethodPost,
			quantityRequestsInSeconds:  50,
			expectedCode:               http.StatusUnauthorized,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "invalidTestExceededControl",
			endPoints:                  "/control",
			methods:                    http.MethodPost,
			quantityRequestsInSeconds:  51,
			expectedCode:               http.StatusTooManyRequests,
			expectedFirstRejectedIndex: 50,
			expectedRejectedQuantity:   1,
		},
		{
			name:                       "validTestBoundaryDevicesIdWithHeaders",
			endPoints:                  "/devices/led1",
			methods:                    http.MethodGet,
			header:                     ht.validAdminHeader,
			quantityRequestsInSeconds:  50,
			expectedCode:               http.StatusOK,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "validTestBoundaryDevicesId",
			endPoints:                  "/devices/led1",
			methods:                    http.MethodGet,
			quantityRequestsInSeconds:  50,
			expectedCode:               http.StatusUnauthorized,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "invalidTestExceededDevicesId",
			endPoints:                  "/devices/led1",
			methods:                    http.MethodGet,
			quantityRequestsInSeconds:  51,
			expectedCode:               http.StatusTooManyRequests,
			expectedFirstRejectedIndex: 50,
			expectedRejectedQuantity:   1,
		},
		{
			name:                       "validTestBoundaryBoardsIDWithHeader",
			endPoints:                  "/boards/esp32_1",
			methods:                    http.MethodGet,
			header:                     ht.validAdminHeader,
			quantityRequestsInSeconds:  50,
			expectedCode:               http.StatusOK,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "validTestBoundaryBoardsID",
			endPoints:                  "/boards/esp32_1",
			methods:                    http.MethodGet,
			quantityRequestsInSeconds:  50,
			expectedCode:               http.StatusUnauthorized,
			expectedFirstRejectedIndex: -1,
			expectedRejectedQuantity:   0,
		},
		{
			name:                       "invalidTestExceededBoardsID",
			endPoints:                  "/boards/esp32_1",
			methods:                    http.MethodGet,
			quantityRequestsInSeconds:  51,
			expectedCode:               http.StatusTooManyRequests,
			expectedFirstRejectedIndex: 50,
			expectedRejectedQuantity:   1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.firstRejectedIndex == 0 {
				tc.firstRejectedIndex = -1
			}
			req := httptest.NewRequest(tc.methods, tc.endPoints, nil)
			req.Header.Set(key, tc.header)

			timeStart := time.Now().UnixMicro()

			for i := 0; i < tc.quantityRequestsInSeconds; i++ {
				// t.Error(tc.firstRejectedIndex)
				w := httptest.NewRecorder()
				server.r.ServeHTTP(w, req)
				code := w.Code
				//  t.Error(code)
				if code == http.StatusTooManyRequests {
					if tc.firstRejectedIndex == -1 {
						tc.firstRejectedIndex = i
					}
					tc.quantityRejected++
				}
				if i+1 == tc.quantityRequestsInSeconds {
					tc.codeInLastRequests = code
				}
			}

			if tc.expectedCode != tc.codeInLastRequests {
				t.Errorf("got codeInLstRequest: %d, expect: %d\n", tc.codeInLastRequests, tc.expectedCode)
			}
			if tc.firstRejectedIndex != tc.expectedFirstRejectedIndex {
				t.Errorf("got fistRejected index: %d, expect: %d\n", tc.firstRejectedIndex, tc.expectedFirstRejectedIndex)
			}
			if tc.expectedRejectedQuantity != tc.quantityRejected {
				t.Errorf("got Rejected quantity: %d, expect: %d\n", tc.quantityRejected, tc.expectedRejectedQuantity)
			}

			timeAfter := time.Now().UnixMicro() - timeStart

			remainder := 1*time.Second.Microseconds() - timeAfter

			if remainder > 0 {
				time.Sleep(time.Duration(remainder) * time.Microsecond)
			}
		})

	}
}

func TestIpRateLimited(t *testing.T) {
	server := initTests(false, true)
	ht := makeHeadersForTests(server)
	time.Sleep(1 * time.Second)

	testCases := []IpRateLimiterTestCase{
		{
			name: "valid_one_IP_valid_quantity",
			ips: []IpBehavior{
				{
					ip: "190.168.0.1:12345",
					behavior: map[string]int{
						"/logs":    5,
						"/devices": 4,
					},
				},
			},
			header:                 ht.validAdminHeader,
			method:                 http.MethodGet,
			expectOkRequests:       9,
			expectQuantityRejected: 0,
		},
		{
			name: "valid_one_IP_boundary_quantity",
			ips: []IpBehavior{
				{
					ip: "190.168.0.1:12345",
					behavior: map[string]int{
						"/logs":    3,
						"/devices": 4,
						"/boards":  3,
					},
				},
			},
			header:                 ht.validAdminHeader,
			method:                 http.MethodGet,
			expectOkRequests:       10,
			expectQuantityRejected: 0,
		},
		{
			name: "inValid_one_IP_exceeded_quantity",
			ips: []IpBehavior{
				{
					ip: "190.168.0.1:12345",
					behavior: map[string]int{
						"/logs":           8,
						"/devices":        1,
						"/boards":         1,
						"/boards/esp32_1": 1,
						"/devices/led1":   1,
					},
				},
			},
			header:                 ht.validAdminHeader,
			method:                 http.MethodGet,
			expectOkRequests:       10,
			expectQuantityRejected: 2,
		},
		{
			name: "inValid_one_IP_extra_exceeded_quantity",
			ips: []IpBehavior{
				{
					ip: "190.168.0.1:12345",
					behavior: map[string]int{
						"/logs":           1000,
						"/devices":        1000,
						"/boards":         1000,
						"/boards/esp32_1": 1000,
						"/login":          1000,
					},
				},
			},
			header:                 ht.validAdminHeader,
			method:                 http.MethodGet,
			expectOkRequests:       10,
			expectQuantityRejected: 4990,
		},
		{
			name: "valid_two_IP_valid_quantity",
			ips: []IpBehavior{
				{
					ip: "190.168.0.1:12345",
					behavior: map[string]int{
						"/logs":    5,
						"/devices": 4,
					},
				},
				{
					ip: "192.168.0.1:12345",
					behavior: map[string]int{
						"/devices":        2,
						"/boards":         3,
						"/boards/esp32_1": 3,
						"/logs":           1,
					},
				},
			},
			header:                 ht.validAdminHeader,
			method:                 http.MethodGet,
			expectOkRequests:       18,
			expectQuantityRejected: 0,
		},
		{
			name: "valid_two_IP_boundary_quantity",
			ips: []IpBehavior{
				{
					ip: "190.168.0.1:12345",
					behavior: map[string]int{
						"/logs":          1,
						"/devices":       3,
						"/devices/led1":  2,
						"/logs?jwtID=1":  2,
						"/logs?logsID=1": 2,
					},
				},
				{
					ip: "192.168.0.1:12345",
					behavior: map[string]int{
						"/devices":        2,
						"/boards":         3,
						"/boards/esp32_1": 3,
						"/logs":           1,
						"/devices/led1":   1,
					},
				},
			},
			header:                 ht.validAdminHeader,
			method:                 http.MethodGet,
			expectOkRequests:       20,
			expectQuantityRejected: 0,
		},
		{
			name: "invalid_two_IP_exceeded_in_first_quantity",
			ips: []IpBehavior{
				{
					ip: "190.168.0.1:12345",
					behavior: map[string]int{
						"/logs":          2,
						"/devices":       3,
						"/devices/led1":  2,
						"/logs?jwtID=1":  2,
						"/logs?logsID=1": 2,
					},
				},
				{
					ip: "192.168.0.1:12345",
					behavior: map[string]int{
						"/devices":        2,
						"/boards":         3,
						"/boards/esp32_1": 3,
						"/logs":           1,
						"/devices/led1":   1,
					},
				},
			},
			header:                 ht.validAdminHeader,
			method:                 http.MethodGet,
			expectOkRequests:       20,
			expectQuantityRejected: 1,
		},
		{
			name: "invalid_two_IP_exceeded_in_second_quantity",
			ips: []IpBehavior{
				{
					ip: "190.168.0.1:12345",
					behavior: map[string]int{
						"/logs":          1,
						"/devices":       3,
						"/devices/led1":  2,
						"/logs?jwtID=1":  2,
						"/logs?logsID=1": 2,
					},
				},
				{
					ip: "192.168.0.1:12345",
					behavior: map[string]int{
						"/devices":        2,
						"/boards":         3,
						"/boards/esp32_1": 5,
						"/logs":           1,
						"/devices/led1":   1,
					},
				},
			},
			header:                 ht.validAdminHeader,
			method:                 http.MethodGet,
			expectOkRequests:       20,
			expectQuantityRejected: 2,
		},
		{
			name: "invalid_two_IP_exceeded_extra_exceeded_in_first_quantity",
			ips: []IpBehavior{
				{
					ip: "190.168.0.1:12345",
					behavior: map[string]int{
						"/logs":          10,
						"/devices":       10,
						"/devices/led1":  10,
						"/logs?jwtID=1":  10,
						"/logs?logsID=1": 10,
					},
				},
				{
					ip: "192.168.0.1:12345",
					behavior: map[string]int{
						"/devices":        2,
						"/boards":         3,
						"/boards/esp32_1": 3,
						"/logs":           1,
						"/devices/led1":   1,
					},
				},
			},
			header:                 ht.validAdminHeader,
			method:                 http.MethodGet,
			expectOkRequests:       20,
			expectQuantityRejected: 40,
		},
		{
			name: "invalid_two_IP_extra_exceeded_quantity",
			ips: []IpBehavior{
				{
					ip: "190.168.0.1:12345",
					behavior: map[string]int{
						"/logs":          10,
						"/devices":       25,
						"/devices/led1":  30,
						"/logs?jwtID=1":  55,
						"/logs?logsID=1": 120,
					},
				},
				{
					ip: "192.168.0.1:12345",
					behavior: map[string]int{
						"/devices": 120,
						"/boards":  120,
						"/logs":    50,
					},
				},
			},
			header:                 ht.validAdminHeader,
			method:                 http.MethodGet,
			expectOkRequests:       20,
			expectQuantityRejected: 510,
		},
		{
			name: "valid_three_IP_valid_quantity",
			ips: []IpBehavior{
				{
					ip: "190.168.0.1:12345",
					behavior: map[string]int{
						"/logs":          2,
						"/devices":       2,
						"/devices/led1":  2,
						"/logs?jwtID=1":  2,
						"/logs?logsID=1": 2,
					},
				},
				{
					ip: "192.168.0.1:12345",
					behavior: map[string]int{
						"/devices": 10,
					},
				},
				{
					ip: "[:::1]:1234",
					behavior: map[string]int{
						"/boards/esp32_1": 10,
					},
				},
			},
			header:                 ht.validAdminHeader,
			method:                 http.MethodGet,
			expectOkRequests:       30,
			expectQuantityRejected: 0,
		},
		{
			name: "invalid_three_IP_exceeded_quantity",
			ips: []IpBehavior{
				{
					ip: "190.168.0.1:12345",
					behavior: map[string]int{
						"/logs":          2,
						"/devices":       3,
						"/devices/led1":  2,
						"/logs?jwtID=1":  2,
						"/logs?logsID=1": 2,
					},
				},
				{
					ip: "192.168.0.1:12345",
					behavior: map[string]int{
						"/devices":        10,
						"/logs":           2,
						"/boards/esp32_1": 2,
					},
				},
				{
					ip: "[:::1]:1234",
					behavior: map[string]int{
						"/boards/esp32_1": 10,
						"/boards":         5,
					},
				},
			},
			header:                 ht.validAdminHeader,
			method:                 http.MethodGet,
			expectOkRequests:       30,
			expectQuantityRejected: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, v := range tc.ips {
				for k, v2 := range v.behavior {
					req := httptest.NewRequest(tc.method, k, nil)
					req.RemoteAddr = v.ip
					req.Header.Set(key, tc.header)
					for i := 0; i < v2; i++ {
						w := httptest.NewRecorder()
						server.r.ServeHTTP(w, req)
						code := w.Code
						if code == http.StatusOK {
							tc.quantityOkRequests++
						}
						if code == http.StatusTooManyRequests {
							tc.quantityRejected++
						}
					}
				}
			}

			timeEnd := time.Now().UnixMicro()

			if tc.quantityOkRequests != tc.expectOkRequests {
				t.Errorf("got quantity ok requests: %d, expect: %d\n", tc.quantityOkRequests, tc.expectOkRequests)
			}
			if tc.quantityRejected != tc.expectQuantityRejected {
				t.Errorf("got quantity rejected: %d, expect: %d\n", tc.quantityRejected, tc.expectQuantityRejected)
			}
			timeAfter := time.Now().UnixMicro() - timeEnd

			remainder := 1*time.Second.Microseconds() - timeAfter

			if remainder > 0 {
				time.Sleep(time.Duration(remainder) * time.Microsecond)
			}

			// time.Sleep(1 * time.Second)
		})
	}

}
