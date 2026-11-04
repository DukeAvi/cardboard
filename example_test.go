package cardboard_test

import (
	"fmt"

	"github.com/DukeAvi/cardboard"
)

func ExampleNewTableWithDeck() {
	deck := append(cardboard.StandardDeck(), cardboard.StandardDeck()...)
	table, err := cardboard.NewTableWithDeck(deck)
	if err != nil {
		panic(err)
	}
	first, _ := table.Common().At(0)
	second, _ := table.Common().At(52)
	fmt.Println("same face:", first.SameFace(second), "same card:", first.ID == second.ID)

	if err := table.Deal(6, 13); err != nil {
		panic(err)
	}
	fmt.Printf("players: %d, common: %d\n", table.PlayerCount(), table.Common().Len())

	hand, _ := table.Hand(0)
	discard, _ := cardboard.NewPile()
	if err := cardboard.Transfer(hand, 0, discard, 0); err != nil {
		panic(err)
	}
	// Collect caller-owned piles before resetting the table.
	if err := cardboard.Transfer(discard, 0, table.Common(), 0); err != nil {
		panic(err)
	}
	table.Reset()
	fmt.Printf("cards after reset: %d, old hand: %d\n", table.TotalCards(), hand.Len())

	// Output:
	// same face: true same card: false
	// players: 6, common: 26
	// cards after reset: 104, old hand: 0
}
