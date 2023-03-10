// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/gorilla/websocket"
	"github.com/kl09/counter-api"
	"sync"
)

// Ensure, that EventWSConnectorMock does implement domain.EventWSConnector.
// If this is not the case, regenerate this file with moq.
var _ domain.EventWSConnector = &EventWSConnectorMock{}

// EventWSConnectorMock is a mock implementation of domain.EventWSConnector.
//
//	func TestSomethingThatUsesEventWSConnector(t *testing.T) {
//
//		// make and configure a mocked domain.EventWSConnector
//		mockedEventWSConnector := &EventWSConnectorMock{
//			AddConnectionFunc: func(ws *websocket.Conn) error {
//				panic("mock out the AddConnection method")
//			},
//			RemoveConnectionFunc: func(ws *websocket.Conn) error {
//				panic("mock out the RemoveConnection method")
//			},
//		}
//
//		// use mockedEventWSConnector in code that requires domain.EventWSConnector
//		// and then make assertions.
//
//	}
type EventWSConnectorMock struct {
	// AddConnectionFunc mocks the AddConnection method.
	AddConnectionFunc func(ws *websocket.Conn) error

	// RemoveConnectionFunc mocks the RemoveConnection method.
	RemoveConnectionFunc func(ws *websocket.Conn) error

	// calls tracks calls to the methods.
	calls struct {
		// AddConnection holds details about calls to the AddConnection method.
		AddConnection []struct {
			// Ws is the ws argument value.
			Ws *websocket.Conn
		}
		// RemoveConnection holds details about calls to the RemoveConnection method.
		RemoveConnection []struct {
			// Ws is the ws argument value.
			Ws *websocket.Conn
		}
	}
	lockAddConnection    sync.RWMutex
	lockRemoveConnection sync.RWMutex
}

// AddConnection calls AddConnectionFunc.
func (mock *EventWSConnectorMock) AddConnection(ws *websocket.Conn) error {
	callInfo := struct {
		Ws *websocket.Conn
	}{
		Ws: ws,
	}
	mock.lockAddConnection.Lock()
	mock.calls.AddConnection = append(mock.calls.AddConnection, callInfo)
	mock.lockAddConnection.Unlock()
	if mock.AddConnectionFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.AddConnectionFunc(ws)
}

// AddConnectionCalls gets all the calls that were made to AddConnection.
// Check the length with:
//
//	len(mockedEventWSConnector.AddConnectionCalls())
func (mock *EventWSConnectorMock) AddConnectionCalls() []struct {
	Ws *websocket.Conn
} {
	var calls []struct {
		Ws *websocket.Conn
	}
	mock.lockAddConnection.RLock()
	calls = mock.calls.AddConnection
	mock.lockAddConnection.RUnlock()
	return calls
}

// RemoveConnection calls RemoveConnectionFunc.
func (mock *EventWSConnectorMock) RemoveConnection(ws *websocket.Conn) error {
	callInfo := struct {
		Ws *websocket.Conn
	}{
		Ws: ws,
	}
	mock.lockRemoveConnection.Lock()
	mock.calls.RemoveConnection = append(mock.calls.RemoveConnection, callInfo)
	mock.lockRemoveConnection.Unlock()
	if mock.RemoveConnectionFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.RemoveConnectionFunc(ws)
}

// RemoveConnectionCalls gets all the calls that were made to RemoveConnection.
// Check the length with:
//
//	len(mockedEventWSConnector.RemoveConnectionCalls())
func (mock *EventWSConnectorMock) RemoveConnectionCalls() []struct {
	Ws *websocket.Conn
} {
	var calls []struct {
		Ws *websocket.Conn
	}
	mock.lockRemoveConnection.RLock()
	calls = mock.calls.RemoveConnection
	mock.lockRemoveConnection.RUnlock()
	return calls
}
