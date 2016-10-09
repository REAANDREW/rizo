package tcp

import (
	"bufio"
	"fmt"
	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var (
	//TestServer ...
	TestServer *RequestRecordingServer
	TestPort   = 6000
)

func URLForTestServer(path string) string {
	return fmt.Sprintf("http://localhost:%d%s", TestPort, path)
}

var _ = Describe("RequestRecordingServer", func() {

	BeforeEach(func() {
		TestServer = New(TestPort)
		TestServer.Start()
		time.Sleep(1000 * time.Millisecond)
	})

	AfterEach(func() {
		TestServer.Clear()
		TestServer.Stop()
	})

	Describe("Find", func() {
		var client *Client

		Describe("Single Request", func() {
			var request string

			BeforeEach(func() {
				client = NewClient("localhost", TestPort)
			})

			It("Body", func() {
				request = "*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n"
				client.Send(request)
				Expect(TestServer.Find(RequestWithBody(request))).To(Equal(true))
			})
		})
	})
})

type Client struct {
	host string
	port int
}

func NewClient(host string, port int) *Client {
	return &Client{
		host: host,
		port: port,
	}
}

func (instance *Client) Send(data string) string {
	conn, e := net.Dial("tcp", fmt.Sprintf(instance.host+":%v", instance.port))
	if e != nil {
		panic(e)
	}
	fmt.Println("connected")
	// fmt.Fprintf(conn, data)
	conn.Write([]byte(data))
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)
	return message
}
