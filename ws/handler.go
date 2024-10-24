package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coder/websocket"

	"github.com/Tutuacs/pkg/logs"
	"github.com/Tutuacs/pkg/routes"
	"github.com/Tutuacs/pkg/types"
)

var i int64 = 0

func init() {
	wsHandler = nil
	topic_Func = &[]topic_function{
		{topic: "/hello", function: HandleHello},
		{topic: "/unknown", function: HandleUnknown},
	}

	if len(*topic_Func) > 0 {
		logs.MessageLog("Initializing WebSocket Server Topics List...")
	}
}

func NewWsHandler() *WsHandler {
	if wsHandler == nil {
		wsHandler = &WsHandler{
			conns: make(map[*websocket.Conn]int64),
		}
	}
	return wsHandler
}

func (h *WsHandler) BuildRoutes(router routes.Route) {
	router.NewRoute(routes.ANY, "/ws", h.onConnectWs)
}

func (h *WsHandler) onConnectWs(w http.ResponseWriter, r *http.Request) {
	opts := &websocket.AcceptOptions{
		OriginPatterns: []string{"*"}, // Allow all origins, or specify particular origins like "127.0.0.1:5500"
	}
	conn, err := websocket.Accept(w, r, opts)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}

	h.conns[conn] = i
	i++

	conn.Write(r.Context(), websocket.MessageText, []byte("Connected to /ws"))

	go h.readLoop(r.Context(), conn)

}

func (h *WsHandler) readLoop(ctx context.Context, conn *websocket.Conn) {
	defer func() {
		fmt.Println("Cleaning up connection:", h.conns[conn])
		delete(h.conns, conn)
		conn.Close(websocket.StatusNormalClosure, "Connection closed by server")
	}()

	for {
		// Attempt to read message from WebSocket connection
		msgType, buff, err := conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway ||
				websocket.CloseStatus(err) == websocket.StatusAbnormalClosure {
				fmt.Println("Connection closed normally:", h.conns[conn])
				break
			}
			fmt.Println("Error reading:", err)
			break
		}

		if msgType != websocket.MessageText {
			fmt.Println("Ignoring non-text message")
			continue
		}

		// Deserialize the message to identify the topic
		var msg types.Message
		err = json.Unmarshal(buff, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		h.handleMessage(ctx, conn, msg)
	}
}

func (h *WsHandler) handleMessage(ctx context.Context, conn *websocket.Conn, msg types.Message) {
	for _, tf := range *topic_Func {
		if tf.topic == msg.Topic {
			tf.function(ctx, conn, msg)
			return
		}
	}
}

func HandleHello(ctx context.Context, conn *websocket.Conn, data types.Message) {
	fmt.Println("Handling /hello with data:", data)

	helloMsg := types.Message{Topic: "/hello", Data: "simple hello to client"}
	NewWsHandler().sendMessage(ctx, conn, helloMsg)
}

func HandleUnknown(ctx context.Context, conn *websocket.Conn, data types.Message) {
	fmt.Println("Handling /unknown with data:", data)

	broadcastMsg := types.Message{}
	NewWsHandler().BroadcastMessage(ctx, broadcastMsg)
}

func (h *WsHandler) BroadcastMessage(ctx context.Context, msg types.Message) {
	for conn := range h.conns {
		msg := types.Message{Topic: "/hello unknown", Data: h.conns[conn]}
		err := h.sendMessage(ctx, conn, msg)
		if err != nil {
			fmt.Println("Error broadcasting message to:", h.conns[conn], err)
		}
	}
}

func (h *WsHandler) sendMessage(ctx context.Context, conn *websocket.Conn, msg types.Message) error {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = conn.Write(ctx, websocket.MessageText, messageBytes)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
