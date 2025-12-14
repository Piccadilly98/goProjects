package status_worker

import (
	"context"
	"fmt"
	"log"
	"time"

	database "github.com/Piccadilly98/goProjects/intellectHome2.0/src/storage/dataBase"
)

const (
	StatusRegistred = "registred"
	StatusActive    = "active"
	StatusLost      = "lost"
	StatusNotActive = "offline"
)

type StatusWorker struct {
	db             *database.DataBase
	intervalUpdate time.Duration
	offlineTime    time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
	updateChan     chan string
	errChan        chan error
}

func MakeStatusWorker(db *database.DataBase, intervalInSecond time.Duration, offlineTime time.Duration) *StatusWorker {
	//сделать размер буфера настраивым
	ctx, cancel := context.WithCancel(context.Background())
	return &StatusWorker{
		db:             db,
		intervalUpdate: intervalInSecond,
		offlineTime:    offlineTime,
		ctx:            ctx,
		cancel:         cancel,
		updateChan:     make(chan string, 50),
		errChan:        db.ErrChan(),
	}
}

func (sw *StatusWorker) Start() {
	go sw.startMarkBoardStateBeforeUpdate()
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-sw.ctx.Done():
				log.Println("status worker stop")
				return
			case <-ticker.C:
				sw.UpdateStatuses()
			}
		}
	}()
}

func (sw *StatusWorker) ErrrorChan() chan error {
	return sw.errChan
}

func (sw *StatusWorker) UpdateChan() chan string {
	return sw.updateChan
}

func (sw *StatusWorker) startMarkBoardStateBeforeUpdate() {
	for {
		select {
		case <-sw.ctx.Done():
			return
		case id := <-sw.updateChan:
			if !sw.db.IsConnect() {
				time.Sleep(sw.intervalUpdate * 2)
				if !sw.db.IsConnect() {
					log.Printf("Status worker: not update status board: %s, DB error\n", id)
					continue
				}
			}
			//закинуть в отдельный метод в бд
			_, err := sw.db.GetPointer().Exec(`
		UPDATE boards b
		SET board_state = $1, updated_date = NOW()
		FROM boardInfo bi
		WHERE b.board_id = $2
		AND bi.updated_date IS NOT NULL;`,
				StatusActive, id)
			if err != nil {
				log.Println(err.Error())
				select {
				case sw.errChan <- err:
				default:
				}
				continue
			}
			log.Printf("set active status in board_id: %s\n", id)
		}
	}
}

func (sw *StatusWorker) UpdateStatuses() {
	if !sw.db.IsConnect() {
		time.Sleep(sw.intervalUpdate * 2)
		if !sw.db.IsConnect() {
			return
		}
	}
	sw.proccesingStatusLost()
	sw.proccesingStatusOffline()
}

func (sw *StatusWorker) proccesingStatusLost() {
	intervalStr := fmt.Sprintf("%d seconds", int(sw.intervalUpdate.Seconds()))
	offlineTime := fmt.Sprintf("%d seconds", int(sw.offlineTime.Seconds()))
	//закинуть в отдельный метод в бд
	res, err := sw.db.GetPointer().Exec(
		`UPDATE boards b
		SET board_state = $1, updated_date = NOW()
		FROM boardInfo bi
		WHERE b.board_id = bi.board_id
		AND bi.updated_date <= NOW() - $2::interval
		AND bi.updated_date >= NOW() - $3::interval
		AND b.board_state != $4;`,
		StatusLost, intervalStr, offlineTime, StatusLost)
	if err != nil {
		select {
		case sw.errChan <- err:
		default:
		}
		return
	}
	aff, err := res.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		select {
		case sw.errChan <- err:
		default:
		}
		return
	}
	if aff > 0 {
		log.Printf("set lost status in %d row\n", aff)
	}
}

func (sw *StatusWorker) proccesingStatusOffline() {
	offlineTime := fmt.Sprintf("%d seconds", int(sw.offlineTime.Seconds()))
	//закинуть в отдельный метод в бд
	res, err := sw.db.GetPointer().Exec(`
	UPDATE boards b
		SET board_state = $1, updated_date = NOW()
		FROM boardInfo bi
		WHERE b.board_id = bi.board_id
		AND bi.updated_date <= NOW() - $2::interval
		AND b.board_state != $3
`, StatusNotActive, offlineTime, StatusNotActive)

	if err != nil {
		select {
		case sw.errChan <- err:
		default:
		}
		return
	}

	aff, err := res.RowsAffected()
	if err != nil {
		select {
		case sw.errChan <- err:
		default:
		}
		return
	}
	if aff > 0 {
		log.Printf("set offline status in %d row\n", aff)
	}
}

func (sw *StatusWorker) Cancel() {
	sw.cancel()
}
