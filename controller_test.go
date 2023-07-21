package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// Cities:
	// CityA north=CityB south=CityC
	// CityB east=CityD south=CityA
	// CityC north=CityA west=CityD
	// CityD west=CityB east=CityC
	//
	// Aliens:
	// Alien0 in CityA
	// Alien1 in CityB
	// Alien2 in CityC
	// Alien3 in CityD

	dummyAppCfg = &DummyAppConfig{
		AlienCount: 4,
		MaxMoves:   500,
		Map: map[string][]interface{}{
			"A": []interface{}{
				map[string]string{"north": "B"},
				map[string]string{"south": "C"},
			},
			"B": []interface{}{
				map[string]string{"east": "D"},
				map[string]string{"south": "A"},
			},
			"C": []interface{}{
				map[string]string{"north": "A"},
				map[string]string{"west": "D"},
			},
			"D": []interface{}{
				map[string]string{"west": "B"},
				map[string]string{"east": "C"},
			},
		},
		AlienLocations: map[string]int{
			"A": 0,
			"B": 1,
			"C": 2,
			"D": 3,
		},
	}
)

func TestCtrl_DestroyAlien(t *testing.T) {
	app := NewDummyApp(dummyAppCfg)
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
	app := NewDummyApp(dummyAppCfg)
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
	app := NewDummyApp(dummyAppCfg)
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
	app := NewDummyApp(dummyAppCfg)
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
	app := NewDummyApp(dummyAppCfg)
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
	app := NewDummyApp(dummyAppCfg)
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
