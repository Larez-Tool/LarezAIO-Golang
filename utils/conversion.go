package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

func ReadLineToArray(path string) []string  {
	inFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error() + `: ` + path)
		return []string{}
	}

	defer func(inFile *os.File) {
		err := inFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(inFile)

	var lines []string

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func ArrayToWriteFile(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	w := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			return err
		}
	}

	return w.Flush()
}

func CleanText(text string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}, text)
}

func BatchActions(a []string, c int) [][]string {
	r := (len(a) + c - 1) / c
	b := make([][]string, r)
	lo, hi := 0, c
	for i := range b {
		if hi > len(a) {
			hi = len(a)
		}
		b[i] = a[lo:hi:hi]
		lo, hi = hi, hi+c
	}
	return b
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

