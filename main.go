package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor = lipgloss.Color("#00FFFF")
	emptyColor   = lipgloss.Color("#1A1A1A")
	textColor    = lipgloss.Color("#FFFFFF")
	helpColor    = lipgloss.Color("#AAAAAA")
	titleStyle   = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			MarginBottom(1)

	filledStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Background(primaryColor)

	emptyStyle = lipgloss.NewStyle().
			Foreground(emptyColor).
			Background(emptyColor)

	percentStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			MarginTop(1).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(helpColor).
			MarginTop(2).
			Align(lipgloss.Center)

	mainBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 3).
			Width(60).
			Align(lipgloss.Center)

	containerStyle = lipgloss.NewStyle().
			Padding(2, 4).
			Align(lipgloss.Center)
)

type model struct {
	brightness int
	max        int
	quitting   bool
	width      int
	height     int
}

func initialModel() model {
	brightness, max := getBrightness()
	return model{
		brightness: brightness,
		max:        max,
		quitting:   false,
		width:      80,
		height:     24,
	}
}

func getBrightness() (int, int) {
	cmd := exec.Command("brightnessctl", "g")
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to get brightness: %v\n", err)
		return 0, 100
	}
	brightness, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to parse brightness: %v\n", err)
		return 0, 100
	}

	cmd = exec.Command("brightnessctl", "m")
	output, err = cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to get max brightness: %v\n", err)
		return brightness, 100
	}
	max, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to parse max brightness: %v\n", err)
		return brightness, 100
	}

	return brightness, max
}

func getIconStyle(percent int) lipgloss.Style {
	var icon string
	var color lipgloss.Color

	switch {
	case percent <= 10:
		icon = "🌑"                        // New moon
		color = lipgloss.Color("#2C2C2C") // Dark gray
	case percent <= 25:
		icon = "🌒" // Waxing crescent
		color = lipgloss.Color("#4A4A4A")
	case percent <= 40:
		icon = "🌓" // First quarter
		color = lipgloss.Color("#666666")
	case percent <= 55:
		icon = "🌔" // Waxing gibbous
		color = lipgloss.Color("#888888")
	case percent <= 70:
		icon = "🌕"                        // Full moon
		color = lipgloss.Color("#CCCCCC") // Light gray
	case percent <= 85:
		icon = "🌞"                        // Sun with face
		color = lipgloss.Color("#FFD700") // Gold
	default:
		icon = "☀️"                       // Full sun
		color = lipgloss.Color("#FFFF00") // Yellow
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		SetString(icon).
		MarginBottom(1).
		Padding(0, 1)
}

func setBrightness(percent int) error {
	cmd := exec.Command("brightnessctl", "set", fmt.Sprintf("%d%%", percent))
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to set brightness: %v\n", err)
		return err
	}
	return nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "up", "right", "k", "l":
			if m.max == 0 {
				return m, nil
			}
			percent := (m.brightness * 100) / m.max
			newPercent := min(percent+5, 100)
			if err := setBrightness(newPercent); err == nil {
				m.brightness, m.max = getBrightness()
			}

		case "down", "left", "j", "h":
			if m.max == 0 {
				return m, nil
			}
			percent := (m.brightness * 100) / m.max
			newPercent := max(percent-5, 0)
			if err := setBrightness(newPercent); err == nil {
				m.brightness, m.max = getBrightness()
			}

		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			num, err := strconv.Atoi(msg.String())
			if err != nil {
				return m, nil
			}
			newPercent := num * 10
			if newPercent == 0 {
				newPercent = 100
			}
			if err := setBrightness(newPercent); err == nil {
				m.brightness, m.max = getBrightness()
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	percent := 0
	if m.max > 0 {
		percent = (m.brightness * 100) / m.max
	}
	barWidth := 50
	filled := (percent * barWidth) / 100
	empty := barWidth - filled

	title := titleStyle.Render("Screen Brightness")
	icon := getIconStyle(percent).String()

	header := lipgloss.JoinHorizontal(lipgloss.Center, icon, " ", title)
	bar := ""
	for i := 0; i < filled; i++ {
		bar += filledStyle.Render("█")
	}
	for i := 0; i < empty; i++ {
		bar += emptyStyle.Render("░")
	}

	percentText := percentStyle.Render(fmt.Sprintf("%d%%", percent))

	brightnessControl := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		bar,
		percentText,
	)

	mainContent := mainBoxStyle.Render(brightnessControl)

	help := helpStyle.Render("↑/↓: ±5%  •  1-9: 10%-90%  •  0: 100%  •  q/esc: quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		mainContent,
		help,
	)

	styledContent := containerStyle.Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		styledContent,
	)

}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
