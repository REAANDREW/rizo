package rizo

import (
	"net/http"
	"testing"
)

type FakeResponseWriter struct {
	Data []byte
}

func (instance *FakeResponseWriter) Header() http.Header {
	return nil
}

func (instance *FakeResponseWriter) Write(data []byte) (int, error) {
	return 0, nil
}

func (instance *FakeResponseWriter) WriteHeader(int) {

}

func TestHandleGet(t *testing.T) {

	const expectedMessage string = "handled the get"
	responseWriter := &FakeResponseWriter{}
	handler := NewPathHandler()

	handler.Get(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expectedMessage))
	})

	handler.Handle(responseWriter, &http.Request{})
}
