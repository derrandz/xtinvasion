package main

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
	err := app.readMapFromFile("nofile.txt")
	require.NotNil(t, err)

	err = app.readMapFromFile("test_map.txt")
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
