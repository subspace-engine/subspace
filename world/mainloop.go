package world

import "fmt"

type Client struct {
	Queue chan string
	Agent WObject
}

func NewClient(agent WObject) Client {
	return Client{make(chan string), agent}
}


type Sim struct {
	clients []Client
	worlds []WObject
}

func NewSim() Sim {
	return Sim {make([]Client, 0), make([]WObject, 0)}
}

func (sim Sim)AddClient(client Client) {
	sim.clients = append(sim.clients, client)
}

func (sim Sim)AddWorld(world WObject) {
	sim.worlds = append(sim.worlds, world)
}

func ListCommands() string {
	return "Help:\nrelations: print your relations"
}

func PrintRelations(agent WObject) string {
	ret := ""
	for key, val := range agent.Relations {
		ret+=fmt.Sprintf("%s %s\n", key, val.Name())
	}
	return ret
}


func (sim Sim)Run() {
	fmt.Println("Running")
	for {
		for _, client := range sim.clients {
			fmt.Println(len(sim.clients))
			fmt.Println(client)
		command := <- client.Queue
		switch command {
			case "help":
			client.Queue <- ListCommands()
		case "relations":
			client.Queue <- PrintRelations(client.Agent)
		}
	}
	}
}
