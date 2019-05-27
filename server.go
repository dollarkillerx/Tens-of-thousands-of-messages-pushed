package main

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WsServer(w http.ResponseWriter,r *http.Request,p httprouter.Params) {
	conn, e := upgrader.Upgrade(w, r, nil)
	defer conn.Close()
	if e != nil {
		panic(e.Error())
	}


}

func RegisterRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/ws",WsServer)

	return router
}

func main() {
	router := RegisterRouter()

	http.ListenAndServe(":8580",router)
}