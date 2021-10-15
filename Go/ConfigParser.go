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

		//Check if no longer reading selection
		if readingSelection && strings.Contains(line, "=") {
			readingSelection = false
			//Append the new process
			processlist = append(processlist, currProcess)
			//Reset selection but keep other settings
			currProcess.Selection = make([]Rule, 0)
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

				//parsing the Inversion specifier
				var newrule Rule
				if strings.HasPrefix(line, "!") {
					newrule.Inv = true
					line = strings.Split(line, "!")[1]
				} else if strings.Count(line, "!") != 0 {
					CheckError("Invalid rule format")
				} else {
					newrule.Inv = false
				}

				//Split the rule into conditions
				conditions := strings.Split(line, "&")
				//Parse each condition in this rule
				for _, cond := range conditions {
					cond = strings.Trim(cond, " ")

					//Error checking
					if strings.Count(cond, "(") != 1 {
						CheckError("Invalid Condition Format")
					}
					if strings.Count(cond, ")") != 1 {
						CheckError("Invalid Condition Format")
					}
					if strings.Count(cond, "+") != 1 {
						CheckError("Invalid Condition Format")
					}
					if strings.Count(cond, "-") != 1 {
						CheckError("Invalid Condition Format")
					}
					idxOpen := strings.Index(cond, "(")
					idxClose := strings.Index(cond, ")")
					idxPlus := strings.Index(cond, "+")
					idxMinus := strings.Index(cond, "-")
					if !(idxOpen < idxClose && idxClose < idxPlus && idxPlus < idxMinus) {
						CheckError("Invalid Condition Format")
					}

					//Split the conditions into tokens
					cond = strings.ReplaceAll(cond, "(", "#")
					cond = strings.ReplaceAll(cond, ")", "#")
					cond = strings.ReplaceAll(cond, "+", "#")
					cond = strings.ReplaceAll(cond, "-", "#")
					tokens := strings.Split(cond, "#")

					if len(tokens[2]) != 0 {
						CheckError("Invalid Condition Format")
					}

					//Parse the tokens
					var newcond Condition
					newcond.Fmt = tokens[0]
					newcond.Arg, err = strconv.ParseFloat(tokens[1], 64)
					CheckError(err)
					newcond.TolPos, err = strconv.ParseFloat(tokens[3], 64)
					CheckError(err)
					newcond.TolNeg, err = strconv.ParseFloat(tokens[4], 64)
					CheckError(err)

					//Append the new condition to the rule
					newrule.Cond = append(newrule.Cond, newcond)
				}

				//Append the rule to the selection
				currProcess.Selection = append(currProcess.Selection, newrule)
			}
		}
	}
	if readingSelection {
		processlist = append(processlist, currProcess)
	}
	fmt.Println("========")
	// for _, process := range processlist {
	// 	fmt.Println(process)
	// }

	file.Close()
	return processlist
}
