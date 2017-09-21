package game

import (
	"github.com/stojg/vector"
	"log"
)

var playerCount = 0

type Player struct {
	Id int
	Name     string
	Pos      vector.Vector3
	Velocity vector.Vector3
}

func (player *Player) ApplyForce(vector vector.Vector3) {
	player.Velocity.Add(&vector)
}

type State struct {
	Players   []*Player
	StepCount int
}

func (state *State) GetPlayers() []*Player {
	return state.Players
}

func (state *State) AddPlayer(name string) *Player {
	newPlayer := Player{Name: name, Id: playerCount}
	log.Printf("Adding player %v", newPlayer)
	state.Players = append(state.Players, &newPlayer)
	playerCount++
	return &newPlayer
}

func (state *State) RemovePlayer(player *Player) {
	if player != nil {
		for i := 0; i < len(state.Players); i++ {
			if state.Players[i].Id == player.Id {
				log.Printf("Removing player %v", *player)
				state.Players = append(state.Players[:i], state.Players[i+1:]...)
				break
			}
		}
	}
}

func (state *State) RunStep() {
	//log.Printf("Running step %v", state.StepCount)
	for _,player := range state.Players {
		(player).Pos.Add(&player.Velocity)
	}
	state.StepCount++
}

func CreateInstance() *State {
	return &State{
		Players: []*Player{},
	}
}
