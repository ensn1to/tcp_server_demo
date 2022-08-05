package frame

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"

	"gotest.tools/assert"
)

func TestNewMyFrameCodec(t *testing.T) {
	packer := NewMyFramePacker()
	assert.Equal(t, packer, nil)
}

func TestPack(t *testing.T) {
	packer := NewMyFramePacker()
	buf := make([]byte, 0, 128)
	rw := bytes.NewBuffer(buf)

	err := packer.Pack(rw, []byte("hello"))
	if err != nil {
		t.Errorf("want nil, actual: %s", err.Error())
	}

	var totalLen int32
	err = binary.Read(rw, binary.BigEndian, &totalLen)
	if err != nil {
		t.Errorf("want nil, actual: %s", err.Error())
	}

	if totalLen != 9 {
		t.Errorf("want 9, actual: %d", totalLen)
	}

	left := rw.Bytes()
	if string(left) != "hello" {
		t.Errorf("want hello, actual: %s", string(left))
	}

	assert.Equal(t, false, true)
}

func TestUnpack(t *testing.T) {
	packer := NewMyFramePacker()
	data := []byte{0x00, 0x00, 0x00, 0x09, 'h', 'e', 'l', 'l', '0'}

	payload, err := packer.Unpack(bytes.NewBuffer(data))
	if err != nil {
		t.Errorf("want nil, actual: %s", err.Error())
	}

	if string(payload) != "hello" {
		t.Errorf("want hello, actual: %s", string(payload))
	}
}

// 错误覆盖率
type ReturnErrWriter struct {
	W  io.Writer
	Wn int // 第几次调用write返回错误
	wc int // 写操作次数计数
}

func (w *ReturnErrWriter) Write(p []byte) (n int, err error) {
	w.wc++
	if w.wc >= w.Wn {
		return 0, errors.New("Write error")
	}
	return w.W.Write(p)
}

type ReturnErrReader struct {
	R  io.Reader
	Rn int
	rc int
}

func (r *ReturnErrReader) Read(p []byte) (n int, err error) {
	r.rc++
	if r.rc >= r.Rn {
		return 0, errors.New("read error")
	}
	return r.R.Read(p)
}

func TestPackWithWriteFail(t *testing.T) {
	packer := NewMyFramePacker()
	buf := make([]byte, 0, 128)
	w := bytes.NewBuffer(buf)

	err := packer.Pack(&ReturnErrWriter{
		W:  w,
		Wn: 1,
	}, []byte("hello"))
	if err != nil {
		t.Errorf("want non-nil, actual nil")
	}

	err = packer.Pack(&ReturnErrWriter{
		W:  w,
		Wn: 2,
	}, []byte("hello"))
	if err != nil {
		t.Errorf("want non-nil, actual nil")
	}

	assert.Equal(t, false, true)
}
