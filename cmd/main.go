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
	parser.ParseSas7bdat("testdata/all_types.sas7bdat", metadata)
	fmt.Printf("%+v\n", metadata)
}
