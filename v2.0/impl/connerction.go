package impl

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Connection struct {
	wsConn *websocket.Conn
	inChan chan []byte // 存储接收消息
	outChan chan []byte // 存储返送消息
	closeChan chan byte // 关闭通知
}

// 初始化长连接
func InitConnection(wsConn *websocket.Conn) (conn *Connection,err error) {
	conn = &Connection{
		wsConn:wsConn,
		inChan:make(chan []byte,1000),
		outChan:make(chan []byte,1000),
		closeChan:make(chan byte,1), // 定义无阻塞通知
	}

	// 启动读协程
	go conn.readLoop()
	// 启动写协程
	go conn.writeLoop()

	return
}

// 获取消息
func (conn *Connection) ReadMessage() (data []byte,err error) {
	data = <-conn.inChan
	return
}

// 发送消息
func (conn *Connection) WriteMessage(data []byte) (err error) {
	conn.outChan <-data
	return
}

// 关闭链接
func (conn *Connection) Close() {
	// 线程安全,可重入的close
	conn.wsConn.Close()


}

// 内部实现

// 读实现
func (conn *Connection) readLoop() {
	for {
		_, p, err := conn.wsConn.ReadMessage()
		if err != nil {
			conn.Close()
			fmt.Println("read message error !!!")
		}
		// 阻塞在这里,等待inChan有空闲的位置
		select {
		case conn.inChan<-p:
		case <-conn.closeChan:
			conn.Close()
		}
	}
}

// 写实现
func (conn *Connection) writeLoop() {
	for {
		select {
		case data := <-conn.outChan :
			err := conn.wsConn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				conn.Close()
				fmt.Println("writer message error !!!")
			}
		case <-conn.closeChan:
			conn.Close()
		}

	}
}