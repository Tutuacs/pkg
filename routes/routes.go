package routes

import (
	"fmt"
	"net/http"
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

	if method == ANY {
		r.Router.HandleFunc(route, function)
		return
	}

	url := fmt.Sprintf("%s %s", method, route)

	r.Router.HandleFunc(url, function)
}
