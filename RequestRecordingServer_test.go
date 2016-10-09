package rizo

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	//TestServer ...
	TestServer *RequestRecordingServer
	TestPort   = 6000
)

func URLForTestServer(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", TestPort, path)
}

//HTTPRequestDo ...
func HTTPRequestDo(verb string, url string, bodyBuffer io.Reader, changeRequestDelegate func(request *http.Request)) (response *http.Response, body string, err error) {
	client := &http.Client{}
	request, err := http.NewRequest(verb, url, bodyBuffer)
	check(err)
	if changeRequestDelegate != nil {
		changeRequestDelegate(request)
	}
	response, err = client.Do(request)
	check(err)
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err == nil {
		body = string(bodyBytes)
	}
	return
}

var _ = Describe("RequestRecordingServer", func() {

	BeforeEach(func() {
		TestServer = CreateRequestRecordingServer(TestPort)
		TestServer.Start()
	})

	AfterEach(func() {
		TestServer.Clear()
		TestServer.Stop()
	})

	Describe("Find", func() {
		var sampleURL string
		var data string

		Describe("Single Request", func() {
			var request *http.Request

			BeforeEach(func() {
				data = "a=1&b=2"
				sampleURL = URLForTestServer("/Fubar?" + data)
				request, _ = http.NewRequest("GET", sampleURL, bytes.NewBuffer([]byte(data)))
				request.Header.Set("Content-type", "application/json")
				TestServer.Requests = append(TestServer.Requests, RecordedRequest{
					Request: request,
					Body:    data,
				})
			})

			It("Path", func() {
				expectedPath := "/Fubar"
				Expect(TestServer.Find(RequestWithPath(expectedPath))).To(Equal(true))
			})

			It("Method", func() {
				expectedMethod := "GET"
				Expect(TestServer.Find(RequestWithMethod(expectedMethod))).To(Equal(true))
			})

			It("Header", func() {
				Expect(TestServer.Find(RequestWithHeader("Content-type", "application/json"))).To(Equal(true))
			})

			It("Body", func() {
				request, _ = http.NewRequest("POST", sampleURL, bytes.NewBuffer([]byte(data)))
				TestServer.Clear()
				TestServer.Requests = append(TestServer.Requests, RecordedRequest{
					Request: request,
					Body:    data,
				})
				Expect(TestServer.Find(RequestWithBody(data))).To(Equal(true))
			})

			It("Querystring", func() {
				Expect(TestServer.Find(RequestWithQuerystring(data))).To(Equal(true))
			})

			It("Handles multiple predicates", func() {
				expectedPath := "/Fubar"
				expectedMethod := "GET"
				Expect(TestServer.Find(RequestWithPath(expectedPath), RequestWithMethod(expectedMethod))).To(Equal(true))
			})
		})

		Describe("Multiple Requests", func() {
			BeforeEach(func() {
				sampleURL = URLForTestServer("/Fubar")
				request, _ := http.NewRequest("GET", sampleURL, nil)
				request.Header.Set("Content-type", "application/json")
				TestServer.Requests = append(TestServer.Requests, RecordedRequest{
					Request: request,
				})

				data = "a=1&b=2"
				postRequest, _ := http.NewRequest("POST", sampleURL, bytes.NewBuffer([]byte(data)))
				TestServer.Requests = append(TestServer.Requests, RecordedRequest{
					Request: postRequest,
					Body:    data,
				})
			})

			It("Path", func() {
				expectedPath := "/Fubar"
				Expect(TestServer.Find(RequestWithPath(expectedPath))).To(Equal(true))
			})

			It("Method", func() {
				expectedMethod := "GET"
				Expect(TestServer.Find(RequestWithMethod(expectedMethod))).To(Equal(true))
			})

			It("Header", func() {
				Expect(TestServer.Find(RequestWithHeader("Content-type", "application/json"))).To(Equal(true))
			})

			It("Body", func() {
				Expect(TestServer.Find(RequestWithBody(data))).To(Equal(true))
			})

			It("Handles multiple predicates", func() {
				expectedPath := "/Fubar"
				expectedMethod := "GET"
				Expect(TestServer.Find(RequestWithPath(expectedPath), RequestWithMethod(expectedMethod))).To(Equal(true))
			})
		})
	})

	Describe("Response factory", func() {

		AfterEach(func() {
			TestServer.Clear()
		})

		It("Defines the response to be used for the server", func() {
			message := "Hello World"

			TestServer.Use(func(w http.ResponseWriter) {
				_, err := io.WriteString(w, message)
				check(err)
			})

			response, body, err := HTTPRequestDo("GET", fmt.Sprintf("http://localhost:%d", TestPort), nil, nil)

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(string(body)).To(Equal(message))
		})

		It("Clears the response to be used for the server", func() {
			message := "Hello World"
			factory := HTTPResponseFactory(func(w http.ResponseWriter) {
				_, err := io.WriteString(w, message)
				check(err)
			})
			TestServer.Use(factory)
			TestServer.Clear()

			response, body, err := HTTPRequestDo("GET", fmt.Sprintf("http://localhost:%d", TestPort), nil, nil)

			Expect(err).To(BeNil())
			Expect(response.StatusCode).To(Equal(http.StatusOK))
			Expect(string(body)).To(Equal(""))
		})

		It("Defines the response to be used for the server with predicate", func() {
			message := "Hello World"
			factory := HTTPResponseFactory(func(w http.ResponseWriter) {
				_, err := io.WriteString(w, message)
				check(err)
			})

			predicates := []HTTPRequestPredicate{
				RequestWithPath("/talula"),
				RequestWithMethod("POST"),
				RequestWithHeader("Content-Type", "application/json"),
			}

			TestServer.Use(factory).For(predicates...)

			pathMatching := fmt.Sprintf("http://localhost:%d/talula", TestPort)
			verbMatching := "POST"
			responseMatching, bodyMatching, errMatching := HTTPRequestDo(verbMatching, pathMatching, nil, func(request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			})

			pathNonMatching := fmt.Sprintf("http://localhost:%d", TestPort)
			verbNonMatching := "GET"
			responseNonMatching, bodyNonMatching, errNonMatching := HTTPRequestDo(verbNonMatching, pathNonMatching, nil, nil)

			Expect(errMatching).To(BeNil())
			Expect(responseMatching.StatusCode).To(Equal(http.StatusOK))
			Expect(string(bodyMatching)).To(Equal(message))

			Expect(errNonMatching).To(BeNil())
			Expect(responseNonMatching.StatusCode).To(Equal(http.StatusOK))
			Expect(string(bodyNonMatching)).To(Equal(""))
		})
	})
})
