package cardboard

import (
	"slices"
	"testing"
)

func selectionPile(t *testing.T) *Pile {
	t.Helper()
	pile, err := NewPile(
		Card{Rank: Ace, Suit: Spades},
		Card{Rank: Two, Suit: Hearts},
		Card{Rank: Three, Suit: Clubs},
		Card{Rank: Four, Suit: Diamonds},
	)
	if err != nil {
		t.Fatal(err)
	}
	return pile
}

func TestSelectionCarriesAndPlacesCard(t *testing.T) {
	pile := selectionPile(t)
	original := pile.Cards()
	s := NewSelection(pile)
	if err := s.Pick(); err != nil {
		t.Fatal(err)
	}
	for position := 1; position <= 2; position++ {
		if err := s.Move(1); err != nil {
			t.Fatal(err)
		}
		if card, _ := pile.At(position); card != original[0] {
			t.Fatal("carried card did not move with the cursor")
		}
		if s.Selected() != position || s.Picked() != position || !s.Carrying() {
			t.Fatal("selection and highlight did not follow the carried card")
		}
	}
	if err := s.Place(); err != nil {
		t.Fatal(err)
	}
	want := []Card{original[1], original[2], original[0], original[3]}
	if !slices.Equal(pile.Cards(), want) || s.Carrying() || s.Picked() != -1 {
		t.Fatal("placing should retain the new order and clear the carry")
	}
	if err := s.Move(1); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(pile.Cards(), want) {
		t.Fatal("moving an idle cursor changed the pile")
	}
}

func TestSelectionCancelRestoresOrder(t *testing.T) {
	for _, movement := range []int{-2, 1} {
		pile := selectionPile(t)
		original := pile.Cards()
		s := NewSelection(pile)
		_ = s.Move(2)
		_ = s.Pick()
		_ = s.Move(movement)
		if err := s.Cancel(); err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(pile.Cards(), original) || s.Selected() != 2 || s.Carrying() {
			t.Fatal("cancel should restore the original order and selection")
		}
	}
}

func TestSelectionBoundariesAndEmptyPile(t *testing.T) {
	pile := selectionPile(t)
	s := NewSelection(pile)
	_ = s.Move(-100)
	if s.Selected() != 0 {
		t.Fatal("cursor moved before the first card")
	}
	_ = s.Move(100)
	if s.Selected() != pile.Len()-1 {
		t.Fatal("cursor moved past the last card")
	}
	for _, pile := range []*Pile{nil, {}} {
		s := NewSelection(pile)
		if s.Selected() != -1 || s.Picked() != -1 || s.Carrying() {
			t.Fatal("empty pile should have no selection")
		}
		if err := s.Pick(); err == nil {
			t.Fatal("picking from an empty pile should fail")
		}
		if err := s.Move(1); err != nil {
			t.Fatal(err)
		}
		if err := s.Place(); err != nil {
			t.Fatal(err)
		}
		if err := s.Cancel(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSelectionRejectsExternalChangesDuringCarry(t *testing.T) {
	for _, action := range []string{"move", "place", "cancel", "pick"} {
		t.Run(action, func(t *testing.T) {
			pile := selectionPile(t)
			s := NewSelection(pile)
			_ = s.Pick()
			_ = s.Move(1)
			_ = pile.Move(1, 3)
			changed := pile.Cards()
			if s.Picked() != -1 {
				t.Fatal("stale carried card should not be highlighted")
			}
			var err error
			switch action {
			case "move":
				err = s.Move(1)
			case "place":
				err = s.Place()
			case "cancel":
				err = s.Cancel()
			case "pick":
				err = s.Pick()
			}
			if err == nil || s.Carrying() || !slices.Equal(pile.Cards(), changed) {
				t.Fatal("external changes should end the carry without moving another card")
			}
		})
	}
}

func TestSelectionClampsAfterCardsLeavePile(t *testing.T) {
	pile := selectionPile(t)
	other, _ := NewPile()
	s := NewSelection(pile)
	_ = s.Move(3)
	_ = Transfer(pile, 3, other, 0)
	if s.Selected() != 2 {
		t.Fatal("idle selection should follow the shortened pile")
	}
	_ = s.Pick()
	_ = Transfer(pile, 2, other, 0)
	before := pile.Cards()
	if err := s.Cancel(); err == nil || !slices.Equal(pile.Cards(), before) {
		t.Fatal("removing a carried card should prevent cancellation from moving another card")
	}
}
