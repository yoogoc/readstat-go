package readstat

import "time"

type Compression int

const (
	CompressNone = iota
	CompressRows
	CompressBinary
)

type Endian int

const (
	EndianNone = iota
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
	TypeString = iota
	TypeInt8
	TypeInt16
	TypeInt32
	TypeFloat
	TypeDouble
	TypeStringRef
)

type VarTypeClass int

const (
	VarTypeClassString = iota
	VarTypeClassStringNumeric
)

type VarFormatClass int

const (
	VarFormatClassDate = iota
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
