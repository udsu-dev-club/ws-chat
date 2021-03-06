package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type server struct {
	hub
}

func newServer() *server {
	h := newHub()

	return &server{*h}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	log.Printf("connected %s", r.RemoteAddr)

	id := new(string)
	for {
		req := &request{}
		if err := ws.ReadJSON(req); err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Printf("%s closed", *id)
			}
			log.Printf("%s error: %s", *id, err)

			return
		}
		log.Printf("Got request %q", req.Cmd)

		res := &response{
			ID:  req.ID,
			Cmd: req.Cmd,
		}

		if err := s.Switch(id, req, ws); err != nil {
			log.Print(err)
			serr := err.Error()
			res.Error = &serr
		}
		if len(*id) > 0 {
			s.Direct(*id, res, false)
		} else {
			ws.WriteJSON(res)
		}
	}
}

func (s *server) Switch(id *string, req *request, ws *websocket.Conn) error {
	switch req.Cmd {
	case cmdLogin:

		return s.Login(id, req, ws)

	case cmdLogout:

		return s.Logout(id, req, ws)

	case cmdPub:

		return s.Publish(id, req, ws)
	}

	return fmt.Errorf("Unknown command")
}

func (s *server) Login(id *string, req *request, ws *websocket.Conn) error {
	if len(*id) != 0 {

		return fmt.Errorf("Already logined")
	}
	d := &username{}
	if req.Data == nil {

		return fmt.Errorf("Username required")
	}
	if err := json.Unmarshal(*req.Data, d); err != nil {

		return err
	}
	if len(d.Username) == 0 {

		return fmt.Errorf("Username is required")
	}
	if strings.Contains(d.Username, " ") {

		return fmt.Errorf("Username must not contain spaces")
	}
	if err := s.Add(d.Username, ws); err != nil {

		return err
	}
	*id = d.Username

	return nil
}

func (s *server) Logout(id *string, _ *request, _ *websocket.Conn) error {
	if len(*id) == 0 {

		return fmt.Errorf("Not logined yet")
	}
	s.Del(*id)
	*id = ""

	return nil
}

func (s *server) Publish(id *string, req *request, ws *websocket.Conn) error {
	if len(*id) == 0 {

		return fmt.Errorf("Not logined to publish")
	}
	if req.Data == nil {

		return fmt.Errorf("Message required")
	}
	d := &reqPublishData{}
	if err := json.Unmarshal(*req.Data, d); err != nil {

		return err
	}
	msg := &message{
		Author: *id,
		Body:   d.Message,
	}
	msgb, err := json.Marshal(msg)
	if err != nil {

		return err
	}
	msgj := json.RawMessage(msgb)
	res := response{
		ID:   -1,
		Cmd:  cmdPub,
		Data: &msgj,
	}
	if d.Message[0] == '@' {
		receiver := strings.Split(d.Message, " ")[0][1:]
		if s.Direct(receiver, res, false) {
			s.Direct(*id, res, false)

			return nil
		}

		return fmt.Errorf("Receiver not found")
	}

	return s.Broadcast(res, false)
}
