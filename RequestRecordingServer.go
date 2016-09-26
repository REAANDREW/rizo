package rizo

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//RecordedRequest ...
type RecordedRequest struct {
	Request *http.Request
	Body    string
}

//HTTPRequestPredicate ...
type HTTPRequestPredicate func(request RecordedRequest) bool

//HTTPResponseFactory ...
type HTTPResponseFactory func(writer http.ResponseWriter)

//UseWithPredicates ...
type UseWithPredicates struct {
	ResponseFactory   HTTPResponseFactory
	RequestPredicates []HTTPRequestPredicate
}

//RequestRecordingServer ...
type RequestRecordingServer struct {
	Requests []RecordedRequest
	port     int
	server   *httptest.Server
	use      []UseWithPredicates
	lock     *sync.Mutex
}

//CreateRequestRecordingServer ...
func CreateRequestRecordingServer(port int) *RequestRecordingServer {
	return &RequestRecordingServer{
		Requests: []RecordedRequest{},
		port:     port,
		use:      []UseWithPredicates{},
		lock:     &sync.Mutex{},
	}
}

//CreateURL ...
func (instance *RequestRecordingServer) CreateURL(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", instance.port, path)
}

func (instance *RequestRecordingServer) evaluatePredicates(recordedRequest RecordedRequest, w http.ResponseWriter) {
	for _, item := range instance.use {
		if item.RequestPredicates != nil {
			result := instance.Evaluate(recordedRequest, item.RequestPredicates...)
			if result {
				item.ResponseFactory(w)
				return
			}
		} else {
			item.ResponseFactory(w)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

//Start ...
func (instance *RequestRecordingServer) Start() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		check(err)
		recordedRequest := RecordedRequest{
			Request: r,
			Body:    string(body),
		}
		instance.lock.Lock()
		instance.Requests = append(instance.Requests, recordedRequest)
		if instance.use != nil {
			instance.evaluatePredicates(recordedRequest, w)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		instance.lock.Unlock()
	})
	instance.server = httptest.NewUnstartedServer(handler)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(instance.port))
	if err != nil {
		panic(err)
	}
	instance.server.Listener = listener
	instance.server.Start()
}

//Stop ...
func (instance *RequestRecordingServer) Stop() {
	if instance.server != nil {
		instance.server.Close()
		time.Sleep(1 * time.Millisecond)
	}
}

//Clear ...
func (instance *RequestRecordingServer) Clear() {
	instance.lock.Lock()
	instance.Requests = []RecordedRequest{}
	instance.use = []UseWithPredicates{}
	instance.lock.Unlock()
}

//Evaluate ...
func (instance *RequestRecordingServer) Evaluate(request RecordedRequest, predicates ...HTTPRequestPredicate) bool {
	results := make([]bool, len(predicates))
	for index, predicate := range predicates {
		results[index] = predicate(request)
	}
	thing := true
	for _, result := range results {
		if !result {
			thing = false
			break
		}
	}
	return thing
}

//Find ...
func (instance *RequestRecordingServer) Find(predicates ...HTTPRequestPredicate) bool {
	for _, request := range instance.Requests {
		if instance.Evaluate(request, predicates...) {
			return true
		}
	}
	return false
}

//Use ...
func (instance *RequestRecordingServer) Use(factory HTTPResponseFactory) *RequestRecordingServer {
	instance.use = append(instance.use, UseWithPredicates{
		ResponseFactory:   factory,
		RequestPredicates: []HTTPRequestPredicate{},
	})
	return instance
}

//For ...
func (instance *RequestRecordingServer) For(predicates ...HTTPRequestPredicate) {
	index := len(instance.use) - 1
	for _, item := range predicates {
		instance.use[index].RequestPredicates = append(instance.use[index].RequestPredicates, item)
	}
}

//RequestWithPath ...
func RequestWithPath(path string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := r.Request.URL.Path == path
		return result
	})
}

//RequestWithMethod ...
func RequestWithMethod(method string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := r.Request.Method == method
		return result
	})
}

//RequestWithHeader ...
func RequestWithHeader(key string, value string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := r.Request.Header.Get(key) == value
		return result
	})
}

//RequestWithBody ...
func RequestWithBody(value string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := string(r.Body) == value
		return result
	})
}

//RequestWithQuerystring ...
func RequestWithQuerystring(value string) HTTPRequestPredicate {
	return HTTPRequestPredicate(func(r RecordedRequest) bool {
		result := r.Request.URL.RawQuery == value
		return result
	})
}
