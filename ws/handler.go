package ws

import (
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/net/websocket"

	"github.com/Tutuacs/pkg/logs"
	"github.com/Tutuacs/pkg/routes"
)

type WsHandler struct {
	conns map[*websocket.Conn]bool
}

type Message struct {
	Topic string      `json:"topic"`
	Data  interface{} `json:"data"`
}

type topic_function struct {
	topic    string
	function func(ws *websocket.Conn, data Message)
}

var topic_Func []topic_function

func init() {
	logs.MessageLog("Initializing topics on WebSocket...")
	topic_Func = []topic_function{
		{topic: "/hello", function: HandleHello},
	}
}

func NewWsHandler() *WsHandler {
	return &WsHandler{conns: make(map[*websocket.Conn]bool)}
}

func (h *WsHandler) BuildRoutes(router routes.Route) {
	// Defina um único endpoint para a comunicação por WebSocket
	router.NewWS("/ws", websocket.Handler(h.HandleWs))
}

func (h *WsHandler) HandleWs(ws *websocket.Conn) {
	fmt.Println("New connection from Client:", ws.RemoteAddr())
	h.conns[ws] = true

	ws.Write([]byte("Connected to /ws"))
	h.readLoop(ws)
}

func (h *WsHandler) readLoop(ws *websocket.Conn) {
	buff := make([]byte, 1024)
	for {
		n, err := ws.Read(buff)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed:", ws.RemoteAddr())
				break
			}
			fmt.Println("Error reading:", err)
			continue
		}

		// Deserialize the message to identify the topic
		var msg Message
		err = json.Unmarshal(buff[:n], &msg)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		fmt.Printf("Received: Topic: %s, Data: %v\n", msg.Topic, msg.Data)

		// Handle message based on the topic
		h.handleMessage(ws, msg)
	}
}

func (h *WsHandler) handleMessage(ws *websocket.Conn, msg Message) {
	for _, tf := range topic_Func {
		if tf.topic == msg.Topic {
			tf.function(ws, msg)
			return
		}
	}
}

func HandleHello(ws *websocket.Conn, data Message) {
	fmt.Println("Handling /hello with data:", data)
	ws.Write([]byte(fmt.Sprintf("Received data for t1: %v", data)))
}
