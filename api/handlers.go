package api

import (
	"net/http"
	"time"

	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	domain "github.com/kl09/counter-api"
)

type Status struct {
	Message string `json:"message"`
}

type Message struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

func (ro *Router) getKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key, ok := vars["key"]
	if !ok {
		encodeErrorResp(domain.ErrInvalidKey, w, ro.config.logger)
		return
	}

	v := ro.statsCounter.Get(key)

	encodeJSONResponse(w, Message{
		Key:   key,
		Value: v,
	})
}

func (ro *Router) incrementHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key, ok := vars["key"]
	if !ok {
		encodeErrorResp(domain.ErrInvalidKey, w, ro.config.logger)
		return
	}

	ro.statsCounter.Increment(key)

	encodeJSONResponse(w, Status{Message: "done"})
}

func (ro *Router) decrementHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key, ok := vars["key"]
	if !ok {
		encodeErrorResp(domain.ErrInvalidKey, w, ro.config.logger)
		return
	}

	ro.statsCounter.Decrement(key)

	encodeJSONResponse(w, Status{Message: "done"})
}

func (ro *Router) resetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key, ok := vars["key"]
	if !ok {
		encodeErrorResp(domain.ErrInvalidKey, w, ro.config.logger)
		return
	}

	ro.statsCounter.Reset(key)

	encodeJSONResponse(w, Status{Message: "done"})
}

func (ro *Router) subscriberHandler(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := u.Upgrade(w, r, http.Header{})
	if err != nil {
		level.Error(ro.config.logger).Log(
			"msg", "Cannot upgrade connection",
			"err", err,
		)
		return
	}

	err = ro.eventWsConnector.AddConnection(conn)
	if err != nil {
		level.Error(ro.config.logger).Log(
			"msg", "Cannot add connection",
			"err", err,
		)
		return
	}

	for {
		h := conn.PingHandler()
		err = h("ping")
		if err != nil {
			_ = ro.eventWsConnector.RemoveConnection(conn)
			return
		}

		time.Sleep(time.Second)
	}
}
