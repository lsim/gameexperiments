package game

type ShapeType int

const (
	PlayerShape ShapeType = iota
	BulletShape ShapeType = iota
	PlanetShape ShapeType = iota
)

type ShapeIdentifier interface {
	GetType() ShapeType
}
