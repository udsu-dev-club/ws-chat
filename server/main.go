package main

import (
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s := newServer()
	http.Handle("/ws", s)
	log.Print("listen on 8080")
	http.ListenAndServe(":8080", nil)
}
