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

void cscroll() {
int x,y,maxx,maxy;
getyx(stdscr, y, x);
getmaxyx(stdscr,y,x);
if (y+1>=maxy)
wscrl(stdscr,1);
}
*/
import "C"

import "unsafe"

func Init() uint32 {
C.initscr()
C.cbreak()
C.noecho()
C.keypad(C.stdscr, true)
C.scrollok(C.stdscr, true)
C.idlok(C.stdscr, true)
C.clear()
return 0
}

func Terminate() {
C.endwin()
}

func Print(text string) {
var ctext * C.char = C.CString(text)
defer C.free(unsafe.Pointer(ctext))
C.addstr(ctext)

}

func Read() int {
return int(C.getch())
}

func Readln() string {
var cbuffer * C.char = C.creadline()
defer C.free(unsafe.Pointer(cbuffer))
return C.GoString(cbuffer)
}

func KeyUp() int {
return int(C.KEY_UP)
}

func KeyDown() int {
return int(C.KEY_DOWN)
}

func KeyLeft() int {
return int(C.KEY_LEFT)
}

func KeyRight() int {
return int(C.KEY_RIGHT)
}

