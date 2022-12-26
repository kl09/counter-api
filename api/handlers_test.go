package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/kl09/counter-api/internal/mock"
	"github.com/stretchr/testify/require"
)

func makeRequest(h http.Handler, method, url string, b []byte) (*http.Response, []byte, error) {
	srv := httptest.NewServer(h)
	defer srv.Close()

	req, err := http.NewRequestWithContext(context.Background(), method, srv.URL+url, bytes.NewReader(b))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	b, err = io.ReadAll(resp.Body)

	return resp, b, err
}

func TestRouter_GetKey_Handler(t *testing.T) {
	key := "a1"

	r := NewRouter(
		NewHandlerConfig(),
		&mock.StatsCounterMock{
			GetFunc: func(key string) int64 {
				return 100500
			},
		},
		&mock.EventWSConnectorMock{
			AddConnectionFunc: func(ws *websocket.Conn) error {
				return nil
			},
		},
	)
	resp, b, err := makeRequest(
		r.Handler(), "GET", "/v1/"+key, []byte(""),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, "{\"key\":\"a1\",\"value\":100500}\n", string(b))
}

func TestRouter_Increment_Handler(t *testing.T) {
	key := "a1"

	r := NewRouter(
		NewHandlerConfig(),
		&mock.StatsCounterMock{
			IncrementFunc: func(key string) {},
		},
		&mock.EventWSConnectorMock{
			AddConnectionFunc: func(ws *websocket.Conn) error {
				return nil
			},
		},
	)
	resp, b, err := makeRequest(
		r.Handler(), "POST", "/v1/"+key+"/increment", []byte(""),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, "{\"message\":\"done\"}\n", string(b))
}

func TestRouter_Decrement_Handler(t *testing.T) {
	key := "a1"

	r := NewRouter(
		NewHandlerConfig(),
		&mock.StatsCounterMock{
			DecrementFunc: func(key string) {},
		},
		&mock.EventWSConnectorMock{
			AddConnectionFunc: func(ws *websocket.Conn) error {
				return nil
			},
		},
	)
	resp, b, err := makeRequest(
		r.Handler(), "POST", "/v1/"+key+"/decrement", []byte(""),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, "{\"message\":\"done\"}\n", string(b))
}

func TestRouter_Reset_Handler(t *testing.T) {
	key := "a1"

	r := NewRouter(
		NewHandlerConfig(),
		&mock.StatsCounterMock{
			ResetFunc: func(key string) {},
		},
		&mock.EventWSConnectorMock{
			AddConnectionFunc: func(ws *websocket.Conn) error {
				return nil
			},
		},
	)
	resp, b, err := makeRequest(
		r.Handler(), "POST", "/v1/"+key+"/reset", []byte(""),
	)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, 200, resp.StatusCode)
	require.Equal(t, "{\"message\":\"done\"}\n", string(b))
}
