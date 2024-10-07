package main

import (
	"context"
	"encoding/json"
	"sync"

	"golang.org/x/net/websocket"
)

type WebSocketHandler struct {
	Changes chan *Event
	errors  chan error

	sockets map[*websocket.Conn]bool
	m       *sync.Mutex
}

func NewWebSocketHandler() *WebSocketHandler {
	return &WebSocketHandler{
		Changes: make(chan *Event, 16),
		errors:  make(chan error, 1),

		// TODO we never empty the session map for websockets and idk how
		sockets: make(map[*websocket.Conn]bool),
		m:       &sync.Mutex{},
	}
}

func (w *WebSocketHandler) AddSocket(ws *websocket.Conn) {
	w.m.Lock()
	defer w.m.Unlock()
	w.sockets[ws] = true
}

func (w *WebSocketHandler) RemoveSocket(ws *websocket.Conn) {
	w.m.Lock()
	defer w.m.Unlock()
	delete(w.sockets, ws)
}

func (w *WebSocketHandler) Broadcast(e *Event) {
	w.Changes <- e
}

func (w *WebSocketHandler) Open(ctx context.Context) {
	go w.Monitor()

	for {
		if e := <-w.Changes; e != nil {
			eventBytes, err := json.Marshal(e)
			if err != nil {
				w.errors <- err
				return
			}
			for ws := range w.sockets {
				if _, err := ws.Write(eventBytes); err != nil {
					w.RemoveSocket(ws)
					ws.Close()
					continue
				}
			}
		}
	}
}

func (w *WebSocketHandler) Monitor() {
	for err := range w.errors {
		if err != nil {
			panic(err)
		}
	}
}
