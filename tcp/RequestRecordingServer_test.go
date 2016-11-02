package tcp

import (
	"bufio"
	"fmt"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RequestRecordingServer", func() {
	var (
		//TestServer ...
		TestServer *RequestRecordingServer
		TestPort   = 6001
		client     *TCPClient
		request    string
	)
	BeforeSuite(func() {
		TestServer = New(TestPort)
		TestServer.Start()
		time.Sleep(1000 * time.Millisecond)
	})

	AfterSuite(func() {
		TestServer.Stop()
	})

	BeforeEach(func() {
		client = NewClient("localhost", TestPort)
	})

	AfterEach(func() {
		client.Close()
		TestServer.Clear()
	})

	Describe("Find", func() {

		Describe("Single Request", func() {
			It("Body", func() {
				request = "*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n"
				client.Send(request)
				Expect(TestServer.Find(RequestWithBody("mylist\r\n"))).To(Equal(true))
			})

		})
	})

	It("Use", func() {
		messages := []string{}
		TestServer.Use(func(request RecordedRequest, w ResponseWriter) {
			fmt.Printf("%#v\n", w)
			if len(messages) == 0 {
				r := []rune(request.Body)
				switch commandType := r[0]; commandType {
				case ':':
					fmt.Println("RESP Integer")
				case '+':
					fmt.Println("RESP Simple String")
				case '-':
					fmt.Println("RESP Error")
				case '$':
					length := request.Body[1:]
					fmt.Printf("RESP Bulk String %v\n", length)
				case '*':
					length := request.Body[1:]
					fmt.Printf("RESP Array %v\n", length)
				}
			}
			w.Send("Talula\r\n")
		})
		request = "*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n"
		result := client.Send(request)
		Expect(result).To(Equal("Talula\r\n"))
	})
})

type TCPClient struct {
	host string
	port int
	conn net.Conn
}

func NewClient(host string, port int) *TCPClient {
	conn, e := net.Dial("tcp", fmt.Sprintf(host+":%v", port))
	if e != nil {
		panic(e)
	}
	fmt.Println("connected")
	return &TCPClient{
		host: host,
		port: port,
		conn: conn,
	}
}

// func send(url string, data string, done chan bool) string {
// 	done <- true
// }
func (instance *TCPClient) Send(data string) string {
	// fmt.Fprintf(conn, data)
	instance.conn.Write([]byte(data))
	message, _ := bufio.NewReader(instance.conn).ReadString('\n')
	fmt.Print("Message from server: " + message)
	return message
}

func (instance *TCPClient) Close() error {
	return instance.conn.Close()
}
