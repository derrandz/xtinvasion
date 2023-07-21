package simulation

import "fmt"

type StateController struct {
	app *App
}

func (sc *StateController) DestroyAlien(alienID int) error {
	if alienID < 0 {
		return fmt.Errorf("invalid alien ID")
	}

	if alien, exists := sc.app.Aliens[alienID]; !exists {
		return fmt.Errorf("alien %d does not exist", alienID)
	} else {
		delete(sc.app.AlienLocations[alien.CurrentCity], alienID)
		delete(sc.app.Aliens, alienID)
	}

	return nil
}

func (sc *StateController) DestroyCity(cityName string) error {
	city, found := sc.app.WorldMap.Cities[cityName]
	if !found {
		return fmt.Errorf("City not found")
	}

	msg := fmt.Sprintf("City %s has been destroyed by aliens: ", cityName)
	for _, alien := range sc.app.AlienLocations[city] {
		msg += fmt.Sprintf("%d ", alien.ID)
		sc.DestroyAlien(alien.ID)
	}

	fmt.Println(msg)

	delete(sc.app.AlienLocations, city)
	delete(sc.app.WorldMap.Cities, cityName)
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

	_, found := sc.app.Aliens[alien.ID]
	if !found {
		return fmt.Errorf("alien %d does not exist in the world", alien.ID)
	}

	// Move the alien
	_, found = sc.app.AlienLocations[alien.CurrentCity]
	if !found {
		return fmt.Errorf("alien %d did not land in any city", alien.ID)
	}

	neighbour, err := getRandomNeighbor(alien.CurrentCity)
	if err != nil {
		return err
	}

	nextCity, found := sc.app.WorldMap.Cities[neighbour.Name]
	if !found {
		return fmt.Errorf("city %s does not exist in world map", neighbour.Name)
	}

	delete(sc.app.AlienLocations[alien.CurrentCity], alien.ID)
	alien.CurrentCity = nextCity
	alien.Moved++
	if nextCityAliens, found := sc.app.AlienLocations[nextCity]; found {
		nextCityAliens[alien.ID] = alien
	} else {
		sc.app.AlienLocations[nextCity] = AlienSet{alien.ID: alien}
	}

	return nil
}

func (sc *StateController) AreAllAliensDestroyed() bool {
	for _, alien := range sc.app.Aliens {
		if alien != nil {
			return false
		}
	}
	return true
}

func (sc *StateController) IsWorldDestroyed() bool {
	return len(sc.app.WorldMap.Cities) == 0
}

func (sc *StateController) IsAlienMovementLimitReached() bool {
	if len(sc.app.Aliens) == 0 {
		return false
	}

	trapped := 0
	for _, alien := range sc.app.Aliens {
		if alien.IsTrapped() {
			trapped++
		} else if alien != nil && alien.Moved < sc.app.Cfg.MaxMoves {
			return false
		}
	}

	return trapped != len(sc.app.Aliens)
}

func (sc *StateController) AreRemainingAliensTrapped() bool {
	if len(sc.app.Aliens) == 0 {
		return false
	}

	for _, alien := range sc.app.Aliens {
		if alien != nil && !alien.IsTrapped() {
			return false
		}
	}
	return true
}

func (sc *StateController) App() *App {
	return sc.app
}

func NewStateController(app *App) *StateController {
	return &StateController{app: app}
}
