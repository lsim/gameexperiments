package game

import (
	"log"
	"github.com/vova616/chipmunk/vect"
	"github.com/vova616/chipmunk"
	"math"
)

var (
	playerRadius float32 = 25
	playerMass float32 = 1
	planetRadius float32 = 70.0
	planetMass float32 = 100
	gravityStrength vect.Float = 5.0e6
)

type Player struct {
	Id int
	Name     string
	Pos      vect.Vect
	Velocity vect.Vect
}

//func (player *Player) ApplyForce(vector vect.Vect) {
//	player.Velocity.Add(&vect)
//}

type State struct {
	Players      []*Player
	StepCount    int
	Space        *chipmunk.Space
	PlanetShape  *chipmunk.Shape
	PlayerShapes []*chipmunk.Shape
}

//func (state *State) GetPlayers() []*Player {
//	return state.Players
//}

func (state *State) AddPlayer(name string) *Player {
	newPlayer := Player{Name: name, Id: len(state.Players)}
	log.Printf("Adding player %v", newPlayer)
	state.Players = append(state.Players, &newPlayer)
	playerShape := createPlayer(state.Space)
	playerShape.UserData = &newPlayer
	state.PlayerShapes = append(state.PlayerShapes, playerShape)
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
		for i := 0; i < len(state.PlayerShapes); i++ {
			playerShape := state.PlayerShapes[i]
			if playerShape.UserData.(*Player).Id == player.Id {
				state.Space.RemoveBody(playerShape.Body)
				state.PlayerShapes = append(state.PlayerShapes[:i], state.PlayerShapes[i+1:]...)
			}
		}
	}
}

func (state *State) RunStep(fps float32) {
	//log.Printf("Running step %v", state.StepCount)
	//for _,player := range state.Players {
	//	(player).Pos.Add(player.Velocity)
	//}
	//state.StepCount++
	state.Space.Step(vect.Float(1 / fps))
}

func CreateInstance() *State {
	space := chipmunk.NewSpace()
	space.Iterations = 20 // Number of refining iterations when computing collisions
	planet := createPlanet(space)
	return &State{
		Players: []*Player{},
		PlanetShape: planet,
		Space: space,
	}
}

func createPlayer(space *chipmunk.Space) *chipmunk.Shape {
	playerShape := chipmunk.NewCircle(vect.Vector_Zero, playerRadius)
	playerBody := chipmunk.NewBody(vect.Float(playerMass), playerShape.Moment(playerMass))
	startRadius := vect.Float(200)
	startPos := vect.Vect{startRadius, 0}
	playerBody.SetPosition(startPos)
	// This is an attempt at placing the player into an orbit around the central planet - see https://github.com/slembcke/Chipmunk2D/blob/master/demo/Planet.c#L36
	v := vect.Float(math.Sqrt(float64(gravityStrength / startRadius)) / float64(startRadius))
	initialVelocity := vect.Perp(startPos)
	initialVelocity.Mult(v)
	playerBody.SetVelocity(float32(initialVelocity.X), float32(initialVelocity.Y))
	playerBody.UpdateVelocityFunc = planetGravityVelocity
	playerBody.AddShape(playerShape)
	space.AddBody(playerBody)
	return playerShape
}

func createPlanet(space *chipmunk.Space) *chipmunk.Shape {
	shape := chipmunk.NewCircle(vect.Vector_Zero, planetRadius)
	planetBody := chipmunk.NewBody(vect.Float(planetMass), shape.Moment(planetMass))
	planetBody.AddShape(shape)
	planetBody.SetAngularVelocity(0.2)
	space.AddBody(planetBody)
	return shape
}

func planetGravityVelocity(body *chipmunk.Body, ignoredGravity vect.Vect, damping, dt vect.Float) {
	p := body.Position()
	sqDist := vect.LengthSqr(p)
	g := vect.Mult(p, vect.Float(-gravityStrength / (sqDist * vect.Float(math.Sqrt(float64(sqDist))))))
	body.UpdateVelocity(g, damping, dt)
}