package term

/*
#include <conio.h>
*/
import "C"


func Init() uint32 {
return 0
}

func Terminate() {
}

func Print(text string) {
C._puts(C.CString(text))
}

func Read() int {
return int(c._getch())
}
