package world

import (
	"fmt"
	engine "github.com/subspace-engine/subspace/engine"
)

type testWorld struct {
	lastId  int
	actor   chan engine.PlayerAction
	players []*BoomPlayer
}

func NewTestWorld() *testWorld {
	return &testWorld{0, make(chan engine.PlayerAction, 1), make([]*BoomPlayer, 0)}
}

func (w *testWorld) Actor() chan<- engine.PlayerAction {
	return chan<- engine.PlayerAction(w.actor)
}
func (w *testWorld) RequestPlayer() Player {
	p := newPlayer(w.lastId)
	w.lastId++
	w.players = append(w.players, p)
	return p
}

func (w *testWorld) Run(done chan struct{}) {
	for act := range w.actor {
		fmt.Printf("world receives: player %d:, action %s, number of players %d\n", act.Id, act.Action.Desc, len(w.players))
		var msg string
		switch len(w.players) {
		case 1:
			msg = "Nothing happens. Maybe you should phone a friend?"
		case 2:
			msg = "You suddenly remember that you have a student card in your pocket. You quickly take it out and hold it to the sensor of the boom."
			w.notifyAll(fmt.Sprintf("Player %d Triumphantly opens the boom!", act.Id))
			close(done)

		case 3:
			msg = "Instead of worrying about the boom, you have a fasinating conversation about game engines."
		default:
			msg = "It seems pretty crouded. Dealing with chaos never was one of your strengths."
			w.notifyAll("It takes 45 minutes for everyone to agree that a drink is what everyone needs. Right now! You all leave.")
			close(done)
		}
		update := &engine.ServerUpdate{msg}
		fmt.Println("Server sending: %s", update.Desc)
		w.players[act.Id].in <- update
	}
	fmt.Println("update loop exit")
}

func (w *testWorld) notifyAll(s string) {
	u := &engine.ServerUpdate{s}
	for i := range w.players {
		w.players[i].in <- u
	}
}

func (w *testWorld) Updater(id int) <-chan *engine.ServerUpdate {
	return w.players[id].in
}

type BoomPlayer struct {
	in chan *engine.ServerUpdate
	id int
}

func newPlayer(id int) *BoomPlayer {
	return &BoomPlayer{make(chan *engine.ServerUpdate, 1), id}
}

func (p *BoomPlayer) Id() int { return p.id }
