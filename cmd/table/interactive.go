package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/DukeAvi/cardboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
)

type pileView struct {
	selected   int
	picked     int
	pickedFrom int
	picking    bool
	layout     cardboard.Layout
}

type tableModel struct {
	table    *cardboard.Table
	renderer *cardboard.Renderer
	views    []pileView
	focused  int
	width    int
	height   int
	message  string
}

func newTableModel(table *cardboard.Table) *tableModel {
	views := make([]pileView, table.PlayerCount()+1)
	for i := range views {
		views[i].picked = -1
		views[i].pickedFrom = -1
		views[i].layout = cardboard.Overlap
	}
	views[len(views)-1].layout = cardboard.Stack
	return &tableModel{
		table:    table,
		renderer: cardboard.NewRenderer(cardboard.Vintage),
		views:    views,
	}
}

func runInteractive(table *cardboard.Table) error {
	if table == nil || table.PlayerCount() == 0 {
		return fmt.Errorf("interactive mode needs at least one dealt hand")
	}
	if !term.IsTerminal(os.Stdin.Fd()) || !term.IsTerminal(os.Stdout.Fd()) {
		return fmt.Errorf("interactive mode needs a terminal; use --deal for text output")
	}
	_, err := tea.NewProgram(newTableModel(table), tea.WithAltScreen(), tea.WithInput(os.Stdin), tea.WithOutput(os.Stdout)).Run()
	return err
}

func (m *tableModel) Init() tea.Cmd { return nil }

func (m *tableModel) focusedPile() *cardboard.Pile {
	if m.focused == m.table.PlayerCount() {
		return m.table.Common()
	}
	hand, _ := m.table.Hand(m.focused)
	return hand
}

func (m *tableModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = message.Width, message.Height
	case tea.KeyMsg:
		key := message.String()
		if key == "q" || key == "ctrl+c" {
			return m, tea.Quit
		}
		view := &m.views[m.focused]
		pile := m.focusedPile()
		switch key {
		case "tab":
			m.focused = (m.focused + 1) % len(m.views)
			m.message = ""
		case "shift+tab":
			m.focused = (m.focused + len(m.views) - 1) % len(m.views)
			m.message = ""
		case "left", "h":
			if view.layout != cardboard.Stack {
				m.moveSelection(view, pile, -1)
			}
		case "right", "l":
			if view.layout != cardboard.Stack {
				m.moveSelection(view, pile, 1)
			}
		case "enter", " ":
			m.pickOrPut(view, pile)
		case "esc":
			m.cancelPick(view, pile)
		case "f":
			pile.Flip()
			m.message = "Pile flipped"
		case "c":
			if view.layout == cardboard.Grid {
				if m.focused == m.table.PlayerCount() {
					view.layout = cardboard.Stack
				} else {
					view.layout = cardboard.Overlap
				}
			} else {
				view.layout = cardboard.Grid
			}
		case "d":
			if view.picking {
				m.message = "Place or cancel the carried card before drawing"
			} else if m.focused == m.table.PlayerCount() {
				m.message = "Focus a player before drawing"
			} else if err := m.table.Draw(m.focused, 1); err != nil {
				m.message = err.Error()
			} else {
				view.selected = pile.Len() - 1
				m.message = "Drew one card"
			}
		}
	}
	return m, nil
}

func (m *tableModel) moveSelection(view *pileView, pile *cardboard.Pile, change int) {
	target := view.selected + change
	if target < 0 || target >= pile.Len() {
		return
	}
	if view.picking {
		if err := pile.Move(view.picked, target); err != nil {
			m.message = err.Error()
			return
		}
		view.picked = target
		m.message = "Carrying card; press Space or Enter to place it"
	} else {
		m.message = ""
	}
	view.selected = target
}

func (m *tableModel) pickOrPut(view *pileView, pile *cardboard.Pile) {
	if pile.Len() == 0 {
		m.message = "Pile is empty"
		return
	}
	if view.picking {
		view.picking = false
		view.picked = -1
		view.pickedFrom = -1
		m.message = "Card placed"
		return
	}
	if view.layout == cardboard.Stack && pile.Len() > 1 {
		m.message = "Switch to grid view before reordering this pile"
		return
	}
	if view.selected < 0 || view.selected >= pile.Len() {
		view.selected = 0
	}
	view.picking = true
	view.picked = view.selected
	view.pickedFrom = view.selected
	m.message = "Card picked; use the arrows to carry it"
}

func (m *tableModel) cancelPick(view *pileView, pile *cardboard.Pile) {
	if !view.picking {
		return
	}
	if err := pile.Move(view.picked, view.pickedFrom); err != nil {
		m.message = err.Error()
		return
	}
	view.selected = view.pickedFrom
	view.picking = false
	view.picked = -1
	view.pickedFrom = -1
	m.message = "Card returned to its original position"
}

func (m *tableModel) View() string {
	if m.width < 1 || m.height < 1 {
		return ""
	}
	pile := m.focusedPile()
	view := m.views[m.focused]
	name := fmt.Sprintf("Player %d", m.focused+1)
	if m.focused == m.table.PlayerCount() {
		name = "Common"
	}
	face := "face up"
	if pile.FaceDown() {
		face = "face down"
	}
	layout := map[cardboard.Layout]string{cardboard.Overlap: "overlap", cardboard.Grid: "grid", cardboard.Stack: "stack"}[view.layout]
	header := fmt.Sprintf("%s | %d cards | %s | %s | common: %d", name, pile.Len(), face, layout, m.table.Common().Len())
	status := m.message
	if status == "" && pile.Len() > 0 {
		status = fmt.Sprintf("Card %d/%d", view.selected+1, pile.Len())
		if !pile.FaceDown() {
			card, _ := pile.At(view.selected)
			status += ": " + card.Short()
		}
	}
	footer := "arrows: select/carry  Space/Enter: pick/place  Esc: return  Tab: pile  f: flip  c: layout  d: draw  q: quit"
	if m.width < ansi.StringWidth(footer) {
		footer = "arrows select/carry | Space pick/place | Tab pile | f flip\nd draw | Esc return | c layout | q quit"
	}
	reserved := 3 + strings.Count(footer, "\n") + 1
	bodyHeight := cardboard.PixelHeight + 2
	if view.layout == cardboard.Overlap {
		bodyHeight += 2
	}
	requiredHeight := bodyHeight + reserved
	if m.width < cardboard.PixelWidth+2 || m.height < requiredHeight {
		resize := fmt.Sprintf("Resize the terminal to at least 28 columns and %d rows.", requiredHeight)
		return m.fit(header + "\n" + status + "\n\n" + resize + "\n\n" + footer)
	}
	body := "No cards"
	if pile.Len() > 0 {
		body = m.renderer.Pile(pile, cardboard.View{
			Width:    m.width,
			Height:   m.height - reserved,
			Layout:   view.layout,
			Selected: view.selected,
			Picking:  view.picking,
			Picked:   view.picked,
		})
	}
	return m.fit(header + "\n" + status + "\n\n" + body + "\n" + footer)
}

func (m *tableModel) fit(text string) string {
	lines := strings.Split(text, "\n")
	lines = lines[:min(len(lines), m.height)]
	for i, line := range lines {
		lines[i] = ansi.Truncate(line, m.width, "")
	}
	return strings.Join(lines, "\n")
}
