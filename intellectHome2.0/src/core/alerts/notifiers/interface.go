package notifiers

type Notifier interface {
	SentMessage(message string) error
}
