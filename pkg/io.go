package simulation

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// IOController handles all input and output operations.
type IOController struct {
	app *App
}

// ReadMapFromFile reads the world map from a file.
func (io *IOController) ReadMapFromFile() error {
	file, err := os.Open(io.app.Cfg.MapInputFile)
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
			return fmt.Errorf("invalid line: %s", line)
		}

		cityName := cityData[0]

		city := &City{Name: cityName, Neighbours: make(map[string]*City)}
		io.app.State.WorldMap.Cities[cityName] = city
	}

	// Reset scanner to start again from the beginning
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)

	// Populate neighboring cities
	for scanner.Scan() {
		line := scanner.Text()
		cityData := strings.Split(line, " ")

		if len(cityData) < 2 {
			return fmt.Errorf("invalid line: %s", line)
		}

		cityName := cityData[0]
		cityNeighbours := cityData[1:]

		city := io.app.State.WorldMap.Cities[cityName]

		for _, neighbourData := range cityNeighbours {
			neighbour := strings.Split(neighbourData, "=")
			if len(neighbour) != 2 {
				return fmt.Errorf("invalid neighbour data: %s", neighbourData)
			}

			neighbourName := neighbour[1]
			direction := neighbour[0]

			if destCity, found := io.app.State.WorldMap.Cities[neighbourName]; !found {
				return fmt.Errorf("neighbour city %s not found for %s", neighbourName, cityName)
			} else {
				city.Neighbours[direction] = destCity
				destCity.Neighbours[oppositeDirection(direction)] = city
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	io.app.logger.Log("Map read successfully.")
	io.app.logger.Logf("Cities: %d", len(io.app.State.WorldMap.Cities))

	return nil
}

// WriteMapToFile writes the world map to a file in the same format as the input.
func (io *IOController) WriteMapToFile() error {
	file, err := os.Create(io.app.Cfg.MapOutputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	for cityName, city := range io.app.State.WorldMap.Cities {
		line := fmt.Sprintf("%s ", cityName)
		for direction, neighbour := range city.Neighbours {
			line += fmt.Sprintf("%s=%s", direction, neighbour.Name)
		}
		line += "\n"

		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}

// printResult prints the remaining cities and aliens in separate tables.
func (io *IOController) PrintResult() {
	app := io.app

	fmt.Println()
	fmt.Println("+-------------------------- Simulation Result --------------------------+")

	fmt.Println("Remaining Cities:")
	printCities(app.State.WorldMap)

	fmt.Println("\nRemaining Aliens:")
	printAliens(app.State.Aliens)

	fmt.Println("+-----------------------------------------------------------------------+")
	fmt.Println("The resulting map of the world is saved to:", app.Cfg.MapOutputFile)
}

// printCities prints the remaining cities in a table.
func printCities(worldMap *Map) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"City", "Neighbours"})

	for _, city := range worldMap.Cities {
		row := []string{city.Name, ""}
		for direction, neighbour := range city.Neighbours {
			row[1] += fmt.Sprintf("%s=%s ", direction, neighbour.Name)
		}
		table.Append(row)
	}

	table.Render()
}

// printAliens prints the remaining aliens in a table.
func printAliens(aliens AlienSet) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Current City", "Moves"})

	for _, alien := range aliens {
		table.Append([]string{fmt.Sprintf("%d", alien.ID), fmt.Sprintf("%s", alien.CurrentCity.Name), fmt.Sprintf("%d", alien.Moved)})
	}

	table.Render()
}

// NewIOController creates a new IOController.
func NewIOController(app *App) *IOController {
	return &IOController{app: app}
}
