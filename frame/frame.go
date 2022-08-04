package frame

import (
	"encoding/binary"
	"errors"
	"io"
)

type FramePayload []byte

type StreamFramePacker interface {
	Pack(io.Writer, FramePayload) error     // struct data -> frame --> io.Writer
	Unpack(io.Reader) (FramePayload, error) // io.Reader -> frame payload --> struct data
}

var (
	ErrShortWrite = errors.New("short write")
	ErrShortRead  = errors.New("short read")
)

type myFramePacker struct{}

func NewMyFrameCodec() StreamFramePacker {
	return &myFramePacker{}
}

func (c *myFramePacker) Pack(w io.Writer, framePayload FramePayload) error {
	f := framePayload
	var totalLen int32 = int32(len(framePayload)) + 4 // 4 bytes

	err := binary.Write(w, binary.BigEndian, &totalLen)
	if err != nil {
		return err
	}

	// write the frame payload to the outbound stream
	n, err := w.Write([]byte(f))
	if err != nil {
		return err
	}

	// check if all data be writen
	if n != len(framePayload) {
		return ErrShortWrite
	}

	return nil
}

func (c *myFramePacker) Unpack(r io.Reader) (FramePayload, error) {
	var totalLen int32
	if err := binary.Read(r, binary.BigEndian, &totalLen); err != nil {
		return nil, err
	}

	buf := make([]byte, totalLen-4)
	// read all bytes needed, unless happends to EOF or ErrUnexpectedEOF
	// can handle stick bag?
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	if n != int(totalLen-4) {
		return nil, ErrShortRead
	}

	return FramePayload(buf), nil
}
