package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func usage() {
	fmt.Print(`Usage: cat [FILE]
Prints out the file's content. Reads from STDIN if no FILE is specified.
`)
	flag.PrintDefaults()
}

type LineRange struct {
	From int
	To   int
}

func main() {
	var noFileName bool
	var noLineNumbers bool
	var lineNumberPadding int
	var lineRange LineRange

	flag.Usage = usage
	flag.BoolVar(&noFileName, "no-file-name", false, "hide file name")
	flag.BoolVar(&noLineNumbers, "no-line-numbers", false, "hide line numbers")
	flag.IntVar(&lineNumberPadding, "line-number-padding", 4, "line number padding")
	flag.Func("line-range", "print from/to line, e.g. '5-13'", func(str string) error {
		split := strings.Split(str, "-")
		if len(split) != 2 {
			return errors.New("parse error")
		}
		from, err := strconv.Atoi(split[0])
		to, err := strconv.Atoi(split[1])
		if err != nil {
			return errors.New("parse error")
		}
		if from < 1 {
			return errors.New("from < 1")
		}
		if to < from {
			return errors.New("to < from")
		}
		lineRange = LineRange{
			From: from,
			To:   to,
		}
		return nil
	})
	flag.Parse()

	var fileName string
	var file *os.File

	if flag.NArg() > 1 {

	} else if flag.NArg() == 1 {
		var err error
		fileName = flag.Arg(0)
		file, err = os.Open(fileName)
		if err != nil {
			panic(err)
		}
	} else {
		file = os.Stdin
	}

	if !noFileName && file != os.Stdin {
		fmt.Println(fileName)
		for range len(fileName) {
			fmt.Print("-")
		}
		fmt.Print("\n")
	}
	reader := bufio.NewReader(file)
	currLine := 1
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				currLine = 1
				break
			}
			panic(err)
		}
		if lineRange.From != 0 && currLine < lineRange.From {
			currLine++
			continue
		}
		if lineRange.To != 0 && currLine > lineRange.To {
			break
		}
		if !noLineNumbers {
			format := fmt.Sprintf("%%%dd | %%s", lineNumberPadding)
			fmt.Printf(format, currLine, line)
		} else {
			fmt.Print(line)
		}
		currLine++
	}
}
