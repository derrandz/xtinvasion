package main

import "fmt"

type Controller struct {
	app *App
}

func NewController(app *App) *Controller {
	return &Controller{app: app}
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

// MoveAlienToCity moves an alien to a city.
// If the alien is not in the city, it returns an error.
// If the city does not exist, it returns an error.
// If the city is isolated, it returns an error.
// Otherwise, it moves the alien to the city.
// Checking for whether the alien is already in the city is omitted to ease up testing
// and such case would be avoided thanks to the caller's logic
func (cc *Controller) MoveAlienToCity(alienID int, cityName string) error {
	if alienID < 0 {
		return fmt.Errorf("invalid alien ID")
	}

	alien, found := cc.app.Aliens[alienID]
	if !found {
		return fmt.Errorf("alien %d does not exist", alienID)
	}

	nextCity, found := cc.app.WorldMap.Cities[cityName]
	if !found {
		return fmt.Errorf("city %s does not exist", cityName)
	}

	if len(nextCity.Neighbours) == 0 {
		return fmt.Errorf("city %s is isolated, alien %d cannot move", cityName, alienID)
	}

	// Move the alien
	_, found = cc.app.AlienLocations[alien.CurrentCity]
	if !found {
		return fmt.Errorf("alien %d is not in city %s", alienID, alien.CurrentCity.Name)
	}

	delete(cc.app.AlienLocations[alien.CurrentCity], alienID)
	alien.CurrentCity = nextCity
	alien.Moved++
	cc.app.AlienLocations[nextCity][alien.ID] = alien

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
	for _, alien := range cc.app.Aliens {
		if alien != nil && alien.Moved < 10000 {
			return false
		}
	}
	return true
}

func (cc *Controller) AreRemainingAliensTrapped() bool {
	for _, alien := range cc.app.Aliens {
		if alien != nil && len(alien.CurrentCity.Neighbours) > 0 {
			return false
		}
	}
	return true
}
