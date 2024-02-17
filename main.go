package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/g-bolmida/weather-app/helpers"
	"github.com/g-bolmida/weather-app/models"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table    table.Model
	location string
	weather  models.OpenMeteoAPI
}

func initialModel() model {
	ipinfo := helpers.IPInfo()

	latlong := strings.Split(ipinfo.Loc, ",")

	weather := helpers.GetWeather(latlong[0], latlong[1], ipinfo.Timezone)

	columns := []table.Column{
		{Title: "Time", Width: 5},
		{Title: "Temp", Width: 7},
		{Title: "📈", Width: 2},
	}

	var rows []table.Row
	for i, w := range weather.Hourly.Time {
		timestamp, err := strconv.ParseInt(fmt.Sprint(w), 10, 64)
		if err != nil {
			panic(err)
		}
		t := time.Unix(timestamp, 0)

		var emoji string
		if i == 0 {
			emoji = "➖"
		} else if weather.Hourly.Temperature2m[i-1] > weather.Hourly.Temperature2m[i] {
			emoji = "🔻"
		} else if weather.Hourly.Temperature2m[i-1] < weather.Hourly.Temperature2m[i] {
			emoji = "🔺"
		} else {
			emoji = "➖"
		}

		rows = append(rows, table.Row{
			fmt.Sprint(t.UTC().Format("15:04")),
			fmt.Sprintf("%v°F", weather.Hourly.Temperature2m[i]),
			emoji,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(12),
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
	t.SetStyles(s)
	return model{
		table:    t,
		location: ipinfo.City + ", " + ipinfo.Region + ", " + ipinfo.Country,
		weather:  weather,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	view := "Weather for " + m.location + "\n\n"
	view += baseStyle.Render(m.table.View()) + "\n"
	view += "Use Up and Down arrows to scroll\n"
	view += "Press 'q' or 'ctrl+c' to quit\n"
	return view
}

func main() {
	m := initialModel()
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
