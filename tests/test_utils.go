package tests

import (
	logger "github.com/derrandz/xtinvasion/logger"
	simulation "github.com/derrandz/xtinvasion/pkg"
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

	app.MaxMoves = cfg.MaxMoves
	app.Aliens = make(simulation.AlienSet)
	app.AlienLocations = make(map[*simulation.City]simulation.AlienSet)
	app.WorldMap = &simulation.Map{Cities: make(map[string]*simulation.City)}

	for i := 0; i < cfg.AlienCount; i++ {
		app.Aliens[i] = &simulation.Alien{ID: i, Moved: 0}
	}

	for city := range cfg.Map {
		app.WorldMap.Cities[city] = &simulation.City{Name: city, Neighbours: make(map[string]*simulation.City)}
	}

	for city, neighbours := range cfg.Map {
		for _, neighbour := range neighbours {
			neighbourCity := neighbour.(map[string]string)
			for direction, neighbourName := range neighbourCity {
				app.WorldMap.Cities[city].Neighbours[direction] = app.WorldMap.Cities[neighbourName]
			}
		}
	}

	for city, alienIDs := range cfg.AlienLocations {
		app.AlienLocations[app.WorldMap.Cities[city]] = simulation.AlienSet{}
		for _, id := range alienIDs {
			app.Aliens[id].CurrentCity = app.WorldMap.Cities[city]
			app.AlienLocations[app.WorldMap.Cities[city]][id] = app.Aliens[id]
		}
	}

	app.SetController(simulation.NewController(app))
	app.SetLogger(logger.NewStdoutLogger())

	return app
}

func NewEmptyDummyApp() *simulation.App {
	app := &simulation.App{}
	app.Aliens = make(simulation.AlienSet)
	app.WorldMap = &simulation.Map{Cities: make(map[string]*simulation.City)}
	app.AlienLocations = make(map[*simulation.City]simulation.AlienSet)

	app.SetController(simulation.NewController(app))
	app.SetLogger(logger.NewStdoutLogger())

	return app
}
