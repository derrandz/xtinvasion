package simulation

import (
	"fmt"
)

// StateController handles all state operations.
type StateController struct {
	app *App
}

// DestroyAlien destroys an alien and removes it from the city
// it is currently in.
func (sc *StateController) DestroyAlien(alienID int) error {
	if alienID < 0 {
		return fmt.Errorf("invalid alien ID")
	}

	if alien, exists := sc.app.State.Aliens[alienID]; !exists {
		return fmt.Errorf("alien %d does not exist", alienID)
	} else {
		delete(sc.app.State.AlienLocations[alien.CurrentCity], alienID)
		delete(sc.app.State.Aliens, alienID)
	}

	return nil
}

// DestroyCity destroys a city and removes it from the world map
// as well as it destroys all aliens in the city.
func (sc *StateController) DestroyCity(cityName string) error {
	city, found := sc.app.State.WorldMap.Cities[cityName]
	if !found {
		return fmt.Errorf("City not found")
	}

	msg := fmt.Sprintf("City %s has been destroyed by aliens: ", cityName)
	for _, alien := range sc.app.State.AlienLocations[city] {
		msg += fmt.Sprintf("%d ", alien.ID)
		sc.DestroyAlien(alien.ID)
	}

	fmt.Println(msg)

	delete(sc.app.State.AlienLocations, city)
	delete(sc.app.State.WorldMap.Cities, cityName)
	for _, neighbour := range city.Neighbours {
		for dir, neighbourNeighbour := range neighbour.Neighbours {
			if neighbourNeighbour.Name == cityName {
				delete(neighbour.Neighbours, dir)
			}
		}
	}
	return nil
}

// MoveAlienToNextCity moves an alien to a city.
// If the alien is not in the city, it returns an error.
// If the city does not exist, it returns an error.
// If the city is isolated, it returns an error.
// Otherwise, it moves the alien to the city.
func (sc *StateController) MoveAlienToNextCity(alien *Alien) error {
	if alien == nil {
		return fmt.Errorf("alien is nil")
	}

	_, found := sc.app.State.Aliens[alien.ID]
	if !found {
		return fmt.Errorf("alien %d does not exist in the world", alien.ID)
	}

	// Move the alien
	_, found = sc.app.State.AlienLocations[alien.CurrentCity]
	if !found {
		return fmt.Errorf("alien %d did not land in any city", alien.ID)
	}

	neighbour, err := getRandomNeighbor(alien.CurrentCity)
	if err != nil {
		return err
	}

	nextCity, found := sc.app.State.WorldMap.Cities[neighbour.Name]
	if !found {
		return fmt.Errorf("city %s does not exist in world map", neighbour.Name)
	}

	delete(sc.app.State.AlienLocations[alien.CurrentCity], alien.ID)
	alien.CurrentCity = nextCity
	alien.Moved++
	if nextCityAliens, found := sc.app.State.AlienLocations[nextCity]; found {
		nextCityAliens[alien.ID] = alien
	} else {
		sc.app.State.AlienLocations[nextCity] = AlienSet{alien.ID: alien}
	}

	return nil
}

// AreAllAliensDestroyed returns true if all aliens are destroyed.
func (sc *StateController) AreAllAliensDestroyed() bool {
	for _, alien := range sc.app.State.Aliens {
		if alien != nil {
			return false
		}
	}
	return true
}

// IsWorldDestroyed returns true if all cities are destroyed.
func (sc *StateController) IsWorldDestroyed() bool {
	return len(sc.app.State.WorldMap.Cities) == 0
}

// IsAlienMovementLimitReached returns true if all aliens have reached
// the maximum number of moves.
// This method does not count trapped aliens.
func (sc *StateController) IsAlienMovementLimitReached() bool {
	if len(sc.app.State.Aliens) == 0 {
		return false
	}

	trapped := 0
	for _, alien := range sc.app.State.Aliens {
		if alien.IsTrapped() {
			trapped++
		} else if alien != nil && alien.Moved < sc.app.Cfg.MaxMoves {
			return false
		}
	}

	return trapped != len(sc.app.State.Aliens)
}

// AreRemainingAliensTrapped returns true if all remaining aliens are trapped.
func (sc *StateController) AreRemainingAliensTrapped() bool {
	if len(sc.app.State.Aliens) == 0 {
		return false
	}

	for _, alien := range sc.app.State.Aliens {
		if alien != nil && !alien.IsTrapped() {
			return false
		}
	}
	return true
}

// App returns the app.
func (sc *StateController) App() *App {
	return sc.app
}

// NewStateController creates a new state controller.
func NewStateController(app *App) *StateController {
	return &StateController{app: app}
}
