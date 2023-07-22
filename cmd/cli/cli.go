package main

import (
	"fmt"
	"os"

	simulation "github.com/derrandz/xtinvasion/pkg"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "app"}

	app := simulation.NewApp()

	// Add a start command
	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the simulation",
		Run:   app.Start,
	}

	// Define flags
	app.DefineFlags(startCmd)
	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
