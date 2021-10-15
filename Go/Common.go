package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	TERM_RESET   = "\033[m"
	TERM_GREY    = "\033[1;30m"
	TERM_RED     = "\033[1;31m"
	TERM_GREEN   = "\033[1;32m"
	TERM_YELLOW  = "\033[1;33m"
	TERM_BLUE    = "\033[1;34m"
	TERM_MAGENTA = "\033[1;35m"
	TERM_CYAN    = "\033[1;36m"
	TERM_WHITE   = "\033[1;37m"
)

func CheckError(err interface{}) {
	if err != nil {
		fmt.Println(TERM_RED+"Error:"+TERM_RESET, err)
		os.Exit(-1)
	}
}

func CreateFile(filepath string) *os.File {
	directory := filepath[:strings.LastIndex(filepath, "/")]
	if len(directory) > 0 {
		os.MkdirAll(directory, 0755)
	}
	oFile, err := os.Create(filepath)
	CheckError(err)
	return oFile
}
