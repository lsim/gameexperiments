package game

import (
	"log"
	"github.com/vova616/chipmunk/vect"
	"github.com/vova616/chipmunk"
	"math"
)

var (
	playerMass float32 = 1
	planetRadius float32 = 70.0
	planetMass float32 = 100
	gravityStrength vect.Float = 5.0e6
)

type Player struct {
	Id    int
	Name  string
	Shape *chipmunk.Shape
}

func (player *Player) AddThrust() {
	if player != nil {
		// Need to create a vector pointing in the direction indicated by the angle
		direction := vect.FromAngle(player.Shape.Body.Angle())
		direction.Mult(50.0)
		player.Shape.Body.AddForce(float32(direction.X), float32(direction.Y))
	}
}

func (player *Player) Rotate(amount float32) {
	if player != nil {
		player.Shape.Body.AddAngularVelocity(amount)
	}
}

type State struct {
	Players      []*Player
	StepCount    int
	Space        *chipmunk.Space
	PlanetShape  *chipmunk.Shape
}

var idCounter int = 0

func (state *State) AddPlayer(name string) *Player {
	playerShape := createPlayerShape(state.Space)
	newPlayer := Player{Name: name, Id: idCounter, Shape: playerShape}
	idCounter += 1
	log.Printf("Adding player %v", newPlayer)
	state.Players = append(state.Players, &newPlayer)
	return &newPlayer
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

func (state *State) RunStep(fps float32) {
	state.Space.Step(vect.Float(1 / fps))
}

func CreateInstance() *State {
	space := chipmunk.NewSpace()
	space.Iterations = 250 // Number of refining iterations when computing collisions
	planet := createPlanet(space)
	return &State{
		PlanetShape: planet,
		Space: space,
	}
}

func createPlayerShape(space *chipmunk.Space) *chipmunk.Shape {
	vertices := chipmunk.Vertices{
		vect.Vect{ -5, 5},
		vect.Vect{ 5, 0},
		vect.Vect{ -5, -5},
		//vect.Vect{ -3, 0},
	}
	playerShape := chipmunk.NewPolygon(vertices, vect.Vect{})
	//playerShape := chipmunk.NewCircle(vect.Vector_Zero, playerRadius)
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
	planetBody := chipmunk.NewBodyStatic() //(vect.Float(planetMass), shape.Moment(planetMass))
	planetBody.AddShape(shape)
	//planetBody.SetAngularVelocity(0.2)
	space.AddBody(planetBody)
	return shape
}

func planetGravityVelocity(body *chipmunk.Body, ignoredGravity vect.Vect, damping, dt vect.Float) {
	p := body.Position()
	sqDist := vect.LengthSqr(p)
	g := vect.Mult(p, vect.Float(-gravityStrength / (sqDist * vect.Float(math.Sqrt(float64(sqDist))))))
	body.UpdateVelocity(g, damping, dt)
}