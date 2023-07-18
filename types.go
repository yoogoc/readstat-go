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
}
