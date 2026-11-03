# Cardboard

A Go library for rendering playing cards, shuffling decks, and dealing hands.

![Cards rendered in the terminal](cards-preview.png)

## Give it a shot

Requires Go 1.24 or newer and a terminal with truecolor and Unicode support.

```sh
go run ./cmd/table
go run ./cmd/table --gallery
go run ./cmd/table --swatch
go run ./cmd/table --back
```

The default command prints a sample hand. `--gallery` shows all 52 cards, two jokers, and the card back.

## Deal some cards

Deal 13 cards each to two players and print the hands as text:

```sh
go run ./cmd/table --deal --players 2 --cards 13
```

Add `--jokers` to use 54 cards. The output lists each player's hand and the remaining common pile. `S`, `H`, `D`, and `C` stand for spades, hearts, diamonds, and clubs.

For the interactive view:

```sh
go run ./cmd/table --interactive --players 2 --cards 13 --jokers
```

`--interactive` replaces the old `--visual` flag. It needs a terminal; use `--deal` when redirecting output.

Choose a card size at startup: `1` for small, `2` for medium, or `3` for large. Enter picks medium. The choice lasts for that run.

| Key | Action |
| --- | --- |
| Left / right, or h / l | Select a card |
| Tab / Shift+Tab | Switch between players and the common pile |
| c | Collapse or expand the cards |
| s | Change card size |
| q / Ctrl+C | Quit |

Collapsed hands show the left edge of each card. The selected card has a gold outline and stays fully visible, along with the last card in the visible group. Expanded hands use as many columns as fit. Both views follow terminal resizing and keep the selection on screen. Each group remembers its selected card.

## Go get it

```sh
go get github.com/DukeAvi/cardboard
```

```go
package main

import (
	"fmt"

	cards "github.com/DukeAvi/cardboard"
)

func main() {
	renderer := cards.NewRenderer(cards.Vintage)
	fmt.Println(renderer.Card(cards.Card{Rank: 1, Suit: cards.Spades}))
}
```

The renderer returns strings you can print or use in a terminal UI.

## Use the deck

```go
pack := cards.NewPack().Jokers() // Leave off .Jokers() for a standard 52-card deck.
if err := pack.Deal(2, 13); err != nil {
	log.Fatal(err)
}

playerOne := pack.Hands[0]
playerTwo := pack.Hands[1]
common := pack.Common // 28 cards left with jokers, 26 without.
```

`NewPack()` shuffles the deck. `Deal(players, cardsEach)` deals one card at a time around the table and stores each hand separately. The rest stay in `Common`. Calling `Jokers()` again won't add duplicates.

Players and cards per player must be positive, and there must be enough cards for everyone. A failed deal leaves the pack alone. A pack can only be dealt once; create a new one for another round.

`Hands` and `Common` are ordinary slices. The game owns any changes made after dealing. Draw/discard rules, scoring, and Rummy meld detection aren't implemented yet.

`Deck()` still returns a fresh, ordered slice of the 52 standard cards. A `Card` has a rank from 1 (ace) to 13 (king) and a suit: `Spades`, `Hearts`, `Diamonds`, or `Clubs`. Printed jokers use suit `Joker` and rank 1 or 2. `Red()` reports the card's color.

## Render a hand

```go
renderer := cards.NewRenderer(cards.Vintage)
renderer.CardSize = cards.Medium

fmt.Println(renderer.Hand(pack.Hands[0], 0, 80))
fmt.Println(renderer.Grid(pack.Hands[0], 0, 80))
fmt.Println(renderer.Back())
fmt.Println(renderer.Joker(1))
```

`Hand` overlaps the cards; `Grid` shows full cards. Both take a selected card index and a width in terminal columns. Pass `-1` for no selection. `HandWindow` returns the visible start and end indexes, with the end excluded.

| Size | Width × height in terminal cells |
| --- | --- |
| `cards.Small` | 12 × 9 |
| `cards.Medium` | 18 × 14 |
| `cards.Large` (renderer default) | 26 × 20 |

All sizes use the same artwork. Smaller sizes blend the original pixels; fine details will soften because fewer terminal cells are available. Rank and suit labels stay readable. Frames and spacing take extra room beyond these card dimensions.

You can pass your own `Theme` to `NewRenderer` to change the palette. Create a new renderer when changing themes.

## Tests

```sh
go test ./...
```

The deal tests cover card counts, dealing order, jokers, invalid requests, and keeping hands separate.
