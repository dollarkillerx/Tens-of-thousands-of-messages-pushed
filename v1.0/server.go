package main

import (
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter,r *http.Request,p httprouter.Params)  {
	conn, e := upgrader.Upgrade(w, r, nil) // 握手应答
	defer conn.Close()
	if e != nil {
		panic(e.Error())
	}

	// websocket.Conn
	for {
		// 支持数据体
		// Text,Binary
		_, data, e := conn.ReadMessage()
		if e != nil {
			panic(e.Error())
		}
		// 发送消息
		e = conn.WriteMessage(websocket.TextMessage, data)
		if e != nil {
			panic(e.Error())
		}

	}
}

func RegisterRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET("/ws",wsHandler)

	return router
}

func main() {
	//http.HandleFunc("/ws",wsHandler)

	router := RegisterRouter()

	http.ListenAndServe(":8580",router)
}