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

func TestCtrl_MoveAlienToNextCity(t *testing.T) {
	app := NewDummyApp(dummyAppCfg)
	controller := NewController(app)

	// move ailen when nil
	err := controller.MoveAlienToNextCity(nil)
	require.NotNil(t, err)

	// move alien when alien not registered in app's alien set
	alien := &Alien{ID: 100}
	err = controller.MoveAlienToNextCity(alien)
	require.NotNil(t, err)

	// move alien when alien is not in any city (alienLocation does not have this alien)
	app.Aliens[100] = alien
	err = controller.MoveAlienToNextCity(alien)
	require.NotNil(t, err)

	// move alien when alien's city is nil
	alien.CurrentCity = nil
	err = controller.MoveAlienToNextCity(alien)
	require.NotNil(t, err)

	// move alien when alien's city has no neighbours
	newCity := &City{Name: "Atlantis"}
	app.AlienLocations[newCity] = AlienSet{}
	alien.CurrentCity = newCity
	err = controller.MoveAlienToNextCity(alien)
	require.NotNil(t, err)

	// move alien when alien's city has neighbours but neighbour does not exist in app's world map
	// move alien
	newCity.Neighbours = map[string]*City{"north": {Name: "Asgard"}}
	err = controller.MoveAlienToNextCity(alien)
	require.NotNil(t, err)

	// move alien when everything is fine
	alien = app.Aliens[0]
	err = controller.MoveAlienToNextCity(alien)
	require.Nil(t, err)

	assert.True(t, app.Aliens[0].CurrentCity == app.WorldMap.Cities["B"] || app.Aliens[0].CurrentCity == app.WorldMap.Cities["C"])
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
	for i := 0; i < 500; i++ {
		controller.MoveAlienToNextCity(app.Aliens[0])
		controller.MoveAlienToNextCity(app.Aliens[1])
		controller.MoveAlienToNextCity(app.Aliens[2])
		controller.MoveAlienToNextCity(app.Aliens[3])
	}

	isAlienMvmtReached = controller.IsAlienMovementLimitReached()
	assert.True(t, isAlienMvmtReached)
}

func TestCtrl_IsAlienMovementLimitReached_SomeTrappedAliens(t *testing.T) {
	// This config creates a map
	// with 4 cities and 2 aliens.
	// Alien0 is free to move from CityA to CityB to CityC
	// Alien1 is trapped in CityD
	// With this config, the movement limit should be reached if alien0 has reached it
	// discarding the movement limit for Alien1 because it's trapped
	appCfg := &DummyAppConfig{
		AlienCount: 2,
		MaxMoves:   500,
		Map: map[string][]interface{}{
			"A": []interface{}{
				map[string]string{"north": "B"},
				map[string]string{"south": "C"},
			},
			"B": []interface{}{
				map[string]string{"south": "A"},
			},
			"C": []interface{}{
				map[string]string{"north": "A"},
			},
			"D": []interface{}{},
		},
		AlienLocations: map[string]int{
			"A": 0,
			"D": 1,
		},
	}
	app := NewDummyApp(appCfg)
	controller := NewController(app)

	isAlienMvmtReached := controller.IsAlienMovementLimitReached()
	assert.False(t, isAlienMvmtReached)

	// Move all aliens 10,000 times.
	for i := 0; i < 500; i++ {
		err := controller.MoveAlienToNextCity(app.Aliens[0])
		require.Nil(t, err)

		err = controller.MoveAlienToNextCity(app.Aliens[1]) // won't move, will return error
		require.NotNil(t, err)
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
