package game

import (
	"github.com/vova616/chipmunk"
	"github.com/vova616/chipmunk/vect"
	"math"
)

func CreatePlanet(space *chipmunk.Space) *chipmunk.Shape {
	shape := chipmunk.NewCircle(vect.Vector_Zero, planetRadius)
	planetBody := chipmunk.NewBodyStatic()
	planetBody.AddShape(shape)
	space.AddBody(planetBody)
	planetBody.CallbackHandler = &PlanetCollisionCallback{}
	return shape
}

func PlanetGravityVelocity(body *chipmunk.Body, _ vect.Vect, damping, dt vect.Float) {
	p := body.Position()
	sqDist := vect.LengthSqr(p)
	g := vect.Mult(p, vect.Float(-gravityStrength/(sqDist*vect.Float(math.Sqrt(float64(sqDist))))))
	body.UpdateVelocity(g, damping, dt)
}

type PlanetCollisionCallback struct {
	BaseCollisionCallback
}

func (cb *PlanetCollisionCallback) CollisionEnter(arbiter *chipmunk.Arbiter) bool {
	inboundShape := arbiter.ShapeB
	collider, ok := inboundShape.UserData.(Collider)
	if ok {
		collider.CollideWith(PlanetShape)
	}
	return false
}

