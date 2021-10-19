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
	Print(args ...interface{})
	Close()
}

type PrinterThread struct {
	CurrentStr string
	UpdateChan chan string
}

func (p PrinterThread) Print(args ...interface{}) {
	p.UpdateChan <- fmt.Sprintf(args[0].(string), args[1:]...)
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
			RefreshLines = 0
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
		//Only, display every so often
		_, more = <-PrinterTicker.C

		//Acquire Threads lock
		ThreadPrintsMutex.Lock()

		if len(ThreadPrints) > 0 {
			//Do prints!
			PrintMutex.Lock()

			//Overwrite previous section
			if RefreshLines > 0 {
				fmt.Print(TERM_UP(RefreshLines))
			}

			//Print each thread's string
			PrintLine("")
			PrintLine("==========")
			for _, thread := range ThreadPrints {
				PrintLine("| " + thread.CurrentStr)
			}
			PrintLine("==========")
			PrintLine("")

			//Count the number of lines that have been printed
			newRefreshLines := len(ThreadPrints) + 4

			//Account for when a thread has been remove
			if newRefreshLines < RefreshLines {
				for i := 0; i < RefreshLines-newRefreshLines; i++ {
					PrintLine("")
					fmt.Print(TERM_UP(1))
				}
			}
			RefreshLines = newRefreshLines

			//Done Prints
			PrintMutex.Unlock()
		}

		var newThreadPrints = make([]PrinterThread, 0)

		for i, thread := range ThreadPrints {
			select {
			case str, more := <-thread.UpdateChan:
				if more {
					ThreadPrints[i].CurrentStr = str
					newThreadPrints = append(newThreadPrints, ThreadPrints[i])
				}
			default:
				newThreadPrints = append(newThreadPrints, ThreadPrints[i])
			}
		}
		ThreadPrints = newThreadPrints

		ThreadPrintsMutex.Unlock()
	}
}

func InitPrinter() {
	TermWidth = getWidth()
	LogChan = make(chan string, 100)
	PrintMutex = new(sync.Mutex)
	ThreadPrintsMutex = new(sync.Mutex)
	PrinterDone = make(chan bool)

	PrinterTicker = time.NewTicker(100 * time.Millisecond)

	go GlobalPrint()
	go PrintProgress()
}

func NewPrinter() Printer {
	newPrinter := PrinterThread{CurrentStr: "", UpdateChan: make(chan string, 100)}

	ThreadPrintsMutex.Lock()
	ThreadPrints = append(ThreadPrints, newPrinter)
	ThreadPrintsMutex.Unlock()

	return newPrinter
}

func Print(args ...interface{}) {
	LogChan <- fmt.Sprintf(args[0].(string), args[1:]...)
}

func ClosePrinter() {
	PrinterTicker.Stop()
	Print(strings.Repeat("\n", RefreshLines))
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
