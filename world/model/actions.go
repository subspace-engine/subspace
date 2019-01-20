package model

import "fmt"

type Action struct {
	Target *Thing
	Tag string
	Dobj *Thing
	Iobj *Thing
}

type Acter interface {
	RegisterAction(string, func(Action)int)
	Act(Action)
}

type ActionManager struct {
	actions map[string][]func(Action)int
}

func MakeActionManager() * ActionManager {
	return &ActionManager{make(map[string][]func(Action)int)}
}

func (self*ActionManager)RegisterAction(tag string, response func (Action)int) {
	val, prs :=self.actions[tag]
	if ! prs {
		val =make([]func(Action)int, 0)
	}
	val = append(val, response)
	self.actions[tag]=val
	fmt.Println(val)
}

func (self*ActionManager)Act(action Action) {
	val, prs :=self.actions[action.Tag]
	if !prs {
		return
	}
	for _,response := range val {
		if response(action)!=0 {
			return
		}
	}
}

