package term

/*
#include <conio.h>
*/
import "C"


func initfunc() uint32 {
return 0
}

func termfunc() {
}

func printfunc(text string) {
C._puts(C.CString(text))
}

func readfunc() int {
return int(c._getch())
}
