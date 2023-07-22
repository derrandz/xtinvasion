package tests

import (
	simulation "github.com/derrandz/xtinvasion/pkg"
	logger "github.com/derrandz/xtinvasion/pkg/logger"
)

// DummyAppConfig is a dummy config for testing.
type DummyAppConfig struct {
	AlienCount     int
	MaxMoves       int
	Map            map[string][]interface{} // [cityName, [{direction, neighbour1}, {direction, neighbour2}, ...]
	AlienLocations map[string][]int
}

// NewDummyApp(cfg *DummyAppConfig) *app.App creates a dummy app for testing.
func NewDummyApp(cfg *DummyAppConfig) *simulation.App {
	app := simulation.NewApp()

	app.Cfg.MaxMoves = cfg.MaxMoves
	app.State = &simulation.AppState{
		Aliens:         make(simulation.AlienSet),
		AlienLocations: make(map[*simulation.City]simulation.AlienSet),
		WorldMap:       &simulation.Map{Cities: make(map[string]*simulation.City)},
	}

	for i := 0; i < cfg.AlienCount; i++ {
		app.State.Aliens[i] = &simulation.Alien{ID: i, Moved: 0}
	}

	for city := range cfg.Map {
		app.State.WorldMap.Cities[city] = &simulation.City{Name: city, Neighbours: make(map[string]*simulation.City)}
	}

	for city, neighbours := range cfg.Map {
		for _, neighbour := range neighbours {
			neighbourCity := neighbour.(map[string]string)
			for direction, neighbourName := range neighbourCity {
				app.State.WorldMap.Cities[city].Neighbours[direction] = app.State.WorldMap.Cities[neighbourName]
			}
		}
	}

	for city, alienIDs := range cfg.AlienLocations {
		app.State.AlienLocations[app.State.WorldMap.Cities[city]] = simulation.AlienSet{}
		for _, id := range alienIDs {
			app.State.Aliens[id].CurrentCity = app.State.WorldMap.Cities[city]
			app.State.AlienLocations[app.State.WorldMap.Cities[city]][id] = app.State.Aliens[id]
		}
	}

	app.SetStateController(simulation.NewStateController(app))
	app.SetIOController(simulation.NewIOController(app))
	app.SetLogger(logger.NewStdoutLogger())

	return app
}

func NewEmptyDummyApp() *simulation.App {
	app := &simulation.App{
		State: &simulation.AppState{
			Aliens:         make(simulation.AlienSet),
			WorldMap:       &simulation.Map{Cities: make(map[string]*simulation.City)},
			AlienLocations: make(map[*simulation.City]simulation.AlienSet),
		},
		Cfg: &simulation.AppCfg{},
	}

	app.SetStateController(simulation.NewStateController(app))
	app.SetIOController(simulation.NewIOController(app))
	app.SetLogger(logger.NewStdoutLogger())

	return app
}
