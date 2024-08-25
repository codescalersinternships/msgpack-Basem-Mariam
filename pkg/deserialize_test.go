package msgpack

import (
	"bytes"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsgPackDeserializeIntegers(t *testing.T) {

	MsgPackTests := []struct {
		input []byte
		want  any
	}{
		{[]byte{0xd0, 0x07}, int8(7)},
		{[]byte{0xd1, 0x07, 0xd0}, int16(2000)},
		{[]byte{0xd2, 0x00, 0x1e, 0x84, 0x80}, int32(2000000)},
		{[]byte{0xd3, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1e, 0x84, 0x80}, int64(2000000)},
	}

	for _, tt := range MsgPackTests {
		got, err := Deserialize(bytes.NewReader(tt.input))
		if err != nil {
			t.Errorf("error %v", err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v want %v", got, tt.want)
		}
	}
}

func TestMsgPackDeserializeUnsignedIntegers(t *testing.T) {

	MsgPackTests := []struct {
		input []byte
		want  any
	}{
		{[]byte{0xcc, 0x7f}, uint8(127)},
		{[]byte{0xcd, 0x01, 0xf4}, uint16(500)},
		{[]byte{0xce, 0x00, 0x1e, 0x84, 0x80}, uint32(2000000)},
		{[]byte{0xcf, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1e, 0x84, 0x80}, uint64(2000000)},
	}

	for _, tt := range MsgPackTests {
		got, err := Deserialize(bytes.NewReader(tt.input))
		if err != nil {
			t.Errorf("error %v", err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v want %v", got, tt.want)
		}
	}
}

func TestMsgPackDeserializeBool(t *testing.T) {

	MsgPackTests := []struct {
		input []byte
		want  any
	}{
		{[]byte{0xc2}, false},
		{[]byte{0xc3}, true},
	}

	for _, tt := range MsgPackTests {

		got, err := Deserialize(bytes.NewReader(tt.input))
		assert.Equal(t, nil, err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v want %v", got, tt.want)
		}
	}
}

func TestMsgPackDeserializeNil(t *testing.T) {
	got, err := Deserialize(bytes.NewReader([]byte{0xc0}))
	assert.Equal(t, nil, err)
	if !reflect.DeepEqual(got, nil) {
		t.Errorf("got %v want nil", got)
	}
}

func TestMsgPackDeserializeFloat(t *testing.T) {
	MsgPackTests := []struct {
		input []byte
		want  any
	}{
		{[]byte{0xca, 0x40, 0x48, 0xf5, 0xc3}, float32(3.14)},
		{[]byte{0xca, 0xbd, 0xcc, 0xcc, 0xcd}, float32(-0.1)},
		{[]byte{0xca, 0x7f, 0x80, 0x00, 0x00}, float32(math.Inf(1))},
		{[]byte{0xca, 0xff, 0x80, 0x00, 0x00}, float32(math.Inf(-1))},
		{[]byte{0xcb, 0x3f, 0xb9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}, float64(0.1)},
		{[]byte{0xcb, 0xbf, 0xb9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}, float64(-0.1)},
	}

	for _, tt := range MsgPackTests {

		got, err := Deserialize(bytes.NewReader(tt.input))
		assert.Equal(t, nil, err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v want %v", got, tt.want)
		}
	}
}

func TestMsgPackDeserializeString(t *testing.T) {

	MsgPackTests := []struct {
		input []byte
		want  any
	}{
		{[]byte{0xa5, 'h', 'e', 'l', 'l', 'o'}, "hello"},
		{[]byte{RAW8, 0x05, 'h', 'e', 'l', 'l', 'o'}, "hello"},
		{[]byte{RAW16, 0x00, 0x05, 'w', 'o', 'r', 'l', 'd'}, "world"},
		{[]byte{RAW32, 0x00, 0x00, 0x00, 0x05, 't', 'e', 's', 't', '1'}, "test1"},
		{[]byte{RAW8, 0x00}, ""},
		{[]byte{RAW16, 0x00, 0x00}, ""},
		{[]byte{RAW32, 0x00, 0x00, 0x00, 0x00}, ""},
	}

	for _, tt := range MsgPackTests {
		got, err := Deserialize(bytes.NewReader(tt.input))
		assert.Equal(t, nil, err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v want %v", got, tt.want)
		}
	}
}

func TestMsgPackDeserializeArray(t *testing.T) {

	MsgPackTests := []struct {
		input []byte
		want  any
	}{
		{[]byte{ARRAY16, 0x00, 0x03, 0xd0, 0x01, 0xd0, 0x02, 0xd0, 0x03}, []any{int8(1), int8(2), int8(3)}},
		{[]byte{ARRAY32, 0x00, 0x00, 0x00, 0x03, 0xd0, 0x01, 0xd0, 0x02, 0xd0, 0x03}, []any{int8(1), int8(2), int8(3)}},
	}

	for _, tt := range MsgPackTests {

		got, err := Deserialize(bytes.NewReader(tt.input))
		assert.Equal(t, nil, err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v want %v", got, tt.want)
		}
	}
}

func TestMsgPackDeserializeMap(t *testing.T) {
	MsgPackTests := []struct {
		input []byte
		want  any
	}{
		{
			[]byte{MAP16, 0x00, 0x02, 0xd0, 0x01, 0xd0, 0x10, 0xd0, 0x02, 0xd0, 0x20},
			map[any]any{int8(1): int8(16), int8(2): int8(32)},
		},
		{
			[]byte{MAP32, 0x00, 0x00, 0x00, 0x02, 0xd0, 0x03, 0xd0, 0x30, 0xd0, 0x04, 0xd0, 0x40},
			map[any]any{int8(3): int8(48), int8(4): int8(64)},
		},
		{
			[]byte{MAP16, 0x00, 0x02, 0xd9, 0x03, 'k', 'e', 'y', 0xd0, 0x05, 0xd9, 0x03, 'a', 'b', 'c', 0xd0, 0x06},
			map[any]any{"key": int8(5), "abc": int8(6)},
		},
		{
			[]byte{MAP32, 0x00, 0x00, 0x00, 0x02, 0xd9, 0x03, 'x', 'y', 'z', 0xd1, 0x00, 0x7f, 0xd9, 0x03, 'p', 'q', 'r', 0xd1, 0x01, 0x01},
			map[any]any{"xyz": int16(127), "pqr": int16(257)},
		},
	}

	for _, tt := range MsgPackTests {

		got, err := Deserialize(bytes.NewReader(tt.input))
		assert.Equal(t, nil, err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	}
}
