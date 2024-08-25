package msgpack

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

type Value struct {
	typ string
	val any
}

type MessagePacker struct {
	reader *bufio.Reader
}

func NewMessagePacker(rd io.Reader) *MessagePacker {
	return &MessagePacker{reader: bufio.NewReader(rd)}
}

func (m *MessagePacker) deserialize() (Value, error) {
	_type, err := m.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case INT8, INT16, INT32, INT64:
		return m.deserializeInteger(_type)
	case NIL:
		return m.deserializeNil()
	case FALSE, TRUE:
		return m.deserializeBool(_type)
	case FLOAT, DOUBLE:
		return m.deserializeFloat(_type)
	case RAW8, RAW16, RAW32:
		return m.deserializeString(_type)
	case UINT8, UINT16, UINT32, UINT64:
		return m.deserializeUnsignedInteger(_type)
	case ARRAY16, ARRAY32:
		return m.deserializeArray(_type)
	case MAP16, MAP32:
		return m.deserializeMap(_type)
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, errors.New("invalid type")
	}
}

func (m *MessagePacker) readBytes(n int) (line []byte, err error) {
	for i := 0; i < n; i++ {
		b, err := m.reader.ReadByte()
		if err != nil {
			return nil, err
		}
		line = append(line, b)
	}
	return line, nil
}

func (m *MessagePacker) deserializeInteger(_type byte) (Value, error) {
	v := Value{}
	v.typ = "integer"
	switch _type {
	case INT8:
		number, err := m.reader.ReadByte()
		if err != nil {
			return v, err
		}
		if _type == INT8 {
			v.val = int8(number)
		}
		return v, nil
	case INT16:
		number, err := m.readBytes(2)
		if err != nil {
			return v, err
		}
		v.val = int16(binary.BigEndian.Uint16(number))
		return v, nil
	case INT32:
		number, err := m.readBytes(4)
		if err != nil {
			return v, err
		}
		v.val = int32(binary.BigEndian.Uint32(number))
		return v, nil
	case INT64:
		number, err := m.readBytes(8)
		if err != nil {
			return v, err
		}
		v.val = int64(binary.BigEndian.Uint64(number))
		return v, nil
	default:
		return v, fmt.Errorf("failed to decode the integer")
	}
}

func (m *MessagePacker) deserializeUnsignedInteger(_type byte) (Value, error) {
	v := Value{}
	v.typ = "Unsigned-integer"
	switch _type {
	case UINT8:
		number, err := m.reader.ReadByte()
		if err != nil {
			return v, err
		}
		v.val = uint8(number)
		return v, nil
	case UINT16:
		number, err := m.readBytes(2)
		if err != nil {
			return v, err
		}
		v.val = binary.BigEndian.Uint16(number)
		return v, nil
	case UINT32:
		number, err := m.readBytes(4)
		if err != nil {
			return v, err
		}
		v.val = binary.BigEndian.Uint32(number)
		return v, nil
	case UINT64:
		number, err := m.readBytes(8)
		if err != nil {
			return v, err
		}
		v.val = binary.BigEndian.Uint64(number)
		return v, nil
	default:
		return v, fmt.Errorf("failed to decode the integer")
	}
}

func (m *MessagePacker) deserializeNil() (Value, error) {
	return Value{typ: "nil", val: nil}, nil
}

func (m *MessagePacker) deserializeBool(_type byte) (Value, error) {
	if _type == FALSE {
		return Value{typ: "boolen", val: false}, nil
	}
	return Value{typ: "boolen", val: true}, nil
}

func (m *MessagePacker) deserializeFloat(_type byte) (Value, error) {
	v := Value{}
	v.typ = "float"
	switch _type {
	case FLOAT:
		number, err := m.readBytes(4)
		if err != nil {
			return v, err
		}
		Uint32Bits := binary.BigEndian.Uint32(number)
		v.val = math.Float32frombits(Uint32Bits)
		return v, nil
	case DOUBLE:
		number, err := m.readBytes(8)
		if err != nil {
			return v, err
		}
		Uint64Bits := binary.BigEndian.Uint64(number)
		v.val = math.Float64frombits(Uint64Bits)
		return v, nil
	default:
		return v, errors.New("failed to decode float")
	}
}

func (m *MessagePacker) deserializeString(_type byte) (Value, error) {
	v := Value{}
	v.typ = "string"

	var length int
	switch _type {
	case RAW8:
		lenByte, err := m.reader.ReadByte()
		if err != nil {
			return v, err
		}
		length = int(lenByte)
	case RAW16:
		var len16 uint16
		err := binary.Read(m.reader, binary.BigEndian, &len16)
		if err != nil {
			return v, err
		}
		length = int(len16)
	case RAW32:
		var len32 uint32
		err := binary.Read(m.reader, binary.BigEndian, &len32)
		if err != nil {
			return v, err
		}
		length = int(len32)
	default:
		return v, errors.New("failed to decode string")
	}

	data, err := m.readBytes(length)
	if err != nil {
		return v, err
	}

	v.val = string(data)
	return v, nil
}

func (m *MessagePacker) deserializeArray(_type byte) (Value, error) {
	v := Value{}
	v.typ = "array"
	var arr []any
	var length int

	switch {
	case _type >= FIXARRAY && _type <= FIXARRAY|0x0F:
		length = int(_type & 0x0F)
	case _type == ARRAY16:
		var len16 uint16
		err := binary.Read(m.reader, binary.BigEndian, &len16)
		if err != nil {
			return v, err
		}
		length = int(len16)
	case _type == ARRAY32:
		var len32 uint32
		err := binary.Read(m.reader, binary.BigEndian, &len32)
		if err != nil {
			return v, err
		}
		length = int(len32)
	default:
		return v, errors.New("failed to decode array")
	}

	for i := 0; i < length; i++ {
		element, err := m.deserialize()
		if err != nil {
			return v, err
		}
		arr = append(arr, element.val)
	}

	v.val = arr
	return v, nil
}

func (m *MessagePacker) deserializeMap(_type byte) (Value, error) {
	v := Value{}
	v.typ = "map"
	mapData := make(map[any]any)
	var length int

	switch {
	case _type >= 0x80 && _type <= 0x8F:
		length = int(_type & 0x0F)
	case _type == MAP16:
		var len16 uint16
		err := binary.Read(m.reader, binary.BigEndian, &len16)
		if err != nil {
			return v, err
		}
		length = int(len16)
	case _type == MAP32:
		var len32 uint32
		err := binary.Read(m.reader, binary.BigEndian, &len32)
		if err != nil {
			return v, err
		}
		length = int(len32)
	default:
		return v, errors.New("failed to decode map")
	}

	for i := 0; i < length; i++ {
		key, err := m.deserialize()
		if err != nil {
			return v, err
		}

		value, err := m.deserialize()
		if err != nil {
			return v, err
		}

		mapData[key.val] = value.val
	}

	v.val = mapData
	return v, nil
}
