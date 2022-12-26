package domain

import "github.com/gorilla/websocket"

// Event - is an event for each action.
type Event struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

// EventSender sends Event about new operations.
//
//go:generate moq -stub -out internal/mock/event_sender.go -pkg mock . EventSender
type EventSender interface {
	Send(event Event) error
}

// EventWSConnector saved ws connection to the pool.
//
//go:generate moq -stub -out internal/mock/event_ws_connector.go -pkg mock . EventWSConnector
type EventWSConnector interface {
	AddConnection(ws *websocket.Conn) error
	RemoveConnection(ws *websocket.Conn) error
}

// StatsCounter calc stats.
//
//go:generate moq -stub -out internal/mock/stats_counter.go -pkg mock . StatsCounter
type StatsCounter interface {
	Increment(key string)
	Decrement(key string)
	Reset(key string)
	Get(key string) int64
}
