package readstat

import (
	"fmt"
	"time"
	"unsafe"
)

/*
#cgo darwin  LDFLAGS: -liconv
#cgo freebsd LDFLAGS: -liconv
#cgo windows LDFLAGS: -liconv
#include <iconv.h>
#include <stdlib.h>
#include <errno.h>
#include "readstat_sas.h"

int goMetadataHandler(readstat_metadata_t *metadata, void *ctx);
int goVariableHandler(int index, readstat_variable_t *variable, char *val_labels, void *ctx);
int goValueHandler(int obs_index, readstat_variable_t *variable, readstat_value_t value, void *ctx);
*/
import "C"

type Parser struct {
	parser *C.readstat_parser_t
}

//export goMetadataHandler
func goMetadataHandler(metadata *C.readstat_metadata_t, ctx unsafe.Pointer) C.int {
	target := (*Metadata)(ctx)
	target.RowCount = int64(C.readstat_get_row_count(metadata))
	target.VarCount = int64(C.readstat_get_var_count(metadata))
	target.CreationTime = time.Unix(int64(int32(C.readstat_get_creation_time(metadata))), 0)
	target.ModifiedTime = time.Unix(int64(int32(C.readstat_get_modified_time(metadata))), 0)
	target.FileFormatVersion = int64(C.readstat_get_file_format_version(metadata))
	target.Compression = Compression(C.readstat_get_compression(metadata))
	target.Endianness = Endian(C.readstat_get_endianness(metadata))
	target.TableName = C.GoString(C.readstat_get_table_name(metadata))
	target.FileLabel = C.GoString(C.readstat_get_file_label(metadata))
	target.FileEncoding = C.GoString(C.readstat_get_file_encoding(metadata))
	target.Is64bit = int(C.readstat_get_file_format_is_64bit(metadata)) == 1
	return C.READSTAT_OK
}

//export goVariableHandler
func goVariableHandler(index C.int, variable *C.readstat_variable_t, val_labels *C.char, ctx unsafe.Pointer) C.int {
	target := (*Metadata)(ctx)
	if target.Vars == nil {
		target.Vars = make(map[int]VarMetadata)
	}
	v := VarMetadata{}
	v.VarName = C.GoString(C.readstat_variable_get_name(variable))
	v.VarType = VarType(C.readstat_variable_get_type(variable))
	v.VarTypeClass = VarTypeClass(C.readstat_variable_get_type_class(variable))
	v.VarLabel = C.GoString(C.readstat_variable_get_label(variable))
	v.VarFormat = C.GoString(C.readstat_variable_get_format(variable))
	v.VarFormatClass = getFormat(v.VarFormat)

	target.Vars[int(index)] = v

	return C.READSTAT_OK
}

//export goValueHandler
func goValueHandler(obsIndex C.int, variable *C.readstat_variable_t, value C.readstat_value_t, ctx unsafe.Pointer) C.int {
	target := (*Data)(ctx)
	varIndex := int(C.readstat_variable_get_index(variable))
	valueType := VarType(C.readstat_value_type(value))
	isMissing := int(C.readstat_value_is_system_missing(value)) == 1
	target.Vars[varIndex] = append(target.Vars[varIndex], getValue(value, valueType, isMissing, target.Metadata.Vars[varIndex].VarFormatClass))
	return C.READSTAT_OK
}

func getValue(value C.readstat_value_t, valueType VarType, missing bool, varFormatClass *VarFormatClass) *Var {
	v := &Var{}
	// Set Type
	switch valueType {
	case TypeString, TypeStringRef:
		v.Type = ValueTypeString
	case TypeInt8:
		v.Type = ValueTypeI8
	case TypeInt16:
		v.Type = ValueTypeI16
	case TypeInt32:
		v.Type = ValueTypeI32
	case TypeFloat:
		v.Type = ValueTypeF32
	case TypeDouble:
		if varFormatClass == nil {
			v.Type = ValueTypeF64
		} else {
			switch *varFormatClass {
			case VarFormatClassDate:
				v.Type = ValueTypeDate
			case VarFormatClassDateTime:
				v.Type = ValueTypeDateTime
			case VarFormatClassDateTimeWithMilliseconds:
				v.Type = ValueTypeDateTimeWithMilliseconds
			case VarFormatClassDateTimeWithMicroseconds:
				v.Type = ValueTypeDateTimeWithMicroseconds
			case VarFormatClassDateTimeWithNanoseconds:
				v.Type = ValueTypeDateTimeWithNanoseconds
			case VarFormatClassTime:
				v.Type = ValueTypeTime
			default:
				fmt.Println("unknown varFormatClass", varFormatClass)
			}
		}
	default:
		fmt.Println("unknown valueType", valueType)
	}
	if missing {
		return v
	}

	// Set Value
	switch valueType {
	case TypeString, TypeStringRef:
		vv := C.GoString(C.readstat_string_value(value))
		v.StringValue = &vv
	case TypeInt8:
		vv := int8(C.readstat_int8_value(value))
		v.I8Value = &vv
	case TypeInt16:
		vv := int16(C.readstat_int16_value(value))
		v.I16Value = &vv
	case TypeInt32:
		vv := int32(C.readstat_int32_value(value))
		v.I32Value = &vv
	case TypeFloat:
		vv := float32(C.readstat_float_value(value))
		v.F32Value = &vv
	case TypeDouble:
		vv := float64(C.readstat_double_value(value))
		if varFormatClass == nil {
			v.F64Value = &vv
		} else {
			vvv := time.Unix(int64(vv), 0)
			switch *varFormatClass {
			case VarFormatClassDate:
				v.DateValue = &vvv
			case VarFormatClassDateTime:
				v.DateTimeValue = &vvv
			case VarFormatClassDateTimeWithMilliseconds:
				v.DateTimeWithMillisecondsValue = &vvv
			case VarFormatClassDateTimeWithMicroseconds:
				v.DateTimeWithMicrosecondsValue = &vvv
			case VarFormatClassDateTimeWithNanoseconds:
				v.DateTimeWithNanosecondsValue = &vvv
			case VarFormatClassTime:
				v.TimeValue = &vvv
			}
		}
	}

	return v
}

func ParserInit() *Parser {
	p := &Parser{C.readstat_parser_init()}
	return p
}

// ReadstatSetMetadataHandler TODO conv error
func (rsp *Parser) ReadstatSetMetadataHandler() {
	C.readstat_set_metadata_handler(rsp.parser, (*[0]byte)(unsafe.Pointer(C.goMetadataHandler)))
}

// ReadstatSetVariableHandler TODO conv error
func (rsp *Parser) ReadstatSetVariableHandler() {
	C.readstat_set_variable_handler(rsp.parser, (*[0]byte)(unsafe.Pointer(C.goVariableHandler)))
}

// ReadstatSetValueHandler TODO conv error
func (rsp *Parser) ReadstatSetValueHandler() {
	C.readstat_set_value_handler(rsp.parser, (*[0]byte)(unsafe.Pointer(C.goValueHandler)))
}

// ReadstatSetRowLimit TODO conv error
func (rsp *Parser) ReadstatSetRowLimit(limit int32) {
	C.readstat_set_row_limit(rsp.parser, C.long(limit))
}

// ReadstatSetRowOffset TODO conv error
func (rsp *Parser) ReadstatSetRowOffset(offset int32) {
	C.readstat_set_row_offset(rsp.parser, C.long(offset))
}

func (rsp *Parser) ParseSas7bdatMetadata(path string, md *Metadata) {
	pathb := []byte(path)
	C.readstat_parse_sas7bdat(rsp.parser, (*C.char)(unsafe.Pointer(&pathb[0])), unsafe.Pointer(md))
}

func (rsp *Parser) ParseSas7bdatData(path string, data *Data) {
	pathb := []byte(path)
	C.readstat_parse_sas7bdat(rsp.parser, (*C.char)(unsafe.Pointer(&pathb[0])), unsafe.Pointer(data))
}

func (rsp *Parser) Close() {
	C.readstat_parser_free(rsp.parser)
}
