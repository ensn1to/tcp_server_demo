package elk

import (
	"encoding/json"
	"fmt"
)

var Logger *ElkLogger

type ElkLogger struct {
	Logger *LogStash
}

func NewElkLogger(hostname string, port int, timeout int) *ElkLogger {
	return &ElkLogger{Logger: newLogStash(hostname, port, timeout)}
}

type ElkMsg struct {
	Service string
	Level   string
	Msg     string
}

func (log *ElkLogger) Errorf(format string, args ...interface{}) {
	msg := &ElkMsg{
		Service: "tcp_server_demo",
		Level:   "error",
		Msg:     fmt.Sprintf(format, args...),
	}
	var (
		bytes []byte
		err   error
	)
	if bytes, err = json.Marshal(msg); err != nil {
		fmt.Println("elk msg marshal error")
		return
	}
	log.Logger.Write(string(bytes))
}

func (log *ElkLogger) Infof(format string, args ...interface{}) {
	msg := &ElkMsg{
		Service: "tcp_server_demo",
		Level:   "info",
		Msg:     fmt.Sprintf(format, args...),
	}
	var (
		bytes []byte
		err   error
	)
	if bytes, err = json.Marshal(msg); err != nil {
		fmt.Println("elk msg marshal error")
		return
	}
	fmt.Println("Infof", msg)
	log.Logger.Write(string(bytes))
}

func (log *ElkLogger) Warnf(format string, args ...interface{}) {
	msg := &ElkMsg{
		Service: "tcp_server_demo",
		Level:   "warn",
		Msg:     fmt.Sprintf(format, args...),
	}
	var (
		bytes []byte
		err   error
	)
	if bytes, err = json.Marshal(msg); err != nil {
		fmt.Println("elk msg marshal error")
		return
	}
	log.Logger.Write(string(bytes))
}
