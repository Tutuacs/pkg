package adapters

import (
	"net/http"

	"golang.org/x/net/websocket"
)

func WebSocketGuarded(handler websocket.Handler, guard http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Executa o guard para verificar a autenticação
		guard(w, r)

		// Agora que o guard passou, inicia o WebSocket se a conexão não foi negada
		if w.Header().Get("X-Unauthorized") == "" { // Exemplo para verificar se a resposta foi rejeitada
			handler.ServeHTTP(w, r)
		}
	}
}
