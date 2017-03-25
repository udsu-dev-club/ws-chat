package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	dlr := websocket.DefaultDialer
	conn, _, err := dlr.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	t.Run("1.Logout", func(t *testing.T) {
		require.NoError(t, conn.WriteJSON(&request{
			ID:  1,
			Cmd: cmdLogout,
		}))
		res := &response{}
		require.NoError(t, conn.ReadJSON(res))
		assert.EqualValues(t, 1, res.ID)
		assert.Equal(t, cmdLogout, res.Cmd)
		require.NotNil(t, res.Error)
	})
	t.Run("2.Login", func(t *testing.T) {
		d := &reqLoginData{
			Username: "Ivan",
		}
		dd, err := json.Marshal(d)
		require.NoError(t, err)
		dj := json.RawMessage(dd)
		require.NoError(t, conn.WriteJSON(&request{
			ID:   2,
			Cmd:  cmdLogin,
			Data: &dj,
		}))
		res := &response{}
		require.NoError(t, conn.ReadJSON(res))
		require.Nil(t, res.Error)
		assert.EqualValues(t, 2, res.ID)
		assert.Equal(t, cmdLogin, res.Cmd)
	})
}
