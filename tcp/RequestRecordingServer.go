package tcp

import (
	"sync"

	"fmt"
	"github.com/firstrow/tcp_server"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//RecordedRequest ...
type RecordedRequest struct {
	//maybe want the local/remote addresses
	Body string
}

//ResponseWriter ...
type ResponseWriter interface {
	Send(m string) error
	SendBytes(b []byte) error
}

//RequestPredicate ...
type RequestPredicate func(request RecordedRequest) bool

//ResponseFactory ...
type ResponseFactory func(writer ResponseWriter)

//UseWithPredicates ...
type UseWithPredicates struct {
	ResponseFactory   ResponseFactory
	RequestPredicates []RequestPredicate
}

//RequestRecordingServer ...
type RequestRecordingServer struct {
	Requests []RecordedRequest
	port     int
	server   TCPServer
	use      []UseWithPredicates
	lock     *sync.Mutex
}

//New ...
func New(port int) *RequestRecordingServer {
	instance := RequestRecordingServer{
		Requests: []RecordedRequest{},
		port:     port,
		use:      []UseWithPredicates{},
		lock:     &sync.Mutex{},
	}
	server := tcp_server.New(fmt.Sprintf("localhost:%v", port))

	server.OnNewClient(func(c *tcp_server.Client) {
		// new client connected
		// lets send some message
		fmt.Println("New client")
		c.Send("Hello")
	})
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		// new message received
		fmt.Println("message recieved")
		// c.Send(":48293\r\n")
		recordedRequest := RecordedRequest{
			Body: message,
		}
		instance.lock.Lock()
		instance.Requests = append(instance.Requests, recordedRequest)
		if instance.use != nil {
			instance.evaluatePredicates(recordedRequest, c)
		} else {
			c.Send(message)
		}
		instance.lock.Unlock()
	})
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		// connection with client lost
		fmt.Println("Client connection lost")
	})
	instance.server = server

	return &instance
}

//Start ...
func (instance *RequestRecordingServer) Start() {
	go instance.server.Listen()
}

//Stop ...
func (instance *RequestRecordingServer) Stop() {
	if instance.server != nil {
		// 	instance.server.Close()
		time.Sleep(1 * time.Millisecond)
	}
}

//Clear ...
func (instance *RequestRecordingServer) Clear() {
	instance.lock.Lock()
	instance.Requests = []RecordedRequest{}
	instance.lock.Unlock()
}

//Use ...
func (instance *RequestRecordingServer) Use(factory ResponseFactory) *RequestRecordingServer {
	instance.use = append(instance.use, UseWithPredicates{
		ResponseFactory:   factory,
		RequestPredicates: []RequestPredicate{},
	})
	return instance
}

//For ...
func (instance *RequestRecordingServer) For(predicates ...RequestPredicate) {
	index := len(instance.use) - 1
	for _, item := range predicates {
		instance.use[index].RequestPredicates = append(instance.use[index].RequestPredicates, item)
	}
}

//Evaluate ...
func (instance *RequestRecordingServer) Evaluate(request RecordedRequest, predicates ...RequestPredicate) bool {
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
func (instance *RequestRecordingServer) Find(predicates ...RequestPredicate) bool {
	for _, request := range instance.Requests {
		if instance.Evaluate(request, predicates...) {
			return true
		}
	}
	return false
}

//RequestWithBody ...
func RequestWithBody(value string) RequestPredicate {
	return RequestPredicate(func(r RecordedRequest) bool {
		result := string(r.Body) == value
		return result
	})
}

func (instance *RequestRecordingServer) evaluatePredicates(recordedRequest RecordedRequest, c *tcp_server.Client) {
	for _, item := range instance.use {
		if item.RequestPredicates != nil {
			result := instance.Evaluate(recordedRequest, item.RequestPredicates...)
			if result {
				item.ResponseFactory(c)
				return
			}
		} else {
			item.ResponseFactory(c)
			return
		}
	}
}
