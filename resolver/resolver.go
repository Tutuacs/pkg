package resolver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

func GetBody(r *http.Request, response interface{}) (err error) {

	if r.Body == nil {
		return fmt.Errorf("missing body")
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(response)

	return
}

func GetParam(r *http.Request, name string) (param string) {
	param = r.PathValue(name)
	return
}

func GetQueryParam(r *http.Request, name string) (param string) {
	param = r.URL.Query().Get(name)
	return
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return strings.Split(tokenAuth, " ")[1]
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func WriteResponse(w http.ResponseWriter, status int, response interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	return encoder.Encode(response)
}

func MakeSseRoute(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Connection", "keep-alive")

	// Flush
	flush, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	flush.Flush()
}

func GetLoggedUser(r *http.Request, key string) any {
	return r.Context().Value(key)
}

var Validate = validator.New()
