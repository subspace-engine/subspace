package ui

import (
	"github.com/nsf/termbox-go"
)

type cursor struct {
	x int
	y int
	maxX int
	maxY int
	lastMouseX int
	lastMouseY int
	isMouseClicked bool
}

var cur cursor;

func InitializeTermbox() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc)
	// termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse) Mouse mode disabled by default
	maxX, maxY := termbox.Size()
	SetCursor(0,0)
	SetLastMouse(0,0)
	setCursorBounds(maxX, maxY)
}


func RunTermboxDemo() {
	InitializeTermbox()
	Println("Welcome to Subspace")
	termbox.Flush()
	TerminalWriteLoop()
}

func Println(message string) {
	Print(message)
	StartAtNextLine()
}

func Print(message string) {
	for i := 0 ; i < len(message) ; i++ {
		WriteChar(rune(message[i]))
	}
}

func SetCursor(x, y int) {
	cur.x = x
	cur.y = y
	termbox.SetCursor(x, y)
}

func SetLastMouse(x, y int) {
	cur.lastMouseX = x
	cur.lastMouseY = y
}

func setCursorBounds(maxX, maxY int) {
	cur.maxX = maxX
	cur.maxY = maxY
}

/*
 * Alters the cursor's position by the given amount
 */
func ModifyCursor(x, y int) {
	SetCursor(cur.x + x, cur.y + y)
}

func StartAtNextLine() {
	SetCursor(0, cur.y + 1)
}

func EndAtPrevLine() {
	SetCursor(cur.maxX-1, cur.y - 1)
}

func IncrementCursorWrapping() {
	if (cur.x < cur.maxX-1) {
		ModifyCursor(1, 0)
	} else {
		StartAtNextLine()
	}
}

func DecrementCursorWrapping() {
	if (cur.x > 0) {
		ModifyCursor(-1, 0)
	} else {
		EndAtPrevLine()
	}
}

func UpCursorBounded() {
	if (cur.y > 0) {
		ModifyCursor(0, -1)
	}
}

func DownCursorBounded() {
	if (cur.y < cur.maxY-1) {
		ModifyCursor(0, 1)
	}
}

func WriteChar(char rune) {
	termbox.SetCell(cur.x, cur.y,
		char,
		termbox.ColorGreen, termbox.ColorDefault)
	IncrementCursorWrapping()
}

func DrawRectangle(x1, y1, x2, y2 int) {
	if (x2 < x1) {
		x1, x2 = x2, x1
	}
	if (y2 < y1) {
		y1, y2 = y2, y1
	}

	for j := y1; j <= y2; j++ {
		for i :=x1 ; i <= x2; i++ {
			termbox.SetCell(i, j, 'X', termbox.ColorDefault, termbox.ColorGreen)
		}
	}


}

func TerminalWriteLoop() {
	loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyArrowLeft:
				DecrementCursorWrapping()
			case termbox.KeyArrowRight:
				IncrementCursorWrapping()
			case termbox.KeyArrowDown:
				DownCursorBounded()
			case termbox.KeyArrowUp:
				UpCursorBounded()
			default:
				if ev.Ch > 20 {
					WriteChar(ev.Ch)
				}
			}
		case termbox.EventMouse:
			switch ev.Key {
			case termbox.MouseLeft:
				if (!cur.isMouseClicked) {
					SetCursor(ev.MouseX, ev.MouseY)
					SetLastMouse(ev.MouseX, ev.MouseY)
					cur.isMouseClicked = true
				}
			case termbox.MouseRelease:
				DrawRectangle(cur.lastMouseX, cur.lastMouseY, ev.MouseX, ev.MouseY)
				cur.isMouseClicked = false
				SetCursor(ev.MouseX, ev.MouseY)
			}
		case termbox.EventResize:
			setCursorBounds(ev.Width, ev.Height)
			if (cur.x >= cur.maxX) {
				SetCursor(cur.maxX - 1, cur.y)
			}
			if (cur.y >= cur.maxY) {
				SetCursor(cur.x, cur.maxY - 1)
			}
			termbox.Flush()
		}
		termbox.Flush()

	}

	defer termbox.Close()
}