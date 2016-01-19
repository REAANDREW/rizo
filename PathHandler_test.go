package rizo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var responseWriter = &FakeResponseWriter{}

func TestHandleGet(t *testing.T) {

	//Arrange
	const expectedMessage string = "handled the GET"
	defer responseWriter.Reset()
	handler := NewPathHandler()
	handler.Get(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expectedMessage))
	})

	//Act
	handler.Handle(responseWriter, &http.Request{
		Method: "GET",
	})

	//Assert
	assert.Equal(t, string(responseWriter.Data), expectedMessage)
}

func TestHandlePost(t *testing.T) {

	//Arrange
	const expectedMessage string = "handled the POST"
	defer responseWriter.Reset()
	handler := NewPathHandler()

	handler.Post(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expectedMessage))
	})

	//Act
	handler.Handle(responseWriter, &http.Request{
		Method: "POST",
	})

	//Assert
	assert.Equal(t, string(responseWriter.Data), expectedMessage)

}

func TestHandlePut(t *testing.T) {

	//Arrange
	const expectedMessage string = "handled the PUT"
	defer responseWriter.Reset()
	handler := NewPathHandler()

	handler.Put(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expectedMessage))
	})

	//Act
	handler.Handle(responseWriter, &http.Request{
		Method: "PUT",
	})

	//Assert
	assert.Equal(t, string(responseWriter.Data), expectedMessage)

}
