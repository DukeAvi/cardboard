package cardboard

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func assertRenderSize(t *testing.T, rendered string, width, height int) {
	t.Helper()
	rows := strings.Split(rendered, "\n")
	if len(rows) != height {
		t.Fatalf("height %d, want %d", len(rows), height)
	}
	for row, line := range rows {
		if got := ansi.StringWidth(line); got != width {
			t.Fatalf("row %d has width %d, want %d", row, got, width)
		}
	}
}

func TestCardRenderer(t *testing.T) {
	renderer := NewRenderer(Vintage)
	seen := make(map[string]bool)
	for _, card := range StandardDeck() {
		rendered := renderer.Card(card)
		assertRenderSize(t, rendered, PixelWidth, PixelHeight)
		if seen[rendered] {
			t.Fatalf("duplicate card rendering: %v", card)
		}
		seen[rendered] = true
	}
	for _, card := range []Card{{Rank: Ace, Suit: Joker}, {Rank: Two, Suit: Joker}, {}} {
		assertRenderSize(t, renderer.Card(card), PixelWidth, PixelHeight)
	}
	assertRenderSize(t, renderer.Back(), PixelWidth, PixelHeight)
	if renderer.Card(Card{Rank: Queen, Suit: Hearts, ID: 1}) != renderer.Card(Card{Rank: Queen, Suit: Hearts, ID: 2}) {
		t.Fatal("physical card IDs should not change the artwork")
	}
}

func TestFaceDownPileDoesNotLeakCards(t *testing.T) {
	one, _ := NewPile(Card{Rank: Ace, Suit: Spades}, Card{Rank: King, Suit: Hearts})
	two, _ := NewPile(Card{Rank: Two, Suit: Clubs}, Card{Rank: Queen, Suit: Diamonds})
	one.SetFaceDown(true)
	two.SetFaceDown(true)
	renderer := NewRenderer(Vintage)
	for _, layout := range []Layout{Overlap, Grid, Stack} {
		view := View{Width: 80, Height: 50, Layout: layout, Selected: 0}
		if renderer.Pile(one, view) != renderer.Pile(two, view) {
			t.Fatalf("%v layout exposed hidden card contents", layout)
		}
	}
}

func TestPileLayoutsStayWithinBounds(t *testing.T) {
	pile, _ := NewPile(StandardDeck()[:13]...)
	renderer := NewRenderer(Vintage)
	for _, layout := range []Layout{Overlap, Grid, Stack} {
		for _, width := range []int{PixelWidth + 2, 60, 100} {
			view := View{Width: width, Height: 50, Layout: layout, Selected: 6, Picking: true, Picked: 2}
			rendered := renderer.Pile(pile, view)
			if rendered == "" {
				t.Fatalf("%v layout unexpectedly empty at width %d", layout, width)
			}
			for row, line := range strings.Split(rendered, "\n") {
				if got := ansi.StringWidth(line); got != width {
					t.Fatalf("%v row %d has width %d, want %d", layout, row, got, width)
				}
			}
		}
	}
	if renderer.Pile(pile, View{Width: PixelWidth + 1}) != "" {
		t.Fatal("layout should be empty when a framed card cannot fit")
	}
	if renderer.Pile(pile, View{Width: 100, Height: PixelHeight + 1, Layout: Grid}) != "" {
		t.Fatal("layout should be empty when a framed card cannot fit vertically")
	}
}

func TestOverlapWindowUsesFramedCardWidth(t *testing.T) {
	start, end := visibleWindow(13, View{Width: 52, Selected: -1})
	if start != 0 || end != 7 {
		t.Fatalf("unselected window is %d:%d, want 0:7", start, end)
	}

	start, end = visibleWindow(13, View{Width: 54, Selected: 6})
	if start != 5 || end != 7 {
		t.Fatalf("54-column selected window is %d:%d, want 5:7", start, end)
	}

	start, end = visibleWindow(13, View{Width: 60, Selected: 6})
	if start != 5 || end != 8 {
		t.Fatalf("60-column selected window is %d:%d, want 5:8", start, end)
	}
}
