package tests

import (
	"testing"

	simulation "github.com/derrandz/xtinvasion/pkg"
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
		AlienLocations: map[string][]int{
			"A": {0},
			"B": {1},
			"C": {2},
			"D": {3},
		},
	}
)

func TestCtrl_DestroyAlien(t *testing.T) {
	app := NewDummyApp(dummyAppCfg)
	ctrl := app.StateController()

	// Destroy an alien that exists.
	err := ctrl.DestroyAlien(0)
	require.Nil(t, err)
	assert.Equal(t, 3, len(app.Aliens))

	// Destroy an alien that doesn't exist.
	err = ctrl.DestroyAlien(6)
	require.NotNil(t, err)
}

func TestCtrl_DestroyCity(t *testing.T) {
	app := NewDummyApp(dummyAppCfg)
	ctrl := app.StateController()

	// Destroy a city that doesn't exist.
	err := ctrl.DestroyCity("Atlantis")
	require.NotNil(t, err)

	// Destroy a city that exists.
	err = ctrl.DestroyCity("A")
	require.Nil(t, err)
	_, exists := app.WorldMap.Cities["A"]
	require.False(t, exists)
}

func TestCtrl_MoveAlienToNextCity(t *testing.T) {
	app := NewDummyApp(dummyAppCfg)
	ctrl := app.StateController()

	// move ailen when nil
	err := ctrl.MoveAlienToNextCity(nil)
	require.NotNil(t, err)

	// move alien when alien not registered in app's alien set
	alien := &simulation.Alien{ID: 100}
	err = ctrl.MoveAlienToNextCity(alien)
	require.NotNil(t, err)

	// move alien when alien is not in any city (alienLocation does not have this alien)
	app.Aliens[100] = alien
	err = ctrl.MoveAlienToNextCity(alien)
	require.NotNil(t, err)

	// move alien when alien's city is nil
	alien.CurrentCity = nil
	err = ctrl.MoveAlienToNextCity(alien)
	require.NotNil(t, err)

	// move alien when alien's city has no neighbours
	newCity := &simulation.City{Name: "Atlantis"}
	app.AlienLocations[newCity] = simulation.AlienSet{}
	alien.CurrentCity = newCity
	err = ctrl.MoveAlienToNextCity(alien)
	require.NotNil(t, err)

	// move alien when alien's city has a nil neighbour
	alien.CurrentCity.Neighbours = map[string]*simulation.City{"north": nil}
	err = ctrl.MoveAlienToNextCity(alien)
	require.NotNil(t, err)

	// move alien when alien's city has neighbours but neighbour does not exist in app's world map
	// move alien
	newCity.Neighbours = map[string]*simulation.City{"north": {Name: "Asgard"}}
	err = ctrl.MoveAlienToNextCity(alien)
	require.NotNil(t, err)

	// move alien when everything is fine
	alien = app.Aliens[0]
	err = ctrl.MoveAlienToNextCity(alien)
	require.Nil(t, err)

	assert.True(t, app.Aliens[0].CurrentCity == app.WorldMap.Cities["B"] || app.Aliens[0].CurrentCity == app.WorldMap.Cities["C"])
}

func TestCtrl_AreAllAliensDestroyed(t *testing.T) {
	app := NewDummyApp(dummyAppCfg)
	ctrl := app.StateController()

	areAllAliensDestroyed := ctrl.AreAllAliensDestroyed()
	assert.False(t, areAllAliensDestroyed)

	// Destroy all aliens.
	for id := range app.Aliens {
		err := ctrl.DestroyAlien(id)
		require.Nil(t, err)
	}

	areAllAliensDestroyed = ctrl.AreAllAliensDestroyed()
	assert.True(t, areAllAliensDestroyed)
	assert.True(t, len(app.Aliens) == 0)
}

func TestCtrl_IsAlienMovementLimitReached(t *testing.T) {
	t.Run("No trapped aliens", func(t *testing.T) {
		app := NewDummyApp(dummyAppCfg)
		ctrl := app.StateController()

		isAlienMvmtReached := ctrl.IsAlienMovementLimitReached()
		assert.False(t, isAlienMvmtReached)

		// Move all aliens 10,000 times.
		for i := 0; i < 500; i++ {
			ctrl.MoveAlienToNextCity(app.Aliens[0])
			ctrl.MoveAlienToNextCity(app.Aliens[1])
			ctrl.MoveAlienToNextCity(app.Aliens[2])
			ctrl.MoveAlienToNextCity(app.Aliens[3])
		}

		isAlienMvmtReached = ctrl.IsAlienMovementLimitReached()
		assert.True(t, isAlienMvmtReached)
	})

	t.Run("Some aliens are trapped", func(t *testing.T) {
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
			AlienLocations: map[string][]int{
				"A": {0},
				"D": {1},
			},
		}
		app := NewDummyApp(appCfg)
		ctrl := app.StateController()

		isAlienMvmtReached := ctrl.IsAlienMovementLimitReached()
		assert.False(t, isAlienMvmtReached)

		// Move all aliens 10,000 times.
		for i := 0; i < 500; i++ {
			err := ctrl.MoveAlienToNextCity(app.Aliens[0])
			require.Nil(t, err)

			err = ctrl.MoveAlienToNextCity(app.Aliens[1]) // won't move, will return error
			require.NotNil(t, err)
		}

		isAlienMvmtReached = ctrl.IsAlienMovementLimitReached()
		assert.True(t, isAlienMvmtReached)
	})
}

func TestCtrl_IsWorldDestroyed(t *testing.T) {
	app := NewDummyApp(dummyAppCfg)
	ctrl := app.StateController()

	isWorldDestroyed := ctrl.IsWorldDestroyed()
	assert.False(t, isWorldDestroyed)

	// Destroy all cities.
	for _, city := range app.WorldMap.Cities {
		err := ctrl.DestroyCity(city.Name)
		require.Nil(t, err)
	}

	isWorldDestroyed = ctrl.IsWorldDestroyed()

	assert.True(t, isWorldDestroyed)
	assert.True(t, len(app.WorldMap.Cities) == 0)
	assert.True(t, len(app.AlienLocations) == 0)
	assert.True(t, len(app.Aliens) == 0)
}

func TestCtrl_AreRemainingAliensTrapped(t *testing.T) {
	t.Run("Single alien trapped", func(t *testing.T) {
		app := NewDummyApp(dummyAppCfg)
		ctrl := app.StateController()

		areRemainingAliensTrapped := ctrl.AreRemainingAliensTrapped()
		assert.False(t, areRemainingAliensTrapped)

		// Destroy all cities except one.
		for _, city := range app.WorldMap.Cities {
			if city.Name != "A" {
				err := ctrl.DestroyCity(city.Name)
				require.Nil(t, err)
			}
		}

		areRemainingAliensTrapped = ctrl.AreRemainingAliensTrapped()
		assert.True(t, areRemainingAliensTrapped)
		assert.True(t, len(app.WorldMap.Cities) > 0)
	})
	t.Run("Multiple aliens trapped", func(t *testing.T) {
		appCfg := &DummyAppConfig{
			AlienCount: 4,
			MaxMoves:   500,
			Map: map[string][]interface{}{
				"A": []interface{}{},
				"B": []interface{}{},
				"C": []interface{}{},
				"D": []interface{}{},
			},
			AlienLocations: map[string][]int{
				"A": {0},
				"B": {1},
				"C": {2},
				"D": {3},
			},
		}
		app := NewDummyApp(appCfg)
		ctrl := app.StateController()

		areRemainingAliensTrapped := ctrl.AreRemainingAliensTrapped()

		assert.True(t, areRemainingAliensTrapped)
		assert.True(t, len(app.WorldMap.Cities) > 0)
	})
}
