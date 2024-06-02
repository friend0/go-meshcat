package internal

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var writeWait = 10 * time.Second
var pongWait = 60 * time.Second
var pingPeriod = (pongWait * 9) / 10

// Define the WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	WriteBufferPool: &sync.Pool{},
}

func wsprocess(ws *websocket.Conn) {
	// defer ws.Close()
	for {

		err := ws.WriteMessage(websocket.BinaryMessage, []byte("Hello!"))
		if err != nil {
			log.Println("Write Err: ", err)
		}
		// Read message from WebSocket
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read:", err, messageType)
			break
		}
		log.Printf("Received: %s", message)

		// Echo the message back to the client
		err = ws.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Write:", err)
			break
		}
	}
}

// func ping(ws *websocket.Conn, done chan struct{}) {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
// 			if err := ws.P(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
// 				log.Println("pingers", err)
// 			}
// 		case <-done:
// 			fmt.Println("no more pingers")
// 			return
// 		}
// 	}
// }

func (s *Server) WsHandler() echo.HandlerFunc {
	// Define a handler for WebSocket connections
	return func(c echo.Context) error {
		// Upgrade the HTTP connection to a WebSocket connection
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Println("Upgrade:", err)
			return err
		}
		ws.SetWriteDeadline(time.Now().Add(pongWait))
		ws.SetPongHandler(func(string) error { ws.SetWriteDeadline(time.Now().Add(pongWait)); fmt.Println("pongers"); return nil })
		ws.SetPingHandler(func(appData string) error {
			ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait))
			return nil
		})
		s.WS = ws
		go wsprocess(s.WS)
		// go ping(s.WS, stdoutDone)
		return nil
	}
}
