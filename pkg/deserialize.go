package msgpack

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"math"
)

func Deserialize(reader io.Reader) (any, error) {
	bufReader := bufio.NewReader(reader)
	_type, err := bufReader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch {
	case _type == INT8 || _type == INT16 || _type == INT32 || _type == INT64:
		return DeserializeInteger(bufReader, _type)
	case _type == NIL:
		return DeserializeNil()
	case _type == FALSE || _type == TRUE:
		return DeserializeBool(_type)
	case _type == FLOAT || _type == DOUBLE:
		return DeserializeFloat(bufReader, _type)
	case (_type >= FIXRAW && _type <= 0xbf) || _type == RAW8 || _type == RAW16 || _type == RAW32:
		return DeserializeString(bufReader, _type)
	case _type == UINT8 || _type == UINT16 || _type == UINT32 || _type == UINT64:
		return DeserializeUnsignedInteger(bufReader, _type)
	case (_type >= FIXARRAY && _type <= 0x9F) || _type == ARRAY16 || _type == ARRAY32:
		return DeserializeArray(bufReader, _type)
	case (_type >= FIXMAP && _type <= 0x8F) || _type == MAP16 || _type == MAP32:
		return DeserializeMap(bufReader, _type)
	default:
		return nil, errors.New("invalid type")
	}
}

func readBytes(reader *bufio.Reader, n int) (line []byte, err error) {
	for i := 0; i < n; i++ {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		line = append(line, b)
	}
	return line, nil
}

func DeserializeInteger(reader *bufio.Reader, _type byte) (any, error) {
	switch _type {
	case INT8:
		number, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		return int8(number), nil
	case INT16:
		number, err := readBytes(reader, 2)
		if err != nil {
			return nil, err
		}
		return int16(binary.BigEndian.Uint16(number)), nil
	case INT32:
		number, err := readBytes(reader, 4)
		if err != nil {
			return nil, err
		}
		return int32(binary.BigEndian.Uint32(number)), nil
	case INT64:
		number, err := readBytes(reader, 8)
		if err != nil {
			return nil, err
		}
		return int64(binary.BigEndian.Uint64(number)), nil
	default:
		return nil, errors.New("failed to decode the integer")
	}
}

func DeserializeUnsignedInteger(reader *bufio.Reader, _type byte) (any, error) {
	switch _type {
	case UINT8:
		number, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		return uint8(number), nil
	case UINT16:
		number, err := readBytes(reader, 2)
		if err != nil {
			return nil, err
		}
		return binary.BigEndian.Uint16(number), nil
	case UINT32:
		number, err := readBytes(reader, 4)
		if err != nil {
			return nil, err
		}
		return binary.BigEndian.Uint32(number), nil
	case UINT64:
		number, err := readBytes(reader, 8)
		if err != nil {
			return nil, err
		}
		return binary.BigEndian.Uint64(number), nil
	default:
		return nil, errors.New("failed to decode the integer")
	}
}

func DeserializeNil() (any, error) {
	return nil, nil
}

func DeserializeBool(_type byte) (any, error) {
	if _type == FALSE {
		return false, nil
	}
	return true, nil
}

func DeserializeFloat(reader *bufio.Reader, _type byte) (any, error) {
	switch _type {
	case FLOAT:
		number, err := readBytes(reader, 4)
		if err != nil {
			return nil, err
		}
		Uint32Bits := binary.BigEndian.Uint32(number)
		return math.Float32frombits(Uint32Bits), nil
	case DOUBLE:
		number, err := readBytes(reader, 8)
		if err != nil {
			return nil, err
		}
		Uint64Bits := binary.BigEndian.Uint64(number)
		return math.Float64frombits(Uint64Bits), nil
	default:
		return nil, errors.New("failed to decode float")
	}
}

func DeserializeString(reader *bufio.Reader, _type byte) (any, error) {
	var length int
	switch {
	case _type >= FIXRAW && _type <= 0xbf:
		length = int(_type & 0x1f)
		if length > 31 {
			return nil, errors.New("invalid string length")
		}
	case _type == RAW8:
		lenByte, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		length = int(lenByte)
	case _type == RAW16:
		var len16 uint16
		err := binary.Read(reader, binary.BigEndian, &len16)
		if err != nil {
			return nil, err
		}
		length = int(len16)
	case _type == RAW32:
		var len32 uint32
		err := binary.Read(reader, binary.BigEndian, &len32)
		if err != nil {
			return nil, err
		}
		length = int(len32)
	default:
		return nil, errors.New("failed to decode string")
	}

	data, err := readBytes(reader, length)
	if err != nil {
		return nil, err
	}

	return string(data), nil
}

func DeserializeArray(reader *bufio.Reader, _type byte) (any, error) {
	var arr []any
	var length int

	switch {
	case _type >= FIXARRAY && _type <= 0x9F:
		length = int(_type & 0x0F)
	case _type == ARRAY16:
		var len16 uint16
		err := binary.Read(reader, binary.BigEndian, &len16)
		if err != nil {
			return nil, err
		}
		length = int(len16)
	case _type == ARRAY32:
		var len32 uint32
		err := binary.Read(reader, binary.BigEndian, &len32)
		if err != nil {
			return nil, err
		}
		length = int(len32)
	default:
		return nil, errors.New("failed to decode array")
	}

	for i := 0; i < length; i++ {
		element, err := Deserialize(reader)
		if err != nil {
			return nil, err
		}
		arr = append(arr, element)
	}

	return arr, nil
}

func DeserializeMap(reader *bufio.Reader, _type byte) (any, error) {
	mapData := make(map[any]any)
	var length int

	switch {
	case _type >= FIXMAP && _type <= 0x8F:
		length = int(_type & 0x0F)
	case _type == MAP16:
		var len16 uint16
		err := binary.Read(reader, binary.BigEndian, &len16)
		if err != nil {
			return nil, err
		}
		length = int(len16)
	case _type == MAP32:
		var len32 uint32
		err := binary.Read(reader, binary.BigEndian, &len32)
		if err != nil {
			return nil, err
		}
		length = int(len32)
	default:
		return nil, errors.New("failed to decode map")
	}

	for i := 0; i < length; i++ {
		key, err := Deserialize(reader)
		if err != nil {
			return nil, err
		}

		value, err := Deserialize(reader)
		if err != nil {
			return nil, err
		}

		mapData[key] = value
	}

	return mapData, nil
}
