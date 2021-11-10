package simplepeg

import "testing"

func TestHello(t *testing.T) {
	v := Hello()

	if v != "World" {
		t.Error("Expected World, got", v)
	}
}
