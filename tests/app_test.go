package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// test_map.txt contains
// A north=B south=C
// B east=D south=A
// C north=A west=D
// D west=B east=C
func TestReadMapFromFile(t *testing.T) {
	app := NewEmptyDummyApp()
	err := app.ReadMapFromFile("nofile.txt")
	require.NotNil(t, err)

	err = app.ReadMapFromFile("test_map.txt")
	require.Nil(t, err)

	assert.Equal(t, 4, len(app.WorldMap.Cities), "Expected 4 cities, got", len(app.WorldMap.Cities))

	for _, city := range app.WorldMap.Cities {
		assert.Equal(t, 2, len(city.Neighbours), "Expected 2 neighbours for", city.Name, "got", len(city.Neighbours))
	}

	assert.Equal(t,
		"B",
		app.WorldMap.Cities["A"].Neighbours["north"].Name,
		"Expected B for A north, got", app.WorldMap.Cities["A"].Neighbours["north"].Name)
	assert.Equal(t,
		"C",
		app.WorldMap.Cities["A"].Neighbours["south"].Name,
		"Expected C for A south, got", app.WorldMap.Cities["A"].Neighbours["south"].Name)

	assert.Equal(t,
		"D",
		app.WorldMap.Cities["B"].Neighbours["east"].Name,
		"Expected D for B east, got", app.WorldMap.Cities["B"].Neighbours["east"].Name)
	assert.Equal(t,
		"A",
		app.WorldMap.Cities["B"].Neighbours["south"].Name,
		"Expected A for B south, got", app.WorldMap.Cities["B"].Neighbours["south"].Name)

	assert.Equal(t,
		"A",
		app.WorldMap.Cities["C"].Neighbours["north"].Name,
		"Expected A for C north, got", app.WorldMap.Cities["C"].Neighbours["north"].Name)
	assert.Equal(t,
		"D",
		app.WorldMap.Cities["C"].Neighbours["west"].Name,
		"Expected D for C west, got", app.WorldMap.Cities["C"].Neighbours["west"].Name)

	assert.Equal(t,
		"B",
		app.WorldMap.Cities["D"].Neighbours["west"].Name,
		"Expected B for D west, got", app.WorldMap.Cities["D"].Neighbours["west"].Name)
	assert.Equal(t,
		"C",
		app.WorldMap.Cities["D"].Neighbours["east"].Name,
		"Expected C for D east, got", app.WorldMap.Cities["D"].Neighbours["east"].Name)
}

func TestPopulateMapWithAliens(t *testing.T) {
	app := NewEmptyDummyApp()
	prefilledApp := NewDummyApp(dummyAppCfg)

	app.Aliens = prefilledApp.Aliens
	app.WorldMap = prefilledApp.WorldMap

	app.PopulateMapWithAliens()

	populatedCities := len(app.AlienLocations)
	assert.LessOrEqual(t, populatedCities, len(app.WorldMap.Cities))

	for city, aliens := range app.AlienLocations {
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

		app := NewDummyApp(cfg)
		ctrl := app.Controller()

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
				"A": []interface{}{
					map[string]string{"north": "B"},
				},
				"B": []interface{}{
					map[string]string{"south": "A"},
				},
			},
			AlienLocations: map[string][]int{
				"A": {0, 1},
				"B": {2, 3},
			},
		}
		app := NewDummyApp(cfg)
		ctrl := app.Controller()

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
				"A": []interface{}{
					map[string]string{"north": "B"},
				},
				"B": []interface{}{
					map[string]string{"south": "A"},
				},
				"C": []interface{}{
					map[string]string{"west": "D"},
				},
				"D": []interface{}{
					map[string]string{"east": "C"},
				},
			},
			AlienLocations: map[string][]int{
				"A": {0},
				"C": {1},
			},
		}

		app := NewDummyApp(cfg)
		ctrl := app.Controller()

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
				"A": []interface{}{},
				"C": []interface{}{},
			},
			AlienLocations: map[string][]int{
				"A": {0},
				"C": {1},
			},
		}
		app := NewDummyApp(cfg)
		ctrl := app.Controller()

		app.Run()

		assert.False(t, ctrl.AreAllAliensDestroyed())
		assert.False(t, ctrl.IsWorldDestroyed())
		assert.False(t, ctrl.IsAlienMovementLimitReached())
		assert.True(t, ctrl.AreRemainingAliensTrapped())
	})
}
