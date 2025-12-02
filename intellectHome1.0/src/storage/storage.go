package storage

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
)

const (
	StatusOFF = "off"
	StatusON  = "on"
)

type Storage struct {
	logs       map[int]string
	logsLength int
	boardData  map[string]*models.DataBoard
	deviceData map[string]*models.Device_data
	rolesList  []string
	mtx        sync.Mutex
}

func MakeStorage(roles ...string) *Storage {
	storage := &Storage{logs: make(map[int]string),
		logsLength: 0,
		boardData:  make(map[string]*models.DataBoard),
		deviceData: make(map[string]*models.Device_data),
		mtx:        sync.Mutex{}}
	storage.rolesList = append(storage.rolesList, roles...)
	storage.AddNewBoard("esp32_1")
	storage.AddNewDeviceId("led1", "esp32_1")
	storage.AddNewDeviceId("led2", "esp32_1")
	return storage
}

func (s *Storage) UpdateStatusDevice(id, status string) bool {
	if status != StatusOFF && status != StatusON {
		return false
	}
	s.mtx.Lock()
	v, ok := s.deviceData[id]
	if ok {
		s.deviceData[id] = &models.Device_data{ID: id, Status: status, BoadrId: v.BoadrId}
		s.mtx.Unlock()
		return true
	}
	s.mtx.Unlock()
	return false
}

func (s *Storage) UpdateBoardDevice(id string, board string) {
	s.deviceData[id].BoadrId = board
}

func (s *Storage) PrintLogs() {
	s.mtx.Lock()
	str := ""
	keys := make([]int, 0)
	for k := range s.logs {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, v := range keys {
		str += fmt.Sprintf("%d: %s\n", v, s.logs[v])
	}
	fmt.Println(str)
	s.mtx.Unlock()
}

func (s *Storage) PrintDataBoards() {
	s.mtx.Lock()
	db := ""
	for k, v := range s.boardData {
		db += fmt.Sprintf("BoardId: %s\nInfo: %s", k, v)
	}
	s.mtx.Unlock()
	fmt.Println(db)
}

func (s *Storage) PrintDataDevice() {
	s.mtx.Lock()
	dd := ""
	for k, v := range s.deviceData {
		dd += fmt.Sprintf("DeviceId: %s -- Status: %s", k, v)
	}
	s.mtx.Unlock()
	fmt.Println(dd)
}

func (s *Storage) AddNewDeviceId(id string, boadrID string) {
	s.mtx.Lock()
	s.deviceData[id] = &models.Device_data{ID: id, Status: StatusOFF, BoadrId: boadrID}
	s.mtx.Unlock()
}

func (s *Storage) AddNewBoard(id string) {
	s.mtx.Lock()
	s.boardData[id] = &models.DataBoard{BoardId: id, TimeAdded: time.Now()}
	s.mtx.Unlock()
}

func (s *Storage) GetAllDevicesStatusJson() []*models.JSONResponse {
	responses := make([]*models.JSONResponse, 0)
	s.mtx.Lock()
	for _, v := range s.deviceData {
		responses = append(responses, &models.JSONResponse{ID: v.ID, Status: v.Status, BoadrId: v.BoadrId})
	}
	s.mtx.Unlock()
	return responses
}

func (s *Storage) GetAllDevicesStatusString() string {
	s.mtx.Lock()
	str := ""
	for _, v := range s.deviceData {
		str += v.String()
	}
	s.mtx.Unlock()
	return str
}

func (s *Storage) AddNewBoardInfo(db *models.DataBoard) bool {
	s.mtx.Lock()
	_, ok := s.boardData[db.BoardId]
	if !ok {
		s.mtx.Unlock()
		return false
	}
	db.TimeAdded = s.boardData[db.BoardId].TimeAdded
	s.boardData[db.BoardId] = db
	s.boardData[db.BoardId].TimeUpload = time.Now()
	s.mtx.Unlock()
	return true
}

func (s *Storage) GetDeviceInfo(id string) (models.Device_data, error) {
	s.mtx.Lock()
	v, ok := s.deviceData[id]
	s.mtx.Unlock()
	if ok {
		return *v, nil
	}
	return models.Device_data{}, fmt.Errorf("invalid device id")
}

func (s *Storage) GetDeviceStatus(id string) string {
	s.mtx.Lock()
	v, ok := s.deviceData[id]
	s.mtx.Unlock()
	if ok {
		return v.Status
	}
	return ""
}

func (s *Storage) GetBoardInfo(id string) (models.DataBoard, error) {
	s.mtx.Lock()
	v, ok := s.boardData[id]
	s.mtx.Unlock()
	if ok {
		return *v, nil
	}
	return models.DataBoard{}, fmt.Errorf("invalid boardID: %s", id)
}

func (s *Storage) GetAllBoardsInfo() []models.DataBoard {
	res := make([]models.DataBoard, 0)

	s.mtx.Lock()
	for _, v := range s.boardData {
		res = append(res, *v)
	}
	s.mtx.Unlock()
	return res
}

func (s *Storage) GetAllLogs() string {
	s.mtx.Lock()
	str := ""
	keys := make([]int, 0)
	for k := range s.logs {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, v := range keys {
		str += fmt.Sprintf("%d: %s\n", v, s.logs[v])
	}
	s.mtx.Unlock()
	return str
}
func (s *Storage) NewLog(r *http.Request, claims *models.ClaimsJSON, httpCode int, errors string, attentions ...string) {
	userID := "not contains"
	var jwtID int64 = 0
	role := "not contains"
	if claims != nil {
		login, err := claims.GetSubject()
		if err != nil {
			userID = claims.Role
		}
		userID = login
		jwtID = claims.TokenID
		role = claims.Role
	}
	attentionsStr := ""
	for _, v := range attentions {
		attentionsStr += v
		attentionsStr += "	"
	}
	s.mtx.Lock()
	s.logsLength++
	if s.logsLength > 200 {
		for s.logsLength >= 0 {
			delete(s.logs, s.logsLength)
			s.logsLength--
		}
		s.logsLength = 1
	}
	s.logs[s.logsLength] = fmt.Sprintf("Time %v -- IP: %s -- URL: %s -- Method: %s -- UserID: %s -- TokenID: %v -- Role: %s --  Code: %d -- Errors: %s -- Attentions: %s",
		time.Now().String(),
		r.RemoteAddr,
		r.URL.Path,
		r.Method,
		userID,
		jwtID,
		role,
		httpCode,
		errors,
		attentionsStr)
	s.mtx.Unlock()
}

func (s *Storage) GetLog(ID int) (string, error) {
	s.mtx.Lock()
	v, ok := s.logs[ID]
	s.mtx.Unlock()
	if ok {
		return v, nil
	}
	return "", fmt.Errorf("invalid id: %d", ID)
}

func (s *Storage) GetAllLogsJSON() []models.LogsJSON {
	res := make([]models.LogsJSON, 0)
	s.mtx.Lock()
	keys := make([]int, 0)
	for k := range s.logs {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		res = append(res, models.LogsJSON{LogsID: k, LogsInfo: s.logs[k]})
	}
	s.mtx.Unlock()
	return res
}

func (s *Storage) GetLogJson(ID int) models.LogsJSON {
	res, err := s.GetLog(ID)
	if err != nil {
		return models.LogsJSON{}
	}
	logJSON := models.LogsJSON{LogsID: ID, LogsInfo: res}
	return logJSON
}

func (s *Storage) GetLogsJWTIDJSON(id string) []models.LogsJSON {
	logs := make([]models.LogsJSON, 0)
	targetStr := fmt.Sprintf("TokenID: %s", id)
	s.mtx.Lock()
	for k, v := range s.logs {
		if strings.Contains(v, targetStr) {
			logs = append(logs, models.LogsJSON{LogsID: k, LogsInfo: v})
		}
	}
	s.mtx.Unlock()
	return logs
}

func (s *Storage) GetLogsJWTIDString(id string) string {
	targetStr := fmt.Sprintf("TokenID: %s", id)
	str := ""
	s.mtx.Lock()
	for k, v := range s.logs {
		if strings.Contains(v, targetStr) {
			str += fmt.Sprintf("%d: %s", k, v)
		}
	}
	s.mtx.Unlock()
	return str
}
func (s *Storage) GetMaxIDLogs() int {
	s.mtx.Lock()
	i := s.logsLength
	s.mtx.Unlock()
	return i
}
func (s *Storage) CheckIdDevice(id string) bool {
	s.mtx.Lock()
	_, ok := s.deviceData[id]
	s.mtx.Unlock()
	return ok
}

func (s *Storage) GetAllRoles() []string {
	s.mtx.Lock()
	roles := make([]string, len(s.rolesList))
	n := copy(roles, s.rolesList)
	s.mtx.Unlock()
	if n != len(s.rolesList) {
		return nil
	}
	return roles
}
