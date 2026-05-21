package states

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rdlucas2/doomscroller/internal/ui"
)

type IntroState struct {
	tick     int
	phase    int // 0=logo 1=subtitle 2=loading 3=done
	progress int
	done     bool
}

// IntroTickMsg is exported so the engine can pattern-match on it.
type IntroTickMsg struct{}

func IntroTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return IntroTickMsg{}
	})
}

func (s *IntroState) Init() tea.Cmd {
	return IntroTick()
}

func (s *IntroState) Done() bool {
	return s.done
}

func (s *IntroState) Update(msg tea.Msg) (done bool) {
	switch msg.(type) {
	case IntroTickMsg:
		s.tick++
		switch {
		case s.tick < 10:
			s.phase = 0
		case s.tick < 20:
			s.phase = 1
		case s.tick < 50:
			s.phase = 2
			s.progress = (s.tick - 20) * 100 / 30
		default:
			s.phase = 3
			s.done = true
		}
	}
	return s.done
}

func (s *IntroState) NextCmd() tea.Cmd {
	if s.done {
		return nil
	}
	return IntroTick()
}

func (s *IntroState) View(width, height int) string {
	lines := make([]string, 0, height)

	logo := []string{
		`  ██████╗  ██████╗  ██████╗ ███╗   ███╗`,
		`  ██╔══██╗██╔═══██╗██╔═══██╗████╗ ████║`,
		`  ██║  ██║██║   ██║██║   ██║██╔████╔██║`,
		`  ██║  ██║██║   ██║██║   ██║██║╚██╔╝██║`,
		`  ██████╔╝╚██████╔╝╚██████╔╝██║ ╚═╝ ██║`,
		`  ╚═════╝  ╚═════╝  ╚═════╝ ╚═╝     ╚═╝`,
		``,
		`  ███████╗ ██████╗ ██████╗  ██████╗ ██╗     ██╗     ███████╗██████╗ `,
		`  ██╔════╝██╔════╝██╔══██╗██╔═══██╗██║     ██║     ██╔════╝██╔══██╗`,
		`  ███████╗██║     ██████╔╝██║   ██║██║     ██║     █████╗  ██████╔╝`,
		`  ╚════██║██║     ██╔══██╗██║   ██║██║     ██║     ██╔══╝  ██╔══██╗`,
		`  ███████║╚██████╗██║  ██║╚██████╔╝███████╗███████╗███████╗██║  ██║`,
		`  ╚══════╝ ╚═════╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝`,
	}

	topPad := (height - 20) / 2
	for i := 0; i < topPad; i++ {
		lines = append(lines, "")
	}

	if s.phase >= 0 {
		for _, l := range logo {
			lines = append(lines, ui.TitleStyle.Render(l))
		}
	}

	lines = append(lines, "")

	if s.phase >= 1 {
		sub := ui.SubtitleStyle.Render("~ A JRPG Roguelike ~")
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(sub))
	}

	lines = append(lines, "")
	lines = append(lines, "")

	if s.phase >= 2 {
		barWidth := 30
		filled := s.progress * barWidth / 100
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
		barStr := fmt.Sprintf("[%s] %d%%", bar, s.progress)
		loading := fmt.Sprintf("Loading...  %s", ui.XPStyle.Render(barStr))
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(loading))
	}

	if s.phase >= 3 {
		lines = append(lines, "")
		press := ui.SuccessStyle.Render("Press any key to continue...")
		lines = append(lines, lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(press))
	}

	return strings.Join(lines, "\n")
}
