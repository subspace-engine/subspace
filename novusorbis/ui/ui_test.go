package ui_test

import (
	// "testing"
	"github.com/subspace-engine/subspace/novusorbis/ui"
)

func ExampleUi_Print() {
	io := ui.NewInputOutput()
	io.Print("Test String 1, ")
	io.Print("Test String 2")

	// Output: Test String 1, Test String 2
}


func ExampleUi_Println() {
	io := ui.NewInputOutput()
	io.Println("Test String 1,")
	io.Println("Test String 2")
	// Output: Test String 1,
	// Test String 2
}