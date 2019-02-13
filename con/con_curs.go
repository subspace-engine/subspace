// +build darwin dragonfly freebsd linux netbsd openbsd

package con

/*
#cgo LDFLAGS: -lcurses
#include <curses.h>
#include <stdlib.h>

char * creadline() {
char * buffer = malloc(4000);
echo();
wgetnstr(stdscr, buffer, 4000);
noecho();
return buffer;
}

void cscroll() {
int x,y,maxx,maxy;
getyx(stdscr, y, x);
getmaxyx(stdscr,y,x);
if (y+1>=maxy)
wscrl(stdscr,1);
}
*/
import "C"

import "fmt"
import "unsafe"

type CursesConsole struct {
}

type CursesEventProc struct {
	parent  *CursesConsole
	keyDown func(rune)
	keyUp   func(rune)
}

func (proc *CursesEventProc) SetKeyDown(keydown func(rune)) {
	proc.keyDown = keydown
}

func (proc *CursesEventProc) SetKeyUp(keyup func(rune)) {
	proc.keyUp = keyup
}

func (proc *CursesEventProc) Pump() {
	proc.parent.nodelay(true)
	key := proc.parent.ReadKey()
	proc.parent.nodelay(false)
	if key != -1 {
		proc.keyDown(key)
	}
}

func MakeTextConsole() *CursesConsole {
	C.initscr()
	C.cbreak()
	C.noecho()
	C.keypad(C.stdscr, true)
	C.scrollok(C.stdscr, true)
	C.idlok(C.stdscr, true)
	C.clear()
	return &CursesConsole{}
}

func (*CursesConsole) Destroy() {
	C.endwin()
}

func (*CursesConsole) Print(text string) {
	/*	var ctext * C.char = C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
		C.addstr(ctext)
			C.refresh()
	*/
	fmt.Print(text)
}

func (self *CursesConsole) Println(text string) {
	self.Print(text + "\n")
}

func (*CursesConsole) ReadKey() rune {
	return rune(C.getch())
}

func (*CursesConsole) ReadLine() string {
	var cbuffer *C.char = C.creadline()
	defer C.free(unsafe.Pointer(cbuffer))
	return C.GoString(cbuffer)
}

func (*CursesConsole) Map() Keymap {
	km := Keymap{}
	km.KeyUp = rune(C.KEY_UP)
	km.KeyDown = rune(C.KEY_DOWN)
	km.KeyLeft = rune(C.KEY_LEFT)
	km.KeyRight = rune(C.KEY_RIGHT)
	return km
}

func (*CursesConsole) nodelay(delay bool) {
	if delay {
		C.nodelay(C.stdscr, true)
	} else {
		C.nodelay(C.stdscr, false)
	}
}

func (self *CursesConsole) MakeEventProc() EventProc {
	proc := CursesEventProc{}
	proc.parent = self
	proc.keyDown = func(key rune) {
	}
	proc.keyUp = func(key rune) {
	}
	return &proc
}
