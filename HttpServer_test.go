package rizo

import (
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const ServerPort = 5000

var _ = Describe("HTTP Server", func() {
	It("Supports a route with a GET method", func() {

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

		Expect(string(content)).To(Equal(ResponseMessage))
	})
})
