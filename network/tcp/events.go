package tcp

import "github.com/loveandpeople/hive.go/events"

type tcpServerEvents struct {
	Start    *events.Event
	Shutdown *events.Event
	Connect  *events.Event
	Error    *events.Event
}
