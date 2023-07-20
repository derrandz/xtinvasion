package main

import "fmt"

type Controller struct {
	app *App
}

func NewController(app *App) *Controller {
	return &Controller{app: app}
}

func (cc *Controller) DestroyAlien(alienID int) error {
	if alienID < 0 || alienID >= len(cc.app.Aliens) {
		return fmt.Errorf("Invalid alien ID.")
	}

	alien := cc.app.Aliens[alienID]
	if alien != nil {
		aliensInCity := cc.app.AlienLocations[alien.CurrentCity]
		aliensInCity = append(aliensInCity[:findAlienIndex(aliensInCity, alienID)], aliensInCity[findAlienIndex(aliensInCity, alienID)+1:]...)
		cc.app.AlienLocations[alien.CurrentCity] = aliensInCity
		cc.app.Aliens[alienID] = nil
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

	delete(cc.app.WorldMap.Cities, cityName)
	return nil
}

func (cc *Controller) MoveAlienToCity(alienID int, cityName string) error {
	if alienID < 0 || alienID >= len(cc.app.Aliens) {
		return fmt.Errorf("invalid alien ID")
	}

	alien := cc.app.Aliens[alienID]
	if alien == nil {
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
	aliensInCity, found := cc.app.AlienLocations[alien.CurrentCity]
	if !found {
		return fmt.Errorf("alien %d is not in city %s", alienID, alien.CurrentCity.Name)
	}

	aliensInCity = append(aliensInCity[:findAlienIndex(aliensInCity, alienID)], aliensInCity[findAlienIndex(aliensInCity, alienID)+1:]...)
	cc.app.AlienLocations[alien.CurrentCity] = aliensInCity
	cc.app.AlienLocations[nextCity] = append(cc.app.AlienLocations[nextCity], alien)
	alien.CurrentCity = nextCity

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

func (cc *Controller) IsAlienMovementLimitReached() bool {
	for _, alien := range cc.app.Aliens {
		if alien != nil && alien.Moved < 10000 {
			return false
		}
	}
	return true
}
