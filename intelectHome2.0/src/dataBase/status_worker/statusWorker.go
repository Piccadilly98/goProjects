package status_worker

import (
	"context"
	"fmt"
	"log"
	"time"

	database "github.com/Piccadilly98/goProjects/intelectHome2.0/src/dataBase"
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
					continue
				}
			}
			_, err := sw.db.GetPointer().Exec(`
		UPDATE boards b
		SET board_state = $1, updated = NOW()
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
	_, err := sw.db.GetPointer().Exec(
		`UPDATE boards b
		SET board_state = $1, updated = NOW()
		FROM boardInfo bi
		WHERE b.board_id = bi.board_id
		AND bi.updated_date <= NOW() - $2::interval
		AND bi.updated_date >= NOW() - $3::interval
		AND b.board_state != $4;`,
		StatusLost, intervalStr, offlineTime, StatusLost)
	if err != nil {
		log.Println(err.Error())
		select {
		case sw.errChan <- err:
		default:
		}
		return
	}
}

func (sw *StatusWorker) proccesingStatusOffline() {
	offlineTime := fmt.Sprintf("%d seconds", int(sw.offlineTime.Seconds()))
	_, err := sw.db.GetPointer().Exec(`
	UPDATE boards b
		SET board_state = $1, updated = NOW()
		FROM boardInfo bi
		WHERE b.board_id = bi.board_id
		AND bi.updated_date <= NOW() - $2::interval
		AND b.board_state != $3
`, StatusNotActive, offlineTime, StatusNotActive)

	if err != nil {
		log.Println(err.Error())
		select {
		case sw.errChan <- err:
		default:
		}
		return
	}
}

func (sw *StatusWorker) Cancel() {
	sw.cancel()
}
