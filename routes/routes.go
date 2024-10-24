package routes

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

type Route struct {
	Router *http.ServeMux
}

func NewRouter() Route {
	return Route{
		Router: http.NewServeMux(),
	}
}

type Method string

const (
	POST   Method = "POST"
	GET    Method = "GET"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	ANY    Method = ""
)

func (r *Route) NewRoute(method Method, route string, function http.HandlerFunc) {

	url := fmt.Sprintf("%s %s", method, route)

	r.Router.HandleFunc(url, function)
}

func (r *Route) NewWS(topic string, function websocket.Handler) {

	r.Router.Handle(topic, function)
}
