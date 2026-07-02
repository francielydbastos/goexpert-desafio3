package web

import "net/http"

type WebServer struct {
	Handler *OrderHandler
}

func NewWebServer(handler *OrderHandler) *WebServer {
	return &WebServer{Handler: handler}
}

func (s *WebServer) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /order", s.Handler.Create)
	mux.HandleFunc("GET /order", s.Handler.List)
	return http.ListenAndServe(addr, mux)
}
