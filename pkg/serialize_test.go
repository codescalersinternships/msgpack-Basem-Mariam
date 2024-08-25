package msgpack

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestSerialize(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []byte
		wantErr  bool
	}{
		{"nil", nil, []byte{NIL}, false},
		{"bool true", true, []byte{TRUE}, false},
		{"bool false", false, []byte{FALSE}, false},
		{"uint8", uint8(42), []byte{UINT8, 42}, false},
		{"uint16", uint16(42), []byte{UINT16, 0, 42}, false},
		{"uint32", uint32(42), []byte{UINT32, 0, 0, 0, 42}, false},
		{"uint64", uint64(42), []byte{UINT64, 0, 0, 0, 0, 0, 0, 0, 42}, false},
		{"int8", int8(42), []byte{INT8, 42}, false},
		{"int16", int16(42), []byte{INT16, 0, 42}, false},
		{"int32", int32(42), []byte{INT32, 0, 0, 0, 42}, false},
		{"int64", int64(42), []byte{INT64, 0, 0, 0, 0, 0, 0, 0, 42}, false},
		{"float32", float32(42.0), func() []byte {
			buf := new(bytes.Buffer)
			buf.WriteByte(FLOAT)
			err := binary.Write(buf, binary.BigEndian, float32(42.0))
			if err != nil {
				t.Errorf("Error writing float32: %v", err)
			}
			return buf.Bytes()
		}(), false},
		{"float64", float64(42.0), func() []byte {
			buf := new(bytes.Buffer)
			buf.WriteByte(DOUBLE)
			err := binary.Write(buf, binary.BigEndian, float64(42.0))
			if err != nil {
				t.Errorf("Error writing float64: %v", err)
			}
			return buf.Bytes()
		}(), false},
		{"string", "hello", append([]byte{byte(FIXRAW | 5)}, []byte("hello")...), false},
		{"slice", []interface{}{1, 2, 3}, []byte{FIXARRAY | 3, 211, 0, 0, 0, 0, 0, 0, 0, 1, 211, 0, 0, 0, 0, 0, 0, 0, 2, 211, 0, 0, 0, 0, 0, 0, 0, 3}, false},
		{"map", map[interface{}]interface{}{"key": "value"}, []byte{FIXMAP | 1, byte(FIXRAW | 3), 'k', 'e', 'y', byte(FIXRAW | 5), 'v', 'a', 'l', 'u', 'e'}, false},
		{"unsupported type", struct{}{}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Serialize(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(result, tt.expected) {
				t.Errorf("Serialize() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
