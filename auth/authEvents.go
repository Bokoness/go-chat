package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bokoness/go-chat/rdb"
	socketio "github.com/googollee/go-socket.io"
)

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

func CreateAuthEvents(server *socketio.Server) {
	server.OnEvent("/", "auth", func(s socketio.Conn, d string) {
		s.SetContext(d)
		arr := map[string]string{}
		json.Unmarshal([]byte(d), &arr)
		v := Values{map[string]string{
			"auth":    arr["auth"],
			"isAdmin": arr["isAdmin"],
			"nsp":     arr["nsp"],
		}}
		s.SetContext(v)
		nsp := GetToken(s, "nsp")
		server.JoinRoom("/", fmt.Sprintf("%s/chatRoom", nsp), s)
		server.JoinRoom("/", fmt.Sprintf("%s/typing", nsp), s)
		chatData, _ := rdb.Client.LRange("chatRoom", 0, rdb.Client.LLen("chatRoom").Val()).Result()
		room := fmt.Sprintf("%s/chatRoom", nsp)
		server.BroadcastToRoom("/", room, "chatData", chatData)
		server.LeaveRoom("/", "setRoom", s)
	})
}

func SetCtx(s socketio.Conn, t string) {
	arr := map[string]string{}
	json.Unmarshal([]byte(t), &arr)
	v := Values{map[string]string{
		"auth":    arr["auth"],
		"isAdmin": arr["isAdmin"],
		"nsp":     arr["nsp"],
	}}
	s.SetContext(v)
}

func GetToken(s socketio.Conn, v string) string {
	ctx := s.Context()
	c1 := context.Background()
	c2 := context.WithValue(c1, "t", ctx)
	c3 := c2.Value("t").(Values).Get(v)
	return c3
}
