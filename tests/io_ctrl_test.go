package tests

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// test_map.txt contains
// A north=B south=C
// B east=D south=A
// C north=A west=D
// D west=B east=C
func TestIOController_ReadMapFromFile(t *testing.T) {
	app := NewEmptyDummyApp()
	app.Cfg.MapInputFile = "nofile.txt"
	err := app.IOController().ReadMapFromFile()
	require.NotNil(t, err)

	app.Cfg.MapInputFile = "testdata/test_map.txt"
	err = app.IOController().ReadMapFromFile()
	require.Nil(t, err)

	assert.Equal(t, 4, len(app.State.WorldMap.Cities), "Expected 4 cities, got", len(app.State.WorldMap.Cities))

	for _, city := range app.State.WorldMap.Cities {
		assert.Equal(t, 2, len(city.Neighbours), "Expected 2 neighbours for", city.Name, "got", len(city.Neighbours))
	}

	assert.Equal(t,
		"B",
		app.State.WorldMap.Cities["A"].Neighbours["north"].Name,
		"Expected B for A north, got", app.State.WorldMap.Cities["A"].Neighbours["north"].Name)
	assert.Equal(t,
		"C",
		app.State.WorldMap.Cities["A"].Neighbours["south"].Name,
		"Expected C for A south, got", app.State.WorldMap.Cities["A"].Neighbours["south"].Name)

	assert.Equal(t,
		"D",
		app.State.WorldMap.Cities["B"].Neighbours["east"].Name,
		"Expected D for B east, got", app.State.WorldMap.Cities["B"].Neighbours["east"].Name)
	assert.Equal(t,
		"A",
		app.State.WorldMap.Cities["B"].Neighbours["south"].Name,
		"Expected A for B south, got", app.State.WorldMap.Cities["B"].Neighbours["south"].Name)

	assert.Equal(t,
		"A",
		app.State.WorldMap.Cities["C"].Neighbours["north"].Name,
		"Expected A for C north, got", app.State.WorldMap.Cities["C"].Neighbours["north"].Name)
	assert.Equal(t,
		"D",
		app.State.WorldMap.Cities["C"].Neighbours["west"].Name,
		"Expected D for C west, got", app.State.WorldMap.Cities["C"].Neighbours["west"].Name)

	assert.Equal(t,
		"B",
		app.State.WorldMap.Cities["D"].Neighbours["west"].Name,
		"Expected B for D west, got", app.State.WorldMap.Cities["D"].Neighbours["west"].Name)
	assert.Equal(t,
		"C",
		app.State.WorldMap.Cities["D"].Neighbours["east"].Name,
		"Expected C for D east, got", app.State.WorldMap.Cities["D"].Neighbours["east"].Name)
}

func TestIOController_WriteMapToFile(t *testing.T) {
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
	app.Cfg.MapOutputFile = "testdata/test_write_map.txt"
	err := app.IOController().WriteMapToFile()
	require.Nil(t, err)

	file, err := os.Open(app.Cfg.MapOutputFile)
	require.Nil(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read and create all cities first
	for scanner.Scan() {
		line := scanner.Text()
		assert.True(t, line == "A north=B" || line == "B south=A")
	}
}
