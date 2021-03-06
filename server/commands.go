package main

import "encoding/json"

type command string

const (
	cmdLogin  command = "LOGIN"
	cmdLogout command = "LOGOUT"
	cmdPub    command = "PUBLISH"
	cmdUsers  command = "USERS"
)

type request struct {
	ID   int              `json:"id"`
	Cmd  command          `json:"cmd"`
	Data *json.RawMessage `json:"data,omitempty"`
}

type reqLoginData struct {
	Username string `json:"username"`
}

type reqPublishData struct {
	Message string `json:"message"`
}

type response struct {
	ID    int              `json:"id"`
	Cmd   command          `json:"cmd"`
	Error *string          `json:"error,omitempty"`
	Data  *json.RawMessage `json:"data,omitempty"`
}

type resUsersData struct {
	Users []string `json:"users"`
}

type username struct {
	Username string `json:"username"`
}

type message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}
