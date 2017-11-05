package term

func Init() uint32 {
return initfunc()
}

func Terminate() {
termfunc()
}

func Print(text string) {
printfunc(text)
}

func Println(text string) {
Print(text)
Print("\n")
}

func Read() int {
return readfunc()
}

func Readln() string {
return readlinefunc()
}

func Keys() KeyStruct {
return keysfunc()
}
