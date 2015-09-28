package ami

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func fakeReadingAMI(t *testing.T, l *net.TCPListener, expect [][]byte) {
	lconn, err := l.Accept()
	if err != nil {
		t.Fail()
		t.Log("error accepting", err)
		return
	}
	err = lconn.SetDeadline(time.Now().Add(5 * time.Second))
	//TODO
	_, err = lconn.Write([]byte("Asterisk Call Manager/1.3\r\n"))
	if err != nil {
		t.Fail()
		t.Log("error writing", err)
		return
	}
	for _, a := range expect {
		b := make([]byte, 1024)
		n, err := lconn.Read(b)
		if err != nil {
			t.Fail()
			t.Log("error reading:", err)
		}
		b = b[:n]
		if !bytes.Equal(a, b) {
			t.Fail()
			t.Logf("read expected %#v, got %#v", string(a), string(b))
		}
	}
}

func TestSend(t *testing.T) {
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Skip("unable to listen")
	}
	defer l.Close()
	end := make(chan struct{})
	go func() {
		fakeReadingAMI(t, l, [][]byte{
			[]byte("Action: Login" + "\r\n" +
				"ActionID: 1" + "\r\n" +
				"Username: user" + "\r\n" +
				"Secret: secret" + "\r\n" +
				"\r\n"),
		})
		close(end)
	}()
	conn, err := dialTCP(l.Addr().(*net.TCPAddr))
	if err != nil {
		t.Fatal("error connecting:", err)
	}
	_, err = conn.Send(Login{Username: "user", Secret: "secret"})
	if err != nil {
		t.Fail()
		t.Log("error in Send:", err)
	}
	<-end
}
