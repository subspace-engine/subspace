package game_test

import (
	"testing"
	"errors"
	"github.com/subspace-engine/subspace/novusorbis/game"
)

type MockIO struct {
	InIter          int
	Announcements   []string
	CollectedPrints []string
}

func NewMockIO(announcements []string) (*MockIO) {
	collectedPrints := make([]string, 0)
	return &MockIO{InIter:0, Announcements:announcements, CollectedPrints:collectedPrints}
}

func (t *MockIO) Print(s string) {
	t.CollectedPrints = append(t.CollectedPrints, s)
}

func (t *MockIO) Println(s string) {
	t.CollectedPrints = append(t.CollectedPrints, s)
}

func (t *MockIO) Read() (s string, err error){
	err = nil
	if t.InIter < len(t.Announcements) {
		s = t.Announcements[t.InIter]
		t.InIter++
		return
	}
	err = errors.New("Tried to read too many outputs")
	return
}

func TestMockioReading(t *testing.T) {
	mockIO := NewMockIO([]string{"hello"})
	read, err := mockIO.Read()
	if (err != nil || read != "hello") {
		t.Fail()
	}
}

func TestMockioReadingMultiple(t *testing.T) {
	mockIO := NewMockIO([]string{"Hello,", "World!"})
	read1, err := mockIO.Read()
	read2, err := mockIO.Read()
	if (err != nil || read1 != "Hello," || read2 != "World!") {
		t.Fatalf("Expected \"%s\" but got \"%s\", \"%s\"", "Hello, World!", read1, read2)
	}
}

func TestMockioWriting(t *testing.T) {
	mockIO := NewMockIO([]string{})
	mockIO.Print("Hello!")
	outputs := mockIO.CollectedPrints
	if (outputs[0] != "Hello!") {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "Hello!", outputs[0])
	}
}

func TestExitCalled_Error(t *testing.T) {
	mockIo := NewMockIO([]string{"y"})
	g := game.GameManager{InputOutput : mockIo}
	err := g.Exit([]string{})
	outputs := mockIo.CollectedPrints
	if (err == nil) {
		t.Fatalf("Expected an error but got nil. Outputs: %v", outputs)
	}
}

func TestGameManager_Exit(t *testing.T) {
	mockIo := NewMockIO([]string{"q", "y"})
	g := game.GameManager{InputOutput : mockIo}
	g.InitializeCommandsMap()
	g.MainLoop()
	// If this doesn't work we expect a timeout
}