package ws

import (
	"context"

	"github.com/coder/websocket"

	"github.com/Tutuacs/pkg/types"
)

type topic_function struct {
	topic    string
	function func(ctx context.Context, ws *websocket.Conn, data types.Message)
}

var topic_Func *[]topic_function

var wsHandler *WsHandler

type WsHandler struct {
	conns map[*websocket.Conn]int64
}
