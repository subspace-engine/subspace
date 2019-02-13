package main

import "github.com/subspace-engine/subspace/world/model"
import "os"
import "encoding/gob"
import "fmt"

func main() {
	obj := model.MakeBasicThing("kaas", "'n blok kaas")
	obj.SetPassable(false)
	file, err := os.Create("data.bin")
	if err != nil {
		fmt.Println("error creating file")
		return
	}
	encoder := gob.NewEncoder(file)
	encoder.Encode(obj)
	file.Close()
	file, err = os.Open("data.bin")
	if err != nil {
		fmt.Println("unable to read file")
		return
	}
	decoder := gob.NewDecoder(file)
	var thing *model.BasicThing
	thing = model.MakeBasicThing("botter", "")
	decoder.Decode(thing)
	fmt.Printf("Got thing %s, %s\n", thing.Name(), thing.Description())

}
