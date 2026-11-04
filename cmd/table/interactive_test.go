package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/DukeAvi/cardboard"
	tea "github.com/charmbracelet/bubbletea"
)

func key(value rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{value}}
}

func TestPickAndCarry(t *testing.T) {
	table := cardboard.NewTable()
	if err := table.Deal(2, 5); err != nil {
		t.Fatal(err)
	}
	hand, _ := table.Hand(0)
	picked, _ := hand.At(0)
	m := newTableModel(table)
	m.width, m.height = 100, 40
	m.Update(key(' '))
	m.Update(tea.KeyMsg{Type: tea.KeyRight})
	if got, _ := hand.At(1); got != picked {
		t.Fatalf("picked card did not follow the first arrow: %v", hand.Cards())
	}
	m.Update(tea.KeyMsg{Type: tea.KeyRight})
	if got, _ := hand.At(2); got != picked {
		t.Fatalf("picked card did not follow the second arrow: %v", hand.Cards())
	}
	if m.views[0].selected != 2 || m.views[0].picked != 2 {
		t.Fatal("selection and picked position should follow the carried card")
	}
	m.Update(key(' '))
	if got, _ := hand.At(2); got != picked {
		t.Fatalf("picked card moved to %v, want index 2", hand.Cards())
	}
	if m.views[0].picking || m.views[0].picked != -1 {
		t.Fatal("placing should clear picked state")
	}
}

func TestCancelCarryRestoresOriginalOrder(t *testing.T) {
	table := cardboard.NewTable()
	_ = table.Deal(2, 5)
	hand, _ := table.Hand(0)
	original := hand.Cards()
	m := newTableModel(table)
	m.Update(key(' '))
	m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if !reflect.DeepEqual(hand.Cards(), original) {
		t.Fatalf("cancelled carry changed the hand: got %v, want %v", hand.Cards(), original)
	}
	if m.views[0].picking || m.views[0].selected != 0 {
		t.Fatal("cancel should clear picked state and restore the original selection")
	}
}

func TestInteractiveUsesTableState(t *testing.T) {
	table := cardboard.NewTable()
	_ = table.Deal(2, 5)
	m := newTableModel(table)
	m.width, m.height = 100, 40
	hand, _ := table.Hand(0)
	m.Update(key('f'))
	if !hand.FaceDown() {
		t.Fatal("flip key did not update pile state")
	}
	before := table.Common().Len()
	m.Update(key('d'))
	if table.Common().Len() != before-1 || hand.Len() != 6 {
		t.Fatal("draw key did not move a card")
	}
	m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if m.focused != 1 {
		t.Fatal("Tab did not change the focused pile")
	}
}

func TestInteractiveMinimumSizeMatchesOverlapLayout(t *testing.T) {
	table := cardboard.NewTable()
	_ = table.Deal(2, 5)
	m := newTableModel(table)
	m.width, m.height = 28, 28
	if output := m.View(); !strings.Contains(output, "Resize the terminal") {
		t.Fatalf("28-row view should request a resize: %q", output)
	}
	m.height = 29
	if output := m.View(); strings.Contains(output, "Resize the terminal") {
		t.Fatal("29-row view should fit an overlapping hand")
	}
}
