package main

import (
	"fmt"
	"github.com/subspace-engine/subspace/world/model"
)

func main() {
	v := make([]model.Thing, 0, 0)
	p := model.MakePlayer("jannie", "jannie se ma se kind")
	v = append(v, p)
	for _, val := range v {
		fmt.Println(val.Name())
	}
}
