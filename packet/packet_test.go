package packet

import (
	"fmt"
	"testing"
)

func TestSubmitDecode(t *testing.T) {
	s := &Submit{}
	packetPayload := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 'h', 'e', 'l', 'l', '0'}
	err := s.Decode(packetPayload)
	if err != nil {
		t.Errorf("want nil, actual: %s", err.Error())
	}
}

func TestSubmitEncode(t *testing.T) {
	s := &Submit{
		ID:      fmt.Sprintf("%08d", 1),
		Payload: []byte("hello"),
	}
	d, _ := s.Encode()
	if len(d) == 0 {
		t.Errorf("want not 0, actual 0")
	}
}
