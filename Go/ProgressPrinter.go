package main

import (
	"fmt"
	"strings"
	"time"
)

type Printer interface {
	Print(str string)
	Close()
}

type PrinterChan struct {
	UpdateChan chan string
}

var ThreadPrints []string
var ThreadPrintChannels []PrinterChan

var PrinterTicker *time.Ticker

func PrintProgress() {
	for true {
		<-PrinterTicker.C

		if len(ThreadPrints) > 0 {
			fmt.Print(TERM_UP(len(ThreadPrints)))
			for i := 0; i < len(ThreadPrints); i++ {
				fmt.Println(ThreadPrints[i] + strings.Repeat(" ", Max(50-len(ThreadPrints[i]), 0)))
			}
		}

		for i := 0; i < len(ThreadPrintChannels); i++ {
			select {
			case str, more := <-ThreadPrintChannels[i].UpdateChan:
				if more {
					ThreadPrints[i] = str
				} else {
					ThreadPrintChannels = append(ThreadPrintChannels[:i], ThreadPrintChannels[i+1:]...)
					ThreadPrints = append(ThreadPrints[:i], ThreadPrints[i+1:]...)
				}
			default:
			}
		}
	}
}

func InitPrinter() {
	PrinterTicker = time.NewTicker(100 * time.Millisecond)

	go PrintProgress()
}

func NewPrinter() Printer {
	newPrinter := PrinterChan{UpdateChan: make(chan string)}

	ThreadPrints = append(ThreadPrints, "")
	ThreadPrintChannels = append(ThreadPrintChannels, newPrinter)

	return newPrinter
}

func (p PrinterChan) Print(str string) {
	p.UpdateChan <- str
}

func (p PrinterChan) Close() {
	close(p.UpdateChan)
}
