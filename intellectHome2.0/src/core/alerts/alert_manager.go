package alerts

import (
	"fmt"
	"log"
	"slices"
	"sync"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/notifiers"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	board_info_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardInfoRules"
	board_status_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardStatusRules"
	database_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/dataBaseRules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

const (
	defaultMaxNumMessageParallel = 10
	defaultBufferMessage         = 50
)

type AlertsManager struct {
	recipients  []notifiers.Notifier
	mu          sync.RWMutex
	eventBus    *events.EventBus
	rules       []rules.Rule
	ch          chan string
	chanWorkers chan struct{}
}

func NewAlertsManager(eventBus *events.EventBus, ruls []rules.Rule, bufferMessage int, maxMessageParallel int, recipients ...notifiers.Notifier) (*AlertsManager, error) {
	if maxMessageParallel <= 0 {
		maxMessageParallel = defaultMaxNumMessageParallel
	}
	if bufferMessage <= 0 {
		bufferMessage = defaultBufferMessage
	}
	if ruls == nil {
		return nil, fmt.Errorf("no rules, alertsManager no point")
	}
	if slices.Contains(ruls, nil) {
		return nil, fmt.Errorf("invalid values in rules")
	}
	if slices.Contains(recipients, nil) {
		return nil, fmt.Errorf("invalid values in recipients")
	}
	if eventBus == nil {
		return nil, fmt.Errorf("eventBus not may be nil")
	}
	am := &AlertsManager{
		recipients: recipients,
		eventBus:   eventBus,
		rules:      ruls,
		// добавить буффер канала
		ch:          make(chan string, bufferMessage),
		chanWorkers: make(chan struct{}, maxMessageParallel),
	}
	return am, nil
}

func (am *AlertsManager) Start() {
	go am.startSendMessage()
	for i, rule := range am.rules {
		switch v := rule.(type) {
		case *board_info_rules.BoardInfoChecker:
			go am.processingBoardInfoChan(am.eventBus.Subscribe(TopicforBoardInfoChecker, NameForBoardInfoChecker), v)
		case *board_status_rules.BoardStatusChecker:
			go am.processingBoardStatusChan(am.eventBus.Subscribe(TopicForBoardStatusChecker, NameForBoardStatusChecker), v)
		case *database_rules.DataBaseChecker:
			am.processingDataBase(v)
		default:
			log.Printf("unexpected type in rules index: %d\n", i)
		}
	}
}

func (am *AlertsManager) startSendMessage() {
	for alert := range am.ch {
		am.chanWorkers <- struct{}{}
		go func(al string) {
			am.mu.RLock()
			notif := make([]notifiers.Notifier, len(am.recipients))
			copy(notif, am.recipients)
			am.mu.RUnlock()
			for _, notifier := range notif {
				err := notifier.SentMessage(al)
				if err != nil {
					log.Printf("Sent message: %s\nnot complete, %s\n", al, err.Error())
					continue
				}
			}
			<-am.chanWorkers
		}(alert)
	}
}

func (am *AlertsManager) processingBoardInfoChan(sub *events.TopicSubscriberOut, b *board_info_rules.BoardInfoChecker) {
	for event := range sub.Chan {
		al, err := b.Check(event)
		if err != nil {
			log.Println(err)
			continue
		}
		if al != nil {
			st := fmt.Sprintf("%s\nBoardID: %s\nMessage: %s", al.Type, *al.BoardID, al.Data)
			select {
			case am.ch <- st:
			default:
				log.Printf("message: %s\nrecipient not ready, deleted\n", st)
			}
		}
	}
}

func (am *AlertsManager) processingBoardStatusChan(sub *events.TopicSubscriberOut, s *board_status_rules.BoardStatusChecker) {
	for event := range sub.Chan {
		al, err := s.Check(event)
		if err != nil {
			log.Println(err)
			continue
		}
		if al != nil {
			st := fmt.Sprintf("%s\nBoardID: %s\nMessage: %s", al.Type, *al.BoardID, al.Data)
			select {
			case am.ch <- st:
			default:
				log.Printf("message: %s\nrecipient not ready, deleted\n", st)
			}
		}
	}
}

func (am *AlertsManager) processingDataBase(d *database_rules.DataBaseChecker) {
	go am.processingDataBaseError(d, am.eventBus.Subscribe(events.TopicErrorsDB, "alert_manager"))
	go am.processingDataBaseStatus(d, am.eventBus.Subscribe(events.TopicDataBaseStatus, "alert_manager"))
}

func (am *AlertsManager) processingDataBaseStatus(d *database_rules.DataBaseChecker, sub *events.TopicSubscriberOut) {
	for event := range sub.Chan {
		al, err := d.Check(event)
		if err != nil {
			log.Println(err)
			continue
		}
		if al != nil {
			str := fmt.Sprintf("%s\n%s\n", al.Type, al.Data)
			select {
			case am.ch <- str:
			default:
				log.Printf("message: %s\nrecipient not ready, deleted\n", str)
			}
		}
	}
}

func (am *AlertsManager) processingDataBaseError(d *database_rules.DataBaseChecker, sub *events.TopicSubscriberOut) {
	for event := range sub.Chan {
		al, err := d.Check(event)
		if err != nil {
			log.Println(err)
			continue
		}
		if al != nil {
			str := fmt.Sprintf("%s\n%s\n", al.Type, al.Data)
			select {
			case am.ch <- str:
			default:
				log.Printf("message: %s\nrecipient not ready, deleted\n", str)
			}
		}
	}
}
