package cardboard

import (
	"fmt"
	"math/rand/v2"
)

// Table owns the common pile and all player hands.
type Table struct {
	common     *Pile
	hands      []*Pile
	jokerCount int
}

// NewTable starts with an ordered standard deck in a face-down common pile.
func NewTable() *Table {
	return &Table{common: newPile(StandardDeck(), true)}
}

func (t *Table) Common() *Pile {
	if t == nil {
		return nil
	}
	return t.common
}

func (t *Table) PlayerCount() int {
	if t == nil {
		return 0
	}
	return len(t.hands)
}

func (t *Table) Hand(player int) (*Pile, error) {
	if t == nil || player < 0 || player >= len(t.hands) {
		return nil, fmt.Errorf("player %d is out of range", player)
	}
	return t.hands[player], nil
}

func (t *Table) Hands() []*Pile {
	if t == nil {
		return nil
	}
	return append([]*Pile(nil), t.hands...)
}

func (t *Table) TotalCards() int {
	if t == nil {
		return 0
	}
	total := t.common.Len()
	for _, hand := range t.hands {
		total += hand.Len()
	}
	return total
}

// AddJokers adds up to two distinguishable jokers before cards are dealt.
func (t *Table) AddJokers(count int) error {
	if t == nil {
		return fmt.Errorf("table must not be nil")
	}
	if len(t.hands) != 0 {
		return fmt.Errorf("jokers must be added before dealing")
	}
	if count < 0 || t.jokerCount+count > 2 {
		return fmt.Errorf("a table can contain at most two jokers")
	}
	for i := 0; i < count; i++ {
		t.jokerCount++
		t.common.cards = append(t.common.cards, jokerCard(t.jokerCount))
	}
	return nil
}

func (t *Table) Shuffle() {
	if t != nil {
		t.common.Shuffle()
	}
}

func (t *Table) ShuffleWith(random *rand.Rand) error {
	if t == nil {
		return fmt.Errorf("table must not be nil")
	}
	return t.common.ShuffleWith(random)
}

// Deal gives cards to each player one at a time around the table.
func (t *Table) Deal(players, cardsEach int) error {
	if t == nil {
		return fmt.Errorf("table must not be nil")
	}
	if len(t.hands) != 0 {
		return fmt.Errorf("cards have already been dealt; reset the table first")
	}
	if players < 1 || cardsEach < 1 {
		return fmt.Errorf("players and cards per player must both be positive")
	}
	if players > t.common.Len()/cardsEach {
		return fmt.Errorf("not enough cards: have %d for %d players with %d cards each", t.common.Len(), players, cardsEach)
	}

	hands := make([]*Pile, players)
	for player := range hands {
		hands[player] = newPile(nil, false)
	}
	next := 0
	for round := 0; round < cardsEach; round++ {
		for player := range hands {
			hands[player].cards = append(hands[player].cards, t.common.cards[next])
			next++
		}
	}
	t.common.cards = t.common.cards[next:]
	t.hands = hands
	return nil
}

// Draw moves count cards from the common pile into a player's hand.
func (t *Table) Draw(player, count int) error {
	hand, err := t.Hand(player)
	if err != nil {
		return err
	}
	if count < 1 {
		return fmt.Errorf("draw count must be positive")
	}
	if count > t.common.Len() {
		return fmt.Errorf("cannot draw %d cards from a common pile of %d", count, t.common.Len())
	}
	hand.cards = append(hand.cards, t.common.cards[:count]...)
	t.common.cards = t.common.cards[count:]
	return nil
}

// Reset collects every hand into the common pile without shuffling it.
func (t *Table) Reset() {
	if t == nil {
		return
	}
	for _, hand := range t.hands {
		t.common.cards = append(t.common.cards, hand.cards...)
	}
	t.hands = nil
	t.common.faceDown = true
}
