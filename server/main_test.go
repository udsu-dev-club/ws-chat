package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

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
	connIvan, _, err := dlr.Dial("ws://localhost:8080/ws", nil)
	require.NoError(t, err)

	t.Run("1.Logout Ivan", func(t *testing.T) {
		require.NoError(t, connIvan.WriteJSON(&request{
			ID:  1,
			Cmd: cmdLogout,
		}))
		res := &response{}
		connIvan.SetReadDeadline(time.Now().Add(time.Second))
		require.NoError(t, connIvan.ReadJSON(res))
		assert.EqualValues(t, 1, res.ID)
		assert.Equal(t, cmdLogout, res.Cmd)
		require.NotNil(t, res.Error)
	})
	t.Run("2.Login Ivan", func(t *testing.T) {
		d := &username{
			Username: "Ivan",
		}
		dd, err := json.Marshal(d)
		require.NoError(t, err)
		dj := json.RawMessage(dd)
		require.NoError(t, connIvan.WriteJSON(&request{
			ID:   2,
			Cmd:  cmdLogin,
			Data: &dj,
		}))
		res := &response{}
		connIvan.SetReadDeadline(time.Now().Add(time.Second))
		require.NoError(t, connIvan.ReadJSON(res))
		require.Nil(t, res.Error)
		assert.EqualValues(t, 2, res.ID)
		assert.Equal(t, cmdLogin, res.Cmd)
	})

	connPetr, _, err := dlr.Dial("ws://localhost:8080/ws", nil)
	require.NoError(t, err)

	t.Run("3.Login Petr", func(t *testing.T) {
		d := &username{
			Username: "Petr",
		}
		dd, err := json.Marshal(d)
		require.NoError(t, err)
		dj := json.RawMessage(dd)
		require.NoError(t, connPetr.WriteJSON(&request{
			ID:   3,
			Cmd:  cmdLogin,
			Data: &dj,
		}))
		{
			res := &response{}
			connPetr.SetReadDeadline(time.Now().Add(time.Second))
			require.NoError(t, connPetr.ReadJSON(res))
			require.Nil(t, res.Error)
			assert.EqualValues(t, 3, res.ID)
			assert.Equal(t, cmdLogin, res.Cmd)
		}
		{
			res := &response{}
			connIvan.SetReadDeadline(time.Now().Add(time.Second))
			require.NoError(t, connIvan.ReadJSON(res))
			require.Nil(t, res.Error)
			assert.EqualValues(t, -1, res.ID)
			assert.Equal(t, cmdLogin, res.Cmd)
			d := &username{}
			require.NoError(t, json.Unmarshal(*res.Data, d))
			assert.Equal(t, "Petr", d.Username)
		}
	})
	t.Run("4.Login Petr", func(t *testing.T) {
		connPetr2, _, err := dlr.Dial("ws://localhost:8080/ws", nil)
		require.NoError(t, err)
		d := &username{
			Username: "Petr",
		}
		dd, err := json.Marshal(d)
		require.NoError(t, err)
		dj := json.RawMessage(dd)
		require.NoError(t, connPetr2.WriteJSON(&request{
			ID:   4,
			Cmd:  cmdLogin,
			Data: &dj,
		}))
		res := &response{}
		connPetr2.SetReadDeadline(time.Now().Add(time.Second))
		require.NoError(t, connPetr2.ReadJSON(res))
		require.NotNil(t, res.Error)
		assert.EqualValues(t, 4, res.ID)
		assert.Equal(t, cmdLogin, res.Cmd)
		assert.Equal(t, "Already exists", *res.Error)
	})

	connVasiliy, _, err := dlr.Dial("ws://localhost:8080/ws", nil)
	require.NoError(t, err)

	t.Run("5.Login Vasiliy", func(t *testing.T) {
		d := &username{
			Username: "Vasiliy",
		}
		dd, err := json.Marshal(d)
		require.NoError(t, err)
		dj := json.RawMessage(dd)
		require.NoError(t, connVasiliy.WriteJSON(&request{
			ID:   3,
			Cmd:  cmdLogin,
			Data: &dj,
		}))
		{
			res := &response{}
			connVasiliy.SetReadDeadline(time.Now().Add(time.Second))
			require.NoError(t, connVasiliy.ReadJSON(res))
			require.Nil(t, res.Error)
			assert.EqualValues(t, 3, res.ID)
			assert.Equal(t, cmdLogin, res.Cmd)
		}
		{
			res := &response{}
			connIvan.SetReadDeadline(time.Now().Add(time.Second))
			require.NoError(t, connIvan.ReadJSON(res))
			require.Nil(t, res.Error)
			assert.EqualValues(t, -1, res.ID)
			assert.Equal(t, cmdLogin, res.Cmd)
			d := &username{}
			require.NoError(t, json.Unmarshal(*res.Data, d))
			assert.Equal(t, "Vasiliy", d.Username)
		}
		{
			res := &response{}
			connPetr.SetReadDeadline(time.Now().Add(time.Second))
			require.NoError(t, connPetr.ReadJSON(res))
			require.Nil(t, res.Error)
			assert.EqualValues(t, -1, res.ID)
			assert.Equal(t, cmdLogin, res.Cmd)
			d := &username{}
			require.NoError(t, json.Unmarshal(*res.Data, d))
			assert.Equal(t, "Vasiliy", d.Username)
		}
	})

}
