package msgpack


const (
	NIL = 0xc0

	FALSE = 0xc2
	TRUE  = 0xc3

	FLOAT  = 0xca
	DOUBLE = 0xcb

	UINT8  = 0xcc
	UINT16 = 0xcd
	UINT32 = 0xce
	UINT64 = 0xcf
	INT8   = 0xd0
	INT16  = 0xd1
	INT32  = 0xd2
	INT64  = 0xd3

	RAW16   = 0xda
	RAW32   = 0xdb
	ARRAY16 = 0xdc
	ARRAY32 = 0xdd
	MAP16   = 0xde
	MAP32   = 0xdf

	FIXMAP   = 0x80
	FIXARRAY = 0x90
	FIXRAW   = 0xa0

	MAXFIXMAP   = 16
	MAXFIXARRAY = 16
	MAXFIXRAW   = 32

	LEN_INT32 = 4
	LEN_INT64 = 8

	MAX16BIT = 2 << (16 - 1)

	REGULAR_UINT7_MAX  = 2 << (7 - 1)
	REGULAR_UINT8_MAX  = 2 << (8 - 1)
	REGULAR_UINT16_MAX = 2 << (16 - 1)
	REGULAR_UINT32_MAX = 2 << (32 - 1)

	SPECIAL_INT8  = 32
	SPECIAL_INT16 = 2 << (8 - 2)
	SPECIAL_INT32 = 2 << (16 - 2)
	SPECIAL_INT64 = 2 << (32 - 2)
)

