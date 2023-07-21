package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"

	"github.com/spf13/cobra"
)

type City struct {
	Name       string
	Neighbours map[string]*City
}

type Map struct {
	Cities map[string]*City
}

type Alien struct {
	ID          int
	CurrentCity *City
	Moved       int
}

type AlienSet map[int]*Alien

type App struct {
	ctrl *Controller

	Aliens         AlienSet
	AlienLocations map[*City]AlienSet
	WorldMap       *Map

	isStopped int32 // Use int32 for atomic operations
	done      chan struct{}
}

func (a *App) createAliens(numAliens int) {
	for i := 0; i < numAliens; i++ {
		alien := &Alien{ID: i, Moved: 0}
		a.Aliens[i] = alien
	}
	fmt.Println("Created", numAliens, "aliens.", a.Aliens)
}

func (a *App) populateMapWithAliens() {
	for _, alien := range a.Aliens {
		city := a.getRandomCity()
		if location, found := a.AlienLocations[city]; found {
			location[alien.ID] = alien
		} else {
			a.AlienLocations[city] = map[int]*Alien{alien.ID: alien}
		}
		alien.CurrentCity = city
	}
}

func (a *App) getRandomCity() *City {
	var cities []*City
	for _, city := range a.WorldMap.Cities {
		cities = append(cities, city)
	}

	if len(cities) == 0 {
		return nil
	}

	return cities[rand.Intn(len(cities))]
}

func (a *App) parseFlags(cmd *cobra.Command) []any {
	numAliens, _ := cmd.Flags().GetInt("aliens")
	filename, _ := cmd.Flags().GetString("file")

	return []any{numAliens, filename}
}

func (a *App) readMapFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read and create all cities first
	for scanner.Scan() {
		line := scanner.Text()
		cityData := strings.Split(line, " ")

		if len(cityData) < 2 {
			fmt.Printf("Invalid line: %s\n", line)
			continue
		}

		cityName := cityData[0]

		city := &City{Name: cityName, Neighbours: make(map[string]*City)}
		a.WorldMap.Cities[cityName] = city
	}

	// Reset scanner to start again from the beginning
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)

	// Populate neighboring cities
	for scanner.Scan() {
		line := scanner.Text()
		cityData := strings.Split(line, " ")

		if len(cityData) < 2 {
			fmt.Printf("Invalid line: %s\n", line)
			continue
		}

		cityName := cityData[0]
		cityNeighbours := cityData[1:]

		city := a.WorldMap.Cities[cityName]

		for _, neighbourData := range cityNeighbours {
			neighbour := strings.Split(neighbourData, "=")
			if len(neighbour) != 2 {
				fmt.Printf("Invalid neighbour data: %s\n", neighbourData)
				continue
			}

			neighbourName := neighbour[1]
			direction := neighbour[0]

			if destCity, found := a.WorldMap.Cities[neighbourName]; !found {
				fmt.Printf("Neighbour city %s not found for %s\n", neighbourName, cityName)
			} else {
				city.Neighbours[direction] = destCity
				destCity.Neighbours[oppositeDirection(direction)] = city
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	fmt.Println("Map read successfully.")
	fmt.Println("Cities:", len(a.WorldMap.Cities))

	return nil
}

func (a *App) DefineFlags(cmd *cobra.Command) {
	cmd.Flags().IntP("aliens", "a", 5, "Number of aliens")
	cmd.Flags().StringP("file", "f", "map.txt", "Map file")
}

func (a *App) Init(cmd *cobra.Command) {
	a.done = make(chan struct{})
	a.isStopped = 0

	// Read the map from the file and create the cities
	flags := a.parseFlags(cmd)

	numAliens := flags[0].(int)
	filename := flags[1].(string)

	a.Aliens = make(AlienSet, numAliens)
	a.WorldMap = &Map{Cities: make(map[string]*City)}
	a.AlienLocations = make(map[*City]AlienSet)

	// Read the map from the file and create the cities
	a.readMapFromFile(filename)

	// Create aliens and assign them to cities
	a.createAliens(numAliens)

	// Populate the alien locations
	a.populateMapWithAliens()

	// Initialize the queryController and commandController
	a.ctrl = &Controller{app: a}
}

func (a *App) Run() {
	a.PrintState()
	for {
		// Check if the app has been stopped
		if atomic.LoadInt32(&a.isStopped) == 1 {
			break
		}

		// Check if all aliens have been destroyed
		if a.ctrl.AreAllAliensDestroyed() {
			fmt.Println("All aliens have been destroyed.")
			break
		}

		// Check if all aliens have moved 10,000 times
		if a.ctrl.IsAlienMovementLimitReached() {
			fmt.Println("All aliens have moved 10,000 times.")
			break
		}

		if a.ctrl.AreRemainingAliensTrapped() {
			fmt.Println("All remaining aliens are trapped.")
			break
		}

		// Check if any city has two or more aliens and destroy them
		for city := range a.AlienLocations {
			if len(a.AlienLocations[city]) > 1 {
				err := a.ctrl.DestroyCity(city.Name)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		// Move aliens around in the map
		for _, alien := range a.Aliens {
			if alien != nil {
				nextCity, err := getRandomNeighbor(alien.CurrentCity)
				if err != nil {
					fmt.Println(err)
					continue
				}
				err = a.ctrl.MoveAlienToCity(alien.ID, nextCity.Name)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		// Print the current state of the app
		fmt.Println("=====================================")
		a.PrintState()
	}

	// Indicate that the main loop has finished by closing the channel
	close(a.done)
}

func (a *App) PrintState() {
	fmt.Println("Remaining Cities:")
	for _, city := range a.WorldMap.Cities {
		fmt.Printf("%s ", city.Name)
		if len(city.Neighbours) > 0 {
			fmt.Printf("connecting to %v", city.Neighbours)
			var neighbours []string
			for _, neighbour := range city.Neighbours {
				neighbours = append(neighbours, neighbour.Name)
			}
			fmt.Printf("%s\n", strings.Join(neighbours, ", "))
		} else {
			fmt.Printf("isolated\n")
		}
	}

	fmt.Println("\nRemaining Aliens:")
	if len(a.Aliens) > 0 {
		for _, alien := range a.Aliens {
			if alien != nil {
				fmt.Printf("Alien %d at %s, moved %d times\n", alien.ID, alien.CurrentCity.Name, alien.Moved)
			}
		}
	} else {
		fmt.Println("No aliens left.")
	}
}

func (a *App) Stop() {
	atomic.StoreInt32(&a.isStopped, 1)
}

func (a *App) Wait() {
	// Wait for the main loop to finish by waiting for the loopDone channel to be closed
	<-a.done
}

func (a *App) Start(cmd *cobra.Command, args []string) {
	a.Init(cmd)
	a.Run()
	a.PrintState()
}
