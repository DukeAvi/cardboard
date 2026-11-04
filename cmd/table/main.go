package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/DukeAvi/cardboard"
)

func main() {
	gallery := flag.Bool("gallery", false, "print all card faces and the card back")
	back := flag.Bool("back", false, "print the card back")
	deal := flag.Bool("deal", false, "deal cards and print them as text")
	interactive := flag.Bool("interactive", false, "open the interactive table")
	players := flag.Int("players", 2, "number of players")
	cardsEach := flag.Int("cards", 13, "cards per player")
	jokers := flag.Int("jokers", 0, "number of jokers (0 to 2)")
	flag.Parse()

	renderer := cardboard.NewRenderer(cardboard.Vintage)
	if *back {
		fmt.Println(renderer.Back())
		return
	}
	if *gallery {
		showGallery(renderer)
		return
	}

	table, err := makeTable(*players, *cardsEach, *jokers)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if *interactive {
		if err := runInteractive(table); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if *deal {
		showTextDeal(table)
		return
	}

	hand, _ := table.Hand(0)
	fmt.Println(renderer.Pile(hand, cardboard.View{Width: 100, Layout: cardboard.Overlap, Selected: 0}))
}

func makeTable(players, cardsEach, jokers int) (*cardboard.Table, error) {
	table := cardboard.NewTable()
	if err := table.AddJokers(jokers); err != nil {
		return nil, err
	}
	table.Shuffle()
	if err := table.Deal(players, cardsEach); err != nil {
		return nil, err
	}
	return table, nil
}

func showGallery(renderer *cardboard.Renderer) {
	deck := cardboard.StandardDeck()
	deck = append(deck,
		cardboard.Card{Rank: cardboard.Ace, Suit: cardboard.Joker},
		cardboard.Card{Rank: cardboard.Two, Suit: cardboard.Joker},
	)
	pile, _ := cardboard.NewPile(deck...)
	fmt.Println(renderer.Pile(pile, cardboard.View{Width: 86, Layout: cardboard.Grid, Selected: -1}))
	fmt.Println()
	fmt.Println(renderer.Back())
}

func showTextDeal(table *cardboard.Table) {
	for player := 0; player < table.PlayerCount(); player++ {
		hand, _ := table.Hand(player)
		fmt.Printf("Player %d (%d): %s\n", player+1, hand.Len(), textCards(hand.Cards()))
	}
	common := table.Common()
	fmt.Printf("Common (%d): %s\n", common.Len(), textCards(common.Cards()))
}

func textCards(cards []cardboard.Card) string {
	names := make([]string, len(cards))
	for i, card := range cards {
		names[i] = card.Short()
	}
	return strings.Join(names, " ")
}
