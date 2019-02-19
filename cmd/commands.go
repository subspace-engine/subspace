package cmd

import "github.com/subspace-engine/subspace/world/model"
import "github.com/subspace-engine/subspace/con"
import "strings"
import "time"
import "fmt"
import "os/exec"

var layouts map[string]string = map[string]string{
	"colemak": "zxcvbkm,./arstdhneio'qwfpgjluy;[]",
	"dvorak":  ";qjkxbmwvzaoeuidhtns-',.pyfgcrl/=",
	"querty":  "zxcvbnm,./asdfghjkl;'qwertyuio[]"}

type Command struct {
	Key     rune
	Command string
	Action  string
}

func (self Command) CommandVerb() string {
	words := strings.Fields(self.Command)
	if len(words) == 0 {
		return ""
	}
	return words[0]
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
	Layout     string
	Keymap     con.Keymap
}

func MakeCommandParser(console con.Console, player model.Player) *CommandParser {
	layout := getKeyboardLayout()
	fmt.Printf("Layout: %s\n", layout)
	km := console.Map()
	return &CommandParser{make(map[string]Command, 0), make(map[rune]Command, 0), console, player, false, km.KeyEscape, km.KeyF1, ':', layout, km}
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
	self.Con.Println(fmt.Sprintf("%s: Quit the application", self.Keymap.DescribeKey(self.QuitKey)))
	for key, val := range self.Keys {
		self.Con.Println(fmt.Sprintf("%s: %s", self.Keymap.DescribeKey(key), val.Command))
	}
}

func (self *CommandParser) RunParser() {
	self.Parsing = true
	self.Con.Println(fmt.Sprintf("Press %s to quit, and %s for help.", self.Keymap.DescribeKey(self.QuitKey), self.Keymap.DescribeKey(self.HelpKey)))
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

func isUpper(key rune) bool {
	if key >= 'A' && key <= 'Z' {
		return true
	}
	return false
}

func translateKey(key rune, layout string) rune {
	usRunes := []rune(layouts["querty"])
	targetRunes := []rune(layouts[layout])
	if isUpper(key) {
		usRunes = []rune(strings.ToUpper(layouts["querty"]))
		targetRunes = []rune(strings.ToUpper(layouts[layout]))
	}
	for i, _ := range usRunes {
		if usRunes[i] == key {
			return targetRunes[i]
		}
	}
	return key
}

func (self CommandParser) findKey(key rune, command Command, pattern string) bool {
	patternRunes := []rune(pattern)
	if isUpper(key) {
		patternRunes = []rune(strings.ToUpper(pattern))
	}
	for _, val := range patternRunes {
		_, prs := self.Keys[val]
		if !prs {
			self.Keys[val] = command
			return true
		}
	}
	return false
}

func (self *CommandParser) MakeKeyAbsolute(key rune) {
	if self.Layout == "us" {
		return
	}
	command, prs := self.Keys[key]
	if !prs {
		return
	}
	newKey := translateKey(key, self.Layout)
	impedingCommand, prs := self.Keys[newKey]
	delete(self.Keys, key)
	self.Keys[newKey] = command
	if !prs {
		return
	}
	if self.findKey(key, impedingCommand, impedingCommand.CommandVerb()) {
		return
	} else {
		self.findKey(key, impedingCommand, layouts[self.Layout])
	}
}

func getKeyboardLayout() string {
	cmd := exec.Command("localectl")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Unable to exec")
		return "us"
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		words := strings.Fields(line)
		if len(words) < 3 {
			continue
		}
		if words[1] == "Keymap:" {
			layout := strings.ToLower(words[2])
			if len(layout) > 0 {
				return layout
			}
		}
	}
	return "us"
}
