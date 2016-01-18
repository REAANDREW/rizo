package rizo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ServerPort = 5000

func TestAddingRouteToHttpServer(t *testing.T) {
	const ResponseMessage string = "Hello World!"
	server := NewHTTPServer(ServerPort)
	defer server.Stop()
	client := &http.Client{}

	server.Get("/something", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(ResponseMessage))
	})

	server.Start()

	request, _ := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/something", ServerPort), nil)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)

	assert.Equal(t, string(content), ResponseMessage)
}
