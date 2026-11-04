package cardboard

import (
	"math/rand/v2"
	"reflect"
	"testing"
)

func testCards() []Card {
	return []Card{{Rank: Ace, Suit: Spades}, {Rank: Two, Suit: Hearts}, {Rank: Three, Suit: Clubs}, {Rank: Four, Suit: Diamonds}}
}

func TestPileUsesDefensiveCopies(t *testing.T) {
	cards := testCards()
	pile, err := NewPile(cards...)
	if err != nil {
		t.Fatal(err)
	}
	cards[0] = Card{}
	copy := pile.Cards()
	copy[1] = Card{}
	if got, _ := pile.At(0); got != (Card{Rank: Ace, Suit: Spades}) {
		t.Fatal("changing constructor input changed the pile")
	}
	if got, _ := pile.At(1); got != (Card{Rank: Two, Suit: Hearts}) {
		t.Fatal("changing Cards output changed the pile")
	}
	if _, err := NewPile(Card{}); err == nil {
		t.Fatal("invalid card should be rejected")
	}
}

func TestPileMove(t *testing.T) {
	pile, _ := NewPile(testCards()...)
	if err := pile.Move(0, 3); err != nil {
		t.Fatal(err)
	}
	want := []Card{{Rank: Two, Suit: Hearts}, {Rank: Three, Suit: Clubs}, {Rank: Four, Suit: Diamonds}, {Rank: Ace, Suit: Spades}}
	if !reflect.DeepEqual(pile.Cards(), want) {
		t.Fatalf("move right: got %v, want %v", pile.Cards(), want)
	}
	if err := pile.Move(3, 1); err != nil {
		t.Fatal(err)
	}
	want = []Card{{Rank: Two, Suit: Hearts}, {Rank: Ace, Suit: Spades}, {Rank: Three, Suit: Clubs}, {Rank: Four, Suit: Diamonds}}
	if !reflect.DeepEqual(pile.Cards(), want) {
		t.Fatalf("move left: got %v, want %v", pile.Cards(), want)
	}
	before := pile.Cards()
	if err := pile.Move(-1, 2); err == nil || !reflect.DeepEqual(pile.Cards(), before) {
		t.Fatal("failed move changed the pile")
	}
}

func TestTransferIsAtomic(t *testing.T) {
	from, _ := NewPile(testCards()...)
	to, _ := NewPile(Card{Rank: Five, Suit: Clubs})
	if err := Transfer(from, 1, to, 1); err != nil {
		t.Fatal(err)
	}
	if from.Len() != 3 || to.Len() != 2 || to.Cards()[1] != (Card{Rank: Two, Suit: Hearts}) {
		t.Fatalf("unexpected transfer: from=%v to=%v", from.Cards(), to.Cards())
	}
	fromBefore, toBefore := from.Cards(), to.Cards()
	if err := Transfer(from, 99, to, 0); err == nil {
		t.Fatal("invalid transfer should fail")
	}
	if !reflect.DeepEqual(from.Cards(), fromBefore) || !reflect.DeepEqual(to.Cards(), toBefore) {
		t.Fatal("failed transfer changed a pile")
	}
	if err := Transfer(to, 0, to, to.Len()); err != nil || to.Cards()[to.Len()-1] != (Card{Rank: Five, Suit: Clubs}) {
		t.Fatal("same-pile transfer to the end failed")
	}
}

func TestPileFacingAndShuffle(t *testing.T) {
	one, _ := NewPile(StandardDeck()...)
	two, _ := NewPile(StandardDeck()...)
	one.SetFaceDown(true)
	if !one.FaceDown() {
		t.Fatal("pile should be face down")
	}
	one.Flip()
	if one.FaceDown() {
		t.Fatal("pile should be face up")
	}
	randomOne := rand.New(rand.NewPCG(10, 20))
	randomTwo := rand.New(rand.NewPCG(10, 20))
	if err := one.ShuffleWith(randomOne); err != nil {
		t.Fatal(err)
	}
	if err := two.ShuffleWith(randomTwo); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(one.Cards(), two.Cards()) || reflect.DeepEqual(one.Cards(), StandardDeck()) {
		t.Fatal("seeded shuffle should be repeatable and change the order")
	}
}
