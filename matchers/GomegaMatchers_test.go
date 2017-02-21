package matchers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"fmt"
	"github.com/guzzlerio/rizo"
	"net/http"
)

var (
	//TestServer ...
	TestServer *rizo.RequestRecordingServer
	TestPort   = 6000
)

func URLForTestServer(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", TestPort, path)
}

var _ = Describe("GomegaMatchers", func() {
	var (
		request   *http.Request
		sampleURL string
	)
	const (
		data = "a=1&b=2"
	)

	BeforeEach(func() {
		TestServer = rizo.CreateRequestRecordingServer(TestPort)
		TestServer.Start()

		sampleURL = URLForTestServer("/Fubar?" + data)
		request, _ = http.NewRequest("GET", sampleURL, bytes.NewBuffer([]byte(data)))
		request.Header.Set("Content-type", "application/json")
		TestServer.Requests = append(TestServer.Requests, rizo.RecordedRequest{
			Request: request,
			Body:    data,
		})
	})

	AfterEach(func() {
		TestServer.Clear()
		TestServer.Stop()
	})
	Describe("Find", func() {
		It("matches when the predicate is successful", func() {
			Ω(TestServer).Should(Find(rizo.RequestWithPath("/Fubar")))
		})
		It("does not match when the predicate is unsuccessful", func() {
			Ω(TestServer).ShouldNot(Find(rizo.RequestWithPath("/talula")))
		})
	})
})
