package con

import "fmt"

type Keymap struct {
	KeyUp     rune
	KeyDown   rune
	KeyLeft   rune
	KeyRight  rune
	KeyF1     rune
	KeyF2     rune
	KeyF3     rune
	KeyF4     rune
	KeyF5     rune
	KeyF6     rune
	KeyF7     rune
	KeyF8     rune
	KeyF9     rune
	KeyF10    rune
	KeyF11    rune
	KeyF12    rune
	KeyEscape rune
}

func (km Keymap) DescribeKey(key rune) string {
	dict := map[rune]string{
		km.KeyUp:     "up arrow",
		km.KeyDown:   "down arrow",
		km.KeyLeft:   "left arrow",
		km.KeyRight:  "right arrow",
		km.KeyF1:     "f1",
		km.KeyF2:     "f2",
		km.KeyF3:     "f3",
		km.KeyF4:     "f4",
		km.KeyF5:     "f5",
		km.KeyF6:     "f6",
		km.KeyF7:     "f7",
		km.KeyF8:     "f8",
		km.KeyF9:     "f9",
		km.KeyF10:    "f10",
		km.KeyF11:    "f11",
		km.KeyF12:    "f12",
		km.KeyEscape: "escape"}
	val, prs := dict[key]
	if !prs {
		return fmt.Sprintf("%c", key)
	}
	return val
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
