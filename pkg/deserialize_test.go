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
		want  Value
	}{
		{input: []byte{0xd0, 0x07}, want: Value{typ: "integer", val: int8(7)}},
		{input: []byte{0xd1, 0x07, 0xd0}, want: Value{typ: "integer", val: int16(2000)}},
		{input: []byte{0xd2, 0x00, 0x1e, 0x84, 0x80}, want: Value{typ: "integer", val: int32(2000000)}},
		{input: []byte{0xd3, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1e, 0x84, 0x80}, want: Value{typ: "integer", val: int64(2000000)}},
	}

	for _, tt := range MsgPackTests {
		msgpacker := NewMessagePacker(bytes.NewReader(tt.input))
		got, err := msgpacker.deserialize()
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
		want  Value
	}{
		{input: []byte{0xcc, 0x7f}, want: Value{typ: "Unsigned-integer", val: uint8(127)}},
		{input: []byte{0xcd, 0x01, 0xf4}, want: Value{typ: "Unsigned-integer", val: uint16(500)}},
		{input: []byte{0xce, 0x00, 0x1e, 0x84, 0x80}, want: Value{typ: "Unsigned-integer", val: uint32(2000000)}},
		{input: []byte{0xcf, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1e, 0x84, 0x80}, want: Value{typ: "Unsigned-integer", val: uint64(2000000)}},
	}

	for _, tt := range MsgPackTests {
		msgpacker := NewMessagePacker(bytes.NewReader(tt.input))
		got, err := msgpacker.deserialize()
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
		want  Value
	}{
		{input: []byte{0xc2}, want: Value{typ: "boolen", val: false}},
		{input: []byte{0xc3}, want: Value{typ: "boolen", val: true}},
	}

	for _, tt := range MsgPackTests {
		msgpacker := NewMessagePacker(bytes.NewReader(tt.input))
		got, err := msgpacker.deserialize()
		assert.Equal(t, nil, err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v want %v", got, tt.want)
		}
	}
}

func TestMsgPackDeserializeNil(t *testing.T) {

	msgpacker := NewMessagePacker(bytes.NewReader([]byte{0xc0}))
	got, err := msgpacker.deserialize()
	assert.Equal(t, nil, err)
	want := Value{typ: "nil", val: nil}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestMsgPackDeserializeFloat(t *testing.T) {
	MsgPackTests := []struct {
		input []byte
		want  Value
	}{
		{input: []byte{0xca, 0x40, 0x48, 0xf5, 0xc3}, want: Value{typ: "float", val: float32(3.14)}},
		{input: []byte{0xca, 0xbd, 0xcc, 0xcc, 0xcd}, want: Value{typ: "float", val: float32(-0.1)}},
		{input: []byte{0xca, 0x7f, 0x80, 0x00, 0x00}, want: Value{typ: "float", val: float32(math.Inf(1))}},
		{input: []byte{0xca, 0xff, 0x80, 0x00, 0x00}, want: Value{typ: "float", val: float32(math.Inf(-1))}},
		{input: []byte{0xcb, 0x3f, 0xb9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}, want: Value{typ: "float", val: float64(0.1)}},
		{input: []byte{0xcb, 0xbf, 0xb9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}, want: Value{typ: "float", val: float64(-0.1)}},
	}

	for _, tt := range MsgPackTests {
		msgpacker := NewMessagePacker(bytes.NewReader(tt.input))
		got, err := msgpacker.deserialize()
		assert.Equal(t, nil, err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v want %v", got, tt.want)
		}
	}
}

func TestMsgPackDeserializeString(t *testing.T) {

	MsgPackTests := []struct {
		input []byte
		want  Value
	}{
		{input: []byte{0xd9, 0x05, 'h', 'e', 'l', 'l', 'o'}, want: Value{typ: "string", val: "hello"}},
		{input: []byte{0xda, 0x00, 0x05, 'w', 'o', 'r', 'l', 'd'}, want: Value{typ: "string", val: "world"}},
		{input: []byte{0xdb, 0x00, 0x00, 0x00, 0x05, 't', 'e', 's', 't', '1'}, want: Value{typ: "string", val: "test1"}},
		{input: []byte{0xd9, 0x00}, want: Value{typ: "string", val: ""}},
		{input: []byte{0xda, 0x00, 0x00}, want: Value{typ: "string", val: ""}},
		{input: []byte{0xdb, 0x00, 0x00, 0x00, 0x00}, want: Value{typ: "string", val: ""}},
	}

	for _, tt := range MsgPackTests {
		msgpacker := NewMessagePacker(bytes.NewReader(tt.input))
		got, err := msgpacker.deserialize()
		assert.Equal(t, nil, err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v want %v", got, tt.want)
		}
	}
}

func TestMsgPackDeserializeArray(t *testing.T) {

	MsgPackTests := []struct {
		input []byte
		want  Value
	}{
		{input: []byte{ARRAY16, 0x00, 0x03, 0xd0, 0x01, 0xd0, 0x02, 0xd0, 0x03}, want: Value{typ: "array", val: []any{int8(1), int8(2), int8(3)}}},
		{input: []byte{ARRAY32, 0x00, 0x00, 0x00, 0x03, 0xd0, 0x01, 0xd0, 0x02, 0xd0, 0x03}, want: Value{typ: "array", val: []any{int8(1), int8(2), int8(3)}}},
	}

	for _, tt := range MsgPackTests {
		msgpacker := NewMessagePacker(bytes.NewReader(tt.input))
		got, err := msgpacker.deserialize()
		assert.Equal(t, nil, err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v want %v", got, tt.want)
		}
	}
}

func TestMsgPackDeserializeMap(t *testing.T) {
	MsgPackTests := []struct {
		input []byte
		want  Value
	}{
		{
			input: []byte{MAP16, 0x00, 0x02, 0xd0, 0x01, 0xd0, 0x10, 0xd0, 0x02, 0xd0, 0x20},
			want:  Value{typ: "map", val: map[any]any{int8(1): int8(16), int8(2): int8(32)}},
		},
		{
			input: []byte{MAP32, 0x00, 0x00, 0x00, 0x02, 0xd0, 0x03, 0xd0, 0x30, 0xd0, 0x04, 0xd0, 0x40},
			want:  Value{typ: "map", val: map[any]any{int8(3): int8(48), int8(4): int8(64)}},
		},
		{
			input: []byte{MAP16, 0x00, 0x02, 0xd9, 0x03, 'k', 'e', 'y', 0xd0, 0x05, 0xd9, 0x03, 'a', 'b', 'c', 0xd0, 0x06},
			want:  Value{typ: "map", val: map[any]any{"key": int8(5), "abc": int8(6)}},
		},
		{
			input: []byte{MAP32, 0x00, 0x00, 0x00, 0x02, 0xd9, 0x03, 'x', 'y', 'z', 0xd1, 0x00, 0x7f, 0xd9, 0x03, 'p', 'q', 'r', 0xd1, 0x01, 0x01},
			want:  Value{typ: "map", val: map[any]any{"xyz": int16(127), "pqr": int16(257)}},
		},
	}

	for _, tt := range MsgPackTests {
		msgpacker := NewMessagePacker(bytes.NewReader(tt.input))
		got, err := msgpacker.deserialize()
		assert.Equal(t, nil, err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	}
}
