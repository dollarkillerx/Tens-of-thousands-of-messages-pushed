package main

import (
	"Tens-of-thousands-of-messages-pushed/v2.0/impl"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
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

	connection, e := impl.InitConnection(conn)
	if e != nil {
		panic(e.Error())
	}

	go func() {
		for  {
			err := connection.WriteMessage([]byte("heartbeat"))
			if err!=nil{
				return
			}
			time.Sleep(time.Second)
		}
	}()

	for {
		data, _ := connection.ReadMessage()
		connection.WriteMessage(data)
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