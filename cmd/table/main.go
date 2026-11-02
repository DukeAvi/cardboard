// The table command shows card art and can print a shuffled deal as text.
package main

import (
	"flag"
	"fmt"
	"os"
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
	var rows []string

	for start := 0; start < len(pictures); start += columns {
		end := start + columns
		if end > len(pictures) {
			end = len(pictures)
		}

		var row []string
		for i := start; i < end; i++ {
			if len(row) > 0 {
				row = append(row, "  ")
			}
			row = append(row, pictures[i])
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, row...))
	}

	return strings.Join(rows, "\n\n")
}

func textCards(hand []cards.Card) string {
	ranks := []string{"", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	suits := []string{"S", "H", "D", "C"}
	var names []string
	for _, card := range hand {
		if card.Suit == cards.Joker {
			names = append(names, fmt.Sprintf("Joker-%d", card.Rank))
		} else {
			names = append(names, ranks[card.Rank]+suits[card.Suit])
		}
	}
	return strings.Join(names, " ")
}

func visualCards(renderer *cards.Renderer, hand []cards.Card) string {
	var pictures []string
	for _, card := range hand {
		pictures = append(pictures, renderer.Card(card))
	}
	return joinPictures(pictures, 3)
}

func main() {
	showGallery := flag.Bool("gallery", false, "print the complete deck")
	showSwatch := flag.Bool("swatch", false, "print five sample cards")
	showBack := flag.Bool("back", false, "print the card back")
	deal := flag.Bool("deal", false, "deal cards and print plain text")
	interactive := flag.Bool("interactive", false, "browse dealt cards interactively")
	players := flag.Int("players", 2, "number of players for -deal")
	cardsEach := flag.Int("cards", 13, "cards per player for -deal")
	jokers := flag.Bool("jokers", false, "include two jokers in -deal")
	flag.Parse()

	if *deal || *interactive {
		pack := cards.NewPack()
		if *jokers {
			pack.Jokers()
		}
		if err := pack.Deal(*players, *cardsEach); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if *interactive {
			if err := runInteractive(pack); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		for i, hand := range pack.Hands {
			fmt.Printf("Player %d (%d): %s\n", i+1, len(hand), textCards(hand))
		}
		fmt.Printf("Common (%d): %s\n", len(pack.Common), textCards(pack.Common))
		return
	}

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
	fmt.Println(visualCards(renderer, sampleHand()))
}
