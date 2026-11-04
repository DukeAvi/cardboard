package cardboard

import (
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
	if err := table.Deal(2, 3); err != nil {
		t.Fatal(err)
	}
	first, _ := table.Hand(0)
	second, _ := table.Hand(1)
	wantFirst := []Card{{Ace, Spades}, {Three, Spades}, {Five, Spades}}
	wantSecond := []Card{{Two, Spades}, {Four, Spades}, {Six, Spades}}
	if !reflect.DeepEqual(first.Cards(), wantFirst) || !reflect.DeepEqual(second.Cards(), wantSecond) {
		t.Fatalf("cards were not dealt around the table: %v / %v", first.Cards(), second.Cards())
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
