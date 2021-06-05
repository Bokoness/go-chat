package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bokoness/go-chat/auth"
	"github.com/bokoness/go-chat/chat"
	"github.com/bokoness/go-chat/rdb"
	socketio "github.com/googollee/go-socket.io"
)

func main() {
	rdb.CreateCon()
	server := socketio.NewServer(nil)
	auth.CreateAuthEvents(server)
	chat.CreateChatEvents(server)
	chat.CreateWaitingEvents(server)

	server.OnConnect("/", func(s socketio.Conn) error {
		server.JoinRoom("/", "auth", s)
		s.Emit("auth", true)
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		room := rdb.Client.HGet("players-rooms", s.ID())
		rdb.Client.HDel("players-rooms", s.ID())
		rdb.Client.HDel(room.Val()+"/players", s.ID())
		server.BroadcastToRoom("/", room.Val(), "left", s.ID())
		//search exit user from room using redis
		fmt.Println(s.Rooms())
		fmt.Println("closed", reason, s.Namespace())
	})

	go server.Serve()
	defer server.Close()

	fs := http.FileServer(http.Dir("static"))

	http.Handle("/", fs)
	http.Handle("/socket.io/", server)

	p := 8001
	log.Printf("Serving at localhost:%d", p)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", p), nil))
}
