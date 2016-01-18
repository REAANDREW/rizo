package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.WriteHeader(http.StatusCreated)
	} else {
		io.WriteString(w, "Hello world!")
	}
}

type pathHandler struct {
	handlers map[string]http.HandlerFunc
}

func newPathHandler() *pathHandler {
	return &pathHandler{
		handlers: map[string]http.HandlerFunc{},
	}
}

func (instance *pathHandler) Get(handler http.HandlerFunc) {
	instance.handlers["GET"] = handler
}

func (instance *pathHandler) Handle(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)
	if handler, ok := instance.handlers[method]; ok {
		handler(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

//HTTPServer ...
type HTTPServer struct {
	Port     uint
	listener net.Listener
	server   *http.Server
	mux      *http.ServeMux
	paths    map[string]*pathHandler
}

//NewHTTPServer ...
func NewHTTPServer(port uint) *HTTPServer {
	return &HTTPServer{
		Port:  port,
		paths: map[string]*pathHandler{},
	}
}

//Start ...
func (instance *HTTPServer) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", instance.Port))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := instance.paths[r.URL.Path]; !ok {
			w.WriteHeader(http.StatusNotFound)
		} else {
			instance.paths[r.URL.Path].Handle(w, r)
		}
	})

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", instance.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err != nil {
		return err
	}
	instance.listener = l
	instance.mux = mux
	instance.server = s

	go func(listener net.Listener) {
		s.Serve(listener)
	}(l)

	return nil
}

//Stop ...
func (instance *HTTPServer) Stop() {
	if instance.listener != nil {
		instance.listener.Close()
	}
}

//Get ...
func (instance *HTTPServer) Get(path string, handler http.HandlerFunc) {
	if _, ok := instance.paths[path]; !ok {
		instance.paths[path] = newPathHandler()
	}
	instance.paths[path].Get(handler)
}
