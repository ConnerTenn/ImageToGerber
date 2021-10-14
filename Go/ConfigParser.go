package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ParseConfig(filename string) []Rule {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println("Parsing configuration")
	fmt.Println("========")
	scanner := bufio.NewScanner(file)
	readingSelection := false
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Split(line, "#")[0]
		line = strings.Trim(line, " ")

		if len(line) > 0 {
			if !readingSelection {
				//Normal line parse
				keyvalue := strings.Split(line, "=")
				key := keyvalue[0]
				// value := keyvalue[1]
				fmt.Println(key)

				switch key {
				case "Width":
				case "Infile":
				case "Outfile":
				case "Selection":
					readingSelection = true
				default:
					// panic("Unknown Key \"" + key + "\"")
				}
			} else {
				//Reading Selection
				if strings.Contains(line, "=") {
					readingSelection = false
				} else {
					fmt.Println("  " + line)
				}
			}
		}
	}
	fmt.Println("========")

	file.Close()
	return []Rule{}
}
