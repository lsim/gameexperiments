package game

import (
	"log"
	"github.com/vova616/chipmunk/vect"
	"github.com/vova616/chipmunk"
	"time"
)

var (
	playerMass      float32    = 1
	playerLength    float32    = 20
	planetRadius    float32    = 70.0
	gravityStrength vect.Float = 1.0e6
	//gravityStrength vect.Float = 0
)

type State struct {
	Bullets      []*PlayerBullet
	Players      []*Player
	StepCount    int
	Space        *chipmunk.Space
	PlanetShape  *chipmunk.Shape
	PlayerDeaths chan *Player
	BulletDeaths chan *PlayerBullet
}

var idCounter int = 0

func (state *State) AddPlayer(name string) *Player {
	newPlayer := SpawnPlayer(state.Space, name, idCounter, state.getPlayerKiller(), state.getPlayerShootHandler())
	idCounter += 1
	newPlayer.SetInitialState()
	log.Printf("Adding player %v", newPlayer)
	state.Players = append(state.Players, newPlayer)
	return newPlayer
}

func (state *State) RemovePlayer(player *Player) {
	if player != nil {
		state.Space.RemoveBody(player.Shape.Body)
		for i := 0; i < len(state.Players); i++ {
			if state.Players[i].Id == player.Id {
				log.Printf("Removing player %v", *player)
				state.Players = append(state.Players[:i], state.Players[i+1:]...)
				break
			}
		}
	}
}

func (state *State) RemoveBullet(bullet *PlayerBullet) {
	if bullet != nil {
		state.Space.RemoveBody(bullet.Shape.Body)
		//state.Space.RemoveShape(bullet.Shape)
		for i:= 0; i < len(state.Bullets); i++ {
			if state.Bullets[i].Id == bullet.Id {
				state.Bullets = append(state.Bullets[:i], state.Bullets[i+1:]...)
				break
			}
		}
	}
}

func (state *State) getPlayerKiller() func(*Player) {
	return func(player *Player) {
		state.PlayerDeaths <- player
	}
}

func (state *State) getPlayerShootHandler() func(*Player) *PlayerBullet{
	return func(player *Player) *PlayerBullet {
		bulletShape := CreateBulletShape(state.Space, player.Shape)
		newBullet := &PlayerBullet{
			Id:           idCounter,
			Shooter:      player,
			Shape:        bulletShape,
			KillCallback: state.getBulletDeathHandler(),
		}
		idCounter++
		bulletShape.UserData = newBullet
		state.Bullets = append(state.Bullets, newBullet)
		// Automatically kill bullets after a while
		ttlTimer := time.NewTimer(bulletTtlSeconds * time.Second)
		go func() {
			<- ttlTimer.C
			newBullet.KillCallback(newBullet)
		}()
		return newBullet
	}
}

func (state *State) getBulletDeathHandler() func(*PlayerBullet) {
	// Note: This function may run on a different thread
	return func(bullet *PlayerBullet) {
		state.BulletDeaths <- bullet
	}
}

func (state *State) RunStep(fps float32) {
	state.Space.Step(vect.Float(1 / fps))
}

func CreateInstance() *State {
	space := chipmunk.NewSpace()
	space.Iterations = 25
	planetShape := CreatePlanet(space)
	state := &State{
		PlanetShape:  planetShape,
		Space:        space,
		PlayerDeaths: make(chan *Player, 100),
		BulletDeaths: make(chan *PlayerBullet, 100),
	}
	return state
}

