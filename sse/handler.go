package sse

import (
	"net/http"

	"github.com/Tutuacs/pkg/resolver"
	"github.com/Tutuacs/pkg/routes"
)

type SseHandler struct {
	conn map[*http.ResponseWriter]int64
}

var sseHandler *SseHandler

func init() {
	sseHandler = nil
}

func NewSseHandler() *SseHandler {
	if sseHandler == nil {
		sseHandler = &SseHandler{
			conn: make(map[*http.ResponseWriter]int64),
		}
	}
	return sseHandler
}

func (h *SseHandler) BuildRoutes(router routes.Route) {
	router.NewRoute(routes.ANY, "/sse", h.onConnectSse)
}

func (h *SseHandler) onConnectSse(w http.ResponseWriter, r *http.Request) {
	// ! To manage the differente connections synthesizing IDs
	resolver.MakeSseRoute(w)

}
