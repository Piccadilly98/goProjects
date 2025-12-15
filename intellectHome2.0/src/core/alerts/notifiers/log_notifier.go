package notifiers

import "log"

type LogNotifier struct {
}

func (l *LogNotifier) SentMessage(message string) error {
	log.Printf("Alert by log_notifier:\n%s\n", message)
	return nil
}
