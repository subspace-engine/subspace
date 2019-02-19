// +build windows

package con

/*
#include <conio.h>
*/
import "C"

import "fmt"

type WindowsConsole struct {
}

type WindowsEventProc struct {
	parent  *WindowsConsole
	keyDown func(rune)
	keyUp   func(rune)
}

func (proc *WindowsEventProc) SetKeyDown(keydown func(rune)) {
	proc.keyDown = keydown
}

func (proc *WindowsEventProc) SetKeyUp(keyup func(rune)) {
	proc.keyUp = keyup
}

func (proc *WindowsEventProc) Pump() {
	key := proc.parent.ReadKey()
	if key >= 0 {
		proc.keyDown(key)
	}
}

func MakeTextConsole() *WindowsConsole {
	return &WindowsConsole{}
}

func (*WindowsConsole) Destroy() {

}

func (*WindowsConsole) Print(text string) {
	fmt.Print(text)
}

func (self *WindowsConsole) Println(text string) {
	self.Print(text + "\n")
}

func (*WindowsConsole) ReadKey() rune {
	ch := rune(C._getch())
	if ch == 0xe0 {
		ch = ch<<8 + rune(C.getch())
	}
	return ch
}

func (*WindowsConsole) ReadLine() string {
	var s string
	fmt.Scanf("%s\n", &s)
	return s
}

func (*WindowsConsole) Map() Keymap {
	km := Keymap{}
	km.KeyLeft = 0xe000 + 75
	km.KeyRight = 0xe000 + 77
	km.KeyUp = 0xe000 + 72
	km.KeyDown = 0xe000 + 80
	km.KeyEscape = 27
	return km
}

func (self *WindowsConsole) MakeEventProc() EventProc {
	proc := WindowsEventProc{}
	proc.parent = self
	proc.keyDown = func(key rune) {
	}
	proc.keyUp = func(key rune) {
	}
	return &proc
}
