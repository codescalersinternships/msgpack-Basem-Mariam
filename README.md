# msgpack-Basem-Mariam

A Go-compatible library for MsgPack serialization and deserialization

## Features
- Serialization: Convert Go data structures into MsgPack format.
- Deserialization: Convert MsgPack byte sequences back into Go data structures.

## How to Install and Run the Project

1. To use the Bencode package, simply import it into your Go project:
    ```go
    import "github.com/codescalersinternships/msgpack-Basem-Mariam"
    ```
## How to use

Convert a Go data structure into MsgPack format.

```go
    data:= []interface{}{1, 2, 3}
    result, err := Serialize(data)
    if err !=nil{
        //handles error
    }
```

Convert MsgPack byte sequences back into a Go data structure
```go
    data:= []byte{0xa5, 'h', 'e', 'l', 'l', 'o'}
    result, err := Deserialize(bytes.NewReader(data))
    if err !=nil{
        //handles error
    }
```



## Overview

### Serialization: type to format conversion

MessagePack serializers convert MessagePack types into formats as following:

source types | output format
------------ | ---------------------------------------------------------------------------------------
Integer      | int format family (int 8/16/32/64 or uint 8/16/32/64)
Nil          | nil
Boolean      | bool format family (false or true)
Float        | float format family (float 32/64)
String       | str format family (fixstr or str 8/16/32)
Array        | array format family (fixarray or array 16/32)
Map          | map format family (fixmap or map 16/32)

### Deserialization: format to type conversion

MessagePack deserializers convert MessagePack formats into types as following:

source formats                                                       | output type
-------------------------------------------------------------------- | -----------
int 8/16/32/64 and uint 8/16/32/64 | Integer
nil                                                                  | Nil
false and true                                                       | Boolean
float 32/64                                                          | Float
fixstr and str 8/16/32                                               | String
fixarray and array 16/32                                             | Array
fixmap map 16/32                                                     | Map


## Running Tests 
```sh
make test
```
    
