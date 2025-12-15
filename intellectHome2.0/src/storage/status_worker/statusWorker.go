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
	StatusRegistred     = "registred"
	StatusActive        = "active"
	StatusLost          = "lost"
	StatusNotActive     = "offline"
	NameWorkerPublisher = "status_worker"
)

type StatusWorker struct {
	db                    *database.DataBase
	intervalUpdate        time.Duration
	offlineTime           time.Duration
	ctx                   context.Context
	cancel                context.CancelFunc
	EventBus              *events.EventBus
	subBoardInfoUpdate    *events.TopicSubscriberOut
	subBoardStatusUpddate *events.TopicSubscriberOut
	subErrorDB            *events.TopicSubscriberOut
}

func MakeStatusWorker(db *database.DataBase, intervalInSecond time.Duration, offlineTime time.Duration, EventBus *events.EventBus) *StatusWorker {
	//сделать размер буфера настраивым
	ctx, cancel := context.WithCancel(context.Background())
	return &StatusWorker{
		db:                    db,
		intervalUpdate:        intervalInSecond,
		offlineTime:           offlineTime,
		ctx:                   ctx,
		cancel:                cancel,
		EventBus:              EventBus,
		subBoardInfoUpdate:    EventBus.Subscribe(events.TopicBoardInfoUpdate, NameWorkerPublisher),
		subBoardStatusUpddate: EventBus.Subscribe(events.TopicBoardsStatusUpdate, NameWorkerPublisher),
		subErrorDB:            EventBus.Subscribe(events.TopicErrorsDB, NameWorkerPublisher),
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

func (sw *StatusWorker) startMarkBoardStateBeforeUpdate() {
	for {
		select {
		case <-sw.ctx.Done():
			return
		case info := <-sw.subBoardInfoUpdate.Chan:
			if !sw.db.IsConnect() {
				time.Sleep(sw.intervalUpdate * 2)
				if !sw.db.IsConnect() {
					log.Printf("Status worker: not update status board: %s, DB error\n", info.BoardID)
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
				StatusActive, info.BoardID)
			if err != nil {
				log.Println(err.Error())
				err = sw.EventBus.Publish(sw.subErrorDB.Topic, events.Event{
					Type:       sw.subErrorDB.Topic,
					BoardID:    info.BoardID,
					Payload:    err,
					Publisher:  NameWorkerPublisher,
					DatePublic: time.Now(),
				}, sw.subErrorDB.ID)
				if err != nil {
					log.Println(err)
				}
				continue
			}
			err = sw.EventBus.Publish(sw.subBoardStatusUpddate.Topic, events.Event{
				Type:       sw.subBoardStatusUpddate.Topic,
				BoardID:    info.BoardID,
				Payload:    "update status to active",
				Publisher:  NameWorkerPublisher,
				DatePublic: time.Now(),
			}, sw.subBoardStatusUpddate.ID)
			if err != nil {
				log.Println(err)
			}
			log.Printf("set active status in board_id: %s\n", info.BoardID)
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
	rows, err := sw.db.GetPointer().Query(
		`UPDATE boards b
		SET board_state = $1, updated_date = NOW()
		FROM boardInfo bi
		WHERE b.board_id = bi.board_id
		AND bi.updated_date <= NOW() - $2::interval
		AND bi.updated_date >= NOW() - $3::interval
		AND b.board_state != $4
		RETURNING b.board_id;`,
		StatusLost, intervalStr, offlineTime, StatusLost)
	if err != nil {
		log.Println(err.Error())
		err := sw.EventBus.Publish(sw.subErrorDB.Topic, events.Event{
			Type:       sw.subErrorDB.Topic,
			Payload:    err,
			Publisher:  NameWorkerPublisher,
			DatePublic: time.Now(),
		}, sw.subErrorDB.ID)
		if err != nil {
			log.Println(err)
		}
		return
	}
	var updatedBoardIDs []string
	for rows.Next() {
		var boardID string
		if err := rows.Scan(&boardID); err != nil {
			log.Printf("Error scanning board_id: %v", err)
			continue
		}
		updatedBoardIDs = append(updatedBoardIDs, boardID)
	}
	if len(updatedBoardIDs) != 0 {
		for _, boardID := range updatedBoardIDs {
			err := sw.EventBus.Publish(sw.subBoardStatusUpddate.Topic, events.Event{
				Type:       sw.subBoardStatusUpddate.Topic,
				BoardID:    boardID,
				Payload:    "updated status to lost",
				DatePublic: time.Now(),
				Publisher:  NameWorkerPublisher,
			}, sw.subBoardStatusUpddate.ID)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (sw *StatusWorker) proccesingStatusOffline() {
	offlineTime := fmt.Sprintf("%d seconds", int(sw.offlineTime.Seconds()))
	//закинуть в отдельный метод в бд
	rows, err := sw.db.GetPointer().Query(`
	UPDATE boards b
		SET board_state = $1, updated_date = NOW()
		FROM boardInfo bi
		WHERE b.board_id = bi.board_id
		AND bi.updated_date <= NOW() - $2::interval
		AND b.board_state != $3
		RETURNING b.board_id;
`, StatusNotActive, offlineTime, StatusNotActive)

	if err != nil {
		log.Println(err.Error())
		err := sw.EventBus.Publish(sw.subErrorDB.Topic, events.Event{
			Type:       sw.subErrorDB.Topic,
			Payload:    err,
			Publisher:  NameWorkerPublisher,
			DatePublic: time.Now(),
		}, sw.subErrorDB.ID)
		if err != nil {
			log.Println(err)
		}
		return
	}

	var updatedBoardIDs []string
	for rows.Next() {
		var boardID string
		if err := rows.Scan(&boardID); err != nil {
			log.Printf("Error scanning board_id: %v", err)
			continue
		}
		updatedBoardIDs = append(updatedBoardIDs, boardID)
	}
	if len(updatedBoardIDs) != 0 {
		for _, boardID := range updatedBoardIDs {
			err := sw.EventBus.Publish(sw.subBoardStatusUpddate.Topic, events.Event{
				Type:       sw.subBoardStatusUpddate.Topic,
				BoardID:    boardID,
				Payload:    "updated status to offline",
				DatePublic: time.Now(),
				Publisher:  NameWorkerPublisher,
			}, sw.subBoardStatusUpddate.ID)
			if err != nil {
				log.Println(err)
			}
		}
		fmt.Println(updatedBoardIDs)
	}

}

func (sw *StatusWorker) Cancel() {
	sw.cancel()
}
