package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDestroyAlien(t *testing.T) {
	app := NewDummyApp()
	controller := NewController(app)

	// Destroy an alien that exists.
	err := controller.DestroyAlien(0)
	require.Nil(t, err)
	assert.Nil(t, app.Aliens[0])

	// Destroy an alien that doesn't exist.
	err = controller.DestroyAlien(6)
	require.NotNil(t, err)
}

func TestDestroyCity(t *testing.T) {
	app := NewDummyApp()
	controller := NewController(app)

	// Destroy a city that doesn't exist.
	err := controller.DestroyCity("Atlantis")
	require.NotNil(t, err)

	// Destroy a city that exists.
	err = controller.DestroyCity("A")
	require.Nil(t, err)
	_, exists := app.WorldMap.Cities["A"]
	require.False(t, exists)
}

func TestMoveAlienToCity(t *testing.T) {
	app := NewDummyApp()
	controller := NewController(app)

	// Move alien with ID that doesn't exist
	err := controller.MoveAlienToCity(-1, "Bar")
	require.NotNil(t, err)

	// Move alien with ID that doesn't exist
	err = controller.MoveAlienToCity(10, "Bar")
	require.NotNil(t, err)

	// Move alien to city that doesn't exist
	err = controller.MoveAlienToCity(0, "Atlantis")
	require.NotNil(t, err)

	// Move an alien that exists.
	err = controller.MoveAlienToCity(0, "A")
	require.Nil(t, err)
	assert.Equal(t, app.Aliens[0].CurrentCity, app.WorldMap.Cities["A"])

	// Move an alien that no longe exists.
	controller.DestroyAlien(0)
	err = controller.MoveAlienToCity(0, "A")
	require.NotNil(t, err)
}

func TestAreAllAliensDestroyed(t *testing.T) {
	app := NewDummyApp()
	controller := NewController(app)

	areAllAliensDestroyed := controller.AreAllAliensDestroyed()
	assert.False(t, areAllAliensDestroyed)

	// Destroy all aliens.
	controller.DestroyAlien(0)
	controller.DestroyAlien(1)
	controller.DestroyAlien(2)
	controller.DestroyAlien(3)

	areAllAliensDestroyed = controller.AreAllAliensDestroyed()
	assert.True(t, areAllAliensDestroyed)
}

func TestIsAlienMovementLimitReached(t *testing.T) {
	app := NewDummyApp()
	controller := NewController(app)

	isAlienMvmtReached := controller.IsAlienMovementLimitReached()
	assert.False(t, isAlienMvmtReached)

	// Move all aliens 10,000 times.
	for i := 0; i < 10000; i++ {
		controller.MoveAlienToCity(0, "A")
		controller.MoveAlienToCity(1, "B")
		controller.MoveAlienToCity(2, "C")
		controller.MoveAlienToCity(3, "D")
	}

	isAlienMvmtReached = controller.IsAlienMovementLimitReached()
	assert.True(t, isAlienMvmtReached)
}
