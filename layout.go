package cards

func newTile(width, height int) tile {
	t := make(tile, height)
	for y := range t {
		t[y] = make([]cell, width)
		for x := range t[y] {
			t[y][x] = cell{char: " ", fg: felt, bg: felt}
		}
	}
	return t
}

func putTile(dst, src tile, x, y int) {
	for dy, row := range src {
		for dx, c := range row {
			if y+dy >= 0 && y+dy < len(dst) && x+dx >= 0 && x+dx < len(dst[y+dy]) && !c.clear {
				dst[y+dy][x+dx] = c
			}
		}
	}
}

func frameCard(card tile, selected bool) tile {
	width, height := len(card[0])+2, len(card)+2
	t := newTile(width, height)
	if selected {
		for y := range t {
			for x := range t[y] {
				if x == 0 || y == 0 || x == width-1 || y == height-1 {
					t[y][x] = cell{char: " ", fg: gold, bg: gold}
				}
			}
		}
	}
	putTile(t, card, 1, 1)
	return t
}

// HandWindow returns the visible range, keeping the selected card in view.
func (r *Renderer) HandWindow(count, selected, width int) (int, int) {
	w, _ := r.CardSize.Dimensions()
	if count < 1 || width < w {
		return 0, 0
	}
	capacity := 1 + (width-w)/4
	if selected >= 0 && selected < count {
		if width < w+2 {
			return 0, 0
		}
		capacity = 1
		// Leave room for the full selection, its frame, and the final card.
		if width >= 2*w+2 {
			capacity = 2 + (width-2*w-2)/4
		}
	}
	capacity = min(capacity, count)
	start := 0
	if selected >= 0 && selected < count {
		start = min(max(0, selected-capacity/2), count-capacity)
	}
	return start, start + capacity
}

// Hand overlaps cards, exposing four columns of each rank and suit.
// The selected card is raised, outlined, and shown in full.
func (r *Renderer) Hand(hand []Card, selected, width int) string {
	if selected < 0 || selected >= len(hand) {
		selected = -1
	}
	start, end := r.HandWindow(len(hand), selected, width)
	if start == end {
		return ""
	}
	hand = hand[start:end]
	if selected >= 0 {
		selected -= start
	}
	w, h := r.CardSize.Dimensions()
	used := w + (len(hand)-1)*4
	if selected >= 0 {
		used += 2
		if selected < len(hand)-1 {
			used += w - 4
		}
	}
	left := (width - used) / 2
	t := newTile(width, h+4)
	for i, c := range hand {
		if i == selected {
			continue
		}
		x := left + i*4
		if selected >= 0 && i > selected {
			x += w + 2 - 4
		}
		putTile(t, r.face(c), x, 3)
	}
	if selected >= 0 {
		putTile(t, frameCard(r.face(hand[selected]), true), left+selected*4, 0)
	}
	return r.render(t)
}

// Grid shows full cards, using as many columns as fit.
func (r *Renderer) Grid(hand []Card, selected, width int) string {
	w, h := r.CardSize.Dimensions()
	if len(hand) == 0 || width < w+2 {
		return ""
	}
	columns := min(len(hand), (width+1)/(w+3))
	rows := (len(hand) + columns - 1) / columns
	t := newTile(width, rows*(h+3)-1)
	left := (width - (columns*(w+3) - 1)) / 2
	for i, c := range hand {
		putTile(t, frameCard(r.face(c), i == selected), left+(i%columns)*(w+3), (i/columns)*(h+3))
	}
	return r.render(t)
}
