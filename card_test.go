package cardboard

import "testing"

func TestStandardDeck(t *testing.T) {
	deck := StandardDeck()
	if len(deck) != 52 {
		t.Fatalf("deck has %d cards, want 52", len(deck))
	}
	seen := make(map[Card]bool)
	for _, card := range deck {
		if card.ID != 0 {
			t.Fatal("StandardDeck should return unassigned card faces")
		}
		if !card.Valid() {
			t.Fatalf("invalid card in deck: %+v", card)
		}
		if seen[card] {
			t.Fatalf("duplicate card: %v", card)
		}
		seen[card] = true
	}
	if deck[0] != (Card{Rank: Ace, Suit: Spades}) || deck[51] != (Card{Rank: King, Suit: Clubs}) {
		t.Fatalf("unexpected deck order: first=%v last=%v", deck[0], deck[51])
	}
}

func TestCardIdentityDoesNotChangeFace(t *testing.T) {
	one := Card{Rank: Queen, Suit: Hearts, ID: 1}
	two := Card{Rank: Queen, Suit: Hearts, ID: 2}
	if one == two || !one.SameFace(two) {
		t.Fatal("matching faces should still have different physical identities")
	}
	if one.Short() != two.Short() || one.String() != two.String() || one.Red() != two.Red() || one.Valid() != two.Valid() {
		t.Fatal("card identity should not affect face behavior")
	}
	if one.SameFace(Card{Rank: Queen, Suit: Clubs}) {
		t.Fatal("different suits should not match")
	}
}

func TestCardFormattingAndValidation(t *testing.T) {
	tests := []struct {
		card  Card
		valid bool
		short string
		long  string
		isRed bool
	}{
		{Card{Rank: Ace, Suit: Spades}, true, "AS", "Ace of Spades", false},
		{Card{Rank: Ten, Suit: Hearts}, true, "10H", "Ten of Hearts", true},
		{Card{Rank: King, Suit: Diamonds}, true, "KD", "King of Diamonds", true},
		{Card{Rank: Ace, Suit: Joker}, true, "Joker-1", "Black Joker", false},
		{Card{Rank: Two, Suit: Joker}, true, "Joker-2", "Red Joker", true},
		{Card{}, false, "?", "Invalid card", false},
		{Card{Rank: King, Suit: Joker}, false, "?", "Invalid card", false},
	}
	for _, test := range tests {
		if test.card.Valid() != test.valid || test.card.Short() != test.short || test.card.String() != test.long || test.card.Red() != test.isRed {
			t.Fatalf("unexpected card behavior for %+v", test.card)
		}
	}
}
