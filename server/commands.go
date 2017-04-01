package main

import (
	"encoding/json"
	"time"
)

type command string

const (
	cmdLogin  command = "LOGIN"
	cmdLogout command = "LOGOUT"
	cmdPub    command = "PUBLISH"
	cmdUsers  command = "USERS"
	cmdSub    command = "SUBSCRIBE"
)

type request struct {
	ID   int32            `json:"id"`
	Cmd  command          `json:"cmd"`
	Data *json.RawMessage `json:"data"`
}

type reqLoginData struct {
	Username string `json:"username"`
}

type reqPublishData struct {
	Message string `json:"message"`
}

type response struct {
	ID    int32            `json:"id"`
	Cmd   command          `json:"cmd"`
	Error *string          `json:"error"`
	Data  *json.RawMessage `json:"data"`
}

type resUsersData struct {
	Users []string `json:"users"`
}

type username struct {
	Username string `json:"username"`
}

type message struct {
	Timestamp time.Time `json:"timestamp"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
}
