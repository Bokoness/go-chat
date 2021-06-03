package chat

import (
	"fmt"
	"net/http"

	"github.com/bokoness/go-chat/auth"
	socketio "github.com/googollee/go-socket.io"
)

func CreateWaitingEvents(server *socketio.Server) {
	server.OnEvent("/", "enter", func(s socketio.Conn, name string) {
		r := auth.GetToken(s, "room")
		//check if room is open (nsp)
		room := fmt.Sprintf("%s/waitingRoom", r)
		server.BroadcastToRoom("/", room, "enter", name)
	})

	http.HandleFunc("/waiting", func(w http.ResponseWriter, r *http.Request) {
		// red := rdb.Client
		// fmt.Println(r.Method)
		// keys, ok := r.URL.Query()["t"]
		// if !ok || len(keys[0]) < 1 {
		// 	http.Redirect(w, r, "/", 403)
		// 	return
		// }
		// t := keys[0]
		//check if room exists here
		// log.Println("Url Param 'key' is: " + string(t))
		w.WriteHeader(200)
	})
}
