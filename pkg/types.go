package simulation

import "fmt"

// City is a city in the world map
type City struct {
	Name       string
	Neighbours map[string]*City
}

// String returns a string representation of the city
func (c *City) String() string {
	neighboursStr := ""
	for direction, neighbour := range c.Neighbours {
		if neighbour == nil {
			neighboursStr += fmt.Sprintf(",{%s: %s}", direction, "nil")
		} else {
			neighboursStr += fmt.Sprintf(",{%s: %s}", direction, neighbour.Name)
		}
	}

	return fmt.Sprintf("City #{%s, Neighbours; %s}", c.Name, neighboursStr)
}

// Map is the world map
type Map struct {
	Cities map[string]*City
}

// Alien is an alien in the world
type Alien struct {
	ID          int
	CurrentCity *City
	Moved       int
}

// IsTrapped returns true if the alien is trapped in a city
func (a *Alien) IsTrapped() bool {
	if a.CurrentCity == nil {
		panic("alien is not in any city")
	}
	return len(a.CurrentCity.Neighbours) == 0
}

// AlienSet is a set of aliens
type AlienSet map[int]*Alien
