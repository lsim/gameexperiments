package game

import (
	"log"
	"github.com/vova616/chipmunk/vect"
	"github.com/vova616/chipmunk"
	"math"
)

var (
	playerMass float32 = 1
	playerLength float32 = 20
	planetRadius float32 = 70.0
	//planetMass float32 = 100
	gravityStrength vect.Float = 1.0e6
)

type Player struct {
	Id    int
	Name  string
	Shape *chipmunk.Shape
}

func (player *Player) AddThrust(amount float32) {
	if player != nil {
		// Need to create a vector pointing in the direction indicated by the angle
		direction := vect.FromAngle(player.Shape.Body.Angle())
		direction.Mult(vect.Float(amount))
		player.Shape.Body.AddForce(float32(direction.X), float32(direction.Y))
	}
}

func (player *Player) Rotate(amount float32) {
	if player != nil {
		player.Shape.Body.SetAngle(player.Shape.Body.Angle() + vect.Float(amount))
		player.Shape.Body.SetAngularVelocity(0)
	}
}

type State struct {
	Players      []*Player
	StepCount    int
	Space        *chipmunk.Space
	PlanetShape  *chipmunk.Shape
	PlayerDeaths chan *Player
}

var idCounter int = 0

func (state *State) AddPlayer(name string) *Player {
	playerShape := createPlayerShape(state.Space)
	newPlayer := &Player{Name: name, Id: idCounter, Shape: playerShape}
	playerShape.UserData = newPlayer
	idCounter += 1
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

func (state *State) KillPlayer(player *Player) {
	state.PlayerDeaths <- player
	state.RemovePlayer(player)
}

func (state *State) RunStep(fps float32) {
	state.Space.Step(vect.Float(1 / fps))
}

func CreateInstance() *State {
	space := chipmunk.NewSpace()
	space.Iterations = 25 // Number of refining iterations when computing collisions TODO: lower
	planet := createPlanet(space)
	state := &State{
		PlanetShape: planet,
		Space: space,
		PlayerDeaths: make(chan *Player, 100),
	}
	planet.Body.CallbackHandler = &collisionCallback{ state }
	return state
}

func createPlayerShape(space *chipmunk.Space) *chipmunk.Shape {
	shapeScale := vect.Float(playerLength / 2)
	vertices := chipmunk.Vertices{
		vect.Vect{ -1 * shapeScale, 1 * shapeScale},
		vect.Vect{ 1 * shapeScale, 0},
		vect.Vect{ -1 * shapeScale, -1 * shapeScale},
	}

	playerShape := chipmunk.NewPolygon(vertices, vect.Vect{})

	playerBody := chipmunk.NewBody(vect.Float(playerMass), playerShape.Moment(playerMass))
	startRadius := vect.Float(200)
	startPos := vect.Vect{startRadius, 0}
	playerBody.SetPosition(startPos)
	// This places the player into an orbit around the central planet - see https://github.com/slembcke/Chipmunk2D/blob/master/demo/Planet.c#L36
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
	planetBody := chipmunk.NewBodyStatic()
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

type collisionCallback struct {
	state *State
}

func (cb *collisionCallback) CollisionEnter(arbiter *chipmunk.Arbiter) bool    {
	playerShape := arbiter.ShapeB
	player := playerShape.UserData.(*Player)
	cb.state.KillPlayer(player)
	return true
}
func (cb *collisionCallback) CollisionPreSolve(arbiter *chipmunk.Arbiter) bool { return true }
func (cb *collisionCallback) CollisionPostSolve(arbiter *chipmunk.Arbiter)     {}
func (cb *collisionCallback) CollisionExit(arbiter *chipmunk.Arbiter)          {}
