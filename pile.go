package cardboard

import (
	"fmt"
	"math/rand/v2"
)

// Pile is an ordered group of cards. Index zero is the top of the pile.
type Pile struct {
	cards    []Card
	faceDown bool
}

func NewPile(cards ...Card) (*Pile, error) {
	for _, card := range cards {
		if !card.Valid() {
			return nil, fmt.Errorf("invalid card: %+v", card)
		}
	}
	return newPile(cards, false), nil
}

func newPile(cards []Card, faceDown bool) *Pile {
	return &Pile{cards: append([]Card(nil), cards...), faceDown: faceDown}
}

func (p *Pile) Len() int {
	if p == nil {
		return 0
	}
	return len(p.cards)
}

// Cards returns a copy, so changing it cannot corrupt the pile.
func (p *Pile) Cards() []Card {
	if p == nil {
		return nil
	}
	return append([]Card(nil), p.cards...)
}

func (p *Pile) At(index int) (Card, error) {
	if p == nil || index < 0 || index >= len(p.cards) {
		return Card{}, fmt.Errorf("card index %d is out of range", index)
	}
	return p.cards[index], nil
}

func (p *Pile) Top() (Card, error) {
	return p.At(0)
}

func (p *Pile) FaceDown() bool {
	return p != nil && p.faceDown
}

func (p *Pile) SetFaceDown(faceDown bool) {
	if p != nil {
		p.faceDown = faceDown
	}
}

func (p *Pile) Flip() {
	if p != nil {
		p.faceDown = !p.faceDown
	}
}

// Move places the card at from at the final index to.
func (p *Pile) Move(from, to int) error {
	if p == nil || from < 0 || from >= len(p.cards) || to < 0 || to >= len(p.cards) {
		return fmt.Errorf("cannot move card from %d to %d in a pile of %d", from, to, p.Len())
	}
	if from == to {
		return nil
	}
	card := p.cards[from]
	if from < to {
		copy(p.cards[from:to], p.cards[from+1:to+1])
	} else {
		copy(p.cards[to+1:from+1], p.cards[to:from])
	}
	p.cards[to] = card
	return nil
}

func (p *Pile) Shuffle() {
	if p != nil {
		rand.Shuffle(len(p.cards), func(i, j int) { p.cards[i], p.cards[j] = p.cards[j], p.cards[i] })
	}
}

func (p *Pile) ShuffleWith(random *rand.Rand) error {
	if p == nil || random == nil {
		return fmt.Errorf("pile and random source must not be nil")
	}
	random.Shuffle(len(p.cards), func(i, j int) { p.cards[i], p.cards[j] = p.cards[j], p.cards[i] })
	return nil
}

// Transfer moves one card between piles. toIndex may equal the destination length.
func Transfer(from *Pile, fromIndex int, to *Pile, toIndex int) error {
	if from == nil || to == nil {
		return fmt.Errorf("source and destination piles must not be nil")
	}
	if from == to {
		if toIndex == len(from.cards) {
			toIndex--
		}
		return from.Move(fromIndex, toIndex)
	}
	if fromIndex < 0 || fromIndex >= len(from.cards) {
		return fmt.Errorf("source index %d is out of range", fromIndex)
	}
	if toIndex < 0 || toIndex > len(to.cards) {
		return fmt.Errorf("destination index %d is out of range", toIndex)
	}
	card := from.cards[fromIndex]
	from.cards = append(from.cards[:fromIndex], from.cards[fromIndex+1:]...)
	to.cards = append(to.cards, Card{})
	copy(to.cards[toIndex+1:], to.cards[toIndex:])
	to.cards[toIndex] = card
	return nil
}
