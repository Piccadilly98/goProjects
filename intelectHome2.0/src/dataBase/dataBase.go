package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	dto "github.com/Piccadilly98/goProjects/intelectHome2.0/src/DTO"
)

type DataBase struct {
	host      string
	port      string
	username  string
	nameDb    string
	password  string
	db        *sql.DB
	errChan   chan error
	mtx       sync.Mutex
	isConnect atomic.Bool
}

func MakeDataBase(host string, port string, username string, nameDb string, password string) (*DataBase, error) {
	db, err := initDataBase(host, port, username, nameDb, password)
	if db == nil {
		return nil, err
	}
	res := &DataBase{
		host:     host,
		port:     port,
		username: username,
		nameDb:   nameDb,
		password: password,
		db:       db,
		errChan:  make(chan error, 10),
	}
	res.isConnect.Store(true)
	return res, nil
}

func (db *DataBase) processingError(err error) (int, error) {
	if errors.Is(err, context.Canceled) {
		return 0, err
	}
	if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "many clients") ||
		strings.Contains(err.Error(), "bad connection") {
		select {
		case db.errChan <- err:
		default:
		}
		return http.StatusServiceUnavailable, fmt.Errorf("fail connection db")
	}
	if errors.Is(err, sql.ErrConnDone) {
		select {
		case db.errChan <- err:
		default:
		}
		return http.StatusServiceUnavailable, fmt.Errorf("fail connection db")
	}
	return http.StatusInternalServerError, err
}

func initDataBase(host string, port string, username string, nameDb string, password string) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, nameDb))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}

func (db *DataBase) Recover() bool {
	db.mtx.Lock()
	if !db.isConnect.Load() {
		db.isConnect.Store(false)
	}
	defer db.mtx.Unlock()
	if err := db.db.Ping(); err != nil {
		for i := range 5 {
			err := db.db.Ping()
			if err == nil {
				db.isConnect.Store(true)
				return true
			}
			time.Sleep(time.Duration(i+1) * time.Second)
		}
		for range 3 {
			db.db, err = initDataBase(db.host, db.port, db.username, db.nameDb, db.password)
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
			err = db.db.Ping()
			if err == nil {
				db.isConnect.Store(true)
				return true
			}
			time.Sleep(5 * time.Second)
		}
	} else {
		db.isConnect.Store(true)
		return true
	}
	return false
}

func (db *DataBase) GetExistWithId(ctx context.Context, id string) (bool, int, error) {
	var exist bool
	if !db.isConnect.Load() {
		return false, http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}
	err := db.GetPointer().QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM boards WHERE board_id = $1);`, id).Scan(&exist)
	if err != nil {
		code, err := db.processingError(err)
		return false, code, err
	}
	return exist, http.StatusOK, nil
}

func (db *DataBase) GetDtoWithId(ctx context.Context, id string) (*dto.GetBoardDataDto, int, error) {
	if !db.isConnect.Load() {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("fail connection db")
	}
	dto := dto.GetBoardDataDto{}
	err := db.GetPointer().QueryRowContext(ctx, `
	SELECT board_id, name, type, board_state, created_date, updated FROM boards
	WHERE board_id = $1;
	`, id).Scan(&dto.BoardId, &dto.Name, &dto.BoardType, &dto.BoardState, &dto.CreatedDate, &dto.UpdatedDate)
	if err != nil {
		code, err := db.processingError(err)
		return nil, code, err
	}
	return &dto, http.StatusOK, nil
}

func (db *DataBase) GetInfoDtoWithId(ctx context.Context, id string) (*dto.GetBoardInfoDTO, int, error) {
	dto := &dto.GetBoardInfoDTO{}
	if !db.isConnect.Load() {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("fail connection db")
	}
	err := db.GetPointer().QueryRowContext(ctx, `
	SELECT board_id, created_date, updated_date, cpu_temp, avalible_ram, rssi_wifi, total_runtime, ip_address, voltage, total_device, mac_address FROM boardInfo
	WHERE board_id = $1;`, id).Scan(&dto.BoardId, &dto.CreatedDate,
		&dto.UpdatedDate, &dto.CpuTemp, &dto.AvalibleRam,
		&dto.RssiWifi, &dto.TotalRunTime, &dto.IpAddress,
		&dto.Voltage, &dto.TotalDeviceCount, &dto.MacAddress)

	if err != nil {
		code, err := db.processingError(err)
		return nil, code, err
	}

	return dto, http.StatusOK, nil
}

func (db *DataBase) RegistrationBoard(ctx context.Context, id *string, name *string, boardType *string, state *string) (int, error) {
	if !db.isConnect.Load() {
		return http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}
	_, err := db.GetPointer().ExecContext(ctx, `
	INSERT INTO boards(board_id, name, type, board_state)
	VALUES($1, $2, $3, $4);
	`, id, name, boardType, state)
	if err != nil {
		code, err := db.processingError(err)
		return code, err
	}

	if !db.isConnect.Load() {
		return http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}
	_, err = db.GetPointer().ExecContext(ctx, `
	INSERT INTO boardInfo(board_id)
	VALUES($1)`, id)
	if err != nil {
		code, err := db.processingError(err)
		return code, err
	}
	return http.StatusCreated, nil
}

func (db *DataBase) UpdateBoard(ctx context.Context, boardID string, dto *dto.UpdateBoardDto) (int, error) {
	if !db.isConnect.Load() {
		return http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}
	var sets []string
	var args []any

	argNum := 1

	if dto.Name != nil {
		sets = append(sets, fmt.Sprintf("name = $%d", argNum))
		args = append(args, *dto.Name)
		argNum++
	}
	if dto.Type != nil {
		sets = append(sets, fmt.Sprintf("type = $%d", argNum))
		args = append(args, *dto.Type)
		argNum++
	}
	if dto.State != nil {
		sets = append(sets, fmt.Sprintf("board_state = $%d", argNum))
		args = append(args, *dto.State)
		argNum++
	}

	sets = append(sets, "updated = NOW()")

	query := "UPDATE boards SET " + strings.Join(sets, ", ") +
		fmt.Sprintf(" WHERE board_id = $%d", argNum)
	args = append(args, boardID)

	res, err := db.GetPointer().ExecContext(ctx, query, args...)
	if err != nil {
		code, err := db.processingError(err)
		return code, err
	}
	if res != nil {
		if count, _ := res.RowsAffected(); count == 0 {
			return http.StatusInternalServerError, fmt.Errorf("no rows in db, invalid board_id")
		}
	}
	return http.StatusOK, nil
}

func (db *DataBase) UpdateBoardInfo(ctx context.Context, id string, data *dto.UpdateBoardInfo) (int, error) {
	if !db.isConnect.Load() {
		return http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}
	_, err := db.GetPointer().ExecContext(ctx, `
	UPDATE boardInfo
	SET cpu_temp = $1, avalible_ram = $2, rssi_wifi = $3, total_runtime = $4, ip_address = $5, voltage = $6, mac_address = $7, total_device = $8, updated_date = $9
	WHERE board_id = $10
	`, *data.CpuTemp, *data.AvalibleRam, *data.RssiWifi,
		*data.TotalRunTime, *data.IpAddress, *data.Voltage,
		*data.MacAddress, *data.TotalDeviceCount, data.TimeUpload,
		id)

	if err != nil {
		code, err := db.processingError(err)
		return code, err
	}
	return http.StatusOK, nil
}

func (db *DataBase) Close() {
	db.db.Close()
}

func (db *DataBase) GetPointer() *sql.DB {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	return db.db
}

func (db *DataBase) ErrChan() chan error {
	return db.errChan
}

func (db *DataBase) IsConnect() bool {
	return db.isConnect.Load()
}
