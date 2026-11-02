package cards

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func assertSize(t *testing.T, rendered string, width, height int) {
	t.Helper()
	rows := strings.Split(rendered, "\n")
	if len(rows) != height {
		t.Fatalf("height: got %d, want %d", len(rows), height)
	}
	for y, row := range rows {
		if w := ansi.StringWidth(row); w != width {
			t.Fatalf("row %d: width %d, want %d", y, w, width)
		}
	}
}

func TestFullPixelDeck(t *testing.T) {
	r := NewRenderer(Vintage)
	seen := map[string]bool{}
	for _, c := range Deck() {
		face := r.Card(c)
		assertSize(t, face, PixelWidth, PixelHeight)
		if seen[face] {
			t.Fatalf("duplicate rendering: %+v", c)
		}
		seen[face] = true
		if !strings.Contains(face, "\x1b[38;2;") {
			t.Fatal("missing truecolor foreground")
		}
	}
	for _, s := range []string{r.Back(), r.Joker(1), r.Joker(2), r.Card(Card{})} {
		assertSize(t, s, PixelWidth, PixelHeight)
	}
	if r.Joker(1) == r.Joker(2) {
		t.Fatal("jokers should have different ink")
	}
}

func TestHandSizingAndSelection(t *testing.T) {
	r := NewRenderer(Vintage)
	hand := Deck()[:13]
	for _, width := range []int{26, 29, 40, 80, 100, 160} {
		for _, selected := range []int{-1, 0, 6, 12, 99} {
			assertSize(t, r.Hand(hand, selected, width), width, PixelHeight+5)
		}
	}
	if r.Hand(hand, 0, 25) != "" || r.Hand(nil, 0, 80) != "" {
		t.Fatal("empty/undersized hand")
	}
	// Even the smallest viewport must keep the selected card visible.
	for i, c := range hand {
		got := strings.Split(ansi.Strip(r.Hand(hand, i, PixelWidth)), "\n")
		want := strings.Split(ansi.Strip(r.Card(c)), "\n")
		for y := range want {
			if got[y] != want[y] {
				t.Fatalf("selected card %d clipped at row %d", i, y)
			}
		}
	}
}

func TestCourtSpritesFitAndRotate(t *testing.T) {
	for rank, sprite := range courts {
		if len(sprite) != 16 {
			t.Fatalf("rank %d has %d rows", rank, len(sprite))
		}
		for _, row := range sprite {
			if len(row) != 16 {
				t.Fatalf("rank %d sprite row width %d", rank, len(row))
			}
		}
		p := newCardCanvas()
		p.court(rank, ink, red)
		for y := 4; y < 36; y++ {
			for x := 4; x < 22; x++ {
				if p[y][x] != p[39-y][25-x] {
					t.Fatalf("rank %d is not double-ended", rank)
				}
			}
		}
	}
}
