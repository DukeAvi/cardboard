package cards

type Suit uint8

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

// Rank runs from 1 (ace) through 13 (king).
type Card struct {
	Rank int
	Suit Suit
}

var labels = [...]string{"", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

type point struct{ x, y int }

var layouts = map[int][]point{
	2:  {{8, 0}, {8, 8}},
	3:  {{8, 0}, {8, 4}, {8, 8}},
	4:  {{3, 0}, {13, 0}, {3, 8}, {13, 8}},
	5:  {{3, 0}, {13, 0}, {8, 4}, {3, 8}, {13, 8}},
	6:  {{3, 0}, {13, 0}, {3, 4}, {13, 4}, {3, 8}, {13, 8}},
	7:  {{3, 0}, {13, 0}, {8, 2}, {3, 4}, {13, 4}, {3, 8}, {13, 8}},
	8:  {{3, 0}, {13, 0}, {8, 2}, {3, 4}, {13, 4}, {8, 6}, {3, 8}, {13, 8}},
	9:  {{3, 0}, {13, 0}, {3, 2}, {13, 2}, {8, 4}, {3, 6}, {13, 6}, {3, 8}, {13, 8}},
	10: {{3, 0}, {13, 0}, {8, 1}, {3, 2}, {13, 2}, {3, 6}, {13, 6}, {8, 7}, {3, 8}, {13, 8}},
}

func (c Card) Red() bool { return c.Suit == Hearts || c.Suit == Diamonds }

// Deck returns all 52 cards, ordered by suit and then rank.
func Deck() []Card {
	deck := make([]Card, 0, 52)
	for s := Spades; s <= Clubs; s++ {
		for r := 1; r <= 13; r++ {
			deck = append(deck, Card{Rank: r, Suit: s})
		}
	}
	return deck
}
