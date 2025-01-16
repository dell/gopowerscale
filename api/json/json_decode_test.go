package json

import (
	"reflect"
	"testing"
)

func TestLiteralStore(t *testing.T) {
	// Test case for unmarshaling a valid JSON null
	d := &decodeState{
		data: []byte(`null`),
		off:  0,
	}
	var item []byte
	item = []byte("n")
	v := reflect.ValueOf(5)
	d.literalStore(item, v, false)

	item = []byte("t")
	v = reflect.ValueOf(5)
	d.literalStore(item, v, true)

	item = []byte("2")
	x := 5
	v = reflect.ValueOf(&x)
	d.literalStore(item, v, true)

	item = []byte("2.2")
	y := 5.5
	v = reflect.ValueOf(&y)
	d.literalStore(item, v, true)

	item = []byte("2.2")
	var z uint8 = 5
	v = reflect.ValueOf(&z)
	d.literalStore(item, v, true)
}
