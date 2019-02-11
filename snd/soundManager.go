package snd

import "github.com/subspace-engine/subspace/util"
import "sync"
import "math"

var driver ALDriver
var sounds map[string][]int
var mutex = &sync.Mutex{}

func Init() {
	mutex.Lock()
	defer mutex.Unlock()
	driver = MakeALDriver()
	sounds = make(map[string][]int, 0)
}

func Terminate() {
	mutex.Lock()
	defer mutex.Unlock()
	driver.Destroy()
}

func PlaySound(file string) int {
	mutex.Lock()
	defer mutex.Unlock()
	val, prs := sounds[file]
	if !prs {
		sound := driver.LoadSound(file)
		if sound == -1 {
			return -1
		}
		sounds[file] = make([]int, 0, 0)
		sounds[file] = append(sounds[file], sound)
		driver.PlaySource(sound)
		return sound
	}
	for _, sound := range val {
		if !driver.IsPlaying(sound) {
			driver.PlaySource(sound)
			return sound
		}
	}
	sound := driver.LoadSound(file)
	driver.PlaySource(sound)
	sounds[file] = append(sounds[file], sound)
	return sound
}

func SetPosition(sound int, pos util.Vec3) {
	driver.SetPosition(sound, pos.X, pos.Y, pos.Z)
}

func SetListenerPosition(pos util.Vec3) {
	driver.SetListenerPosition(pos.X, pos.Y, pos.Z)
}

func SetListenerDirection(direction float64) {
	pos := util.Vec3{
		math.Sin(direction),
		0,
		-math.Cos(direction)}
	driver.SetListenerOrientation(pos.X, pos.Y, pos.Z, 0, 1, 0)
}

func SetLooping(sound int, looping bool) {
	driver.SetLooping(sound, looping)
}
