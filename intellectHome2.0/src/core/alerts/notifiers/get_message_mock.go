package notifiers

import "fmt"

type getMessageMock struct {
	MsgCh chan string
}

func NewGetMessageMock(bufferSize int) *getMessageMock {
	return &getMessageMock{
		MsgCh: make(chan string, bufferSize),
	}
}

func (g *getMessageMock) SentMessage(message string) error {
	var err error
	select {
	case g.MsgCh <- message:
	default:
		err = fmt.Errorf("message sent fail")
	}
	return err
}
