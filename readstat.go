package readstat

import (
	"fmt"
	"unsafe"
)

/*
#cgo darwin  LDFLAGS: -liconv
#cgo freebsd LDFLAGS: -liconv
#cgo windows LDFLAGS: -liconv
#cgo LDFLAGS: -L./sas
#include <iconv.h>
#include <stdlib.h>
#include <errno.h>
#include "sas/readstat_sas.h"

int metadataHandler(readstat_metadata_t *metadata, void *ctx);
*/
import "C"

type Parser struct {
	parser *C.readstat_parser_t
}

//export metadataHandler
func metadataHandler(metadata *C.readstat_metadata_t, ctx unsafe.Pointer) C.int {
	fmt.Println("xxxx")
	return C.READSTAT_OK
}

func ParserInit() *Parser {
	return &Parser{C.readstat_parser_init()}
}

func (rsp *Parser) ReadstatSetMetadataHandler() {
	C.readstat_set_metadata_handler(rsp.parser, (*[0]byte)(unsafe.Pointer(C.metadataHandler)))
}

func (rsp *Parser) ParseSas7bdat(path string, md *Metadata) {
	pathb := []byte(path)
	C.readstat_parse_sas7bdat(rsp.parser, (*C.char)(unsafe.Pointer(&pathb[0])), unsafe.Pointer(md))
}
