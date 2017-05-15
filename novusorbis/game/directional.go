package game

import (
	"strings"
	"errors"
	"github.com/subspace-engine/subspace/novusorbis/world"
)

func  (g *GameManager) SetUpDirectionMaps() (err error) {
	letterToDirection := map[rune]Direction{'h':HERE,'n':NORTH,'e':EAST,'s':SOUTH,
		'w':WEST,'u':UP,'d':DOWN,'x':CANCEL,'p':SHOW_POSSIBILITIES}
	g.LetterToDirection = letterToDirection

	directionToString := map[Direction]string{HERE : "Here", NORTH : "North", EAST : "East",
		SOUTH : "South", WEST : "West", UP : "Up", DOWN : "Down"}
	g.DirectionToString = directionToString

	err = nil
	return
}


type Direction int

const (
	HERE Direction = iota
	NORTH
	EAST
	SOUTH
	WEST
	UP
	DOWN
	SHOW_POSSIBILITIES
	CANCEL
	RETRY
)


func  (g *GameManager) GetDirection(dirString string) (dir Direction) {
	dirChar := dirString[0]
	dir, isThere := g.LetterToDirection[rune(dirChar)]

	if !isThere {
		dir = RETRY
	}
	return
}

func (g *GameManager) ReadDirection(actionString string, args []string) (dir Direction) {
	if (len(args) <= 1) {
		g.Println("In which direction would you like to " + actionString +
			"? ('c' to cancel or 'p' for possible directions)")
		dirString, _ := g.Read()
		dirString = strings.ToLower(dirString)
		if (strings.HasPrefix(dirString, actionString)) {
			dirString = strings.TrimSpace(dirString[4:])
		}
		dir = g.GetDirection(dirString)
		if (dir == RETRY) {
			return
		}
	} else {
		dirString := args[1]
		dir = g.GetDirection(dirString)
		if (dir == RETRY) {
			return
		}
	}
	return
}

func (g *GameManager) getDirectionalPosition(dir Direction) (pos world.Position, err error) {
	switch dir {
	case HERE: pos = world.Position{0, 0, 0}
	case NORTH: pos = world.Position{0, 1, 0}
	case EAST: pos = world.Position{1, 0, 0}
	case SOUTH: pos = world.Position{0, -1, 0}
	case WEST: pos = world.Position{-1, 0, 0}
	case UP: pos = world.Position{0, 0, 1}
	case DOWN: pos = world.Position{0, 0, -1}
	default: err = errors.New("UnrecognizedDirection")
	}
	return
}

// TODO direction calculation logic is not correct
func (g *GameManager) StateItemsAtSpecifiedDirection(actionString string,
cancellationString string, args []string) (pos world.Position, err error){
	w := g.World
	var name1, name2 string
	var things []world.Mover
	var isClear bool

	dir := g.ReadDirection(actionString, args)

	if dir == SHOW_POSSIBILITIES {
		g.Println("The possible directions are H, N, E, S, W, U, D")
		err = errors.New("Showing possibilities")
		return
	} else if dir == CANCEL {
		g.Println("Cancelled " + cancellationString + ".")
		err = errors.New("Cancelled")
		return
	}

	pos, err = g.getDirectionalPosition(dir);
	if (err != nil) {
		g.Println("Unrecognized direction. The possible directions are H, N, E, S, W, U, D,")
		// TODO let user try again
		return
	}

	name1, name2, things, isClear, _ = w.GetNamesOfTerrainsAndObjects(pos)

	dirName := g.DirectionToString[dir]

	var middleString string
	var endString string
	if (isClear) {
		middleString =  " above " + name2
	}

	if (things != nil) && (len(things) > 0) {
		endString = " with: "
		for i, thing := range things {
			endString += thing.Name()
			if (i < len(things) - 2) {
				endString += ","
			} else if (i == len(things) - 2) {
				endString += ", and "
			}
		}
	}
	g.Println(dirName  + " is " + name1 + middleString + endString + ".")
	err = nil
	return
}

func (g *GameManager) Look(args []string) (err error) {
	actionString := "look"
	cancellationString := "looking"
	g.StateItemsAtSpecifiedDirection(actionString, cancellationString, args)
	return
}

func (g *GameManager) Move(args []string) (err error) {
	actionString := "look"
	cancellationString := "looking"
	pos, err := g.StateItemsAtSpecifiedDirection(actionString, cancellationString, args)
	if (err != nil) {
		err = nil
		return
	}
	g.World.ShiftColonist(pos)
	err = nil
	return
}