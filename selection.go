package cardboard

import (
	"fmt"
	"slices"
)

// Selection tracks a cursor and a carried card in one pile.
// Moving while carrying immediately reorders the pile; Cancel restores its order.
type Selection struct {
	pile       *Pile
	selected   int
	pickedFrom int
	expected   []Card
}

func NewSelection(pile *Pile) *Selection {
	s := &Selection{pile: pile, pickedFrom: -1}
	s.clamp()
	return s
}

// Selected returns -1 when the pile is empty.
func (s *Selection) Selected() int {
	if !s.Carrying() {
		s.clamp()
	}
	return s.selected
}

func (s *Selection) Picked() int {
	if !s.Carrying() || !slices.Equal(s.expected, s.pile.Cards()) {
		return -1
	}
	return s.selected
}

func (s *Selection) Carrying() bool { return s.expected != nil }

// Move changes the selection by delta, stopping at either end of the pile.
func (s *Selection) Move(delta int) error {
	if err := s.checkCarry(); err != nil {
		return err
	}
	s.clamp()
	if s.selected < 0 {
		return nil
	}
	target := s.selected
	switch {
	case delta < -s.selected:
		target = 0
	case delta >= s.pile.Len()-s.selected:
		target = s.pile.Len() - 1
	default:
		target += delta
	}
	if s.Carrying() {
		if err := s.pile.Move(s.selected, target); err != nil {
			return err
		}
		s.expected = s.pile.Cards()
	}
	s.selected = target
	return nil
}

func (s *Selection) Pick() error {
	if err := s.checkCarry(); err != nil {
		return err
	}
	if s.Carrying() {
		return fmt.Errorf("a card is already being carried")
	}
	s.clamp()
	if s.selected < 0 {
		return fmt.Errorf("cannot pick a card from an empty pile")
	}
	s.pickedFrom = s.selected
	s.expected = s.pile.Cards()
	return nil
}

func (s *Selection) Place() error {
	if err := s.checkCarry(); err != nil {
		return err
	}
	s.clearCarry()
	return nil
}

func (s *Selection) Cancel() error {
	if err := s.checkCarry(); err != nil {
		return err
	}
	if !s.Carrying() {
		return nil
	}
	if err := s.pile.Move(s.selected, s.pickedFrom); err != nil {
		return err
	}
	s.selected = s.pickedFrom
	s.clearCarry()
	return nil
}

// External changes end a carry rather than letting an old index move another card.
func (s *Selection) checkCarry() error {
	if s.Carrying() && !slices.Equal(s.expected, s.pile.Cards()) {
		s.clearCarry()
		s.clamp()
		return fmt.Errorf("pile changed while carrying; pick the card again")
	}
	return nil
}

func (s *Selection) clearCarry() {
	s.expected = nil
	s.pickedFrom = -1
}

func (s *Selection) clamp() {
	if s.pile.Len() == 0 {
		s.selected = -1
	} else {
		s.selected = min(max(s.selected, 0), s.pile.Len()-1)
	}
}
