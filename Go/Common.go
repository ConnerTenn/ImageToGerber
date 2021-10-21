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

func Min(a int, b int) int {
	if a <= b {
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

func ProgressBar(val int, minimum int, maximum int) (string, float64) {
	barSize := TermWidth - 30
	percent := (float64(val) - float64(minimum)) / float64(maximum)
	progress := Min(int(float64(barSize)*percent), barSize)
	return "[" + strings.Repeat("=", progress) + strings.Repeat(" ", barSize-progress) + "]", percent * 100.0
}

func PrintProgressBar(name string, color string, val *int, minimum int, maximum int, printer Printer, done chan bool, resp chan bool) {
	ticker := time.NewTicker(100 * time.Millisecond)

	name = fmt.Sprintf("%-16s", name)

	bar, percent := ProgressBar(minimum, minimum, maximum)
	printer.Print(color+name+" %s  %.1f%%"+TERM_RESET, bar, percent)

	go func() {
		run := true
		for run {
			select {
			case <-ticker.C:
				bar, percent := ProgressBar(*val, minimum, maximum)
				printer.Print(color+name+" %s  %.1f%%"+TERM_RESET, bar, percent)
			case <-done:
				ticker.Stop()
				run = false
				bar, percent := ProgressBar(maximum, minimum, maximum)
				printer.Print(color+name+" %s  %.1f%%"+TERM_RESET, bar, percent)
			}
		}
		resp <- true
	}()
}
