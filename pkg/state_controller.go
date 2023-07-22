package simulation

import (
	"fmt"

	"github.com/derrandz/xtinvasion/pkg/logger"
)

// StateController handles all state operations.
type StateController struct {
	app *App

	// printer is used to print required messages
	// we differentiate between the printer and the app's logger
	// as the app's logger logs the app's errors and warnings (to file generally or stdout if specified)
	// while the printer is used to print messages to the user
	// usually to stdout (or to other writers, check cmd/tui/tui.go for example)
	printer *logger.Logger
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

	if sc.printer != nil {
		sc.printer.Log(msg)
	}

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

// SimulationResult returns a string describing the simulation termination reason.
func (sc *StateController) SimulationResult() string {
	if sc.IsWorldDestroyed() {
		return "The world has been destroyed"
	} else if sc.AreAllAliensDestroyed() {
		return "All aliens have been destroyed"
	} else if sc.IsAlienMovementLimitReached() {
		return "Alien movement limit reached"
	} else if sc.AreRemainingAliensTrapped() {
		return "All remaining aliens are trapped"
	} else {
		return "Unknown"
	}
}

// CopyState is a state getter, returns a copy of the state
// made public for testing
func (sc *StateController) CopyState() AppState {
	return *sc.app.State
}

// BroadcastStateChanges broadcasts the state changes to the state channel.
// Non-blocking
func (sc *StateController) BroadcastStateChanges() {
	select {
	case sc.app.stateCh <- sc.CopyState():
	default:
	}
}

// ListenForStateUpdates returns a channel that can be used to listen for state updates
func (sc *StateController) ListenForStateUpdates() chan AppState {
	return sc.app.stateCh
}

// App returns the app.
func (sc *StateController) App() *App {
	return sc.app
}

// SetPrinter sets the printer.
func (sc *StateController) SetPrinter(printer *logger.Logger) {
	sc.printer = printer
}

// NewStateController creates a new state controller.
func NewStateController(app *App) *StateController {
	return &StateController{app: app}
}
