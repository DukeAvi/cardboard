package cardboard

import "testing"

func TestTurnOrder(t *testing.T) {
	turns, err := NewTurnOrder(4)
	if err != nil {
		t.Fatal(err)
	}

	if turns.Current() != 0 {
		t.Fatal("expected player 0")
	}

	next, err := turns.Advance(1)
	if err != nil || next != 1 {
		t.Fatalf("expected player 1, got %d", next)
	}

	next, err = turns.Advance(2)
	if err != nil || next != 3 {
		t.Fatalf("expected player 3, got %d", next)
	}

	next, err = turns.Advance(1)
	if err != nil || next != 0 {
		t.Fatalf("expected player 0, got %d", next)
	}

	turns.Reverse()
	next, err = turns.Advance(1)
	if err != nil || next != 3 {
		t.Fatalf("expected player 3, got %d", next)
	}
}

func TestInvalidTurnOrder(t *testing.T) {
	_, err := NewTurnOrder(0)
	if err == nil {
		t.Fatal("expected error for zero players")
	}

	turns, _ := NewTurnOrder(3)

	_, err = turns.Advance(0)
	if err == nil {
		t.Fatal("expected error for zero steps")
	}

	err = turns.Reset(3)
	if err == nil {
		t.Fatal("expected error for invalid player")
	}
}
