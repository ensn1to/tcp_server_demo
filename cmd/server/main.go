package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/ensn1to/tcp_server_demo/config"
	"github.com/ensn1to/tcp_server_demo/frame"
	"github.com/ensn1to/tcp_server_demo/metrics"
	"github.com/ensn1to/tcp_server_demo/middleware/elk"
	"github.com/ensn1to/tcp_server_demo/packet"
	"github.com/spf13/viper"
)

func init() {
	elk.Logger = elk.NewElkLogger("192.168.0.223", 4560, 10)
}

func main() {
	cfg := prepareRun()
	go func() {
		http.ListenAndServe(cfg.HttpAddr, nil)
	}()

	run(cfg.TcpAddr)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM|syscall.SIGINT)
	<-c
}

func prepareRun() *config.Config {
	config.ReadFromEnv()

	serverMode := viper.GetString("ServerMode")
	consulAddr := viper.GetString("ConsulAddr")

	c, err := config.ReadFromConsul(consulAddr, serverMode)
	if err != nil {
		panic(err)
	}

	return c
}

func run(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	elk.Logger.Infof("server running...")

	for {
		c, err := l.Accept()
		if err != nil {
			elk.Logger.Errorf("accept error: ", err)
			break
		}

		go handlerConn(c)
	}
}

// handleConn的调用结构
// read frame from conn
//     ->frame decode
//       -> handle packet
//         -> packet decode
//         -> packet(ack) encode
//     ->frame(ack) encode
// write ack frame to conn

func handlerConn(c net.Conn) {
	metrics.ClientConnected.Inc()
	defer func() {
		metrics.ClientConnected.Dec()
		c.Close()
	}()

	framePacker := frame.NewMyFramePacker()
	rbf := bufio.NewReader(c)
	wbf := bufio.NewWriter(c)

	for {
		// read frame from the connection
		framePayload, err := framePacker.Unpack(rbf)
		if err != nil {
			elk.Logger.Errorf("handlerConn error: ", err)
			return
		}

		// add 1 recive request
		metrics.ReqRecvTotal.Add(1)

		ackFramePayload, err := handlePacket(framePayload)
		if err != nil {
			elk.Logger.Errorf("handleConn: handlePacket error: ", err)
			return
		}

		// write ack frame to the connetion
		err = framePacker.Pack(wbf, ackFramePayload)
		if err != nil {
			elk.Logger.Errorf("handleConn: framePacker pack error: ", err)
			return
		}

		// add 1 response
		metrics.RspSendTotal.Add(1)
	}
}

func handlePacket(framePayload []byte) (ackFramePayload []byte, err error) {
	var p packet.Packet
	p, err = packet.Decode(framePayload)
	if err != nil {
		elk.Logger.Errorf("handleConn: packet decode error: ", err)
		return
	}

	switch p.(type) {
	case *packet.Submit:
		submit := p.(*packet.Submit)
		elk.Logger.Infof("recv submit: id = %s, payload=%s\n", submit.ID, string(submit.Payload))
		submitAck := &packet.SubmitAck{
			ID:     submit.ID,
			Result: 0,
		}
		packet.SubmitPool.Put(submit) // 将submit对象归还给Pool池
		ackFramePayload, err = packet.Encode(submitAck)
		if err != nil {
			elk.Logger.Errorf("handleConn: packet encode error:", err)
			return nil, err
		}
		return ackFramePayload, nil
	default:
		return nil, fmt.Errorf("unknown packet type")
	}
}
