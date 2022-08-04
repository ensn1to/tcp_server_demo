package frame

// import (
// 	"bytes"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestNewMyFrameCodec(t *testing.T) {
// 	codec := NewMyFrameCodec()
// 	assert.Equal(t, codec, nil)
// }

// func TestPack(t *testing.T) {
// 	packer := NewMyFrameCodec()
// 	buf := make([]byte, 0, 128)
// 	rw := bytes.NewBuffer(buf)

// 	err := packer.Pack(rw, []byte("hello world"))
// 	if err != nil {
// 		t.Errorf("want nil, actual: %s", err.Error())
// 	}
// }
