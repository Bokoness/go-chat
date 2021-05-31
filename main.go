package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	socketio "github.com/googollee/go-socket.io"
)

type UserMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())
		server.JoinRoom("/", "chatRoom", s)
		server.JoinRoom("/", "typing", s)

		server.Adapter(&socketio.RedisAdapterOptions{})
		return nil
	})
	//this handles the recive message event
	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) {
		client.RPush("chatRoom", msg)
		b := client.Expire("chatRoom", time.Duration(time.Hour*2))
		fmt.Println(b)
		usrMsg := msgToJson(msg)
		server.BroadcastToRoom("/", "chatRoom", "msg", usrMsg)
	})

	//this event shows all connection that user is typing, exept the typing user himself
	server.OnEvent("/", "typing", func(s socketio.Conn, usr string) {
		s.Leave("typing")
		server.BroadcastToRoom("/", "typing", "typing", usr)
		s.Join("typing")
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	//handling all sockets in differnt channels
	go server.Serve()
	defer server.Close()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This func should handle the authentication of the socket"))
	})

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
