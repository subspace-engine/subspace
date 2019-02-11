package model

type Action struct {
	Source Thing
	Tag    string
	Dobj   Thing
	Iobj   Thing
}

type Actor interface {
	Act(Action) bool
	RegisterPreaction(string, func(Action) bool)
	RegisterPostaction(string, func(Action) bool)
	RegisterAction(string, func(Action) bool)
}

type ActionManager struct {
	actions     map[string][]func(Action) bool
	preactions  map[string][]func(Action) bool
	postactions map[string][]func(Action) bool
}

func MakeActionManager() Actor {
	return &ActionManager{make(map[string][]func(Action) bool), make(map[string][]func(Action) bool), make(map[string][]func(Action) bool)}
}

func registerAction(actions *(map[string][]func(Action) bool), tag string, response func(Action) bool) {
	val, prs := (*actions)[tag]
	if !prs {
		val = make([]func(Action) bool, 0)
	}
	val = append(val, response)
	(*actions)[tag] = val
}

func (self *ActionManager) RegisterAction(tag string, response func(Action) bool) {
	registerAction(&self.actions, tag, response)
}

func (self *ActionManager) RegisterPreaction(tag string, response func(Action) bool) {
	registerAction(&self.preactions, tag, response)
}

func (self *ActionManager) RegisterPostaction(tag string, response func(Action) bool) {
	registerAction(&self.postactions, tag, response)
}

func act(actions map[string][]func(Action) bool, action Action) bool {
	val, prs := actions[action.Tag]
	if !prs {
		return true // if we don't have actions, we are by default successful
	}
	for _, response := range val {
		if !response(action) {
			return false
		}
	}
	return true
}

func addActionManager(actors map[*ActionManager]struct{}, actor Actor) map[*ActionManager]struct{} {
	if actor == nil {
		return actors
	}
	actionManager, ok := actor.(*ActionManager)
	if !ok {
		return actors
	}
	_, ok = actors[actionManager]
	if !ok {
		actors[actionManager] = struct{}{}
	}
	return actors
}

func addActionManagerFromThing(actors map[*ActionManager]struct{}, thing Thing) map[*ActionManager]struct{} {
	if thing != nil {
		return addActionManager(actors, thing.Actions())
	}
	return actors
}

func (self *ActionManager) Act(action Action) bool {
	actors := make(map[*ActionManager]struct{}, 0)
	actors = addActionManagerFromThing(actors, action.Source)
	actors = addActionManagerFromThing(actors, action.Dobj)
	actors = addActionManagerFromThing(actors, action.Iobj)
	actors = addActionManager(actors, self)
	for i, _ := range actors {

		success := act(i.preactions, action)
		if !success {
			return false
		}
	}
	for i, _ := range actors {
		success := act(i.actions, action)
		if !success {
			return false
		}
	}
	for i, _ := range actors {
		success := act(i.postactions, action)
		if !success {
			return false
		}
	}
	return true
}
