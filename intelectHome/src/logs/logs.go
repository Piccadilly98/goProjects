package logs

import (
	"fmt"
	"net/http"
	"time"
)

type Logs struct {
	records []string
}

func MakeNewLogsInfo() *Logs {
	return &Logs{make([]string, 0)}
}

func newRecords(record string, l *Logs) {
	l.records = append(l.records, record)
}

func (l *Logs) String() string {
	str := ""
	for i, v := range l.records {
		str += fmt.Sprintf("Records %d: %s\n", i+1, v)
	}
	return str
}

func (l *Logs) CreateAndAddRecordPost(r *http.Request, ID, status string, httpCode int) {
	str := fmt.Sprintf("Time: %v\nUrl: %s\nMehod: %v\nBody:\nID: %s\nStatus: %s\nCode: %d", time.Now(), r.URL, r.Method, ID, status, httpCode)

	newRecords(str, l)
}

func (l *Logs) CreateAndAddRecordGet(r *http.Request, body []byte, httpCode int) {
	str := fmt.Sprintf("Time: %v\nUrl: %s\nMehod: %v\nBody: %s\nCode: %d", time.Now(), r.URL, r.Method, string(body), httpCode)
	newRecords(str, l)
}
