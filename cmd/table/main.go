// The table command prints a few static examples of the card art.
package main

import (
	"flag"
	"fmt"
	"strings"

	cards "github.com/DukeAvi/cardboard"
	"github.com/charmbracelet/lipgloss"
)

func sampleHand() []cards.Card {
	return []cards.Card{
		{Rank: 1, Suit: cards.Spades},
		{Rank: 13, Suit: cards.Spades},
		{Rank: 11, Suit: cards.Hearts},
		{Rank: 5, Suit: cards.Hearts},
		{Rank: 12, Suit: cards.Clubs},
		{Rank: 11, Suit: cards.Clubs},
		{Rank: 7, Suit: cards.Clubs},
		{Rank: 13, Suit: cards.Diamonds},
		{Rank: 11, Suit: cards.Diamonds},
		{Rank: 10, Suit: cards.Hearts},
		{Rank: 9, Suit: cards.Spades},
		{Rank: 8, Suit: cards.Spades},
		{Rank: 2, Suit: cards.Hearts},
	}
}

func galleryCards(renderer *cards.Renderer) []string {
	var pictures []string

	for _, card := range cards.Deck() {
		pictures = append(pictures, renderer.Card(card))
	}

	pictures = append(pictures, renderer.Joker(1))
	pictures = append(pictures, renderer.Joker(2))
	pictures = append(pictures, renderer.Back())
	return pictures
}

func joinPictures(pictures []string, columns int) string {
	gap := strings.Repeat("  \n", cards.PixelHeight-1) + "  "
	var rows []string

	for start := 0; start < len(pictures); start += columns {
		end := start + columns
		if end > len(pictures) {
			end = len(pictures)
		}

		var row []string
		for i := start; i < end; i++ {
			if len(row) > 0 {
				row = append(row, gap)
			}
			row = append(row, pictures[i])
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}

	return strings.Join(rows, "\n\n")
}

func main() {
	showGallery := flag.Bool("gallery", false, "print the complete deck")
	showSwatch := flag.Bool("swatch", false, "print five sample cards")
	showBack := flag.Bool("back", false, "print the card back")
	flag.Parse()

	renderer := cards.NewRenderer(cards.Vintage)

	if *showBack {
		fmt.Println(renderer.Back())
		return
	}

	if *showSwatch {
		pictures := []string{
			renderer.Card(cards.Card{Rank: 1, Suit: cards.Spades}),
			renderer.Card(cards.Card{Rank: 11, Suit: cards.Hearts}),
			renderer.Card(cards.Card{Rank: 12, Suit: cards.Clubs}),
			renderer.Card(cards.Card{Rank: 13, Suit: cards.Diamonds}),
			renderer.Back(),
		}
		fmt.Println(joinPictures(pictures, len(pictures)))
		return
	}

	if *showGallery {
		fmt.Println(joinPictures(galleryCards(renderer), 3))
		return
	}

	fmt.Println("THE DECK")
	var pictures []string
	for _, card := range sampleHand() {
		pictures = append(pictures, renderer.Card(card))
	}
	fmt.Println(joinPictures(pictures, 3))
}
