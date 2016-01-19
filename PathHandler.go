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
	instance.Method("GET", handler)
}

//Post ...
func (instance *PathHandler) Post(handler http.HandlerFunc) {
	instance.Method("POST", handler)
}

//Put ...
func (instance *PathHandler) Put(handler http.HandlerFunc) {
	instance.Method("PUT", handler)
}

//Delete ...
func (instance *PathHandler) Delete(handler http.HandlerFunc) {
	instance.Method("DELETE", handler)
}

//Patch ...
func (instance *PathHandler) Patch(handler http.HandlerFunc) {
	instance.Method("PATCH", handler)
}

//Method ...
func (instance *PathHandler) Method(method string, handler http.HandlerFunc) {
	upperCaseMethod := strings.ToUpper(method)
	instance.handlers[upperCaseMethod] = handler
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
