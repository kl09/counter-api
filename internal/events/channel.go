package events

import (
	"container/list"
	"context"
	"encoding/json"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/websocket"
	service "github.com/kl09/counter-api"
)

// EventCh - events stream based on channel. It sends events to all alive websocket connections.
type EventCh struct {
	logger log.Logger

	ch chan service.Event
	// NOTE: linked list is a good option here - because we should add/remove clients from pool of connections often.
	wsConnections *list.List
}

// NewEventCh returns a new instance of EventCh.
func NewEventCh(logger log.Logger) *EventCh {
	return &EventCh{
		logger: logger,
		// TODO: should think about the capacity of channel.
		ch:            make(chan service.Event, 1000),
		wsConnections: list.New(),
	}
}

// AddConnection - adds ws connection to pool of connections.
func (c *EventCh) AddConnection(ws *websocket.Conn) error {
	c.wsConnections.PushBack(ws)

	return nil
}

// RemoveConnection - removes ws connection from the pool of connections.
func (c *EventCh) RemoveConnection(ws *websocket.Conn) error {
	// TODO: to speed up removing connections we ca add cache map[string]listElement.

	return nil
}

// Send - sends a new Event.
func (c *EventCh) Send(event service.Event) error {
	c.ch <- event

	return nil
}

// Listen - listens events.
func (c *EventCh) Listen(ctx context.Context) (err error) {
	level.Info(c.logger).Log("msg", "event listening started")
	defer func() {
		level.Info(c.logger).Log("msg", "event listening finished", "err", err)
	}()

	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		case e := <-c.ch:
			level.Info(c.logger).Log("msg", "got event", "e", e)
			b, err := json.Marshal(e)
			if err != nil {
				level.Error(c.logger).Log("err", err)
				continue
			}

			el := c.wsConnections.Front()
			for {
				if el == nil {
					break
				}

				ws := el.Value.(*websocket.Conn)
				err = ws.WriteMessage(websocket.TextMessage, b)
				if err != nil {
					// looks like connection aborted - remove the ws connection from the list.
					c.wsConnections.Remove(el)
				}

				el = el.Next()
			}
		}
	}
}
