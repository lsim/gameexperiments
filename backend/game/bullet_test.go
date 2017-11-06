package game

import (
	"testing"
	"github.com/vova616/chipmunk/vect"
	"github.com/vova616/chipmunk"
	"math"
)

func TestUnitTesting(t *testing.T) {
	t.Logf("Unittests running - yay!")
}

func TestBulletCreationAndTravel(t *testing.T) {
	space := chipmunk.NewSpace()
	playerShape := CreatePlayerShape(space)
	// Turn the player to shoot along the X-axis in the negative direction
	playerShape.Body.SetAngle(math.Pi)
	playerShape.Update()
	// Create bullet from player
	bulletShape := CreateBulletShape(space, playerShape)

	lower1 := bulletShape.BB.Lower
	space.Step(1)
	lower2 := bulletShape.BB.Lower

	if !floatsEqual(lower1.X - bulletSpeed, lower2.X) {
		t.Error("Bullet did not travel the correct distance")
	}
	if !floatsEqual(lower1.Y, lower2.Y) {
		t.Error("Bullet did not travel along the X-axis")
	}
}

const tolerance = 0.0001
func floatsEqual(a, b vect.Float) bool {
	diff := math.Abs(float64(a - b))
	return diff < tolerance
}