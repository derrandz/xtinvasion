package simulation

import "fmt"

type City struct {
	Name       string
	Neighbours map[string]*City
}

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

type Map struct {
	Cities map[string]*City
}

type Alien struct {
	ID          int
	CurrentCity *City
	Moved       int
}

func (a *Alien) IsTrapped() bool {
	if a.CurrentCity == nil {
		panic("alien is not in any city")
	}
	return len(a.CurrentCity.Neighbours) == 0
}

type AlienSet map[int]*Alien
