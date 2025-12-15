package alerts

import (
	"fmt"
	"log"
	"sync"

	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/notifiers"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules"
	board_info_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardInfoRules"
	board_status_rules "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/alerts/rules/boardStatusRules"
	"github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"
)

/*
Принимаем нотифиеров
bool о том какие уведомления нужны

храним:
получателей
mu
*/

type Alert struct {
	Type    string
	BoardID *string
	Data    string
}

type AlertsManager struct {
	recipients []notifiers.Notifier
	mu         sync.RWMutex
	eventBus   *events.EventBus
	rules      []rules.Rule
	ch         chan string
}

func NewAlertsManager(eventBus *events.EventBus, ruls []rules.Rule, bufferMessage int, recipients ...notifiers.Notifier) *AlertsManager {
	if ruls == nil {
		return nil
	}
	am := &AlertsManager{
		recipients: recipients,
		eventBus:   eventBus,
		rules:      ruls,
		ch:         make(chan string),
	}
	return am
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
	}
}

func (am *AlertsManager) StartSentMessage() {
	for alert := range am.ch {
		am.mu.RLock()
		for _, notifier := range am.recipients {
			err := notifier.SentMessage(alert)
			if err != nil {
				log.Printf("Sent message: %s\nnot complete, %v\n", err.Error())
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
