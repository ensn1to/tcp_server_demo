package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ensn1to/tcp_server_demo/frame"
	"github.com/ensn1to/tcp_server_demo/packet"
	"github.com/lucasepe/codename"
)

func main() {
	var wg sync.WaitGroup
	var num int = 100

	wg.Add(num)
	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()
			startClient(i)
		}(i + 1)
	}
	wg.Wait()
}

func startClient(i int) {
	quit := make(chan struct{})
	done := make(chan struct{})

	c, err := net.Dial("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	defer c.Close()
	fmt.Println("client running...")

	rng, err := codename.DefaultRNG()
	if err != nil {
		panic(err)
	}

	framePacker := frame.NewMyFramePacker()
	go func() {
		for {
			select {
			case <-quit:
				done <- struct{}{}
				return
			default:
			}

			c.SetReadDeadline(time.Now().Add(time.Second * 5))
			ackFramePayload, err := framePacker.Unpack(c)
			if err != nil {
				if e, ok := err.(net.Error); ok {
					if e.Timeout() {
						continue
					}
				}
				panic(err)
			}

			p, err := packet.Decode(ackFramePayload)
			submitAck, ok := p.(*packet.SubmitAck)
			if !ok {
				panic("not submitAck")
			}

			fmt.Printf("[client %d]: the result of submit ack[%s] is %d\n", i, submitAck.ID, submitAck.Result)
		}
	}()

	var counter int
	for {
		counter++
		id := fmt.Sprintf("%08d", counter)
		packetPayload := codename.Generate(rng, 4)
		s := &packet.Submit{
			ID:      id,
			Payload: []byte(packetPayload),
		}
		framePayload, err := packet.Encode(s)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[client %d]: send submit id = %s, payload=%s, frame length = %d\n",
			i, s.ID, s.Payload, len(framePayload)+4)

		err = framePacker.Pack(c, framePayload)
		if err != nil {
			panic(err)
		}

		time.Sleep(1)
		if counter >= 10 {
			quit <- struct{}{}
			<-done

			fmt.Printf("client[%d] exiting...\n", i)
			return
		}
	}
}
