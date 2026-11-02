package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	cards "github.com/DukeAvi/cardboard"
)

func TestGalleryHasCompleteDeck(t *testing.T) {
	renderer := cards.NewRenderer(cards.Vintage)
	pictures := galleryCards(renderer)

	if len(pictures) != 55 {
		t.Fatalf("gallery has %d cards, want 55", len(pictures))
	}
}

func TestSampleHandLayout(t *testing.T) {
	renderer := cards.NewRenderer(cards.Vintage)
	var pictures []string
	for _, card := range sampleHand() {
		pictures = append(pictures, renderer.Card(card))
	}
	rows := strings.Split(joinPictures(pictures, 3), "\n")

	if len(rows) != 5*cards.PixelHeight+4 {
		t.Fatalf("hand has %d rows, want %d", len(rows), 5*cards.PixelHeight+4)
	}

	for i, row := range rows {
		if i%(cards.PixelHeight+1) == cards.PixelHeight {
			continue
		}
		want := 3*cards.PixelWidth + 4
		if i >= 4*(cards.PixelHeight+1) {
			want = cards.PixelWidth
		}
		if ansi.StringWidth(row) != want {
			t.Fatalf("hand row %d has width %d, want %d", i, ansi.StringWidth(row), want)
		}
	}
}
