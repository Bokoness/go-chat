package main

import (
	"context"
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

type Token struct {
	Nsp     string `json:"nsp"`
	Auth    string `json:"auth"`
	IsAdmin bool   `json:"isAdmin"`
}

type Values struct {
	m map[string]string
}

func (v Values) Get(key string) string {
	return v.m[key]
}

func main() {
	red := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	server := socketio.NewServer(nil)
	server.OnConnect("/", func(s socketio.Conn) error {
		fmt.Println("joined")
		server.JoinRoom("/", "setRoom", s)
		// server.JoinRoom("/", "chatRoom", s)
		// server.JoinRoom("/", "typing", s)

		// chatData, _ := red.LRange("chatRoom", 0, red.LLen("chatRoom").Val()).Result()
		// server.BroadcastToRoom("/", "chatRoom", "chatData", chatData)
		return nil
	})

	server.OnEvent("/", "setRoom", func(s socketio.Conn, d string) {
		s.SetContext(d)
		arr := map[string]string{}
		json.Unmarshal([]byte(d), &arr)
		v := Values{map[string]string{
			"auth":    arr["auth"],
			"isAdmin": arr["isAdmin"],
			"nsp":     arr["nsp"],
		}}
		s.SetContext(v)
		nsp := getToken(s, "nsp")
		server.JoinRoom("/", fmt.Sprintf("%s/chatRoom", nsp), s)
		server.JoinRoom("/", fmt.Sprintf("%s/typing", nsp), s)
		server.LeaveRoom("/", "setRoom", s)
	})

	//this handles the recive message event
	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) {
		nsp := getToken(s, "nsp")
		usrMsg := msgToJson(msg)
		red.RPush("chatRoom", msg)
		red.Expire("chatRoom", time.Duration(time.Hour*2))
		room := fmt.Sprintf("%s/chatRoom", nsp)
		server.BroadcastToRoom("/", room, "msg", usrMsg)
	})

	//this event shows all connection that user is typing, exept the typing user himself
	server.OnEvent("/", "typing", func(s socketio.Conn, usr string) {
		room := fmt.Sprintf("%s/chatRoom", getToken(s, "nsp"))
		s.Leave(room)
		server.BroadcastToRoom("/", room, "typing", usr)
		s.Join(room)
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason, s.Namespace())
	})

	//handling all sockets in differnt channels
	go server.Serve()
	defer server.Close()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	//this route will give cookies to user
	http.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("giving cookies")
		// var t Token
		// _ = json.NewDecoder(r.Body).Decode(&t)
		// c := http.Cookie{Name: "nsp", Value: t.Nsp, Path: "/"}
		// http.SetCookie(w, &c)
		// w.WriteHeader(200)
		w.Write([]byte("Good"))
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
func getToken(s socketio.Conn, v string) string {
	ctx := s.Context()
	c := context.Background()
	c2 := context.WithValue(c, "token", ctx)
	c3 := c2.Value("token").(Values).Get(v)
	return c3
}
