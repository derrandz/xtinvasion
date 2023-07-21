package simulation

import "fmt"

type Controller struct {
	app *App
}

func (cc *Controller) DestroyAlien(alienID int) error {
	if alienID < 0 {
		return fmt.Errorf("invalid alien ID")
	}

	if alien, exists := cc.app.Aliens[alienID]; !exists {
		return fmt.Errorf("alien %d does not exist", alienID)
	} else {
		delete(cc.app.AlienLocations[alien.CurrentCity], alienID)
		delete(cc.app.Aliens, alienID)
	}

	return nil
}

func (cc *Controller) DestroyCity(cityName string) error {
	city, found := cc.app.WorldMap.Cities[cityName]
	if !found {
		return fmt.Errorf("City not found")
	}

	msg := fmt.Sprintf("City %s has been destroyed by aliens: ", cityName)
	for _, alien := range cc.app.AlienLocations[city] {
		msg += fmt.Sprintf("%d ", alien.ID)
		cc.DestroyAlien(alien.ID)
	}

	fmt.Println(msg)

	delete(cc.app.AlienLocations, city)
	delete(cc.app.WorldMap.Cities, cityName)
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
func (cc *Controller) MoveAlienToNextCity(alien *Alien) error {
	if alien == nil {
		return fmt.Errorf("alien is nil")
	}

	_, found := cc.app.Aliens[alien.ID]
	if !found {
		return fmt.Errorf("alien %d does not exist in the world", alien.ID)
	}

	// Move the alien
	_, found = cc.app.AlienLocations[alien.CurrentCity]
	if !found {
		return fmt.Errorf("alien %d did not land in any city", alien.ID)
	}

	neighbour, err := getRandomNeighbor(alien.CurrentCity)
	if err != nil {
		return err
	}

	nextCity, found := cc.app.WorldMap.Cities[neighbour.Name]
	if !found {
		return fmt.Errorf("city %s does not exist in world map", neighbour.Name)
	}

	delete(cc.app.AlienLocations[alien.CurrentCity], alien.ID)
	alien.CurrentCity = nextCity
	alien.Moved++
	if nextCityAliens, found := cc.app.AlienLocations[nextCity]; found {
		nextCityAliens[alien.ID] = alien
	} else {
		cc.app.AlienLocations[nextCity] = AlienSet{alien.ID: alien}
	}

	return nil
}

func (cc *Controller) AreAllAliensDestroyed() bool {
	for _, alien := range cc.app.Aliens {
		if alien != nil {
			return false
		}
	}
	return true
}

func (cc *Controller) IsWorldDestroyed() bool {
	return len(cc.app.WorldMap.Cities) == 0
}

func (cc *Controller) IsAlienMovementLimitReached() bool {
	if len(cc.app.Aliens) == 0 {
		return false
	}

	trapped := 0
	for _, alien := range cc.app.Aliens {
		if alien.IsTrapped() {
			trapped++
		} else if alien != nil && alien.Moved < cc.app.MaxMoves {
			return false
		}
	}

	return trapped != len(cc.app.Aliens)
}

func (cc *Controller) AreRemainingAliensTrapped() bool {
	if len(cc.app.Aliens) == 0 {
		return false
	}

	for _, alien := range cc.app.Aliens {
		if alien != nil && !alien.IsTrapped() {
			return false
		}
	}
	return true
}

func (cc *Controller) App() *App {
	return cc.app
}

func NewController(app *App) *Controller {
	return &Controller{app: app}
}
