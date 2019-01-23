package con

type Keymap struct {
	KeyUp    int
	KeyDown  int
	KeyLeft  int
	KeyRight int
}

type EventProc interface {
	SetKeyDown(keydown func(int))
	SetKeyUp(keyup func(int))
	Pump()
}

type Console interface {
	MakeEventProc() EventProc
	ReadKey() int
	ReadLine() string
	Print(string)
	Println(string)
	Map() Keymap
	Destroy()
}
