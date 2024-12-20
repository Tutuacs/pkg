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
	// ! Dont remove this nil initialization
	wsHandler = nil
	// * Its to prevent loose the conns map

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
	// ? Whant to use validation?
	// router.NewRoute(routes.ANY, "/ws", guards.AuthenticatedUrlRoute(h.onConnectWs))
	// TODO: On client use like this: ws = new WebSocket(`ws://localhost:9000/ws?token=${encodeURIComponent(token)}`);
	// ! You can change the way you get the token on guards/middlewares.go
}

func (h *WsHandler) onConnectWs(w http.ResponseWriter, r *http.Request) {
	opts := &websocket.AcceptOptions{
		InsecureSkipVerify: true,          // Skip verifying the origin in production
		OriginPatterns:     []string{"*"}, // Allow all origins, or specify particular origins like "127.0.0.1:5500"
	}

	// ? Using Validation for the WebSocket connection?
	// userLogged := r.Context().Value(guards.UserKey).(*types.User)
	// TODO: get the users info to use on the conns map!

	conn, err := websocket.Accept(w, r, opts)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}

	h.conns[conn] = i
	i++

	h.readLoop(r.Context(), conn)
}

func (h *WsHandler) readLoop(ctx context.Context, conn *websocket.Conn) {
	for {
		msgType, buff, err := conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
				websocket.CloseStatus(err) == websocket.StatusGoingAway ||
				websocket.CloseStatus(err) == websocket.StatusAbnormalClosure {
				fmt.Println("Connection closed normally:", h.conns[conn])
				break
			}
			fmt.Println("Error reading:", err)
			fmt.Println(msgType, buff, err)
			break
		}

		if msgType != websocket.MessageText {
			fmt.Println("Ignoring non-text message")
			continue
		}

		var msg types.WsMessage
		err = json.Unmarshal(buff, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		h.handleMessage(ctx, conn, msg)
	}
	delete(h.conns, conn)
}

func (h *WsHandler) handleMessage(ctx context.Context, conn *websocket.Conn, msg types.WsMessage) {
	for _, tf := range *topic_Func {
		if tf.topic == msg.Topic {
			tf.function(ctx, conn, msg)
			return
		}
	}
}

func (h *WsHandler) BroadcastMessage(ctx context.Context, msg types.WsMessage) {
	for conn := range h.conns {
		msg := types.WsMessage{Topic: "/hello unknown", Data: h.conns[conn]}
		err := h.sendMessage(ctx, conn, msg)
		if err != nil {
			fmt.Println("Error broadcasting message to:", h.conns[conn], err)
		}
	}
}

func (h *WsHandler) sendMessage(ctx context.Context, conn *websocket.Conn, msg types.WsMessage) error {
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

func (h *WsHandler) GetConn(id int64) *websocket.Conn {
	for w, i := range h.conns {
		if i == id {
			return w
		}
	}
	return nil
}

// Recomended to use private functions on topics_functions
// ! Use public functions only if you need to use it on another package

func HandleHello(ctx context.Context, conn *websocket.Conn, data types.WsMessage) {
	fmt.Println("Handling /hello with data:", data)

	helloMsg := types.WsMessage{Topic: "/hello", Data: "simple hello to client"}
	NewWsHandler().sendMessage(ctx, conn, helloMsg)
}

func HandleUnknown(ctx context.Context, conn *websocket.Conn, data types.WsMessage) {
	fmt.Println("Handling /unknown with data:", data)

	broadcastMsg := types.WsMessage{}
	NewWsHandler().BroadcastMessage(ctx, broadcastMsg)
}
