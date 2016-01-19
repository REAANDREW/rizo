package rizo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleGet(t *testing.T) {

	const expectedMessage string = "handled the get"
	responseWriter := &FakeResponseWriter{}
	defer responseWriter.Reset()
	handler := NewPathHandler()

	handler.Get(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expectedMessage))
	})

	handler.Handle(responseWriter, &http.Request{
		Method: "GET",
	})

	assert.Equal(t, string(responseWriter.Data), expectedMessage)

}
