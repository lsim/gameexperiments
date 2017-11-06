package game

import "testing"

func TestPlayerIsCollider(t *testing.T) {
	var player interface{}

	player = &Player{}

	collider, ok := player.(Collider)
	if ok {
		collider.CollideWith(PlanetShape)
	} else {
		t.Fail()
	}
}

