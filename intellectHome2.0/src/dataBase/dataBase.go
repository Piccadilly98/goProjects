package database

// сделать в бд чек статусов: registred or active or offline
// сделать выдачу id, что бы это не было кастомным полем
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

	dto "github.com/Piccadilly98/goProjects/intellectHome2.0/src/DTO"
)

const (
	StatusRegistred = "registred"
	StatusActive    = "active"
	StatusLost      = "lost"
	StatusNotActive = "offline"

	ControllerColumnName         = "name"
	ControllerColumnType         = "type"
	ControllerColumnCreated      = "created_date"
	ControllerColumnUpdated      = "updated_date"
	ControllerColumnPinNumber    = "pin_number"
	BinaryControllerColumnStatus = "status"

	SensorControllerColumnUnit  = "unit"
	SensorControllerColumnValue = "value"
)

// сделать кэш в in memory мапе с board_id
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
		errChan:  make(chan error),
	}
	res.isConnect.Store(true)
	return res, nil
}

// вынести в отдельную сущность error worker
// по возможности убрать что бы бд возвращал код, нужно возвращать ошибку и какую нибудь false/true
// cancel context
// решить проблему с ошибками: неверный запрос даёт нам всегда 500, это не правильно

func (db *DataBase) processingError(err error) (int, error) {
	if strings.Contains(err.Error(), "canceling statement due to user request") ||
		strings.Contains(err.Error(), "context deadline exceeded") || errors.Is(err, context.Canceled) {
		return 0, context.Canceled
	}
	if strings.Contains(err.Error(), "duplicate") {
		return http.StatusBadRequest, err
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
		return nil, err
	}
	return db, nil
}

func (db *DataBase) Recover() bool {
	if !db.isConnect.Load() {
		db.isConnect.Store(false)
	}
	if err := db.db.Ping(); err != nil {
		for i := range 5 {
			err := db.db.Ping()
			if err == nil {
				db.isConnect.Store(true)
				return true
			}
			log.Printf("Connection re-try %d/5\n", i+1)
			time.Sleep(time.Duration(i+(i*3)) * time.Second)
		}
		db.mtx.Lock()
		defer db.mtx.Unlock()
		for i := range 3 {
			db.db, err = initDataBase(db.host, db.port, db.username, db.nameDb, db.password)
			if err != nil {
				log.Printf("Init try %d/3\n", i+1)
				time.Sleep(10 * time.Second)
				continue
			}
			time.Sleep(2 * time.Second)
			err = db.db.Ping()
			if err == nil {
				db.isConnect.Store(true)
				return true
			}
			log.Printf("DB Init try %d/3\n", i+1)
			time.Sleep(10 * time.Second)
		}
	} else {
		db.isConnect.Store(true)
		return true
	}
	return false
}

func (db *DataBase) GetExistWithBoardId(ctx context.Context, id string) (bool, int, error) {
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

func (db *DataBase) GetExistWithControllerId(ctx context.Context, deviceId string) (bool, int, error) {
	var exist bool
	if !db.IsConnect() {
		return false, http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}
	err := db.GetPointer().QueryRowContext(ctx, `
	SELECT EXISTS (
  	SELECT 1
  	FROM boards,
       jsonb_array_elements(controllers->'devices'->'binary')  AS elem
 	 WHERE elem->>'controller_id' = $1

  	UNION ALL

 	 SELECT 1
  	FROM boards,
       jsonb_array_elements(controllers->'devices'->'sensor')  AS elem
  	WHERE elem->>'controller_id' = $1
	);`, deviceId).Scan(&exist)
	if err != nil {
		code, err := db.processingError(err)
		return false, code, err
	}

	return exist, http.StatusOK, nil
}

func (db *DataBase) RegistrationController(ctx context.Context, json []byte, boardID string, binary, sensor bool, controllerID string) (int, error) {
	if !db.IsConnect() {
		return http.StatusServiceUnavailable, fmt.Errorf("fail connection db")
	}
	exist, code, err := db.GetExistWithBoardId(ctx, boardID)
	if err != nil {
		return code, err
	}
	if !exist {
		return http.StatusBadRequest, fmt.Errorf("invalid boardID")
	}
	exist, code, err = db.GetExistWithControllerId(ctx, controllerID)
	if err != nil {
		return code, err
	}
	if exist {
		return http.StatusBadRequest, fmt.Errorf("invalid controllerID")
	}

	query := ""
	if binary {
		query = `UPDATE boards
	SET controllers = jsonb_set(
	controllers,
	'{devices, binary}',
	COALESCE(controllers->'devices'->'binary', '[]') || $1::jsonb
	),updated_date = now()
	WHERE board_id = $2;`
	} else if sensor {
		query = `UPDATE boards
	SET controllers = jsonb_set(
	controllers,
	'{devices, sensor}',
	COALESCE(controllers->'devices'->'sensor', '[]') || $1::jsonb
	),updated_date = now()
	WHERE board_id = $2;`
	}
	_, err = db.GetPointer().ExecContext(ctx, query, json, boardID)
	if err != nil {
		code, err := db.processingError(err)
		return code, err
	}
	return http.StatusCreated, nil
}

// нейминг говна кусок
func (db *DataBase) GetDtoWithId(ctx context.Context, id string) (*dto.GetBoardDataDto, int, error) {
	if !db.isConnect.Load() {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("fail connection db")
	}
	dto := dto.GetBoardDataDto{}
	err := db.GetPointer().QueryRowContext(ctx, `
	SELECT board_id, name, type, board_state, created_date, updated_date FROM boards
	WHERE board_id = $1;
	`, id).Scan(&dto.BoardId, &dto.Name, &dto.BoardType, &dto.BoardState, &dto.CreatedDate, &dto.UpdatedDate)
	if err != nil {
		code, err := db.processingError(err)
		return nil, code, err
	}
	return &dto, http.StatusOK, nil
}

// нейминг говна кусок
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
		if err == sql.ErrNoRows {
			return nil, http.StatusBadRequest, fmt.Errorf("invalid boardID")
		}
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

// постараться разделить метод на 2
func (db *DataBase) UpdateBoard(ctx context.Context, boardID string, dto *dto.UpdateBoardDataDto) (int, error) {
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

	sets = append(sets, "updated_date = NOW()")

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

func (db *DataBase) GetAllBoardsWithConditions(ctx context.Context, state string, boardId string, boardType string, name string) ([]dto.GetBoardDataDto, int, error) {
	if !db.isConnect.Load() {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}
	var sets []string
	var args []any
	var quantityArgs int
	if boardId != "" {
		quantityArgs++
		str := ""
		if quantityArgs == 1 {
			str = fmt.Sprintf("WHERE board_id = $%d", quantityArgs)
		} else {
			str = fmt.Sprintf("board_id = $%d", quantityArgs)
		}
		sets = append(sets, str)
		args = append(args, boardId)
	}
	if state != "" {
		quantityArgs++
		str := ""
		if quantityArgs == 1 {
			str = fmt.Sprintf("WHERE board_state=$%d", quantityArgs)
		} else {
			str = fmt.Sprintf("board_state=$%d", quantityArgs)
		}
		sets = append(sets, str)
		args = append(args, state)
	}
	if boardType != "" {
		quantityArgs++
		str := ""
		if quantityArgs == 1 {
			str = fmt.Sprintf("WHERE type=$%d", quantityArgs)
		} else {
			str = fmt.Sprintf("type=$%d", quantityArgs)
		}
		sets = append(sets, str)
		args = append(args, boardType)
	}
	if name != "" {
		quantityArgs++
		str := ""
		if quantityArgs == 1 {
			str = fmt.Sprintf("WHERE name=$%d", quantityArgs)
		} else {
			str = fmt.Sprintf("name=$%d", quantityArgs)
		}
		sets = append(sets, str)
		args = append(args, name)
	}
	qery := "SELECT board_id, name, type, board_state, created_date, updated_date FROM boards " + strings.Join(sets, " AND ")
	if !db.isConnect.Load() {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}

	res := []dto.GetBoardDataDto{}

	rows, err := db.GetPointer().QueryContext(ctx, qery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, http.StatusOK, nil
		}
		code, err := db.processingError(err)
		return nil, code, err
	}

	for rows.Next() {
		dto := &dto.GetBoardDataDto{}

		err := rows.Scan(&dto.BoardId, &dto.Name, &dto.BoardType, &dto.BoardState, &dto.CreatedDate, &dto.UpdatedDate)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		res = append(res, *dto)
	}
	return res, http.StatusOK, nil
}

// чуть понятнее нейминг
func (db *DataBase) GetControllersByte(ctx context.Context, id string) ([]byte, int, error) {
	if !db.IsConnect() {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}
	res := []byte{}
	err := db.GetPointer().QueryRowContext(ctx, `
	SELECT controllers->'devices' FROM boards
	WHERE board_id = $1`, id).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusBadRequest, fmt.Errorf("invalid controllerID")
		}
		code, err := db.processingError(err)
		return nil, code, err
	}
	return res, http.StatusOK, nil
}

// CONTROLLER UPDATE
func (db *DataBase) GetJSONBuilderArgs(boardID string, dto *dto.ControllerUpdateDTO, controllerID string) ([]string, []any, int) {
	var sets []string
	var args []any

	argNum := 2
	args = append(args, controllerID)
	if dto.Name != nil {
		sets = append(sets, fmt.Sprintf("\n\t'%s', $%d::text", ControllerColumnName, argNum))
		args = append(args, *dto.Name)
		argNum++
	}
	if dto.PinNumber != nil {
		sets = append(sets, fmt.Sprintf("\n\t'%s', $%d::integer", ControllerColumnPinNumber, argNum))
		args = append(args, *dto.PinNumber)
		argNum++
	}
	if dto.Type != nil {
		sets = append(sets, fmt.Sprintf("\n\t'%s', $%d::text", ControllerColumnType, argNum))
		args = append(args, *dto.Type)
		argNum++
	}
	if dto.Status != nil {
		sets = append(sets, fmt.Sprintf("\n\t'%s', $%d::boolean", BinaryControllerColumnStatus, argNum))
		args = append(args, *dto.Status)
		argNum++
	}
	if dto.Value != nil {
		sets = append(sets, fmt.Sprintf("\n\t'%s', $%d::integer", SensorControllerColumnValue, argNum))
		args = append(args, *dto.Value)
		argNum++
	}
	if dto.Unit != nil {
		sets = append(sets, fmt.Sprintf("\n\t'%s', $%d::text", SensorControllerColumnUnit, argNum))
		args = append(args, *dto.Unit)
		argNum++
	}
	sets = append(sets, fmt.Sprintf("\n\t'%s', to_jsonb(NOW())", ControllerColumnUpdated))
	args = append(args, boardID)
	return sets, args, argNum
}

func (db *DataBase) GetQueryToUpdateConroller(sets []string, args []any, argNum int, controllerType string) string {
	beginStr := fmt.Sprintf(`
		UPDATE boards b
		SET controllers = jsonb_set(
		controllers,
		'{devices,%s}',
		 coalesce((
	   select jsonb_agg(
	   case
	   	when elem->>'controller_id' = $1
		then elem||jsonb_build_object(`, controllerType)

	endStr := fmt.Sprintf(`
				)
			else elem
	   		end)
	   		from jsonb_array_elements(b.controllers->'devices'->'%s') elem
	   		), '[]'::jsonb)
		)
	WHERE board_id = $%d;`, controllerType, argNum)
	query := beginStr + strings.Join(sets, ",") + endStr
	return query

}

func (db *DataBase) UpdateControllerData(ctx context.Context, boardID string, dto *dto.ControllerUpdateDTO, controllerType string, controllerID string) ([]byte, int, error) {
	if !db.IsConnect() {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}

	sets, args, argNum := db.GetJSONBuilderArgs(boardID, dto, controllerID)
	query := db.GetQueryToUpdateConroller(sets, args, argNum, controllerType)
	if !db.IsConnect() {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}
	_, err := db.GetPointer().ExecContext(ctx, query, args...)
	if err != nil {
		code, err := db.processingError(err)
		return nil, code, err
	}
	return db.GetControllersInfoWithType(ctx, boardID, controllerType, controllerID)
}

// зачем нам здесь боард айди
// переделать
func (db *DataBase) GetControllersInfoWithType(ctx context.Context, boardID string, controllerType string, controllerID string) ([]byte, int, error) {
	if !db.IsConnect() {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}
	var b []byte
	err := db.GetPointer().QueryRowContext(ctx, `
	SELECT elem FROM boards, jsonb_array_elements(controllers->'devices'->$1) elem
	WHERE board_id = $2 AND elem->>'controller_id' = $3`, controllerType, boardID, controllerID).Scan(&b)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusBadRequest, fmt.Errorf("invalid controllerID")
		}
		code, err := db.processingError(err)
		return nil, code, err
	}
	return b, http.StatusOK, nil
}

// закинуть в метод выше
func (db *DataBase) GetControllerTypeAndBoardID(ctx context.Context, controllerID string) (string, string, int, error) {
	var controllerType *string
	var boardID *string
	if !db.IsConnect() {
		return "", "", http.StatusServiceUnavailable, fmt.Errorf("data base not ready")
	}

	err := db.GetPointer().QueryRowContext(ctx, `
	SELECT board_id, 'binary' AS device_type FROM boards, jsonb_array_elements(controllers->'devices'->'binary') elem
	WHERE elem->>'controller_id' = $1

	UNION ALL

	SELECT board_id, 'sensor' AS device_type FROM boards, jsonb_array_elements(controllers->'devices'->'sensor') elem
	WHERE elem->>'controller_id' = $1;
	`, controllerID).Scan(&boardID, &controllerType)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", http.StatusBadRequest, fmt.Errorf("invalid body")
		}
		code, err := db.processingError(err)
		return "", "", code, err
	}
	return *controllerType, *boardID, http.StatusOK, nil
}

func (db *DataBase) Close() error {
	err := db.db.Close()
	return err
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
