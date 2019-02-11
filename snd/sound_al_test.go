package snd

import "testing"
import "time"

func TestLoading(t *testing.T) {
	driver := MakeALDriver()
	defer driver.Destroy()
	sound1 := driver.LoadSound("/home/rkruger/audio/musicnight.mp3")
	if sound1 < 1 {
		t.Fail()
	}
	sound2 := driver.LoadSound("/home/rkruger/audio/dingdong.wav")
	if sound2 < 1 {
		t.Fail()
	}
	sound3 := driver.LoadSound("/home/rkruger/audio/longesttime.mp3")
	if sound3 < 1 {
		t.Fail()
	}

}

func TestThreadsafe(t *testing.T) {
	driver := MakeALDriver()
	defer driver.Destroy()
	makeSound(driver)
	makeSound(driver)
	makeSound(driver)
	time.Sleep(time.Second)
}

func makeSound(driver ALDriver) {
	sound := driver.LoadSound("/home/rkruger/audio/footstep.ogg")
	driver.PlaySource(sound)
}
