package game

import "github.com/vova616/chipmunk"

type Collider interface {
	CollideWith(shapeType ShapeType)
}

type BaseCollisionCallback struct {}

func (cb *BaseCollisionCallback) CollisionEnter(arbiter *chipmunk.Arbiter) bool    { return true }
func (cb *BaseCollisionCallback) CollisionPreSolve(arbiter *chipmunk.Arbiter) bool { return true }
func (cb *BaseCollisionCallback) CollisionPostSolve(arbiter *chipmunk.Arbiter)     {}
func (cb *BaseCollisionCallback) CollisionExit(arbiter *chipmunk.Arbiter)          {}

