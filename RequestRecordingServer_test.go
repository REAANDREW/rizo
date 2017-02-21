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

		var (
			sampleURL string
			predicate HTTPRequestPredicate
		)
		const (
			data = "a=1&b=2"
		)
		Describe("Single Request", func() {
			var request *http.Request

			BeforeEach(func() {
				sampleURL = URLForTestServer("/Fubar?" + data)
				request, _ = http.NewRequest("GET", sampleURL, bytes.NewBuffer([]byte(data)))
				request.Header.Set("Content-type", "application/json")
				TestServer.Requests = append(TestServer.Requests, RecordedRequest{
					Request: request,
					Body:    data,
				})
			})

			Describe("Path", func() {
				const expectedPath = "/Fubar"
				BeforeEach(func() {
					predicate = RequestWithPath(expectedPath)
				})
				It("provides a String representation", func() {
					Ω(predicate.String()).Should(Equal("WithPath: '/Fubar'"))
				})
				It("can find using the predicate", func() {
					Expect(TestServer.Find(predicate)).To(Equal(true))
				})
			})

			Describe("Method", func() {
				const expectedMethod = "GET"
				BeforeEach(func() {
					predicate = RequestWithMethod(expectedMethod)
				})
				It("provides a String representation", func() {
					Ω(predicate.String()).Should(Equal("WithMethod: 'GET'"))
				})
				It("can find using the predicate", func() {
					Expect(TestServer.Find(predicate)).To(Equal(true))
				})
			})

			Describe("Header", func() {
				BeforeEach(func() {
					predicate = RequestWithHeader("Content-type", "application/json")
				})
				It("provides a String representation", func() {
					Ω(predicate.String()).Should(Equal("WithHeader: 'Content-type':'application/json'"))
				})
				It("can find using the predicate", func() {
					Expect(TestServer.Find(predicate)).To(Equal(true))
				})
			})

			Describe("Body", func() {
				BeforeEach(func() {
					predicate = RequestWithBody(data)
				})
				It("provides a String representation", func() {
					Ω(predicate.String()).Should(Equal("WithBody: 'a=1&b=2'"))
				})
				It("can find using the predicate", func() {
					request, _ = http.NewRequest("POST", sampleURL, bytes.NewBuffer([]byte(data)))
					TestServer.Clear()
					TestServer.Requests = append(TestServer.Requests, RecordedRequest{
						Request: request,
						Body:    data,
					})
					Expect(TestServer.Find(predicate)).To(Equal(true))
				})
			})

			Describe("Querystring", func() {
				BeforeEach(func() {
					predicate = RequestWithQuerystring(data)
				})
				It("provides a String representation", func() {
					Ω(predicate.String()).Should(Equal("WithQuerystring: 'a=1&b=2'"))
				})
				It("can find using the predicate", func() {
					Expect(TestServer.Find(predicate)).To(Equal(true))
				})
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
