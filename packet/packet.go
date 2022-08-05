package packet

import (
	"bytes"
	"fmt"
)

const (
	CommandConn = iota + 0x01
	CommandSubmit
)

const (
	CommandConnAck = iota + 0x80
	CommandSubmitAck
)

type Packet interface {
	Decode([]byte) error     // []byte -> struct
	Encode() ([]byte, error) // struct -> []byte
}

type Submit struct {
	ID      string //
	Payload []byte //
}

func (s *Submit) Decode(data []byte) error {
	s.ID = string(data[:8])
	s.Payload = data[8:]
	return nil
}

func (s *Submit) Encode() ([]byte, error) {
	id := []byte(s.ID[:8])
	packetBytes := [][]byte{id, s.Payload}
	return bytes.Join(packetBytes, nil), nil
}

type SubmitAck struct {
	ID     string
	Result uint8
}

func (s *SubmitAck) Decode(data []byte) error {
	s.ID = string(data[:8])
	s.Result = uint8(data[8])
	return nil
}

func (s *SubmitAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.ID[:8]), {s.Result}}, nil), nil
}

func Decode(packet []byte) (Packet, error) {
	commandID := packet[0]
	pktBody := packet[1:]

	switch commandID {
	case CommandConn:
		return nil, nil
	case CommandConnAck:
		return nil, nil
	case CommandSubmit:
		s := Submit{}
		if err := s.Decode(pktBody); err != nil {
			return nil, err
		}
		return &s, nil
	case CommandSubmitAck:
		s := SubmitAck{}
		if err := s.Decode(pktBody); err != nil {
			return nil, err
		}
		return &s, nil
	default:
		return nil, fmt.Errorf("unknown commandID [%d]", commandID)
	}
}

func Encode(p Packet) ([]byte, error) {
	var commandID uint8
	var pktBody []byte
	var err error

	switch t := p.(type) {
	case *Submit:
		commandID = CommandSubmit
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *SubmitAck:
		commandID = CommandSubmitAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown type [%s]", t)
	}
	return bytes.Join([][]byte{{commandID}, pktBody}, nil), nil
}
