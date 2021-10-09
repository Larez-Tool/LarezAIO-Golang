package utils

import (
	"github.com/fatih/color"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))
var ModeLoop int64 = 0
var ModeReset int64 = 1

func ComplexTextAnim(word string, ms int64, colors []*color.Color, mode int64)  {
	arrPtr := 0
	wordSplit := strings.Split(word, "")
	indexesLen := len(colors) - 1
	rotation := "forward"

	for _, letter := range wordSplit {
		_, _ = colors[arrPtr].Printf(letter)

		switch mode {
			case ModeLoop:
				if rotation == "forward" {
					if arrPtr < indexesLen {
						arrPtr++
					} else if arrPtr == indexesLen {
						rotation = "backward"
						arrPtr--
					}
				} else if rotation == "backward" {
					if arrPtr > 0 {
						arrPtr--
					} else if arrPtr == 0 {
						rotation = "forward"
						arrPtr++
					}
				}
			break

			case ModeReset:
				if arrPtr < indexesLen {
					arrPtr++
				} else if arrPtr == indexesLen {
					rotation = "backward"
				} else if arrPtr == indexesLen {
					arrPtr = 0
				}
			break
		}

		time.Sleep(time.Millisecond * time.Duration(ms))
	}
}

var DefaultColor *color.Color = color.New(color.FgYellow, color.Bold)

func TextAnim(word string, ms int64, color *color.Color)  {
	wordSplit := strings.Split(word, "")

	for _, letter := range wordSplit {
		_, _ = color.Printf(letter)
		time.Sleep(time.Millisecond * time.Duration(ms))
	}
}

func GetStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if e == -1 {
		return
	}
	e += s + e - 1
	return str[s:e]
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func Remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func Chunks(xs []string, chunkSize int) [][]string {
	if len(xs) == 0 {
		return nil
	}
	divided := make([][]string, (len(xs)+chunkSize-1)/chunkSize)
	prev := 0
	i := 0
	till := len(xs) - chunkSize
	for prev < till {
		next := prev + chunkSize
		divided[i] = xs[prev:next]
		prev = next
		i++
	}
	divided[i] = xs[prev:]
	return divided
}