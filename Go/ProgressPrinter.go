package main

import (
	"fmt"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

type Printer interface {
	Print(str string)
	Close()
}

type PrinterThread struct {
	CurrentStr string
	UpdateChan chan string
}

func (p PrinterThread) Print(str string) {
	p.UpdateChan <- str
}

func (p PrinterThread) Close() {
	close(p.UpdateChan)
}

var TermWidth int
var PrintMutex *sync.Mutex
var LogChan chan string
var RefreshLines = 0
var PrinterDone chan bool

var ThreadPrints []PrinterThread
var ThreadPrintsMutex *sync.Mutex

var PrinterTicker *time.Ticker

func FillRestOfLine(taken int) string {
	return strings.Repeat(" ", Max(TermWidth-taken, 0))
}

func PrintLine(str string) {
	fmt.Println(str + FillRestOfLine(len(str)))
}

func GlobalPrint() {
	for true {
		str, more := <-LogChan
		if more {
			PrintMutex.Lock()
			PrintLine(str)
			if RefreshLines > 0 {
				RefreshLines -= 1
			}
			PrintMutex.Unlock()
		} else {
			PrinterDone <- true
			return
		}
	}
}
func PrintProgress() {
	more := true
	for more {
		_, more = <-PrinterTicker.C

		ThreadPrintsMutex.Lock()
		if len(ThreadPrints) > 0 {
			PrintMutex.Lock()

			PrintLine("==========")
			for _, thread := range ThreadPrints {
				PrintLine(thread.CurrentStr)
			}
			PrintLine("==========")
			RefreshLines = len(ThreadPrints) + 2

			if RefreshLines > 0 {
				fmt.Print(TERM_UP(RefreshLines))
			}

			PrintMutex.Unlock()
		}

		for i, thread := range ThreadPrints {
			select {
			case str, more := <-thread.UpdateChan:
				if more {
					ThreadPrints[i].CurrentStr = str
				} else {
					ThreadPrints = append(ThreadPrints[:i], ThreadPrints[i+1:]...)
				}
			default:
			}
		}
		ThreadPrintsMutex.Unlock()
	}
}

func InitPrinter() {
	TermWidth = getWidth()
	LogChan = make(chan string)
	PrintMutex = new(sync.Mutex)
	ThreadPrintsMutex = new(sync.Mutex)
	PrinterDone = make(chan bool)

	PrinterTicker = time.NewTicker(100 * time.Millisecond)

	go GlobalPrint()
	go PrintProgress()
}

func NewPrinter() Printer {
	newPrinter := PrinterThread{CurrentStr: "", UpdateChan: make(chan string)}

	ThreadPrintsMutex.Lock()
	ThreadPrints = append(ThreadPrints, newPrinter)
	ThreadPrintsMutex.Unlock()

	return newPrinter
}

func Print(str string) {
	LogChan <- str
}

func ClosePrinter() {
	PrinterTicker.Stop()
	// Print(strings.Repeat("\n", RefreshLines))
	close(LogChan)
	<-PrinterDone
	close(PrinterDone)
}

//====================================================================
//https://stackoverflow.com/questions/16569433/get-terminal-size-in-go
type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getWidth() int {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return int(ws.Col)
}

//====================================================================
