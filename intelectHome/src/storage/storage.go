package storage

import (
	"fmt"
	"net/http"
	"sync"
	"time"

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
	str := s.logs.String()
	s.mtx.Unlock()
	return str
}
func (s *Storage) NewLog(r *http.Request, body []byte, httpCode int, errors string, a ...any) {
	s.mtx.Lock()
	s.logs.CreateAndAddRecord(r, body, httpCode, errors)
	s.mtx.Unlock()
}

func (s *Storage) CheckIdDevice(id string) bool {
	s.mtx.Lock()
	_, ok := s.deviceData[id]
	s.mtx.Unlock()
	return ok
}
