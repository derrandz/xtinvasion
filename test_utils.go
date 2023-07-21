package main

import "xtinvasion/logger"

// DummyAppConfig is a dummy config for testing.
type DummyAppConfig struct {
	AlienCount     int
	MaxMoves       int
	Map            map[string][]interface{} // [cityName, [{direction, neighbour1}, {direction, neighbour2}, ...]
	AlienLocations map[string]int
}

// NewDummyApp(cfg *DummyAppConfig) *App creates a dummy app for testing.
func NewDummyApp(cfg *DummyAppConfig) *App {
	app := &App{
		Aliens:         make(AlienSet),
		AlienLocations: make(map[*City]AlienSet),
		WorldMap:       &Map{Cities: make(map[string]*City)},
		isStopped:      0,
		done:           make(chan struct{}),
	}

	for i := 0; i < cfg.AlienCount; i++ {
		app.Aliens[i] = &Alien{ID: i, Moved: 0}
	}

	for city := range cfg.Map {
		app.WorldMap.Cities[city] = &City{Name: city, Neighbours: make(map[string]*City)}
	}

	for city, neighbours := range cfg.Map {
		for _, neighbour := range neighbours {
			neighbourCity := neighbour.(map[string]string)
			for direction, neighbourName := range neighbourCity {
				app.WorldMap.Cities[city].Neighbours[direction] = app.WorldMap.Cities[neighbourName]
			}
		}
	}

	for city, alienID := range cfg.AlienLocations {
		app.AlienLocations[app.WorldMap.Cities[city]] = AlienSet{alienID: app.Aliens[alienID]}
		app.Aliens[alienID].CurrentCity = app.WorldMap.Cities[city]
	}

	app.ctrl = NewController(app)
	app.logger = logger.NewStdoutLogger()
	return app
}

func NewEmptyDummyApp() *App {
	app := &App{}
	app.Aliens = make(AlienSet)
	app.WorldMap = &Map{Cities: make(map[string]*City)}
	app.AlienLocations = make(map[*City]AlienSet)
	app.ctrl = NewController(app)
	app.logger = logger.NewStdoutLogger()
	return app
}
