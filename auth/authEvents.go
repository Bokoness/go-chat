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
	Nsp     string `json:"nsp"`
	Auth    string `json:"auth"`
	IsAdmin bool   `json:"isAdmin"`
}

type AdminData struct {
	Nsp   string `json:"nsp"`
	Name  string `json:"name"`
	Token string `json:"token"`
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
		chatRoom := fmt.Sprintf("%s/chatRoom", nsp)
		server.JoinRoom("/", chatRoom, s)
		server.JoinRoom("/", fmt.Sprintf("%s/typing", nsp), s)
		chatData, _ := rdb.Client.LRange(chatRoom, 0, rdb.Client.LLen(chatRoom).Val()).Result()
		server.BroadcastToRoom("/", chatRoom, "chatData", chatData)
		server.LeaveRoom("/", "setRoom", s)
	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		red := rdb.Client
		var b AdminData
		e := json.NewDecoder(r.Body).Decode(&b)
		if e != nil {
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		if !validateToken(b.Token) {
			http.Error(w, e.Error(), http.StatusUnauthorized)
			return
		}
		// res := red.Set(b.Nsp, true, 2*time.Hour)
		var m = make(map[string]interface{})
		m["room"] = b.Nsp
		m["admin"] = b.Name
		m["status"] = "waiting"
		red.HMSet(b.Nsp, m)
		red.Expire(b.Nsp, time.Duration(time.Hour*2))
		w.WriteHeader(200)
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

//TODO: make propper validation
func validateToken(t string) bool {
	return t != ""
}
