package ws

import (
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/net/websocket"

	"github.com/Tutuacs/pkg/logs"
	"github.com/Tutuacs/pkg/routes"
)

type client struct {
	conns  map[*websocket.Conn]bool
	Client map[*websocket.Conn]int64
}

type WsHandler struct {
	client
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
	topic_Func = []topic_function{
		{topic: "/hello", function: HandleHello},
	}
	if len(topic_Func) > 0 {
		logs.MessageLog("Initializing WebSocket Server Topics List...")
	}
}

func NewWsHandler() *WsHandler {
	return &WsHandler{
		client: client{
			conns: make(map[*websocket.Conn]bool),
		},
	}
}

func (h *WsHandler) BuildRoutes(router routes.Route) {
	// Defina um único endpoint para a comunicação por WebSocket
	router.NewWS("/ws", websocket.Handler(h.HandleWs))
}

var i int64 = 0

func (h *WsHandler) HandleWs(ws *websocket.Conn) {
	fmt.Println("New connection from Client:", ws.RemoteAddr())
	h.conns[ws] = true
	h.Client[ws] = i // Increment the user id
	i++

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
	ws.Write([]byte(fmt.Sprintf("Received data for /hello: %v", data)))

	// Exemplo de envio de mensagem para todos os conectados
	handler := NewWsHandler()
	broadcastMsg := Message{Topic: "/hello", Data: "Broadcasting to all clients"}
	handler.sendMessage(ws, broadcastMsg)
}

func (h *WsHandler) BroadcastMessage(msg Message) {
	for conn := range h.conns {
		err := h.sendMessage(conn, msg)
		if err != nil {
			fmt.Println("Error broadcasting message to:", conn.RemoteAddr(), err)
		}
	}
}

func (h *WsHandler) sendMessage(conn *websocket.Conn, msg Message) error {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = conn.Write(messageBytes)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
