package main

import (
	"fmt"

	"github.com/yoogoc/readstat-go"
)

func main() {
	parser := readstat.ParserInit()
	fmt.Printf("%+v\n", parser)
	parser.ReadstatSetMetadataHandler()
	parser.ReadstatSetVariableHandler()
	metadata := &readstat.Metadata{}
	parser.ParseSas7bdatMetadata("testdata/all_types.sas7bdat", metadata)
	fmt.Printf("%+v\n", metadata)
	data := &readstat.Data{Metadata: metadata}
	data.Vars = make([][]*readstat.Var, metadata.VarCount)
	for i := 0; i < int(metadata.VarCount); i++ {
		data.Vars[i] = make([]*readstat.Var, 0, metadata.RowCount)
	}
	dparser := readstat.ParserInit()
	dparser.ReadstatSetValueHandler()
	dparser.ReadstatSetRowLimit(int32(metadata.RowCount))
	dparser.ReadstatSetRowOffset(0)
	dparser.ParseSas7bdatData("testdata/all_types.sas7bdat", data)
	for i, vars := range data.Vars {
		var forp []any
		forp = append(forp, i)
		for _, v := range vars {
			var x any
			if v == nil {
				continue
			}
			switch v.Type {
			case readstat.ValueTypeString:
				if v.StringValue != nil {
					x = *v.StringValue
				}
			case readstat.ValueTypeI8:
				if v.I8Value != nil {
					x = *v.I8Value
				}
			case readstat.ValueTypeI16:
				if v.I16Value != nil {
					x = *v.I16Value
				}
			case readstat.ValueTypeI32:
				if v.I32Value != nil {
					x = *v.I32Value
				}
			case readstat.ValueTypeF32:
				if v.F32Value != nil {
					x = *v.F32Value
				}
			case readstat.ValueTypeF64:
				if v.F64Value != nil {
					x = *v.F64Value
				}
			case readstat.ValueTypeDate:
				if v.DateValue != nil {
					x = *v.DateValue
				}
			case readstat.ValueTypeDateTime:
				if v.DateTimeValue != nil {
					x = *v.DateTimeValue
				}
			case readstat.ValueTypeDateTimeWithMilliseconds:
				if v.DateTimeWithMillisecondsValue != nil {
					x = *v.DateTimeWithMillisecondsValue
				}
			case readstat.ValueTypeDateTimeWithMicroseconds:
				if v.DateTimeWithMicrosecondsValue != nil {
					x = *v.DateTimeWithMicrosecondsValue
				}
			case readstat.ValueTypeDateTimeWithNanoseconds:
				if v.DateTimeWithNanosecondsValue != nil {
					x = *v.DateTimeWithNanosecondsValue
				}
			case readstat.ValueTypeTime:
				if v.TimeValue != nil {
					x = *v.TimeValue
				}
			}

			forp = append(forp, fmt.Sprintf("%s(%v)", v.Type.String(), x))
		}
		fmt.Println(forp...)
	}
}
