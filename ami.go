package ami

import (
	"bufio"
	"net"
	"sync"
)

type ActionID uint64

// An AMI-returned `Response: Error`, containing the value of Message
type AMIErrorResponse string

func (e AMIErrorResponse) Error() string {
	return string(e)
}

type Conn struct {
	c *net.TCPConn
	s *bufio.Scanner

	// Last used ActionID for this connection (lock, inc, take, unlock)
	idmu sync.Mutex
	id   ActionID
}

// Dial sets up an AMI connection to the given address. The address can
// contain a port, and will fall back to port 5038 if none is given.
func Dial(address string) (*Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		tcpAddr, err = net.ResolveTCPAddr("tcp", address+":5038")
	}
	if err != nil {
		return nil, err
	}
	return dialTCP(tcpAddr)
}

func dialTCP(address *net.TCPAddr) (*Conn, error) {
	c := &Conn{}
	var err error
	c.c, err = net.DialTCP("tcp", nil, address)
	if err != nil {
		return nil, err
	}
	c.s = bufio.NewScanner(c.c)
	c.s.Scan()
	if c.s.Text() != "Asterisk Call Manager/1.3" {
		c.Close()
		return nil, AMIErrorResponse("unexpected AMI identification string: " + c.s.Text())
	}
	return c, nil
}

// Send an Action message through this AMI connection. Action should be a
// struct from a named type, with exported fields for each valuein the Action.
// The struct fields should be of string or integer types.
func (c *Conn) Send(action interface{}) (ActionID, error) {
	c.idmu.Lock()
	c.id += 1
	id := c.id
	c.idmu.Unlock()
	b, err := marshalAction(action, id)
	if err != nil {
		return id, err
	}
	//TODO: move this to a central dispatch?
	_, err = c.c.Write(b)
	if err != nil {
		//TODO: close/reconnect/whatever
		return id, err
	}
	return id, nil
}

func (c *Conn) Close() error {
	return c.c.Close()
}
