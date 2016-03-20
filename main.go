package main

import (
	"fmt"

	"github.com/atsman/interviewr-go/db"
	"github.com/atsman/interviewr-go/handlers"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

const (
	Port = "3000"
)

func init() {
	db.Connect()
}

var Socketio_Server *socketio.Server

type Message struct {
	UserId    string `json:"userId"`
	Username  string `json:"username"`
	Time      string `json:"time"`
	Message   string `json:"message"`
	RoomID    string `json:"roomId"`
	UserImage string `json:"userImage"`
}

type CodeSharingState struct {
	Lang   string `json:"lang"`
	Code   string `json:"code"`
	RoomId string `json:"roomId"`
	Cursor Cursor `json:"cursor"`
}

type Cursor struct {
	Column int `json:"column"`
	Row    int `json:"row"`
}

func socketHandler(c *gin.Context) {
	Socketio_Server.On("connection", func(so socketio.Socket) {
		fmt.Println("on connection")

		so.On("joinRoom", func(roomId string) {
			log.Debugf("joinRoom, roomId: %s", roomId)
			so.Join(roomId)
		})

		so.On("sendMessage", func(message Message) {
			log.Debug("sendMessage, message: ", message)
			so.BroadcastTo(message.RoomID, "newMessage", message)
		})

		so.On("sendCode", func(code CodeSharingState) {
			log.Debug("sendCode, codeState:", code)
			so.BroadcastTo(code.RoomId, "receiveCodeChange", code)
		})

		so.On("disconnection", func() {
			fmt.Println("on disconnect")
		})
	})

	Socketio_Server.On("error", func(so socketio.Socket, err error) {
		fmt.Printf("[ WebSocket ] Error : %v", err.Error())
	})

	Socketio_Server.ServeHTTP(c.Writer, c.Request)
}

func main() {
	r := handlers.NewEngine()

	var err error
	Socketio_Server, err = socketio.NewServer(nil)
	if err != nil {
		panic(err)
	}

	r.GET("/socket.io", socketHandler)
	r.POST("/socket.io", socketHandler)
	r.Handle("WS", "/socket.io", socketHandler)
	r.Handle("WSS", "/socket.io", socketHandler)

	port := Port
	r.Run(":" + port)
}
