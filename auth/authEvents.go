package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bokoness/go-chat/rdb"
	socketio "github.com/googollee/go-socket.io"
)

type Token struct {
	Room    string `json:"room"`
	Auth    string `json:"auth"`
	IsAdmin bool   `json:"isAdmin"`
}

type AdminData struct {
	Room      string `json:"room"`
	AdminName string `json:"adminName"`
	Token     string `json:"token"`
}

type Values struct {
	m map[string]string
}

func (v Values) Get(key string) string {
	return v.m[key]
}

func CreateAuthEvents(server *socketio.Server) {
	server.OnEvent("/", "auth", func(s socketio.Conn, d string) {
		red := rdb.Client
		s.SetContext(d)
		arr := map[string]string{}
		json.Unmarshal([]byte(d), &arr)
		v := Values{map[string]string{
			"auth":    arr["auth"],
			"isAdmin": arr["isAdmin"],
			"room":    arr["room"],
		}}
		s.SetContext(v)
		room := GetToken(s, "room")
		fmt.Println(red.HExists(room, "status").Val())
		if red.HExists(room, "status").Val() {
			status := red.HGet(room, "status").Val()
			fmt.Println(status)
			chatRoom := fmt.Sprintf("%s/chatRoom", room)
			server.JoinRoom("/", chatRoom, s)
			server.JoinRoom("/", fmt.Sprintf("%s/typing", room), s)
			chatData, _ := rdb.Client.LRange(chatRoom, 0, rdb.Client.LLen(chatRoom).Val()).Result()
			server.BroadcastToRoom("/", chatRoom, "chatData", chatData)
		}
		server.LeaveRoom("/", "setRoom", s)
	})

	server.OnEvent("/", "joinRoom", func(s socketio.Conn, data string) {
		var r map[string]string
		json.Unmarshal([]byte(data), &r)
		red := rdb.Client
		if !red.HExists(r["room"], "name").Val() {
			fmt.Println("BAD")
			server.LeaveRoom("/", "auth", s)
			return
		} else {
			//add user to room users list
			uCol := r["room"] + "/players"
			red.HSet(uCol, string(s.ID()), r["name"])
			usrs := red.HGetAll(uCol)
			//add user to general users-rooms list
			red.HSet("players-rooms", string(s.ID()), r["room"])
			s.Join(r["room"])
			server.BroadcastToRoom("/", r["room"], "joined", usrs.Val())
		}
	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		red := rdb.Client
		var b AdminData
		e := json.NewDecoder(r.Body).Decode(&b)
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(b)
		if !validateToken(b.Token) {
			http.Error(w, e.Error(), http.StatusUnauthorized)
			return
		}
		if red.HExists(b.Room, "status").Val() {
			http.Error(w, "Room is taken", http.StatusUnauthorized)
			return
		}
		// res := red.Set(b.Nsp, true, 2*time.Hour)
		var m = make(map[string]interface{})
		m["name"] = b.Room
		m["admin"] = b.Room
		m["status"] = "waiting"
		red.HMSet(b.Room, m)
		red.Expire(b.Room, time.Duration(time.Hour*2))
		// w.WriteHeader(200)
		w.WriteHeader(200)
	})
}

func SetCtx(s socketio.Conn, t string) {
	arr := map[string]string{}
	json.Unmarshal([]byte(t), &arr)
	v := Values{map[string]string{
		"auth":    arr["auth"],
		"isAdmin": arr["isAdmin"],
		"room":    arr["room"],
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

//TODO: make propper validation
func validateToken(t string) bool {
	return t != ""
}
