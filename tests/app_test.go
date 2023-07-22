package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestApp_PopulateMapWithAliens(t *testing.T) {
	app := NewEmptyDummyApp()
	prefilledApp := NewDummyApp(dummyAppCfg)

	app.State.Aliens = prefilledApp.State.Aliens
	app.State.WorldMap = prefilledApp.State.WorldMap

	app.PopulateMapWithAliens()

	populatedCities := len(app.State.AlienLocations)
	assert.LessOrEqual(t, populatedCities, len(app.State.WorldMap.Cities))

	for city, aliens := range app.State.AlienLocations {
		assert.GreaterOrEqual(t, len(aliens), 1)
		for _, alien := range aliens {
			assert.Equal(t, city, alien.CurrentCity)
		}
	}
}

func TestApp_Run(t *testing.T) {
	t.Run("All aliens get destroyed but part of the world remains", func(t *testing.T) {
		cfg := &DummyAppConfig{
			AlienCount: 4,
			MaxMoves:   500,
			Map: map[string][]interface{}{
				"A": {
					map[string]string{"north": "B"},
					map[string]string{"south": "C"},
				},
				"B": {
					map[string]string{"east": "D"},
					map[string]string{"south": "A"},
				},
				"C": {
					map[string]string{"north": "A"},
					map[string]string{"west": "D"},
				},
				"D": {
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

		app := NewDummyApp(cfg)
		ctrl := app.StateController()

		app.Run()

		assert.True(t, ctrl.AreAllAliensDestroyed())
		assert.False(t, ctrl.IsWorldDestroyed())
		assert.False(t, ctrl.IsAlienMovementLimitReached())
		assert.False(t, ctrl.AreRemainingAliensTrapped())
	})
	t.Run("All aliens get destroyed and the world is destroyed", func(t *testing.T) {
		cfg := &DummyAppConfig{
			AlienCount: 4,
			MaxMoves:   500,
			Map: map[string][]interface{}{
				"A": {
					map[string]string{"north": "B"},
				},
				"B": {
					map[string]string{"south": "A"},
				},
			},
			AlienLocations: map[string][]int{
				"A": {0, 1},
				"B": {2, 3},
			},
		}
		app := NewDummyApp(cfg)
		ctrl := app.StateController()

		app.Run()

		assert.True(t, ctrl.AreAllAliensDestroyed())
		assert.True(t, ctrl.IsWorldDestroyed())
		assert.False(t, ctrl.IsAlienMovementLimitReached())
		assert.False(t, ctrl.AreRemainingAliensTrapped())
	})
	t.Run("Aliens reach the maximum number of moves", func(t *testing.T) {
		cfg := &DummyAppConfig{
			AlienCount: 2,
			MaxMoves:   500,
			Map: map[string][]interface{}{
				"A": {
					map[string]string{"north": "B"},
				},
				"B": {
					map[string]string{"south": "A"},
				},
				"C": {
					map[string]string{"west": "D"},
				},
				"D": {
					map[string]string{"east": "C"},
				},
			},
			AlienLocations: map[string][]int{
				"A": {0},
				"C": {1},
			},
		}

		app := NewDummyApp(cfg)
		ctrl := app.StateController()

		app.Run()

		assert.False(t, ctrl.AreAllAliensDestroyed())
		assert.False(t, ctrl.IsWorldDestroyed())
		assert.True(t, ctrl.IsAlienMovementLimitReached())
		assert.False(t, ctrl.AreRemainingAliensTrapped())
	})
	t.Run("Aliens get trapped", func(t *testing.T) {
		cfg := &DummyAppConfig{
			AlienCount: 2,
			MaxMoves:   500,
			Map: map[string][]interface{}{
				"A": {},
				"C": {},
			},
			AlienLocations: map[string][]int{
				"A": {0},
				"C": {1},
			},
		}
		app := NewDummyApp(cfg)
		ctrl := app.StateController()

		app.Run()

		assert.False(t, ctrl.AreAllAliensDestroyed())
		assert.False(t, ctrl.IsWorldDestroyed())
		assert.False(t, ctrl.IsAlienMovementLimitReached())
		assert.True(t, ctrl.AreRemainingAliensTrapped())
	})
}

// Test Stop behavior by creating a world with two aliens
// that will never meet but will keep moving until maximum movement limit is reached (50000)
// Stop will be called after 100ms, prior to maximum movement limit being reached.
// This test config allows us to ensure that the world won't destroyed, nor the aliens,
// nor will they be trapped.
func TestApp_Stop(t *testing.T) {
	cfg := &DummyAppConfig{
		AlienCount: 2,
		MaxMoves:   50000,
		Map: map[string][]interface{}{
			"A": {
				map[string]string{"north": "B"},
			},
			"B": {
				map[string]string{"south": "A"},
			},
			"C": {
				map[string]string{"west": "D"},
			},
			"D": {
				map[string]string{"east": "C"},
			},
		},
		AlienLocations: map[string][]int{
			"A": {0},
			"C": {1},
		},
	}

	app := NewDummyApp(cfg)
	ctrl := app.StateController()

	go app.Run()

	time.AfterFunc(100*time.Millisecond, func() {
		app.Stop()
	})

	app.Wait()

	assert.True(t, app.IsStopped())
	assert.False(t, ctrl.AreAllAliensDestroyed())
	assert.False(t, ctrl.IsWorldDestroyed())
	assert.False(t, ctrl.IsAlienMovementLimitReached())
	assert.False(t, ctrl.AreRemainingAliensTrapped())
}
