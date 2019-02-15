package cmd

import "github.com/subspace-engine/subspace/world/model"
import "github.com/subspace-engine/subspace/con"
import "strings"
import "time"

type Command struct {
	Key     rune
	Command string
	Action  string
}

type CommandParser struct {
	Commands map[string]Command
	Keys     map[rune]Command
	Con      con.Console
	Player   model.Player
	Parsing  bool
}

func MakeCommandParser(console con.Console, player model.Player) *CommandParser {
	return &CommandParser{make(map[string]Command, 0), make(map[rune]Command, 0), console, player, false}
}

func (self CommandParser) ParseKey(key rune) {
	command, prs := self.Keys[key]
	if !prs {
		self.Player.Say("Unhandeled key.")
		return
	}
	action := command.Action
	self.Player.Act(model.Action{self.Player, action, nil, nil})
}

func (self CommandParser) ParseCommand(line string) {
	words := strings.Split(line, " ")
	if len(words) == 0 {
		self.Player.Say("I beg your pardon?")
		return
	}
	first := words[0]
	command, prs := self.Commands[first]
	if !prs {
		self.Player.Say("I don't know that command.")
		return
	}
	self.Player.Act(model.Action{self.Player, command.Action, nil, nil})
}

func (self *CommandParser) Start() {
	go self.RunParser()
}

func (self *CommandParser) RunParser() {
	self.Parsing = true
	proc := self.Con.MakeEventProc()
	proc.SetKeyDown(func(key rune) {
		switch key {
		case ':':
			line := self.Con.ReadLine()
			self.ParseCommand(line)
		default:
			self.ParseKey(key)
		}
	})
	for self.Parsing {
		proc.Pump()
		time.Sleep(time.Millisecond * 5)
	}
}

func (self *CommandParser) Stop() {
	self.Parsing = false
}

func (self *CommandParser) AddCommand(key rune, line string, action string) {
	command := Command{key, line, action}
	words := strings.Split(line, " ")
	if len(words) == 0 {
		panic("Trying to add empty command")
	}
	first := words[0]
	self.Commands[first] = command
	self.Keys[key] = command
}
