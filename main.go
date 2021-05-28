package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

type UserMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func main() {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		server.JoinRoom("/", "chatRoom", s)
		return nil
	})

	//this handles the recive message event
	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) {
		usrMsg := msgToJson(msg)
		server.BroadcastToRoom("/", "chatRoom", "msg", usrMsg)
		log.Println("Message recived from client: " + usrMsg.Message)
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.Handle("/socket.io/", server)
	p := 8000
	log.Printf("Serving at localhost:%d", p)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", p), nil))
}

func msgToJson(msg string) UserMessage {
	var usrMsg UserMessage
	json.Unmarshal([]byte(msg), &usrMsg)
	return usrMsg
}
