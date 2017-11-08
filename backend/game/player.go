package game

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"

	"math"
)

type Player struct {
	Id    int
	Name  string
	Shape *chipmunk.Shape
	KillCallback func(player *Player)
	ShootCallback func(player *Player) *PlayerBullet
}

func SpawnPlayer(space *chipmunk.Space, name string, id int, killCallback func(player *Player), shootCallback func(player *Player) *PlayerBullet) *Player {
	playerShape := CreatePlayerShape(space)
	newPlayer := &Player{
		Name:         name,
		Id:           id,
		Shape:        playerShape,
		KillCallback: killCallback,
		ShootCallback: shootCallback,
	}

	playerShape.UserData = newPlayer

	return newPlayer
}

func CreatePlayerShape(space *chipmunk.Space) *chipmunk.Shape {
	shapeScale := vect.Float(PlayerLength / 2)
	vertices := chipmunk.Vertices{
		vect.Vect{-1 * shapeScale, 1 * shapeScale},
		vect.Vect{1 * shapeScale, 0},
		vect.Vect{-1 * shapeScale, -1 * shapeScale},
	}

	playerShape := chipmunk.NewPolygon(vertices, vect.Vect{})

	playerBody := chipmunk.NewBody(vect.Float(PlayerMass), playerShape.Moment(PlayerMass))
	playerBody.UpdateVelocityFunc = PlanetGravityVelocity
	playerBody.AddShape(playerShape)
	space.AddBody(playerBody)
	return playerShape
}

func (player *Player) AddThrust() {
	if player != nil {
		// Need to create a vector pointing in the direction indicated by the angle
		direction := vect.FromAngle(player.Shape.Body.Angle())
		direction.Mult(vect.Float(ThrustFactor))
		player.Shape.Body.AddVelocity(float32(direction.X), float32(direction.Y))
	}
}

func (player *Player) SetInitialState() {
	startPos := vect.Vect{PlayerStartRadius, 0}
	player.Shape.Body.SetPosition(startPos)

	// This places the player into an orbit around the central planet - see https://github.com/slembcke/Chipmunk2D/blob/master/demo/Planet.c#L36
	v := vect.Float(math.Sqrt(float64(GravityStrength/PlayerStartRadius)) / float64(PlayerStartRadius))
	initialVelocity := vect.Perp(startPos)
	initialVelocity.Mult(v)
	player.Shape.Body.SetVelocity(float32(initialVelocity.X), float32(initialVelocity.Y))
}

func (player *Player) Rotate(sign float64) {
	if player != nil {
		player.Shape.Body.SetAngle(player.Shape.Body.Angle() + vect.Float(sign *RotateFactor))
		player.Shape.Body.SetAngularVelocity(0)
	}
}

func (player *Player) Shoot() (*PlayerBullet, bool) {
	// Create a simple box shape - ok
	// Needs to be positioned at the tip of the player's ship - ok
	// Needs to be oriented with the same angle as the ship - ok
	// Needs to have velocity - ok

	// Needs to be broadcast to all players so they can animate the bullet - ok
	// When the bullet hits something (or disappears out of bounds), it will be broadcast to all players - ok
	// When a player joins, he will need to be told which bullets are where (and their heading and velocity)
	// This way the clients animate the bullets and all collision computations are done server side

	// Each bullet will need an id so that its termination can be communicated correctly to clients - ok
	// Each bullet may also need to remain associated with whoever fired it (for scoring purposes) - ok
	// Firing of bullets may need to be rate limited (serverside?) - ok for client side
	if player != nil {
		return player.ShootCallback(player), true
	}
	return nil, false
}

// Collider interface
func (player *Player) CollideWith(shapeType ShapeType) {
	if shapeType == BulletShape || shapeType == PlanetShape {
		player.KillCallback(player)
	}
}
