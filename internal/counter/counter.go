package counter

import (
	"sync"

	service "github.com/kl09/counter-api"
)

// Stats - is a data for storing counts for specific keys.
type Stats struct {
	data        map[string]int64
	eventSender service.EventSender

	mx sync.RWMutex
}

// NewStats returns a new instance of Stats.
func NewStats(eventSender service.EventSender) *Stats {
	return &Stats{
		data:        map[string]int64{},
		eventSender: eventSender,
	}
}

// Increment data for specific key.
func (c *Stats) Increment(key string) {
	var newVal int64

	c.mx.Lock()
	defer func() {
		c.mx.Unlock()
		_ = c.eventSender.Send(service.Event{
			Key:   key,
			Value: newVal,
		})
	}()

	_, ok := c.data[key]
	if !ok {
		newVal = 1
		c.data[key] = newVal
		return
	}

	newVal = c.data[key] + 1
	c.data[key] = newVal
}

// Decrement data for specific key.
func (c *Stats) Decrement(key string) {
	var newVal int64

	c.mx.Lock()
	defer func() {
		c.mx.Unlock()
		_ = c.eventSender.Send(service.Event{
			Key:   key,
			Value: newVal,
		})
	}()

	_, ok := c.data[key]
	if !ok {
		newVal = -1
		c.data[key] = newVal
		return
	}

	newVal = c.data[key] - 1
	c.data[key] = newVal
}

// Reset resets stat for specific key.
func (c *Stats) Reset(key string) {
	c.mx.Lock()
	defer func() {
		c.mx.Unlock()
		_ = c.eventSender.Send(service.Event{
			Key:   key,
			Value: 0,
		})
	}()

	c.data[key] = 0
}

// Get returns counter by key.
func (c *Stats) Get(key string) int64 {
	c.mx.Lock()
	defer c.mx.Unlock()

	v, ok := c.data[key]
	if !ok {
		c.data[key] = 0
		return 0
	}

	return v
}
