package chat

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bokoness/go-chat/auth"
	"github.com/bokoness/go-chat/rdb"
	socketio "github.com/googollee/go-socket.io"
)

func CreateWaitingEvents(server *socketio.Server) {
	server.OnEvent("/", "joined", func(s socketio.Conn, name string) {
		r := auth.GetToken(s, "room")
		//check if room is open (nsp)
		room := fmt.Sprintf("%s/waitingRoom", r)
		server.BroadcastToRoom("/", room, "join", name)
	})

	http.HandleFunc("/waiting", func(w http.ResponseWriter, r *http.Request) {
		e, keys := getUrlParams(r, []string{"r", "t"})
		if e != nil {
			http.Error(w, e.Error(), http.StatusUnauthorized)
			return
		}
		room := keys[0]
		// token := keys[1]
		b := rdb.Client.HExists(room, "name").Val()
		if !b {
			http.Error(w, e.Error(), http.StatusUnauthorized)
			return
		}
		dir, _ := os.Getwd()
		path := filepath.Join(dir, "static", "waitingRoom.html")
		fmt.Println(path)
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, path)
	})
}

func getUrlParams(r *http.Request, params []string) (error, []string) {
	var p []string
	for _, val := range params {
		keys, ok := r.URL.Query()[val]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'key' is missing")
			return errors.New("no such param"), nil
		}
		p = append(p, keys[0])
	}
	return nil, p
}
