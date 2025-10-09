package logs

import (
	"fmt"
	"net/http"
	"time"
)

type Logs struct {
	Records []string
}

func MakeNewLogsInfo() *Logs {
	return &Logs{make([]string, 0)}
}

func newRecords(record string, l *Logs) {
	l.Records = append(l.Records, record)
}

func (l *Logs) String() string {
	str := ""
	for i, v := range l.Records {
		str += fmt.Sprintf("Records %d: %s\n", i+1, v)
	}
	return str
}

func (l *Logs) CreateAndAddRecord(r *http.Request, body []byte, httpCode int, errors string, a ...any) {
	str := fmt.Sprintf("Time: %v\nUrl: %s\nMehod: %v\nBody: %s\nCode: %d\nErrors: %s\n", time.Now(), r.URL, r.Method, string(body), httpCode, errors)
	newRecords(str, l)
}
