package cards

import (
	"fmt"
	"math/rand/v2"
)

// Pack holds each player's hand and the shared cards left to draw.
type Pack struct {
	Hands  [][]Card
	Common []Card

	jokers bool
	dealt  bool
}

// NewPack starts with a shuffled standard deck of 52 cards.
func NewPack() *Pack {
	p := &Pack{Common: Deck()}
	p.shuffle()
	return p
}

// Jokers adds the two printed jokers once and shuffles the common pile.
func (p *Pack) Jokers() *Pack {
	if p.jokers {
		return p
	}
	p.Common = append(p.Common,
		Card{Rank: 1, Suit: Joker},
		Card{Rank: 2, Suit: Joker},
	)
	p.jokers = true
	p.shuffle()
	return p
}

func (p *Pack) shuffle() {
	rand.Shuffle(len(p.Common), func(i, j int) {
		p.Common[i], p.Common[j] = p.Common[j], p.Common[i]
	})
}

// Deal gives each player cardsEach cards, leaving the rest in Common.
func (p *Pack) Deal(players, cardsEach int) error {
	if p.dealt {
		return fmt.Errorf("this pack has already been dealt")
	}
	if players < 1 || cardsEach < 1 {
		return fmt.Errorf("players and cards per player must both be positive")
	}
	// Divide first so very large inputs cannot overflow a multiplication.
	if players > len(p.Common)/cardsEach {
		return fmt.Errorf("not enough cards: have %d for %d players with %d cards each", len(p.Common), players, cardsEach)
	}

	hands := make([][]Card, players)
	for player := range hands {
		hands[player] = make([]Card, 0, cardsEach)
	}

	// Deal one card at a time around the table.
	next := 0
	for round := 0; round < cardsEach; round++ {
		for player := range hands {
			hands[player] = append(hands[player], p.Common[next])
			next++
		}
	}

	p.Hands = hands
	p.Common = p.Common[next:]
	p.dealt = true
	return nil
}
