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

	"github.com/ensn1to/tcp_server_demo/frame"
	"github.com/ensn1to/tcp_server_demo/metrics"
	"github.com/ensn1to/tcp_server_demo/packet"
)

func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	run()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM|syscall.SIGINT)
	<-c
}

func run() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}

	fmt.Println("server running...")

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
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
			fmt.Println("handlerConn error: ", err)
			return
		}

		// add 1 recive request
		metrics.ReqRecvTotal.Add(1)

		ackFramePayload, err := handlePacket(framePayload)
		if err != nil {
			fmt.Println("handleConn: handlePacket error: ", err)
			return
		}

		// write ack frame to the connetion
		err = framePacker.Pack(wbf, ackFramePayload)
		if err != nil {
			fmt.Println("handleConn: framePacker pack error: ", err)
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
		fmt.Println("handleConn: packet decode error: ", err)
		return
	}

	switch p.(type) {
	case *packet.Submit:
		submit := p.(*packet.Submit)
		fmt.Printf("recv submit: id = %s, payload=%s\n", submit.ID, string(submit.Payload))
		submitAck := &packet.SubmitAck{
			ID:     submit.ID,
			Result: 0,
		}
		packet.SubmitPool.Put(submit) // 将submit对象归还给Pool池
		ackFramePayload, err = packet.Encode(submitAck)
		if err != nil {
			fmt.Println("handleConn: packet encode error:", err)
			return nil, err
		}
		return ackFramePayload, nil
	default:
		return nil, fmt.Errorf("unknown packet type")
	}
}
