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
	filledStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#7D56F4"))

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C3C3C")).
			Background(lipgloss.Color("#3C3C3C"))

	percentStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			MarginTop(1).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(2).
			Align(lipgloss.Center)

	brightIconStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)

	fadedIconStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#3C3C3C"))

	iconContainerStyle = lipgloss.NewStyle().
				Align(lipgloss.Center).
				MarginBottom(2)

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

func setBrightness(percent int) error {
	cmd := exec.Command("brightnessctl", "set", fmt.Sprintf("%d%%", percent))
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to set brightness: %v\n", err)
		return err
	}
	return nil
}

func getBrightnessIcon(percent int) string {
	icons := []string{"󱩎", "󱩏", "󱩐", "󱩑", "󱩒", "󱩓", "󱩔", "󱩕", "󱩖", "󰛨"}

	activeIndex := percent / 10
	if activeIndex >= len(icons) {
		activeIndex = len(icons) - 1
	}

	var result []string
	for i, icon := range icons {
		if i <= activeIndex {
			result = append(result, brightIconStyle.Render(icon))
		} else {
			result = append(result, fadedIconStyle.Render(icon))
		}
	}

	return strings.Join(result, " ")
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

	icon := iconContainerStyle.Render(getBrightnessIcon(percent))

	barWidth := 50
	filled := (percent * barWidth) / 100
	empty := barWidth - filled

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
		bar,
		percentText,
	)

	help := helpStyle.Render("↑/↓: ±5%  •  1-9: 10%-90%  •  0: 100%  •  q/esc: quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		icon,
		"",
		brightnessControl,
		"",
		help,
	)

	styledContent := containerStyle.Width(70).Render(content)

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
