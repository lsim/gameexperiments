package game

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
)

const (
	bulletMass       = 1
	bulletWidth      = 3
	bulletHeight     = 1
	bulletSpeed      = 70
	bulletTtlSeconds = 6
)

type PlayerBullet struct {
	Id           int
	Shooter      *Player
	Shape        *chipmunk.Shape
	KillCallback func(*PlayerBullet) // Note: This may run on a different thread
}

func CreateBulletShape(space *chipmunk.Space, shooterShape *chipmunk.Shape) *chipmunk.Shape {
	bulletPos := shooterShape.Body.Position()
	// TODO: The problem with this position is that it is up-to-date whereas the player position on the client is not
	// To fix this, the client will have to do proper predictive simulation (of the planet's gravity and player velocity etc)
	// Client and server will both need to fully simulate the game - and the server should overrule the client when they disagree
	playerOrientation := vect.FromAngle(shooterShape.Body.Angle())
	playerOrientation.Normalize()
	bulletOffset := vect.Mult(playerOrientation, vect.Float(playerLength/2))
	bulletPos.Add(bulletOffset)
	bulletShape := chipmunk.NewBox(vect.Vect{X: bulletWidth / 2}, vect.Float(bulletWidth), vect.Float(bulletHeight))
	bulletBody := chipmunk.NewBody(bulletMass, bulletShape.Moment(bulletMass))
	bulletBody.CallbackHandler = &BulletCollisionCallback{}
	bulletBody.SetPosition(bulletPos)
	bulletBody.SetAngle(shooterShape.Body.Angle())
	// Compute bullet velocity (comes from ship velocity and weapon exit velocity)
	bulletVelocity := playerOrientation
	bulletVelocity.Mult(vect.Float(bulletSpeed))
	bulletVelocity.Add(shooterShape.Body.Velocity())
	bulletBody.SetVelocity(float32(bulletVelocity.X), float32(bulletVelocity.Y))
	bulletShape.IsSensor = true
	bulletBody.AddShape(bulletShape)
	space.AddBody(bulletBody)
	return bulletShape
}

// Collider interface
func (bullet *PlayerBullet) CollideWith(shapeType ShapeType) {
	// Whatever the bullet collides with, it gets consumed
	bullet.KillCallback(bullet)
}

// ShapeIdentifier interface
func (bullet *PlayerBullet) GetType() ShapeType { return BulletShape }

type BulletCollisionCallback struct {
	BaseCollisionCallback
}

func (cb *BulletCollisionCallback) CollisionEnter(arbiter *chipmunk.Arbiter) bool {
	inboundShape := arbiter.ShapeA
	// TODO: Here we want ShapeA - in the planet collider we want shapeB - it seems brittle!
	collider, ok := inboundShape.UserData.(Collider)
	if ok {
		collider.CollideWith(BulletShape)
		return false
	}
	return true
}
