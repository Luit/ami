package ami

import (
	"fmt"
	"net"
	"sync"
)

type (
	Message map[string][]string
)

type Conn interface {
	Do(action interface{}, response interface{}) error
	// List(action interface{}, response interface{}, list interface{}) error

	// Action(action string, m Message) (<-chan Message, error)
	// Subscribe(event string, ch chan<- Message)
	// Unsubscribe(event string, ch chan<- Message)
	Close() error
}

type conn struct {
	c  *net.TCPConn
	id struct {
		id int
		mu sync.Mutex
	}
}

func Dial(addr string) (Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		tcpAddr, err = net.ResolveTCPAddr("tcp", addr+":5038")
	}
	if err != nil {
		return nil, err
	}
	c := &conn{}
	c.c, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *conn) Do(action interface{}, response interface{}) error {
	c.id.mu.Lock()
	c.id.id += 1
	id := c.id.id
	c.id.mu.Unlock()
	b, err := marshalAction(action, id)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

// func (c *conn) Action(action string, m Message) (<-chan Message, error) {
// 	return nil, nil
// }

// func (c *conn) Subscribe(event string, ch chan<- Message)   {}
// func (c *conn) Unsubscribe(event string, ch chan<- Message) {}

func (c *conn) Close() error {
	return c.c.Close()
}
