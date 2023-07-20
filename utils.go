package main

import (
	"fmt"
	"math/rand"
	"time"
)

func getRandomNeighbor(city *City) *City {
	if city == nil || len(city.Neighbours) == 0 {
		fmt.Println("getRandomNeighbor: city is nil or has no neighbours", city, city.Neighbours)
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(city.Neighbours))
	i := 0
	for neighbour := range city.Neighbours {
		if i == index {
			return city.Neighbours[neighbour]
		}
		i++
	}

	// This should not happen, but return nil as a safety measure
	return nil
}

func oppositeDirection(direction string) string {
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

func findAlienIndex(aliens []*Alien, alienID int) int {
	for i, alien := range aliens {
		if alien != nil && alien.ID == alienID {
			return i
		}
	}
	return -1
}
