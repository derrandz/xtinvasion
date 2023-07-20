package main

func NewDummyApp() *App {
	app := &App{
		WorldMap: &Map{
			Cities: make(map[string]*City),
		},
		Aliens:         make([]*Alien, 6),
		AlienLocations: make(map[*City][]*Alien),
	}

	// Create aliens.
	app.Aliens[0] = &Alien{ID: 0, Moved: 0}
	app.Aliens[1] = &Alien{ID: 1, Moved: 0}
	app.Aliens[2] = &Alien{ID: 2, Moved: 0}
	app.Aliens[3] = &Alien{ID: 3, Moved: 0}

	cities := []*City{
		&City{Name: "A"},
		&City{Name: "B"},
		&City{Name: "C"},
		&City{Name: "D"},
	}

	// Create cities.
	//
	// CityA north=CityB south=CityC
	// CityB east=CityD south=CityA
	// CityC north=CityA west=CityD
	// CityD west=CityB east=CityC

	app.WorldMap.Cities["A"] = cities[0]
	app.WorldMap.Cities["A"].Neighbours = map[string]*City{
		"north": cities[1],
		"south": cities[2],
	}

	app.WorldMap.Cities["B"] = cities[1]
	app.WorldMap.Cities["B"].Neighbours = map[string]*City{
		"east":  cities[3],
		"south": cities[0],
	}

	app.WorldMap.Cities["C"] = cities[2]
	app.WorldMap.Cities["C"].Neighbours = map[string]*City{
		"north": cities[0],
		"west":  cities[3],
	}

	app.WorldMap.Cities["D"] = cities[2]
	app.WorldMap.Cities["D"].Neighbours = map[string]*City{
		"west": cities[1],
		"east": cities[2],
	}

	app.AlienLocations[cities[0]] = []*Alien{app.Aliens[0]}
	app.Aliens[0].CurrentCity = cities[0]

	app.AlienLocations[cities[1]] = []*Alien{app.Aliens[1]}
	app.Aliens[1].CurrentCity = cities[1]

	app.AlienLocations[cities[2]] = []*Alien{app.Aliens[2]}
	app.Aliens[2].CurrentCity = cities[2]

	app.AlienLocations[cities[3]] = []*Alien{app.Aliens[3]}
	app.Aliens[3].CurrentCity = cities[3]

	return app
}
