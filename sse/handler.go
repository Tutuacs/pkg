package sse

import "net/http"

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

func (h *SseHandler) BuildRoutes() {
	http.HandleFunc("/sse", h.onConnectSse)
}

func (h *SseHandler) onConnectSse(w http.ResponseWriter, r *http.Request) {
	// ! To manage the differente connections synthesizing IDs

}
