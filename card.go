package cardboard

import "fmt"

type Suit uint8

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
	Joker
)

type Rank uint8

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

type Card struct {
	Rank Rank
	Suit Suit
}

var rankLabels = [...]string{"", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
var rankNames = [...]string{"", "Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}
var suitNames = [...]string{"Spades", "Hearts", "Diamonds", "Clubs"}
var suitLabels = [...]string{"S", "H", "D", "C"}

func (c Card) Valid() bool {
	if c.Suit == Joker {
		return c.Rank == Ace || c.Rank == Two
	}
	return c.Suit <= Clubs && c.Rank >= Ace && c.Rank <= King
}

func (c Card) Red() bool {
	return c.Suit == Hearts || c.Suit == Diamonds || (c.Suit == Joker && c.Rank == Two)
}

func (c Card) Short() string {
	if !c.Valid() {
		return "?"
	}
	if c.Suit == Joker {
		return fmt.Sprintf("Joker-%d", c.Rank)
	}
	return rankLabels[c.Rank] + suitLabels[c.Suit]
}

func (c Card) String() string {
	if !c.Valid() {
		return "Invalid card"
	}
	if c.Suit == Joker {
		if c.Red() {
			return "Red Joker"
		}
		return "Black Joker"
	}
	return rankNames[c.Rank] + " of " + suitNames[c.Suit]
}

// StandardDeck returns a new ordered deck of 52 cards.
func StandardDeck() []Card {
	deck := make([]Card, 0, 52)
	for suit := Spades; suit <= Clubs; suit++ {
		for rank := Ace; rank <= King; rank++ {
			deck = append(deck, Card{Rank: rank, Suit: suit})
		}
	}
	return deck
}

func jokerCard(number int) Card {
	return Card{Rank: Rank(number), Suit: Joker}
}
