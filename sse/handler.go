package sse

import (
	"net/http"

	"github.com/Tutuacs/pkg/resolver"
	"github.com/Tutuacs/pkg/routes"
)

type SseHandler struct {
	conn map[http.ResponseWriter]int64
}

var sseHandler *SseHandler

func init() {
	sseHandler = nil
}

func NewSseHandler() *SseHandler {
	if sseHandler == nil {
		sseHandler = &SseHandler{
			conn: make(map[http.ResponseWriter]int64),
		}
	}
	return sseHandler
}

func (h *SseHandler) BuildRoutes(router routes.Route) {
	router.NewRoute(routes.ANY, "/sse", h.onConnectSse)
	// ? Whant to use validation?
	// router.NewRoute(routes.ANY, "/sse", guards.AuthenticatedUrlRoute(h.onConnectSse))
	// TODO: On client use like this: eventSource = new EventSource(`http://localhost:9000/sse?token=${encodeURIComponent(token)}`);
	// ! You can change the way you get the token on guards/middlewares.go
}

var i int64 = 0

func (h *SseHandler) onConnectSse(w http.ResponseWriter, r *http.Request) {

	// * make this route on SSE using the resolver:
	resolver.MakeSseRoute(w)
	// * Now its a SSE router

	// ? Want to send messages to specific users?
	// ! recomended to use auth validation to get the user info

	// ? Using Validation for the WebSocket connection?
	// userLogged := r.Context().Value(guards.UserKey).(*types.User)
	// TODO: get the users info to use on the conns map!

	h.conn[w] = i
	i++

}

func (h *SseHandler) SendMessage(w http.ResponseWriter, message interface{}) {
	w.Write(message.([]byte))
}

func (h *SseHandler) GetConn(id int64) http.ResponseWriter {
	for w, i := range h.conn {
		if i == id {
			return w
		}
	}
	return nil
}
