package cardboard

import (
	"fmt"
	"math/rand/v2"
)

// Table owns the common pile and all player hands.
// Create one with NewTable or NewTableWithDeck before use.
type Table struct {
	common     *Pile
	hands      []*Pile
	jokerCount int
	lastID     uint64
}

// NewTable starts with an ordered standard deck in a face-down common pile.
func NewTable() *Table {
	table, _ := NewTableWithDeck(StandardDeck())
	return table
}

// NewTableWithDeck copies a nonempty deck into a face-down common pile.
// Repeated faces and any number of printed jokers are allowed. Supplied IDs
// are replaced with IDs unique to this table, starting at one.
func NewTableWithDeck(cards []Card) (*Table, error) {
	if len(cards) == 0 {
		return nil, fmt.Errorf("deck must not be empty")
	}
	for index, card := range cards {
		if !card.Valid() {
			return nil, fmt.Errorf("invalid card at index %d", index)
		}
	}
	table := &Table{common: newPile(cards, true)}
	for index := range table.common.cards {
		table.lastID++
		table.common.cards[index].ID = table.lastID
		if table.common.cards[index].Suit == Joker {
			table.jokerCount++
		}
	}
	return table, nil
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

// AddJokers is a convenience for decks with at most two printed jokers.
// Use NewTableWithDeck to supply decks containing more jokers.
func (t *Table) AddJokers(count int) error {
	if t == nil {
		return fmt.Errorf("table must not be nil")
	}
	if len(t.hands) != 0 {
		return fmt.Errorf("jokers must be added before dealing")
	}
	if count < 0 {
		return fmt.Errorf("joker count must not be negative")
	}
	if count == 0 {
		return nil
	}
	if count > 2 || t.jokerCount > 2-count {
		return fmt.Errorf("AddJokers supports at most two jokers; supply larger decks with NewTableWithDeck")
	}
	if t.lastID > ^uint64(0)-uint64(count) {
		return fmt.Errorf("no card IDs left")
	}
	for i := 0; i < count; i++ {
		number := 1
		for _, card := range t.common.cards {
			if card.Suit == Joker && card.Rank == Ace {
				number = 2
				break
			}
		}
		card := jokerCard(number)
		t.lastID++
		card.ID = t.lastID
		t.jokerCount++
		t.common.cards = append(t.common.cards, card)
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
// Previously returned hand pointers become empty piles.
func (t *Table) Reset() {
	if t == nil {
		return
	}
	for _, hand := range t.hands {
		t.common.cards = append(t.common.cards, hand.cards...)
		hand.cards = nil
	}
	t.hands = nil
	t.common.faceDown = true
}
