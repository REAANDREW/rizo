package rizo

import (
	"net/http"
	"strings"
)

//PathHandler ...
type PathHandler struct {
	handlers map[string]http.HandlerFunc
}

//NewPathHandler ...
func NewPathHandler() *PathHandler {
	return &PathHandler{
		handlers: map[string]http.HandlerFunc{},
	}
}

//Get ...
func (instance *PathHandler) Get(handler http.HandlerFunc) {
	instance.handlers["GET"] = handler
}

//Handle ...
func (instance *PathHandler) Handle(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)
	if handler, ok := instance.handlers[method]; ok {
		handler(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
