package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCtrl_DestroyAlien(t *testing.T) {
	app := NewDummyApp()
	controller := NewController(app)

	// Destroy an alien that exists.
	err := controller.DestroyAlien(0)
	require.Nil(t, err)
	assert.Equal(t, 3, len(app.Aliens))

	// Destroy an alien that doesn't exist.
	err = controller.DestroyAlien(6)
	require.NotNil(t, err)
}

func TestCtrl_DestroyCity(t *testing.T) {
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

func TestCtrl_MoveAlienToCity(t *testing.T) {
	app := NewDummyApp()
	controller := NewController(app)

	// Move alien with ID that's not valid
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
}

func TestCtrl_AreAllAliensDestroyed(t *testing.T) {
	app := NewDummyApp()
	controller := NewController(app)

	areAllAliensDestroyed := controller.AreAllAliensDestroyed()
	assert.False(t, areAllAliensDestroyed)

	// Destroy all aliens.
	for id := range controller.app.Aliens {
		err := controller.DestroyAlien(id)
		require.Nil(t, err)
	}

	fmt.Println(controller.app.Aliens)

	areAllAliensDestroyed = controller.AreAllAliensDestroyed()
	assert.True(t, areAllAliensDestroyed)
}

func TestCtrl_IsAlienMovementLimitReached(t *testing.T) {
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

func TestCtrl_IsWorldDestroyed(t *testing.T) {
	app := NewDummyApp()
	controller := NewController(app)

	isWorldDestroyed := controller.IsWorldDestroyed()
	assert.False(t, isWorldDestroyed)

	// Destroy all cities.
	for _, city := range controller.app.WorldMap.Cities {
		err := controller.DestroyCity(city.Name)
		require.Nil(t, err)
	}

	isWorldDestroyed = controller.IsWorldDestroyed()
	assert.True(t, isWorldDestroyed)
}
