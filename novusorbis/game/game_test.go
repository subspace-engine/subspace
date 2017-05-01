package game_test

import (
	"testing"
	"errors"
)

type InputOutput interface {
	Print(s string)
	Println(s string)
	Read() (s string)
}


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

func TestMockioWriting(t *testing.T) {
	mockIO := NewMockIO([]string{})
	mockIO.Print("Hello!")
	outputs := mockIO.CollectedPrints
	if (outputs[0] != "Hello!") {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "Hello!", outputs[0])
	}
}