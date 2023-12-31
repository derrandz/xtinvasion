package simulation

import (
	"fmt"
	"math/rand"
	"time"
)

// functions made public for testing

// getRandomCity returns a random city from the map
func GetRandomNeighbor(city *City) (*City, error) {
	if city == nil {
		return nil, fmt.Errorf("GetRandomNeighbor: city is nil")
	}

	if len(city.Neighbours) == 0 {
		return nil, fmt.Errorf("GetRandomNeighbor: city has no neighbours or neighbours have been destroyed. city=%s", city.Name)
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(city.Neighbours))
	i := 0
	for neighbour := range city.Neighbours {
		if i == index {
			if city.Neighbours[neighbour] == nil {
				return nil, fmt.Errorf("GetRandomNeighbor: city %s has a nil neighbour", city.Name)
			}
			return city.Neighbours[neighbour], nil
		}
		i++
	}

	// This should not happen, but return nil as a safety measure
	return nil, fmt.Errorf("getRandomNeighbor: could not find a random neighbour")
}

// oppositeDirection returns the opposite direction of the given direction
// used during map creation
func OppositeDirection(direction string) string {
	switch direction {
	case "north":
		return "south"
	case "south":
		return "north"
	case "east":
		return "west"
	case "west":
		return "east"
	default:
		return ""
	}
}

// removeSliceElement removes an element from a slice
// used during state updates
func RemoveSliceElement[T any](slice []T, index int) []T {
	// Check if the index is out of range
	if index < 0 || index >= len(slice) {
		return slice
	}

	// Create a new slice without the element at the given index
	newSlice := append(slice[:index], slice[index+1:]...)
	return newSlice
}

// sleepMs sleeps for the given number of milliseconds
func sleepMS(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
