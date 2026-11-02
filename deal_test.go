package cards

import (
	"reflect"
	"testing"
)

func TestDeal(t *testing.T) {
	tests := []struct {
		name       string
		players    int
		cardsEach  int
		jokers     bool
		commonLeft int
	}{
		{"two rummy hands", 2, 13, false, 26},
		{"whole deck", 4, 13, false, 0},
		{"one card each", 52, 1, false, 0},
		{"with jokers", 4, 13, true, 2},
		{"deal all including jokers", 2, 27, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pack := NewPack()
			expected := Deck()
			if tt.jokers {
				pack.Jokers()
				expected = append(expected, Card{Rank: 1, Suit: Joker}, Card{Rank: 2, Suit: Joker})
			}
			before := append([]Card(nil), pack.Common...)

			if err := pack.Deal(tt.players, tt.cardsEach); err != nil {
				t.Fatal(err)
			}
			if len(pack.Hands) != tt.players || len(pack.Common) != tt.commonLeft {
				t.Fatalf("got %d hands and %d common cards", len(pack.Hands), len(pack.Common))
			}

			seen := make(map[Card]int)
			for player, hand := range pack.Hands {
				if len(hand) != tt.cardsEach {
					t.Fatalf("player %d has %d cards, want %d", player, len(hand), tt.cardsEach)
				}
				for round, card := range hand {
					if card != before[round*tt.players+player] {
						t.Fatal("cards were not dealt around the table in order")
					}
					seen[card]++
				}
			}
			for i, card := range pack.Common {
				if card != before[tt.players*tt.cardsEach+i] {
					t.Fatal("common pile order changed during the deal")
				}
				seen[card]++
			}
			if len(seen) != len(expected) {
				t.Fatalf("got %d unique cards, want %d", len(seen), len(expected))
			}
			for _, card := range expected {
				if seen[card] != 1 {
					t.Fatalf("card %+v appears %d times, want once", card, seen[card])
				}
			}
		})
	}
}

func TestDealInvalidInput(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	for _, input := range [][2]int{{0, 13}, {-1, 13}, {2, 0}, {2, -1}, {5, 13}, {maxInt, 2}, {2, maxInt}} {
		pack := NewPack()
		before := append([]Card(nil), pack.Common...)
		if err := pack.Deal(input[0], input[1]); err == nil {
			t.Fatalf("Deal(%d, %d) should fail", input[0], input[1])
		}
		if pack.Hands != nil || !reflect.DeepEqual(pack.Common, before) {
			t.Fatal("failed deal changed the pack")
		}
		if err := pack.Deal(2, 13); err != nil {
			t.Fatalf("valid deal after an invalid request failed: %v", err)
		}
	}
}

func TestDealCannotReplaceExistingHands(t *testing.T) {
	pack := NewPack()
	if err := pack.Deal(2, 13); err != nil {
		t.Fatal(err)
	}
	hands := make([][]Card, len(pack.Hands))
	for i, hand := range pack.Hands {
		hands[i] = append([]Card(nil), hand...)
	}
	common := append([]Card(nil), pack.Common...)
	if err := pack.Deal(2, 7); err == nil {
		t.Fatal("dealing twice should fail")
	}
	if !reflect.DeepEqual(pack.Hands, hands) || !reflect.DeepEqual(pack.Common, common) {
		t.Fatal("second deal changed existing cards")
	}

	pack.Hands[0][0] = Card{}
	if !reflect.DeepEqual(pack.Hands[1], hands[1]) || !reflect.DeepEqual(pack.Common, common) {
		t.Fatal("changing one hand changed another group")
	}
}

func TestJokersAddedOnlyOnce(t *testing.T) {
	pack := NewPack().Jokers()
	before := append([]Card(nil), pack.Common...)
	if pack.Jokers() != pack || len(pack.Common) != 54 {
		t.Fatal("Jokers should add exactly two cards and return the same pack")
	}
	if !reflect.DeepEqual(pack.Common, before) {
		t.Fatal("adding jokers again should leave the pack alone")
	}
}
