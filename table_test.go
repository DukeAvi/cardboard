package cardboard

import (
	"math/rand/v2"
	"reflect"
	"testing"
)

func TestTableDealAndDraw(t *testing.T) {
	table := NewTable()
	if err := table.AddJokers(2); err != nil {
		t.Fatal(err)
	}
	if err := table.Deal(4, 13); err != nil {
		t.Fatal(err)
	}
	if table.PlayerCount() != 4 || table.Common().Len() != 2 || table.TotalCards() != 54 {
		t.Fatalf("bad deal: players=%d common=%d total=%d", table.PlayerCount(), table.Common().Len(), table.TotalCards())
	}
	for player := 0; player < 4; player++ {
		hand, _ := table.Hand(player)
		if hand.Len() != 13 || hand.FaceDown() {
			t.Fatalf("bad hand for player %d", player)
		}
	}
	if err := table.Draw(0, 2); err != nil {
		t.Fatal(err)
	}
	hand, _ := table.Hand(0)
	if hand.Len() != 15 || table.Common().Len() != 0 || table.TotalCards() != 54 {
		t.Fatal("draw lost or duplicated cards")
	}
}

func TestDealOrder(t *testing.T) {
	table := NewTable()
	deck := table.Common().Cards()
	if err := table.Deal(2, 3); err != nil {
		t.Fatal(err)
	}
	first, _ := table.Hand(0)
	second, _ := table.Hand(1)
	wantFirst := []Card{deck[0], deck[2], deck[4]}
	wantSecond := []Card{deck[1], deck[3], deck[5]}
	if !reflect.DeepEqual(first.Cards(), wantFirst) || !reflect.DeepEqual(second.Cards(), wantSecond) {
		t.Fatalf("cards were not dealt around the table: %v / %v", first.Cards(), second.Cards())
	}
}

func TestCustomDeckOwnershipAndValidation(t *testing.T) {
	for _, deck := range [][]Card{nil, {}, {{Rank: Ace, Suit: Spades}, {}}} {
		if table, err := NewTableWithDeck(deck); err == nil || table != nil {
			t.Fatal("empty or invalid deck should be rejected")
		}
	}
	deck := append(StandardDeck(), StandardDeck()...)
	for i := 0; i < 4; i++ {
		deck = append(deck, Card{Rank: Ace, Suit: Joker, ID: 999})
	}
	table, err := NewTableWithDeck(deck)
	if err != nil {
		t.Fatal(err)
	}
	if table.TotalCards() != 108 || !table.Common().FaceDown() {
		t.Fatal("custom deck should be copied into the face-down common pile")
	}
	cards := table.Common().Cards()
	if !cards[0].SameFace(cards[52]) || cards[0].ID == cards[52].ID {
		t.Fatal("copies of a face need distinct IDs")
	}
	seen := make(map[uint64]bool)
	for _, card := range cards {
		if card.ID == 0 || seen[card.ID] {
			t.Fatalf("missing or repeated ID: %d", card.ID)
		}
		seen[card.ID] = true
	}
	if deck[0].ID != 0 || deck[104].ID != 999 {
		t.Fatal("constructor changed the supplied deck")
	}
	deck[0] = Card{}
	if got, _ := table.Common().Top(); !got.Valid() {
		t.Fatal("changing the supplied deck changed table state")
	}
	if err := table.AddJokers(0); err != nil {
		t.Fatal(err)
	}
	if err := table.AddJokers(1); err == nil || table.TotalCards() != 108 {
		t.Fatal("joker convenience limit should count supplied jokers")
	}
}

func TestCardIdentitiesSurvivePileOperations(t *testing.T) {
	table, err := NewTableWithDeck(append(StandardDeck(), StandardDeck()...))
	if err != nil {
		t.Fatal(err)
	}
	initial := make(map[uint64]Card)
	for _, card := range table.Common().Cards() {
		initial[card.ID] = card
	}
	check := func() {
		t.Helper()
		seen := make(map[uint64]bool)
		piles := append(table.Hands(), table.Common())
		for _, pile := range piles {
			for _, card := range pile.Cards() {
				if card.ID == 0 || seen[card.ID] || initial[card.ID] != card {
					t.Fatalf("card was lost, duplicated or changed: %+v", card)
				}
				seen[card.ID] = true
			}
		}
		if len(seen) != len(initial) {
			t.Fatalf("found %d identities, want %d", len(seen), len(initial))
		}
	}
	if err := table.ShuffleWith(rand.New(rand.NewPCG(1, 2))); err != nil {
		t.Fatal(err)
	}
	if err := table.Deal(4, 13); err != nil {
		t.Fatal(err)
	}
	check()
	if err := table.Draw(0, 1); err != nil {
		t.Fatal(err)
	}
	one, _ := table.Hand(0)
	two, _ := table.Hand(1)
	if err := Transfer(one, 0, two, two.Len()); err != nil {
		t.Fatal(err)
	}
	if err := two.Move(0, two.Len()-1); err != nil {
		t.Fatal(err)
	}
	check()
	table.Reset()
	check()
	if one.Len() != 0 || two.Len() != 0 {
		t.Fatal("reset left cards in old hand pointers")
	}
	if err := Transfer(one, 0, table.Common(), 0); err == nil {
		t.Fatal("an old hand should not be able to duplicate a collected card")
	}
	if err := table.Deal(2, 13); err != nil {
		t.Fatal(err)
	}
	check()
}

func TestAddedJokersGetNewIdentities(t *testing.T) {
	table, err := NewTableWithDeck([]Card{{Rank: Two, Suit: Joker}})
	if err != nil {
		t.Fatal(err)
	}
	if err := table.AddJokers(1); err != nil {
		t.Fatal(err)
	}
	cards := table.Common().Cards()
	if cards[0].ID == cards[1].ID || cards[1].ID == 0 || cards[1].Rank != Ace {
		t.Fatal("added joker should have a new ID and the missing face")
	}
	standard := NewTable()
	before := standard.Common().Cards()
	if err := standard.AddJokers(2); err != nil {
		t.Fatal(err)
	}
	cards = standard.Common().Cards()
	if cards[52].ID <= before[51].ID || cards[53].ID <= cards[52].ID {
		t.Fatal("added joker IDs should follow existing IDs")
	}
}

func TestJokerFailuresLeaveIdentitiesAlone(t *testing.T) {
	table := NewTable()
	before := table.Common().Cards()
	for _, count := range []int{-1, int(^uint(0) >> 1)} {
		if err := table.AddJokers(count); err == nil || !reflect.DeepEqual(table.Common().Cards(), before) {
			t.Fatal("invalid joker count changed the deck")
		}
	}
	table.lastID = ^uint64(0)
	if err := table.AddJokers(1); err == nil || !reflect.DeepEqual(table.Common().Cards(), before) {
		t.Fatal("exhausted ID range should fail without changing the deck")
	}
}

func TestFailedTableOperationsLeaveStateAlone(t *testing.T) {
	for _, request := range [][2]int{{0, 13}, {2, 0}, {5, 13}} {
		table := NewTable()
		before := table.Common().Cards()
		if err := table.Deal(request[0], request[1]); err == nil {
			t.Fatalf("Deal(%d, %d) should fail", request[0], request[1])
		}
		if table.PlayerCount() != 0 || !reflect.DeepEqual(table.Common().Cards(), before) {
			t.Fatal("failed deal changed the table")
		}
	}
	table := NewTable()
	_ = table.Deal(2, 13)
	before := table.Common().Cards()
	if err := table.Draw(0, 100); err == nil || !reflect.DeepEqual(table.Common().Cards(), before) {
		t.Fatal("failed draw changed the common pile")
	}
	if err := table.AddJokers(1); err == nil {
		t.Fatal("adding jokers after dealing should fail")
	}
}

func TestResetAndHandsCopy(t *testing.T) {
	table := NewTable()
	_ = table.AddJokers(2)
	_ = table.Deal(2, 13)
	hands := table.Hands()
	hands[0] = nil
	if hand, _ := table.Hand(0); hand == nil {
		t.Fatal("changing Hands output changed the table")
	}
	table.Common().SetFaceDown(false)
	table.Reset()
	if table.PlayerCount() != 0 || table.Common().Len() != 54 || !table.Common().FaceDown() {
		t.Fatalf("bad reset: players=%d common=%d", table.PlayerCount(), table.Common().Len())
	}
	if err := table.Deal(3, 10); err != nil {
		t.Fatalf("table should deal again after reset: %v", err)
	}
}

func TestJokerLimit(t *testing.T) {
	table := NewTable()
	if err := table.AddJokers(1); err != nil {
		t.Fatal(err)
	}
	if err := table.AddJokers(1); err != nil {
		t.Fatal(err)
	}
	before := table.Common().Cards()
	if err := table.AddJokers(1); err == nil || !reflect.DeepEqual(table.Common().Cards(), before) {
		t.Fatal("third joker should fail without changing the table")
	}
}
