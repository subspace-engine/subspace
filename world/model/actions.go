package model

type Action struct {
	Source Thing
	Tag    string
	Dobj   Thing
	Iobj   Thing
}

type Actor interface {
	Act(Action)
	RegisterPreaction(string, func(Action) int)
	RegisterPostaction(string, func(Action) int)
	RegisterAction(string, func(Action) int)
}

type ActionManager struct {
	actions     map[string][]func(Action) int
	preactions  map[string][]func(Action) int
	postactions map[string][]func(Action) int
}

func MakeActionManager() Actor {
	return &ActionManager{make(map[string][]func(Action) int), make(map[string][]func(Action) int), make(map[string][]func(Action) int)}
}

func registerAction(actions *(map[string][]func(Action) int), tag string, response func(Action) int) {
	val, prs := (*actions)[tag]
	if !prs {
		val = make([]func(Action) int, 0)
	}
	val = append(val, response)
	(*actions)[tag] = val
}

func (self *ActionManager) RegisterAction(tag string, response func(Action) int) {
	registerAction(&self.actions, tag, response)
}

func (self *ActionManager) RegisterPreaction(tag string, response func(Action) int) {
	registerAction(&self.preactions, tag, response)
}

func (self *ActionManager) RegisterPostaction(tag string, response func(Action) int) {
	registerAction(&self.postactions, tag, response)
}

func act(actions map[string][]func(Action) int, action Action) {
	val, prs := actions[action.Tag]
	if !prs {
		return
	}
	for _, response := range val {
		if response(action) != 0 {
			return
		}
	}
}

func (self *ActionManager) Act(action Action) {
	act(self.preactions, action)
	act(self.actions, action)
	act(self.postactions, action)
}
