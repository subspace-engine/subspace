package ui

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

type InputOutput interface {
	Print(s string)
	Println(s string)
	Read() (s string)
}

type ui struct {
	reader *bufio.Reader
}

func NewInputOutput() (io InputOutput) {
	reader := bufio.NewReader(os.Stdin)
	io = &ui{reader}
	return
}

func (u *ui) Print(message string) {
	fmt.Print(message)
}

func (u *ui) Println(message string) {
	fmt.Println(message)
}

func (u *ui) Read() (message string) {
	message , _ = u.reader.ReadString('\n')
	message = strings.TrimSpace(message)
	return
}