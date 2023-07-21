package simulation

import (
	"fmt"
	"math/rand"
	"strings"
	"sync/atomic"

	"github.com/derrandz/xtinvasion/pkg/logger"

	"github.com/spf13/cobra"
)

type AppCfg struct {
	MaxMoves      int
	MapInputFile  string
	LogFile       string
	MapOutputFile string
}

type App struct {
	logger *logger.Logger

	stateCtrl *StateController
	ioCtrl    *IOController

	Cfg *AppCfg

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

func (a *App) PopulateMapWithAliens() {
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

func (a *App) DefineFlags(cmd *cobra.Command) {
	cmd.Flags().IntP("aliens", "a", 5, "Number of aliens")
	cmd.Flags().IntP("max_moves", "m", 10000, "Max number of moves allowed for each alien")
	cmd.Flags().StringP("input", "i", "data/map.txt", "Map input file")
	cmd.Flags().StringP("output", "l", "output/map.txt", "Map output file")
	cmd.Flags().StringP("log", "o", "output/stdout.log", "Log file")
}

func (a *App) parseFlags(cmd *cobra.Command) []any {
	numAliens, _ := cmd.Flags().GetInt("aliens")
	maxMoves, _ := cmd.Flags().GetInt("max_moves")
	inputFilename, _ := cmd.Flags().GetString("input")
	outputFilename, _ := cmd.Flags().GetString("output")
	logfile, _ := cmd.Flags().GetString("log")

	return []any{numAliens, maxMoves, inputFilename, outputFilename, logfile}
}

func (a *App) Init(cmd *cobra.Command) {
	a.done = make(chan struct{})
	a.isStopped = 0

	// Read the map from the file and create the cities
	flags := a.parseFlags(cmd)

	// store configuration
	a.Cfg = &AppCfg{
		MaxMoves:      flags[1].(int),
		MapInputFile:  flags[2].(string),
		MapOutputFile: flags[3].(string),
		LogFile:       flags[4].(string),
	}

	// Initialize the logger
	if a.Cfg.LogFile == "" {
		a.logger = logger.NewStdoutLogger()
	} else {
		if loggr, err := logger.NewFileLogger(a.Cfg.LogFile); err != nil {
			fmt.Printf("error creating logger: %v", err)
			panic(err)
		} else {
			a.logger = loggr
		}
	}

	// Initialize the state and io controllers
	a.stateCtrl = &StateController{app: a}
	a.ioCtrl = &IOController{app: a}

	// Initialize the map and aliens (state)
	numAliens := flags[0].(int)

	a.Aliens = make(AlienSet, numAliens)
	a.WorldMap = &Map{Cities: make(map[string]*City)}
	a.AlienLocations = make(map[*City]AlienSet)

	// Read the map from the file and create the cities
	if err := a.ioCtrl.ReadMapFromFile(a.Cfg.MapInputFile); err != nil {
		a.logger.Logf("error: %v", err)
		panic(err)
	}

	// Create aliens and assign them to cities
	a.createAliens(numAliens)

	// Populate the alien locations
	a.PopulateMapWithAliens()
}

func (a *App) Run() {
	for {
		// Check if the app has been stopped
		if atomic.LoadInt32(&a.isStopped) == 1 {
			break
		}

		// Check if all aliens have been destroyed
		if a.stateCtrl.AreAllAliensDestroyed() {
			a.logger.Log("All aliens have been destroyed.")
			break
		}

		// Check if all aliens have moved 10,000 times
		if a.stateCtrl.IsAlienMovementLimitReached() {
			a.logger.Log("All aliens have moved 10,000 times.")
			break
		}

		if a.stateCtrl.AreRemainingAliensTrapped() {
			a.logger.Log("All remaining aliens are trapped.")
			break
		}

		// Check if any city has two or more aliens and destroy them
		for city := range a.AlienLocations {
			if len(a.AlienLocations[city]) > 1 {
				err := a.stateCtrl.DestroyCity(city.Name)
				if err != nil {
					a.logger.Logf("error: %v", err)
				}
			}
		}

		// Move aliens around in the map
		for _, alien := range a.Aliens {
			if alien != nil {
				err := a.stateCtrl.MoveAlienToNextCity(alien)
				if err != nil && !strings.Contains(err.Error(), "getRandomNeighbor: city has no neighbours") {
					a.logger.Logf("error: %v", err)
				}
			}
		}
	}

	// Indicate that the main loop has finished by closing the channel
	close(a.done)
}

func (a *App) Stop() {
	atomic.StoreInt32(&a.isStopped, 1)
}

func (a *App) Wait() {
	// Wait for the main loop to finish by waiting for the loopDone channel to be closed
	<-a.done
}

func (a *App) SaveResult() {
	a.ioCtrl.WriteMapToFile()
	a.ioCtrl.PrintResult()
}

func (a *App) IsStopped() bool {
	return atomic.LoadInt32(&a.isStopped) == 1
}

func (a *App) Start(cmd *cobra.Command, args []string) {
	a.Init(cmd)
	a.Run()
	a.SaveResult()
}

// Setter for testing
func (a *App) SetStateController(stateCtrl *StateController) {
	a.stateCtrl = stateCtrl
}

// Getter for testing
func (a *App) StateController() *StateController {
	return a.stateCtrl
}

// Getter for testing
func (a *App) IOController() *IOController {
	return a.ioCtrl
}

// Setter for testing
func (a *App) SetIOController(ioCtrl *IOController) {
	a.ioCtrl = ioCtrl
}

func (a *App) SetLogger(logger *logger.Logger) {
	a.logger = logger
}

func NewApp() *App {
	app := &App{
		done:      make(chan struct{}),
		isStopped: 0,
		Cfg:       &AppCfg{},
	}
	return app
}
