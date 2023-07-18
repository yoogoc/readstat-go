package main

import (
	"fmt"

	"github.com/yoogoc/readstat-go"
)

func main() {
	parser := readstat.ParserInit()
	fmt.Printf("%+v\n", parser)
	parser.ReadstatSetMetadataHandler()
	metadata := &readstat.Metadata{}
	parser.ParseSas7bdat("testdata/all_types.sas7bdat", metadata)
}
