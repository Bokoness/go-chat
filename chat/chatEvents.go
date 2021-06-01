package chat

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bokoness/go-chat/auth"
	"github.com/bokoness/go-chat/rdb"
	socketio "github.com/googollee/go-socket.io"
)

type UserMessage struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func CreateChatEvents(server *socketio.Server) {
	red := rdb.CreateCon()
	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) {
		nsp := auth.GetToken(s, "nsp")
		usrMsg := MsgToJson(msg)
		red.RPush("chatRoom", msg)
		red.Expire("chatRoom", time.Duration(time.Hour*2))
		room := fmt.Sprintf("%s/chatRoom", nsp)
		server.BroadcastToRoom("/", room, "msg", usrMsg)
	})

	//this event shows all connection that user is typing, exept the typing user himself
	server.OnEvent("/", "typing", func(s socketio.Conn, usr string) {
		room := fmt.Sprintf("%s/chatRoom", auth.GetToken(s, "nsp"))
		s.Leave(room)
		server.BroadcastToRoom("/", room, "typing", usr)
		s.Join(room)
	})
}

func MsgToJson(msg string) UserMessage {
	var usrMsg UserMessage
	json.Unmarshal([]byte(msg), &usrMsg)
	return usrMsg
}