package main

import "fmt"
import "github.com/subspace-engine/subspace/snd"
import "github.com/subspace-engine/subspace/util"

func main() {
	snd.Init()
	sound := snd.PlaySound("footstep.ogg")
	snd.SetPosition(sound, util.Vec3{1.0, 0.0, 0.0})
	snd.PlaySound("footstep.ogg")
	fmt.Scanf("\n")
}
