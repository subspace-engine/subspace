package cmd

import "github.com/subspace-engine/subspace/world/model"
import "github.com/subspace-engine/subspace/con"
import "strings"
import "time"
import "fmt"

type Command struct {
	Key     rune
	Command string
	Action  string
}

type CommandParser struct {
	Commands   map[string]Command
	Keys       map[rune]Command
	Con        con.Console
	Player     model.Player
	Parsing    bool
	QuitKey    rune
	HelpKey    rune
	CommandKey rune
}

func MakeCommandParser(console con.Console, player model.Player) *CommandParser {
	return &CommandParser{make(map[string]Command, 0), make(map[rune]Command, 0), console, player, false, 'q', 'h', ':'}
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
	if len(line) == 0 {
		self.Player.Say("I beg your pardon?")
		return
	}
	words := strings.Split(line, " ")
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

func (self CommandParser) DisplayHelp() {
	self.Con.Println("Help:")
	self.Con.Println(fmt.Sprintf("%c: Quit the application", self.QuitKey))
	for key, val := range self.Keys {
		self.Con.Println(fmt.Sprintf("%c: %s", key, val.Command))
	}
}

func (self *CommandParser) RunParser() {
	self.Parsing = true
	self.Con.Println(fmt.Sprintf("Press %c to quit, and %c for help.", self.QuitKey, self.HelpKey))
	proc := self.Con.MakeEventProc()
	proc.SetKeyDown(func(key rune) {
		switch key {
		case self.CommandKey:
			self.Con.Print(":")
			line := self.Con.ReadLine()
			self.ParseCommand(line)
		case self.QuitKey:
			self.Con.Println("Quitting...")
			self.Parsing = false
		case self.HelpKey:
			self.DisplayHelp()
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
