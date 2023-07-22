package tests

import (
	"testing"

	simulation "github.com/derrandz/xtinvasion/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRandomNeighbor(t *testing.T) {
	t.Run("nil city", func(t *testing.T) {
		neighbour, err := simulation.GetRandomNeighbor(nil)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "city is nil")

		assert.Nil(t, neighbour)
	})

	t.Run("city with nil neighbours", func(t *testing.T) {
		city := &simulation.City{Name: "A", Neighbours: make(map[string]*simulation.City)}
		city.Neighbours = map[string]*simulation.City{
			"B": nil,
		}

		neighbour, err := simulation.GetRandomNeighbor(city)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), " has a nil neighbour")

		assert.Nil(t, neighbour)
	})

	t.Run("city with no neighbours", func(t *testing.T) {
		city := &simulation.City{Name: "A", Neighbours: make(map[string]*simulation.City)}
		city.Neighbours = map[string]*simulation.City{}

		neighbour, err := simulation.GetRandomNeighbor(city)
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "GetRandomNeighbor: city has no neighbours")

		assert.Nil(t, neighbour)
	})

	t.Run("city with neighbour", func(t *testing.T) {
		city := &simulation.City{Name: "A", Neighbours: make(map[string]*simulation.City)}
		city.Neighbours = map[string]*simulation.City{
			"B": {Name: "B"},
			"C": {Name: "C"},
			"D": {Name: "D"},
		}

		neighbour, err := simulation.GetRandomNeighbor(city)
		require.Nil(t, err)

		assert.NotNil(t, neighbour)
		assert.Contains(t, []string{"B", "C", "D"}, neighbour.Name)
	})

}

func TestOppositeDirection(t *testing.T) {
	t.Run("north", func(t *testing.T) {
		assert.Equal(t, "south", simulation.OppositeDirection("north"))
	})

	t.Run("south", func(t *testing.T) {
		assert.Equal(t, "north", simulation.OppositeDirection("south"))
	})

	t.Run("east", func(t *testing.T) {
		assert.Equal(t, "west", simulation.OppositeDirection("east"))
	})

	t.Run("west", func(t *testing.T) {
		assert.Equal(t, "east", simulation.OppositeDirection("west"))
	})

	t.Run("invalid", func(t *testing.T) {
		assert.Equal(t, "", simulation.OppositeDirection("invalid"))
	})
}

func TestRemoveSliceElement(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		slice := []int{}
		slice = simulation.RemoveSliceElement(slice, 0)

		assert.Equal(t, []int{}, slice)
	})

	t.Run("out of range", func(t *testing.T) {
		slice := []int{1, 2, 3}
		slice = simulation.RemoveSliceElement(slice, 3)

		assert.Equal(t, []int{1, 2, 3}, slice)
	})

	t.Run("in range", func(t *testing.T) {
		slice := []int{1, 2, 3}
		slice = simulation.RemoveSliceElement(slice, 1)

		assert.Equal(t, []int{1, 3}, slice)
	})
}
