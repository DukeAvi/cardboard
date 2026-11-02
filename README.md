# Cardboard

A small Go library for playing cards and terminal rendering. This is an early prototype.

![Cards rendered in the terminal](cards-preview.png)

## Try it

Requires Go 1.24 or newer and a terminal with truecolor and Unicode support.

```sh
go run ./cmd/table
go run ./cmd/table --gallery
```

## Use it

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

The renderer returns strings that can also be used in a TUI. It includes a standard deck, card backs, jokers, and a hand layout.
