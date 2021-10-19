package main

import (
	"fmt"
	"os"
	"strings"
	"time"
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

func TERM_UP(count int) string {
	if count > 0 {
		return fmt.Sprintf("\033[%dF", count)
	}
	return ""
}

func Max(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

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

func ProgressBar(val int, minimum int, maximum int) string {
	barSize := TermWidth - 40
	progress := (barSize * (val - minimum)) / maximum
	return "[" + strings.Repeat("=", progress) + strings.Repeat(" ", barSize-progress) + "]"
}

func PrintProgressBar(name string, color string, val *int, minimum int, maximum int, printer Printer, done chan bool, resp chan bool) {
	ticker := time.NewTicker(100 * time.Millisecond)
	tStart := time.Now()

	name = fmt.Sprintf("%-16s", name)

	bar := ProgressBar(minimum, minimum, maximum)
	printer.Print(color+name+" %s Time:%v"+TERM_RESET, bar, "---")

	go func() {
		run := true
		for run {
			select {
			case t := <-ticker.C:
				bar := ProgressBar(*val, minimum, maximum)
				printer.Print(color+name+" %s Time:%v"+TERM_RESET, bar, t.Sub(tStart).Truncate(10*time.Millisecond))
			case <-done:
				ticker.Stop()
				run = false
				bar := ProgressBar(maximum, minimum, maximum)
				printer.Print(color+name+" %s Time:%v"+TERM_RESET, bar, time.Now().Sub(tStart).Truncate(10*time.Millisecond))
			}
		}
		resp <- true
	}()
}
