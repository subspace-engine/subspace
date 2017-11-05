// +build darwin dragonfly freebsd linux netbsd openbsd

package term

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
*/
import "C"

import "unsafe"

func initfunc() uint32 {
C.initscr()
C.cbreak()
C.noecho()
C.keypad(C.stdscr, true)
C.clear()
return 0
}

func termfunc() {
C.endwin()
}

func printfunc(text string) {
var ctext * C.char = C.CString(text)
defer C.free(unsafe.Pointer(ctext))
C.addstr(ctext)
}

func readfunc() int {
return int(C.getch())
}

func readlinefunc() string {
var cbuffer * C.char = C.creadline()
defer C.free(unsafe.Pointer(cbuffer))
return C.GoString(cbuffer)
}
