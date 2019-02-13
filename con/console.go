package con

type Keymap struct {
	KeyUp    rune
	KeyDown  rune
	KeyLeft  rune
	KeyRight rune
}

type EventProc interface {
	SetKeyDown(keydown func(rune))
	SetKeyUp(keyup func(rune))
	Pump()
}

type Console interface {
	MakeEventProc() EventProc
	ReadKey() rune
	ReadLine() string
	Print(string)
	Println(string)
	Map() Keymap
	Destroy()
}
