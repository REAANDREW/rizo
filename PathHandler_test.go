package rizo

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var responseWriter = &FakeResponseWriter{}

var _ = Describe("PathHandler", func() {
	var (
		responseWriter *FakeResponseWriter
		handler        *PathHandler
	)

	BeforeEach(func() {
		responseWriter = &FakeResponseWriter{}
		handler = NewPathHandler()
	})

	AfterEach(func() {
		responseWriter.Reset()
	})

	It("Configures a GET handler", func() {
		const expectedMessage string = "handled the GET"

		handler.Get(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(expectedMessage))
		})

		handler.Handle(responseWriter, &http.Request{
			Method: "GET",
		})

		Expect(string(responseWriter.Data)).To(Equal(expectedMessage))
	})

	It("Configures a POST handler", func() {
		const expectedMessage string = "handled the POST"

		handler.Post(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(expectedMessage))
		})

		handler.Handle(responseWriter, &http.Request{
			Method: "POST",
		})

		Expect(string(responseWriter.Data)).To(Equal(expectedMessage))
	})

	It("Configures a PUT handler", func() {
		const expectedMessage string = "handled the PUT"

		handler.Put(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(expectedMessage))
		})

		handler.Handle(responseWriter, &http.Request{
			Method: "PUT",
		})

		Expect(string(responseWriter.Data)).To(Equal(expectedMessage))
	})

	It("Configures a DELETE handler", func() {
		const expectedMessage string = "handled the DELETE"

		handler.Delete(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(expectedMessage))
		})

		handler.Handle(responseWriter, &http.Request{
			Method: "DELETE",
		})

		Expect(string(responseWriter.Data)).To(Equal(expectedMessage))
	})

	It("Configures a PATCH handler", func() {
		const expectedMessage string = "handled the PATCH"

		handler.Patch(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(expectedMessage))
		})

		handler.Handle(responseWriter, &http.Request{
			Method: "PATCH",
		})

		Expect(string(responseWriter.Data)).To(Equal(expectedMessage))
	})

	It("Configures any other Method", func() {
		const expectedMessage string = "handled the BOOM"

		handler.Method("BOOM", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(expectedMessage))
		})

		handler.Handle(responseWriter, &http.Request{
			Method: "BOOM",
		})

		Expect(string(responseWriter.Data)).To(Equal(expectedMessage))
	})
})
