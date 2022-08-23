package elk

import (
	"fmt"
	"net"
	"time"
)

type LogStash struct {
	hostname string
	port     int
	Conn     *net.TCPConn
	TimeOut  int
}

func newLogStash(hostname string, port int, timeout int) *LogStash {
	ls := &LogStash{
		hostname: hostname,
		port:     port,
		Conn:     nil,
		TimeOut:  timeout,
	}
	var err error
	ls.Conn, err = ls.Connect()
	if err != nil {
		panic(err)
	}

	return ls
}

func (l *LogStash) Connect() (*net.TCPConn, error) {
	var conn *net.TCPConn
	service := fmt.Sprintf("%s:%d", l.hostname, l.port)
	addr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		return nil, err
	}
	conn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	if conn != nil {
		l.Conn = conn
		l.Conn.SetKeepAlive(true)
		l.Conn.SetKeepAlivePeriod(time.Duration(5) * time.Second)
		l.setTimeouts()
	}
	return conn, nil
}

func (l *LogStash) setTimeouts() {
	deadline := time.Now().Add(time.Duration(l.TimeOut) * time.Millisecond)
	_ = l.Conn.SetDeadline(deadline)
	_ = l.Conn.SetReadDeadline(deadline)
	_ = l.Conn.SetWriteDeadline(deadline)
}

// Write. message: json
func (l *LogStash) Write(message string) (err error) {
	message = fmt.Sprintf("%s\n", message)
	if l.Conn != nil {
		fmt.Printf("write message: %s\n", message)
		_, err = l.Conn.Write([]byte(message))
		if err != nil {
			l.Connect()
			return
		} else {
			l.setTimeouts()
			return nil
		}
	}
	return
}
