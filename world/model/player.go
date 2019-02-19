package model

type Player interface {
	MobileThing
	RegisterPrintFunc(func(text string))
	Print(text string)
}

type BasicPlayer struct {
	MobileThing
	printFunc func(text string)
}

func MakePlayer(name string, description string) Player {
	p := &BasicPlayer{MakeMobileThing(name, description), nil}
	p.SetType(p.Type() + ":Player")
	return p
}

func (self *BasicPlayer) RegisterPrintFunc(printFunc func(text string)) {
	self.printFunc = printFunc
}

// the print func called by the world
func (self *BasicPlayer) Print(text string) {
	if self.printFunc != nil {
		self.printFunc(text)
	}
}
