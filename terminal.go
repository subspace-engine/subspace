package main

import (
	"github.com/nsf/termbox-go"
)

func runTermbox() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	numX := 0
	numY := 0
	loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyArrowUp:
				numY -= 1
			case termbox.KeyArrowDown:
				numY += 1
			case termbox.KeyArrowLeft:
				numX -= 1
			case termbox.KeyArrowRight:
				numX += 1
			}
		}
		termbox.SetCell(numX, numY, rune(48+((numX + numY)%10)), termbox.ColorGreen, termbox.ColorDefault)
		termbox.Flush()
	}

	defer termbox.Close()
}
