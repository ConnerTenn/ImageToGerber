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
var RefreshLines = 0

var ThreadPrints []PrinterThread
var ThreadPrintsMutex *sync.Mutex

var PrinterTicker *time.Ticker
var PrinterDone chan bool
var AllPrintThreadsDone chan bool

func FillRestOfLine(taken int) string {
	return strings.Repeat(" ", Max(TermWidth-taken, 0))
}

func PrintLine(str string) {
	printlength := 0
	escapeSequ := false
	for _, chr := range str {
		//Check for escape sequence
		if chr == '\033' {
			escapeSequ = true
		} else if escapeSequ {
			//Check for end of escape sequence
			if ('a' < chr && chr < 'z') || ('A' < chr && chr < 'Z') {
				escapeSequ = false
			}
		} else {
			printlength++
		}
	}
	fmt.Println(str + FillRestOfLine(printlength))
}

func PrintProgress() {
	done := false
	requestDone := false
	for !done {
		//Only, display every so often
		select {
		case _ = <-PrinterTicker.C:
		case requestDone = <-PrinterDone:
		}

		//Acquire Threads lock
		ThreadPrintsMutex.Lock()

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

		{
			//Do prints!
			PrintMutex.Lock()

			//Print each thread's string
			PrintLine("")
			PrintLine(TERM_CYAN + strings.Repeat("=", TermWidth) + TERM_RESET)
			for _, thread := range ThreadPrints {
				PrintLine(TERM_CYAN + "| " + TERM_RESET + thread.CurrentStr)
			}
			PrintLine(TERM_CYAN + strings.Repeat("=", TermWidth) + TERM_RESET)
			PrintLine("")

			//Count the number of lines that have been printed
			newRefreshLines := len(ThreadPrints) + 4

			//Account for when a thread has been removed
			if newRefreshLines < RefreshLines {
				for i := 0; i < RefreshLines-newRefreshLines; i++ {
					PrintLine("")
				}
				for i := 0; i < RefreshLines-newRefreshLines; i++ {
					fmt.Print(TERM_UP(1))
				}
			}

			RefreshLines = newRefreshLines

			//Overwrite previous section
			if RefreshLines > 0 {
				fmt.Print(TERM_UP(RefreshLines))
			}

			//Done Prints
			PrintMutex.Unlock()
		}

		if len(ThreadPrints) == 0 && requestDone {
			AllPrintThreadsDone <- true
			done = true
		}

		ThreadPrintsMutex.Unlock()
	}
}

func InitPrinter() {
	TermWidth = getWidth()
	PrintMutex = new(sync.Mutex)
	ThreadPrintsMutex = new(sync.Mutex)

	PrinterTicker = time.NewTicker(100 * time.Millisecond)
	PrinterDone = make(chan bool)
	AllPrintThreadsDone = make(chan bool)

	// go GlobalPrint()
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
	str := fmt.Sprintf(args[0].(string), args[1:]...)

	PrintMutex.Lock()
	PrintLine(str)
	RefreshLines = 0
	PrintMutex.Unlock()
}

func ClosePrinter() {
	//Stop the ticker
	PrinterTicker.Stop()
	//Request done
	PrinterDone <- true
	close(PrinterDone)
	//Wait till it actually finishes
	<-AllPrintThreadsDone
	close(AllPrintThreadsDone)

	//Flush till end of thread print section
	Print(strings.Repeat("\n", RefreshLines))
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
