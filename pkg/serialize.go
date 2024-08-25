// package msgpack is a package that provides serialization and deserialization of data in the MessagePack format.
package msgpack

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// serialzieBool serializes a boolean value
func serializeBool(b bool) ([]byte, error) {

	buf := new(bytes.Buffer)
	if b {
		buf.WriteByte(TRUE)
	} else {
		buf.WriteByte(FALSE)
	}

	return buf.Bytes(), nil
}

// serializeNil serializes a nil value
func serializeNil() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(NIL)
	return buf.Bytes(), nil
}

// serializeUint8 serializes a uint8 value
func serializeUint8(n uint8) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(UINT8)
	buf.WriteByte(n)
	return buf.Bytes(), nil
}

// serializeUint16 serializes a uint16 value
func serializeUint16(n uint16) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(UINT16)
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes(), nil
}

// serializeUint32 serializes a uint32 value
func serializeUint32(n uint32) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(UINT32)
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes(), nil
}

// serializeUint64 serializes a uint64 value
func serializeUint64(n uint64) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(UINT64)
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes(), nil
}

// serializeInt8 serializes an int8 value
func serializeInt8(n int8) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(INT8)
	buf.WriteByte(byte(n))
	return buf.Bytes(), nil
}

// serializeInt16 serializes an int16 value
func serializeInt16(n int16) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(INT16)
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes(), nil
}

// serializeInt32 serializes an int32 value
func serializeInt32(n int32) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(INT32)
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes(), nil
}

// serializeInt64 serializes an int64 value
func serializeInt64(n int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(INT64)
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes(), nil
}

// serializeFloat32 serializes a float32 value
func serializeFloat32(n float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(FLOAT)
	binary.Write(buf, binary.BigEndian, float32(n))
	return buf.Bytes(), nil
}

// serializeFloat64 serializes a float64 value
func serializeFloat64(n float64) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte(DOUBLE)
	binary.Write(buf, binary.BigEndian, float64(n))
	return buf.Bytes(), nil
}

// serializeString serializes a string value
func serializeString(s string) ([]byte, error) {
	buf := new(bytes.Buffer)
	strLen := len(s)

	if strLen <= 31 {
		buf.WriteByte(byte(FIXRAW | strLen)) // FixStr
	} else if strLen <= 255 {
		buf.WriteByte(RAW8) // str8
		buf.WriteByte(byte(strLen))
	} else if strLen <= 65535 {
		buf.WriteByte(RAW16) // str16
		binary.Write(buf, binary.BigEndian, uint16(strLen))
	} else if strLen <= 4294967295 {
		buf.WriteByte(RAW32) // str32
		binary.Write(buf, binary.BigEndian, uint32(strLen))
	} else {
		return nil, fmt.Errorf("string too long to serialize")
	}

	buf.WriteString(s)
	return buf.Bytes(), nil
}

// serializeArray serializes an array value
func serializeArray(arr interface{}) ([]byte, error) {
	val := reflect.ValueOf(arr)

	buf := new(bytes.Buffer)
	length := val.Len()

	if length <= 15 {
		buf.WriteByte(byte(FIXARRAY | length)) // FixArray
	} else if length <= 65535 {
		buf.WriteByte(ARRAY16) // array16
		binary.Write(buf, binary.BigEndian, uint16(length))
	} else if length <= 4294967295 {
		buf.WriteByte(ARRAY32) // array32
		binary.Write(buf, binary.BigEndian, uint32(length))
	} else {
		return nil, fmt.Errorf("array too long to serialize")
	}

	for i := 0; i < length; i++ {
		elem, err := Serialize(val.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		buf.Write(elem)
	}

	return buf.Bytes(), nil
}

// serializeMap serializes a map value
func serializeMap(m interface{}) ([]byte, error) {
	val := reflect.ValueOf(m)

	buf := new(bytes.Buffer)
	length := val.Len()

	if length <= 15 {
		buf.WriteByte(byte(FIXMAP | length)) // FixMap
	} else if length <= 65535 {
		buf.WriteByte(MAP16) // map16
		binary.Write(buf, binary.BigEndian, uint16(length))
	} else if length <= 4294967295 {
		buf.WriteByte(MAP32) // map32
		binary.Write(buf, binary.BigEndian, uint32(length))
	} else {
		return nil, fmt.Errorf("map too long to serialize")
	}

	for _, key := range val.MapKeys() {
		keySerialized, err := Serialize(key.Interface())
		if err != nil {
			return nil, err
		}
		valueSerialized, err := Serialize(val.MapIndex(key).Interface())
		if err != nil {
			return nil, err
		}
		buf.Write(keySerialized)
		buf.Write(valueSerialized)
	}

	return buf.Bytes(), nil
}

// Serialize can be given a value of any type and will serialize it to a byte slice
func Serialize(v interface{}) ([]byte, error) {

	if v == nil {
		return serializeNil()
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.Bool:
		return serializeBool(v.(bool))
	case reflect.Uint8:
		return serializeUint8(v.(uint8))
	case reflect.Uint16:
		return serializeUint16(v.(uint16))
	case reflect.Uint32:
		return serializeUint32(v.(uint32))
	case reflect.Uint64:
		return serializeUint64(v.(uint64))
	case reflect.Uint:
		return serializeUint64(uint64(v.(uint)))
	case reflect.Int8:
		return serializeInt8(v.(int8))
	case reflect.Int16:
		return serializeInt16(v.(int16))
	case reflect.Int32:
		return serializeInt32(v.(int32))
	case reflect.Int64:
		return serializeInt64(v.(int64))
	case reflect.Int:
		return serializeInt64(int64(v.(int)))
	case reflect.Float32:
		return serializeFloat32(v.(float32))
	case reflect.Float64:
		return serializeFloat64(v.(float64))
	case reflect.String:
		return serializeString(v.(string))
	case reflect.Slice, reflect.Array:
		return serializeArray(v)
	case reflect.Map:
		return serializeMap(v)
	default:
		fmt.Print(reflect.TypeOf(v).Kind())
		return nil, fmt.Errorf("unsupported type")
	}
}
