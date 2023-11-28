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
			forp = append(forp, fmt.Sprintf("%s(%v)\n", v.Type.String(), x))
		}
		fmt.Println(forp...)
	}
}
