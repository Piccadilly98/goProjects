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

type AlertsManager struct {
	recipients []notifiers.Notifier
	mu         sync.RWMutex
	eventBus   *events.EventBus
	rules      []rules.Rule
	ch         chan string
}

func NewAlertsManager(eventBus *events.EventBus, ruls []rules.Rule, bufferMessage int, recipients ...notifiers.Notifier) (*AlertsManager, error) {
	if ruls == nil {
		return nil, fmt.Errorf("no rules, alertsManager no point")
	}
	if slices.Contains(ruls, nil) {
		return nil, fmt.Errorf("invalid values in rules")
	}
	am := &AlertsManager{
		recipients: recipients,
		eventBus:   eventBus,
		rules:      ruls,
		ch:         make(chan string),
	}
	return am, nil
}

func (am *AlertsManager) Start() {
	go am.StartSentMessage()
	for _, rule := range am.rules {
		if b, ok := rule.(*board_info_rules.BoardInfoChecker); ok {
			go am.processingBoardInfoChan(am.eventBus.Subscribe(TopicforBoardInfoChecker, NameForBoardInfoChecker), b)
		}
		if s, ok := rule.(*board_status_rules.BoardStatusChecker); ok {
			go am.processingBoardStatusChan(am.eventBus.Subscribe(TopicForBoardStatusChecker, NameForBoardStatusChecker), s)
		}
		if d, ok := rule.(*database_rules.DataBaseChecker); ok {
			go am.processingDataBase(d)
		}
	}
}

func (am *AlertsManager) StartSentMessage() {
	for alert := range am.ch {
		am.mu.RLock()
		for _, notifier := range am.recipients {
			err := notifier.SentMessage(alert)
			if err != nil {
				log.Printf("Sent message: %s\nnot complete, %s\n", alert, err.Error())
				continue
			}
		}
		am.mu.RUnlock()
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
				log.Printf("message: %s\nmessage recipient not ready, deleted\n", st)
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
				log.Printf("message: %s\nmessage recipient not ready, deleted\n", st)
			}
		}
	}
}

func (am *AlertsManager) processingDataBase(d *database_rules.DataBaseChecker) {
	if d.ErrorCheck() {
		go am.processingDataBaseError(d, am.eventBus.Subscribe(events.TopicErrorsDB, "data_base_error_checker"))
	}
	if d.StatusCheck() {
		go am.processingDataBaseStatus(d, am.eventBus.Subscribe(events.TopicDataBaseStatus, "data_base_status_checker"))
	}
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
				log.Printf("message: %s\nmessage recipient not ready, deleted\n", str)
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
			select {
			case am.ch <- al.Data:
			default:
				log.Printf("message: %s\nmessage recipient not ready, deleted\n", al.Data)
			}
		}
	}
}
