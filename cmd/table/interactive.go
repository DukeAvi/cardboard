package main

import (
	"fmt"
	"os"
	"strings"

	cards "github.com/DukeAvi/cardboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
)

type table struct {
	pack         *cards.Pack
	renderer     *cards.Renderer
	cardSize     cards.CardSize
	choosingSize bool
	collapsed    bool
	group        int
	selected     []int
	width        int
	height       int
}

func runInteractive(pack *cards.Pack) error {
	if !term.IsTerminal(os.Stdin.Fd()) || !term.IsTerminal(os.Stdout.Fd()) {
		return fmt.Errorf("interactive mode needs a terminal; use --deal for text output")
	}
	m := table{
		pack:         pack,
		renderer:     cards.NewRenderer(cards.Vintage),
		cardSize:     cards.Medium,
		choosingSize: true,
		collapsed:    true,
		selected:     make([]int, len(pack.Hands)+1),
	}
	_, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout)).Run()
	return err
}

func (m table) Init() tea.Cmd {
	return nil
}

func (m table) hand() []cards.Card {
	if m.group == len(m.pack.Hands) {
		return m.pack.Common
	}
	return m.pack.Hands[m.group]
}

func (m table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		key := msg.String()
		if key == "q" || key == "ctrl+c" {
			return m, tea.Quit
		}
		if m.choosingSize {
			switch key {
			case "1":
				m.cardSize = cards.Small
			case "2":
				m.cardSize = cards.Medium
			case "3":
				m.cardSize = cards.Large
			case "enter":
			default:
				return m, nil
			}
			m.renderer.CardSize = m.cardSize
			m.choosingSize = false
			return m, nil
		}
		switch key {
		case "left", "h":
			if m.selected[m.group] > 0 {
				m.selected[m.group]--
			}
		case "right", "l":
			if m.selected[m.group] < len(m.hand())-1 {
				m.selected[m.group]++
			}
		case "tab":
			m.group = (m.group + 1) % len(m.selected)
		case "shift+tab":
			m.group = (m.group + len(m.selected) - 1) % len(m.selected)
		case "c":
			m.collapsed = !m.collapsed
		case "s":
			m.choosingSize = true
		}
	}
	return m, nil
}

func (m table) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	if m.choosingSize {
		return m.fit(fmt.Sprintf("Choose card size\n\n1  small   12 x 9\n2  medium  18 x 14\n3  large   26 x 20\n\nEnter: %s   q: quit", m.cardSize))
	}

	hand := m.hand()
	selected := m.selected[m.group]
	name := fmt.Sprintf("Player %d", m.group+1)
	if m.group == len(m.pack.Hands) {
		name = "Common"
	}
	mode := "expanded"
	if m.collapsed {
		mode = "collapsed"
	}
	header := fmt.Sprintf("%s (%d cards) | Common: %d | %s | %s\n", name, len(hand), len(m.pack.Common), m.cardSize, mode)
	footer := "left/right: browse  Tab: player/common  c: collapse/expand  s: size  q: quit"
	if m.width < 70 {
		footer = "left/right: browse | Tab: group\nc: fold | s: size | q: quit"
	}
	if m.width < 30 {
		footer = "←/→ Tab c s q"
	}
	textHeight := 5 + strings.Count(footer, "\n")
	if len(hand) == 0 {
		return m.fit(header + "\nNo cards here.\n\n" + footer)
	}

	w, h := m.cardSize.Dimensions()
	neededHeight := h + 2
	if m.collapsed {
		neededHeight = h + 4
	}
	if m.width < w+2 || m.height < neededHeight+textHeight {
		return m.fit(header + "\nResize the terminal or press s for smaller cards.\n\n" + footer)
	}

	start, end := 0, len(hand)
	var body string
	if m.collapsed {
		start, end = m.renderer.HandWindow(len(hand), selected, m.width)
		body = m.renderer.Hand(hand, selected, m.width)
	} else {
		columns := (m.width + 1) / (w + 3)
		rows := (m.height - textHeight + 1) / (h + 3)
		pageSize := columns * rows
		start = selected / pageSize * pageSize
		end = min(start+pageSize, len(hand))
		body = m.renderer.Grid(hand[start:end], selected-start, m.width)
	}
	status := fmt.Sprintf("Card %d/%d: %s | Showing %d-%d | %d before, %d after", selected+1, len(hand), textCards(hand[selected:selected+1]), start+1, end, start, len(hand)-end)
	return m.fit(header + status + "\n\n" + body + "\n\n" + footer)
}

func (m table) fit(view string) string {
	lines := strings.Split(view, "\n")
	lines = lines[:min(len(lines), m.height)]
	for i, line := range lines {
		lines[i] = ansi.Truncate(line, m.width, "")
	}
	return strings.Join(lines, "\n")
}
