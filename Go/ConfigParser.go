package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Process struct {
	BoardWidth float64
	Infile     string
	Outfile    string
	Types      []string
	Selection  []Rule
}

func ParseConfig(filename string) []Process {
	file, err := os.Open(filename)
	CheckError(err)

	fmt.Println("Parsing configuration")
	fmt.Println("========")

	//State variables
	readingSelection := false
	var currProcess Process
	var processlist []Process

	scanner := bufio.NewScanner(file)
	//Loop through each line
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Split(line, "#")[0] //Trim the comments
		line = strings.Trim(line, " ")     //Trim the whitespace

		if readingSelection && strings.Contains(line, "=") {
			readingSelection = false
			processlist = append(processlist, currProcess)
		}

		if len(line) > 0 {
			if !readingSelection {
				//Normal line parse
				keyvalue := strings.Split(line, "=")
				key := keyvalue[0]
				value := keyvalue[1]
				fmt.Println(key)

				switch key {
				case "Width":
					value := strings.Split(value, "mm")
					if len(value) != 2 {
						CheckError("Invalid Width format")
					}
					boardwidth, err := strconv.ParseFloat(value[0], 64)
					CheckError(err)
					currProcess.BoardWidth = boardwidth
				case "Infile":
					currProcess.Infile = strings.Trim(value, "\"")
				case "Outfile":
					currProcess.Outfile = strings.Trim(value, "\"")
				case "Selection":
					currProcess.Types = strings.Split(value, ",")
					readingSelection = true
				default:
					CheckError("Unknown Key: " + key)
				}
			} else {
				//Reading Selection
				fmt.Println("  " + line)
			}
		}
	}
	if readingSelection {
		processlist = append(processlist, currProcess)
	}
	fmt.Println("========")
	fmt.Println(processlist)

	file.Close()
	return processlist
}
