package main

import (
	"fmt"
	"sync"

	"encoding/json"

	"github.com/gorilla/websocket"
)

// hub is a map of websocket connections to addressed on it
type hub struct {
	sync.RWMutex
	m map[string]chan interface{}
}

func newHub() *hub {
	return &hub{
		m: map[string]chan interface{}{},
	}
}

// Add new connection to map and create inbox
func (h *hub) Add(id string, ws *websocket.Conn) error {
	h.Lock()

	defer h.Unlock()

	if _, ok := h.m[id]; ok {
		return fmt.Errorf("Already exists")
	}

	d := &username{
		Username: id,
	}
	bd, err := json.Marshal(d)
	if err != nil {
		return err
	}
	rd := json.RawMessage(bd)

	res := &response{
		ID:    -1,
		Cmd:   cmdLogin,
		Error: nil,
		Data:  &rd,
	}

	h.Broadcast(res, true)

	ch := make(chan interface{}, 10)

	go func() {
		for m := range ch {
			ws.WriteJSON(m)
		}
	}()

	h.m[id] = ch

	for u := range h.m {
		d := &username{
			Username: u,
		}
		bd, err := json.Marshal(d)
		if err != nil {
			return err
		}
		rd := json.RawMessage(bd)
		res := &response{
			ID:    -1,
			Cmd:   cmdLogin,
			Error: nil,
			Data:  &rd,
		}
		h.Direct(id, res, true)
	}

	return nil
}

// Del connection from map and close inbox
func (h *hub) Del(id string) {
	h.Lock()

	defer h.Unlock()

	if ch, ok := h.m[id]; ok {
		delete(h.m, id)
		close(ch)
	}
}

// Broadcast message ot all inboxes
func (h *hub) Broadcast(m interface{}, locked bool) error {
	if !locked {
		h.RLock()

		defer h.RUnlock()
	}

	for _, ch := range h.m {
		ch <- m
	}

	return nil
}

// Direct send message. Return false if receiver is not found
func (h *hub) Direct(id string, m interface{}, locked bool) bool {
	if !locked {
		h.RLock()

		defer h.RUnlock()
	}

	if ch, ok := h.m[id]; ok {
		ch <- m

		return true
	}

	return false
}
