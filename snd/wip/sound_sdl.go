package main

import "github.com/veandco/go-sdl2/mix"
import "fmt"

func main() {
	a := mix.Init(0)
	b := mix.OpenAudio(44100, 0x8010, 2, 2048)
	if a != nil || b != nil {
		fmt.Println("Error opening audio")
		fmt.Println(a)
		fmt.Println(b)
		return
	}
	chunk, err := mix.LoadWAV("/home/rkruger/audio/out.wav")
	if err != nil {
		return
	}
	c, err := chunk.Play(-1, 0)
	fmt.Println(c)
	mix.SetPosition(c, 0, 1)
	//	mix.SetPanning(c, 0, 1)
	fmt.Scanf("\n")
	mix.CloseAudio()
	mix.Quit()
}
