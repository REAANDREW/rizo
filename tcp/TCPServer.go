package tcp

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type TCPServer struct {
	addr    *net.TCPAddr
	service *Service
}

func (s *TCPServer) Start() {
	listener, err := net.ListenTCP("tcp", s.addr)
	if nil != err {
		log.Fatalln(err)
	}
	log.Println("listening on", listener.Addr())

	go s.service.Serve(listener)
}

// Make a new Service.
func NewServer(addr string, port int) *TCPServer {
	s := &Service{
		ch:        make(chan bool),
		waitGroup: &sync.WaitGroup{},
	}
	s.waitGroup.Add(1)
	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:%v", addr, port))
	if nil != err {
		log.Fatalln(err)
	}
	server := &TCPServer{
		service: s,
		addr:    laddr,
	}
	return server
}

// An uninteresting service.
type Service struct {
	ch        chan bool
	waitGroup *sync.WaitGroup
}

// Accept connections and spawn a goroutine to serve each one.  Stop listening
// if anything is received on the service's channel.
func (s *Service) Serve(listener *net.TCPListener) {
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			log.Println("stopping listening on", listener.Addr())
			listener.Close()
			return
		default:
		}
		listener.SetDeadline(time.Now().Add(1e9))
		conn, err := listener.AcceptTCP()
		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		}
		log.Println(conn.RemoteAddr(), "connected")
		s.waitGroup.Add(1)
		go s.serve(conn)
	}
}

// Stop the service by closing the service's channel.  Block until the service
// is really stopped.
func (s *Service) Stop() {
	close(s.ch)
	s.waitGroup.Wait()
}

// Serve a connection by reading and writing what was read.  That's right, this
// is an echo service.  Stop reading and writing if anything is received on the
// service's channel but only after writing what was read.
func (s *Service) serve(conn *net.TCPConn) {
	defer conn.Close()
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			log.Println("disconnecting", conn.RemoteAddr())
			return
		default:
		}
		conn.SetDeadline(time.Now().Add(1e9))
		buf := make([]byte, 4096)
		if _, err := conn.Read(buf); nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
			return
		}
		if _, err := conn.Write(buf); nil != err {
			log.Println(err)
			return
		}
	}
}
