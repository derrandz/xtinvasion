package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	simulation "github.com/derrandz/xtinvasion/pkg"
	"github.com/derrandz/xtinvasion/pkg/logger"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func awaitStateUpdates(sub <-chan simulation.AppState) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

func awaitActivityUpdates(sub <-chan string) tea.Cmd {
	return func() tea.Msg {
		return <-sub
	}
}

// ChannelWriter is an implementation of io.Writer that sends bytes to a channel.
type ChannelWriter struct {
	ch chan string
}

// NewChannelWriter creates a new ChannelWriter with the given channel.
func NewChannelWriter(ch chan string) *ChannelWriter {
	return &ChannelWriter{
		ch: ch,
	}
}

// Write writes bytes to the channel.
func (w *ChannelWriter) Write(p []byte) (n int, err error) {
	w.ch <- string(p)
	return len(p), nil
}

func (w *ChannelWriter) Chan() <-chan string {
	return w.ch
}

type model struct {
	aliensTable   table.Model
	citiesTable   table.Model
	activityTable table.Model

	sub chan simulation.AppState

	activityCh chan string
}

func isAlienTrapped(alien *simulation.Alien) string {
	if alien.IsTrapped() {
		return "Yes"
	}

	return "No"
}

func (m *model) handleStateUpdate(msg tea.Msg) {
	appState := msg.(simulation.AppState)

	// Update aliens table
	newAlienRows := make([]table.Row, 0)
	for _, alien := range appState.Aliens {
		newAlienRows = append(newAlienRows, table.Row{
			fmt.Sprintf("%d", len(newAlienRows)+1),
			alien.CurrentCity.Name,
			fmt.Sprintf("%d", alien.Moved),
			isAlienTrapped(alien),
		})
	}
	m.aliensTable.SetRows(newAlienRows)

	// Update cities table
	newCityRows := make([]table.Row, 0)
	for _, city := range appState.WorldMap.Cities {
		neighbours := ""
		for direction, neighbour := range city.Neighbours {
			neighbours += fmt.Sprintf("%s=%s, ", direction, neighbour.Name)
		}

		cityAliens := ""
		for _, alien := range appState.AlienLocations[city] {
			cityAliens += fmt.Sprintf("%d, ", alien.ID)
		}

		newCityRows = append(newCityRows, table.Row{
			fmt.Sprintf("%d", len(newCityRows)+1),
			city.Name,
			neighbours,
			cityAliens,
		})
	}

	m.citiesTable.SetRows(newCityRows[:4])
}

func (m *model) handleActivityUpdate(msg string) {
	m.activityTable.SetRows(
		append(m.activityTable.Rows(), table.Row{string(msg)}),
	)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		awaitActivityUpdates(m.activityCh),
		tea.Batch(awaitStateUpdates(m.sub), tickCmd()),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.aliensTable.Focused() {
				m.aliensTable.Blur()
				m.citiesTable.Focus()
			} else if m.citiesTable.Focused() {
				m.aliensTable.Blur()
				m.citiesTable.Blur()
				m.activityTable.Focus()
			} else {
				m.aliensTable.Focus()
				m.citiesTable.Blur()
				m.activityTable.Blur()
			}

			return m, cmd
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.aliensTable.SelectedRow()[1]),
			)
		}

	case simulation.AppState:
		m.handleStateUpdate(msg)
		return m, tea.Batch(cmd, tea.Batch(awaitStateUpdates(m.sub), tickCmd()))

	case string:
		m.handleActivityUpdate(msg)
		return m, tea.Batch(cmd, awaitActivityUpdates(m.activityCh))
	}

	if m.aliensTable.Focused() {
		m.aliensTable, cmd = m.aliensTable.Update(msg)
	} else if m.citiesTable.Focused() {
		m.citiesTable, cmd = m.citiesTable.Update(msg)
	} else {
		m.activityTable, cmd = m.activityTable.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.JoinVertical(
				lipgloss.Left,
				baseStyle.Render(m.aliensTable.View()),
			),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.JoinVertical(
				lipgloss.Right,
				baseStyle.Render(m.citiesTable.View()),
			),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			baseStyle.Render(m.activityTable.View()),
		),
	)
}

func runTUI(app *simulation.App) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		aliensColumns := []table.Column{
			{Title: "ID", Width: 10},
			{Title: "Current City", Width: 15},
			{Title: "Moved", Width: 5},
			{Title: "Is Trapped ?", Width: 10},
		}
		aliensRows := []table.Row{}

		citiesColumns := []table.Column{
			{Title: "ID", Width: 10},
			{Title: "City", Width: 10},
			{Title: "Neighbours", Width: 50},
			{Title: "Aliens", Width: 10},
		}
		citiesRows := []table.Row{}

		activityColumns := []table.Column{
			{Title: "Activity", Width: 101},
		}
		activityRows := []table.Row{}

		at := table.New(
			table.WithColumns(aliensColumns),
			table.WithRows(aliensRows),
			table.WithFocused(true),
			table.WithHeight(7),
		)

		ct := table.New(
			table.WithColumns(citiesColumns),
			table.WithRows(citiesRows),
			table.WithFocused(true),
			table.WithHeight(7),
		)

		act := table.New(
			table.WithColumns(activityColumns),
			table.WithRows(activityRows),
			table.WithFocused(true),
			table.WithHeight(7),
		)

		s := table.DefaultStyles()
		s.Header = s.Header.
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false)
		s.Selected = s.Selected.
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false)
		at.SetStyles(s)
		ct.SetStyles(s)
		act.SetStyles(s)

		sub := make(chan simulation.AppState)
		m := model{
			at,
			ct,
			act,
			sub,
			make(chan string),
		}

		go func() {
			activityWriter := NewChannelWriter(m.activityCh)
			activityPrinter := logger.NewLogger(activityWriter)
			app.Init(cmd)
			app.StateController().SetPrinter(activityPrinter)
			app.Run()
		}()

		go func() {
			<-app.Ready()
			for {
				select {
				case state := <-app.StateController().ListenForStateUpdates():
					sub <- state
				}
			}
		}()

		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	}
}

func main() {
	var rootCmd = &cobra.Command{Use: "app"}

	app := simulation.NewApp()

	// Add a start command
	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the simulation with terminal UI",
		Run:   runTUI(app),
	}

	// Define flags
	app.DefineFlags(startCmd)
	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
