package storage

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Piccadilly98/goProjects/intelectHome/src/logs"
	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
)

const (
	StatusOFF = "off"
	StatusON  = "on"
)

type Storage struct {
	logs       *logs.Logs
	boardData  map[string]*models.DataBoard
	deviceData map[string]*models.Device_data
	mtx        sync.Mutex
}

func MakeStorage() *Storage {
	storage := &Storage{logs: logs.MakeNewLogsInfo(),
		boardData:  make(map[string]*models.DataBoard),
		deviceData: make(map[string]*models.Device_data),
		mtx:        sync.Mutex{}}
	storage.boardData["esp32_1"] = &models.DataBoard{BoardId: "esp32_1"}
	storage.deviceData["led1"] = &models.Device_data{ID: "led1", Status: "off"}
	return storage
}

func (s *Storage) UpdateStatusDevice(id, status string) bool {
	if status != StatusOFF && status != StatusON {
		return false
	}
	s.mtx.Lock()
	_, ok := s.deviceData[id]
	if ok {
		s.deviceData[id] = &models.Device_data{ID: id, Status: status}
		s.mtx.Unlock()
		return true
	}
	s.mtx.Unlock()
	return false
}

func (s *Storage) PrintLogs() {
	s.mtx.Lock()
	l := s.logs.String()
	s.mtx.Unlock()
	fmt.Println(l)
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

func (s *Storage) AddNewDeviceId(id string) {
	s.mtx.Lock()
	s.deviceData[id] = &models.Device_data{ID: id, Status: StatusOFF}
	s.mtx.Unlock()
}

func (s *Storage) AddNewBoard(id string) {
	s.mtx.Lock()
	s.boardData[id] = &models.DataBoard{BoardId: id}
	s.mtx.Unlock()
}

func (s *Storage) GetAllDevicesStatusJson() []*models.JSONResponse {
	responses := make([]*models.JSONResponse, 0)
	s.mtx.Lock()
	for _, v := range s.deviceData {
		responses = append(responses, &models.JSONResponse{ID: v.ID, Status: v.Status})
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

func (s *Storage) GetDeviceStatus(id string) string {
	s.mtx.Lock()
	v, ok := s.deviceData[id]
	s.mtx.Unlock()
	if ok {
		return v.Status
	}
	return ""
}

func (s *Storage) GetBoardInfo(id string) string {
	s.mtx.Lock()
	v, ok := s.boardData[id]
	s.mtx.Unlock()
	if ok {
		return v.String()
	}
	return ""
}

func (s *Storage) NewLogPost(r *http.Request, ID, status string, httpCode int) {
	s.mtx.Lock()
	s.logs.CreateAndAddRecordPost(r, ID, status, httpCode)
	s.mtx.Unlock()
}

func (s *Storage) NewLogGet(r *http.Request, body []byte, httpCode int) {
	s.mtx.Lock()
	s.logs.CreateAndAddRecordGet(r, body, httpCode)
	s.mtx.Unlock()
}

func (s *Storage) CheckIdDevice(id string) bool {
	s.mtx.Lock()
	_, ok := s.deviceData[id]
	s.mtx.Unlock()
	return ok
}
