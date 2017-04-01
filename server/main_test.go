package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"path/filepath"
	"strings"

	"encoding/json"

	"io/ioutil"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var users = []string{"Ivan", "Petr", "Vasiliy"}

type Case struct {
	Requests  map[string]json.RawMessage `json:"requests"`
	Responses map[string]json.RawMessage `json:"responses"`
}

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s := newServer()
	http.Handle("/ws", s)
	hs := http.Server{
		Addr: ":8080",
	}
	go hs.ListenAndServe()

	status := m.Run()

	hs.Shutdown(context.Background())

	os.Exit(status)
}

func TestWorkflow(t *testing.T) {
	var err error

	dlr := websocket.DefaultDialer
	conns := map[string]*websocket.Conn{}
	for _, u := range users {
		conns[u], _, err = dlr.Dial("ws://localhost:8080/ws", nil)
		require.NoError(t, err)
	}
	cases, err := filepath.Glob(filepath.Join("test_data", "*"))
	require.NoError(t, err)
	for _, c := range cases {
		name := strings.Split(filepath.Base(c), ".")[0]
		t.Run(name, func(t *testing.T) {
			cc := &Case{}
			cb, err := ioutil.ReadFile(c)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(cb, cc))
			for u, d := range cc.Requests {
				require.NoError(t, conns[u].WriteJSON(d))
			}
			for u, d := range cc.Responses {
				conns[u].SetReadDeadline(time.Now().Add(time.Second))
				res := json.RawMessage{}
				require.NoError(t, conns[u].ReadJSON(&res))
				assert.JSONEq(t, string(d), string(res))
			}
		})
	}
}
