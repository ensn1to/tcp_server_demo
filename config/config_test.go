package config

import (
	"testing"
)

func TestReadFromConsul(t *testing.T) {
	c, err := ReadFromConsul("192.168.0.223:8500", "dev")
	if err != nil {
		t.Logf("wanted nil, actual error: %s", err.Error())
	}

	t.Logf("http_addr%s\n", c.HttpAddr)
}
