package sse

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Tutuacs/pkg/resolver"
	"github.com/Tutuacs/pkg/routes"
)

type SseHandler struct {
	conns map[http.ResponseWriter]int64
	close chan http.ResponseWriter // Channel to handle closed connections
}

var sseHandler *SseHandler

func init() {
	sseHandler = nil
}

func NewSseHandler() *SseHandler {
	if sseHandler == nil {
		sseHandler = &SseHandler{
			conns: make(map[http.ResponseWriter]int64),
			close: make(chan http.ResponseWriter),
		}
		// Start a goroutine to listen for closed connections
		go sseHandler.cleanupConnections()
	}
	return sseHandler
}

func (h *SseHandler) BuildRoutes(router routes.Route) {
	router.NewRoute(routes.ANY, "/sse", h.onConnectSse)
	router.NewRoute(routes.GET, "/sse/message", h.NewMessage)
}

var i int64 = 0

func (h *SseHandler) onConnectSse(w http.ResponseWriter, r *http.Request) {
	resolver.MakeSseRoute(w)

	h.conns[w] = i
	i++

	fmt.Println("New connection")

	// Listen and send messages until the connection is closed
	for i := 0; i <= 100; i++ {
		err := h.SendMessage(w, fmt.Sprintf("Hello %d", i))
		if err != nil {
			fmt.Println("Client disconnected")
			h.close <- w // Notify the cleanup goroutine of the closed connection
			return       // Exit the loop and stop sending messages
		}
		fmt.Println("Message sent")
		time.Sleep(5 * time.Second)
	}
	// ! the client will reconnect after 3 seconds by default if the connection get closed
}

func (h *SseHandler) NewMessage(w http.ResponseWriter, r *http.Request) {
	message := resolver.GetQueryParam(r, "message")
	if message == "" {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}
	err := h.SendBroadcastMessage(message)
	if err != nil {
		fmt.Printf("Error sending message: %v\n", err)
	}

}

func (h *SseHandler) SendMessage(w http.ResponseWriter, message string) error {
	// TODO: Write the message starting with "data: %v\n\n"
	data := fmt.Sprintf("data: %s\n\n", message)
	_, err := w.Write([]byte(data))
	if err != nil {
		// If an error occurs, remove the connection from the conns map
		// Probably the user is not connected anymore
		// * you can remove the channel that verifyes the connection if you do that
		delete(h.conns, w)
		return fmt.Errorf("error writing to a closed connection: %v", err)
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}

func (h *SseHandler) SendBroadcastMessage(message string) (erro error) {
	data := fmt.Sprintf("data: %s\n\n", message)
	erro = nil
	for w := range h.conns {
		_, err := w.Write([]byte(data))
		if err != nil {
			erro = err
			continue
		}
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}

	}
	return
}

// cleanupConnections listens for closed connections and removes them from the conns map
func (h *SseHandler) cleanupConnections() {
	for w := range h.close {
		delete(h.conns, w)
		fmt.Println("Connection removed from conns map")
	}
}

func (h *SseHandler) GetConn(id int64) http.ResponseWriter {
	for w, i := range h.conns {
		if i == id {
			return w
		}
	}
	return nil
}
