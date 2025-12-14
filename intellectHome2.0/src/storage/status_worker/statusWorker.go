package status_worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
	database "github.com/Piccadilly98/goProjects/intellectHome2.0/src/storage/dataBase"
)

const (
	StatusRegistred = "registred"
	StatusActive    = "active"
	StatusLost      = "lost"
	StatusNotActive = "offline"
	NamePublisher   = "status_worker"
)

type StatusWorker struct {
	db             *database.DataBase
	intervalUpdate time.Duration
	offlineTime    time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
	EventBus       *events.EventBus
	errChan        chan error
}

func MakeStatusWorker(db *database.DataBase, intervalInSecond time.Duration, offlineTime time.Duration, EventBus *events.EventBus) *StatusWorker {
	//сделать размер буфера настраивым
	ctx, cancel := context.WithCancel(context.Background())
	return &StatusWorker{
		db:             db,
		intervalUpdate: intervalInSecond,
		offlineTime:    offlineTime,
		ctx:            ctx,
		cancel:         cancel,
		EventBus:       EventBus,
		errChan:        db.ErrChan(),
	}
}

func (sw *StatusWorker) Start() {
	subBoardUploadInfo := sw.EventBus.Subscribe(events.TopicBoardUploadInfo, NamePublisher)
	subBoardChangeStatus := sw.EventBus.Subscribe(events.TopicBoardsEventStatus, NamePublisher)
	go sw.startMarkBoardStateBeforeUpdate(subBoardUploadInfo, subBoardChangeStatus)
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

func (sw *StatusWorker) startMarkBoardStateBeforeUpdate(subBoardUploadInfo *events.TopicSubscriberOut, subBoardChangeStatus *events.TopicSubscriberOut) {
	for {
		select {
		case <-sw.ctx.Done():
			return
		case sub := <-subBoardUploadInfo.Chan:
			log.Printf("Sent event by: %s\n", sub.Publisher)
			if !sw.db.IsConnect() {
				time.Sleep(sw.intervalUpdate * 2)
				if !sw.db.IsConnect() {
					log.Printf("Status worker: not update status board: %s, from publisher: %s. DB error\n", sub.BoardID, sub.Publisher)
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
				StatusActive, sub.BoardID)
			if err != nil {
				log.Println(err.Error())
				select {
				case sw.errChan <- err:
				default:
				}
				continue
			}
			log.Printf("set active status in board_id: %s\n", sub.BoardID)
			err = sw.EventBus.Publish(subBoardChangeStatus.Topic, events.Event{
				Type:       subBoardChangeStatus.Topic,
				BoardID:    sub.BoardID,
				Payload:    fmt.Sprintf("change status to %s after upload", StatusActive),
				Publisher:  NamePublisher,
				DatePublic: time.Now(),
			}, subBoardChangeStatus.ID)
			if err != nil {
				log.Printf("Erro in publish statusWorker: %s\n", err.Error())
			} else {
				log.Printf("publish event in topic: %s, publisher: %s\n", subBoardChangeStatus.Topic, NamePublisher)
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
