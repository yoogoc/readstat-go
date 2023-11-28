package readstat

import (
	"time"
)

type Compression int

const (
	CompressNone Compression = iota
	CompressRows
	CompressBinary
)

type Endian int

const (
	EndianNone Compression = iota
	EndianLittle
	EndianBig
)

type Metadata struct {
	RowCount          int64
	VarCount          int64
	CreationTime      time.Time
	ModifiedTime      time.Time
	FileFormatVersion int64
	Compression       Compression
	Endianness        Endian
	TableName         string
	FileLabel         string
	FileEncoding      string
	Is64bit           bool
	Vars              map[int]VarMetadata
}

type VarType int

const (
	TypeString VarType = iota
	TypeInt8
	TypeInt16
	TypeInt32
	TypeFloat
	TypeDouble
	TypeStringRef
)

type VarTypeClass int

const (
	VarTypeClassString VarTypeClass = iota
	VarTypeClassStringNumeric
)

type VarFormatClass int

const (
	VarFormatClassDate VarFormatClass = iota
	VarFormatClassDateTime
	VarFormatClassDateTimeWithMilliseconds
	VarFormatClassDateTimeWithMicroseconds
	VarFormatClassDateTimeWithNanoseconds
	VarFormatClassTime
)

type VarMetadata struct {
	VarName        string
	VarType        VarType
	VarTypeClass   VarTypeClass
	VarLabel       string
	VarFormat      string
	VarFormatClass *VarFormatClass
}

type Data struct {
	Metadata *Metadata
	Vars     [][]*Var
}

//go:generate go run -mod=mod golang.org/x/tools/cmd/stringer -type=ValueType
type ValueType int

const (
	ValueTypeString ValueType = iota
	ValueTypeI8
	ValueTypeI16
	ValueTypeI32
	ValueTypeF32
	ValueTypeF64
	ValueTypeDate
	ValueTypeDateTime
	ValueTypeDateTimeWithMilliseconds
	ValueTypeDateTimeWithMicroseconds
	ValueTypeDateTimeWithNanoseconds
	ValueTypeTime
)

type Var struct {
	Type  ValueType
	Value any
}
