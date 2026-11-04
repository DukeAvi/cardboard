package cardboard

type Layout uint8

const (
	Overlap Layout = iota
	Grid
	Stack
)

// View describes how a pile should be rendered. Picked is used only when Picking is true.
type View struct {
	Width    int
	Height   int
	Layout   Layout
	Selected int
	Picking  bool
	Picked   int
}

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

func frameCard(card tile, color int, visible bool) tile {
	width, height := len(card[0])+2, len(card)+2
	t := newTile(width, height)
	if visible {
		for y := range t {
			for x := range t[y] {
				if x == 0 || y == 0 || x == width-1 || y == height-1 {
					t[y][x] = cell{char: " ", fg: color, bg: color}
				}
			}
		}
	}
	putTile(t, card, 1, 1)
	return t
}

// Pile renders a pile according to its face state and the requested layout.
func (r *Renderer) Pile(pile *Pile, view View) string {
	if pile == nil || pile.Len() == 0 || view.Width < PixelWidth+2 {
		return ""
	}
	minimumHeight := PixelHeight + 2
	if view.Layout == Overlap {
		minimumHeight = PixelHeight + 4
	}
	if view.Height > 0 && view.Height < minimumHeight {
		return ""
	}
	if view.Selected < 0 || view.Selected >= pile.Len() {
		view.Selected = -1
	}
	if !view.Picking || view.Picked < 0 || view.Picked >= pile.Len() {
		view.Picking = false
		view.Picked = -1
	}
	card := func(index int) tile {
		if pile.faceDown {
			return r.back()
		}
		return r.face(pile.cards[index])
	}
	switch view.Layout {
	case Grid:
		return r.grid(pile.Len(), view, card)
	case Stack:
		return r.stack(view, card(0))
	default:
		return r.overlap(pile.Len(), view, card)
	}
}

func highlight(view View, index int) (int, bool) {
	if view.Picking && index == view.Picked {
		return red, true
	}
	return gold, index == view.Selected
}

func visibleWindow(count int, view View) (int, int) {
	if count < 1 || view.Width < PixelWidth+2 {
		return 0, 0
	}
	frameWidth := PixelWidth + 2
	capacity := min(1+(view.Width-frameWidth)/4, count)
	for capacity > 0 {
		start := 0
		if view.Selected >= 0 {
			start = min(max(0, view.Selected-capacity/2), count-capacity)
		}
		used := frameWidth + (capacity-1)*4
		if view.Selected >= start && view.Selected < start+capacity-1 {
			used += frameWidth - 4
		}
		if used <= view.Width {
			return start, start + capacity
		}
		capacity--
	}
	return 0, 0
}

func (r *Renderer) overlap(count int, view View, card func(int) tile) string {
	start, end := visibleWindow(count, view)
	if start == end {
		return ""
	}
	shown := end - start
	selected := view.Selected - start
	used := PixelWidth + 2 + (shown-1)*4
	if selected >= 0 && selected < shown-1 {
		used += PixelWidth - 2
	}
	left := (view.Width - used) / 2
	t := newTile(view.Width, PixelHeight+4)
	for i := 0; i < shown; i++ {
		if i == selected {
			continue
		}
		x := left + i*4
		if selected >= 0 && i > selected {
			x += PixelWidth + 2 - 4
		}
		color, visible := highlight(view, start+i)
		putTile(t, frameCard(card(start+i), color, visible), x, 2)
	}
	if selected >= 0 && selected < shown {
		color, visible := highlight(view, view.Selected)
		putTile(t, frameCard(card(view.Selected), color, visible), left+selected*4, 0)
	}
	return r.render(t)
}

func (r *Renderer) grid(count int, view View, card func(int) tile) string {
	columns := min(count, (view.Width+1)/(PixelWidth+3))
	if columns < 1 {
		return ""
	}
	rowsAvailable := (count + columns - 1) / columns
	if view.Height > 0 {
		rowsAvailable = max(1, min(rowsAvailable, (view.Height+1)/(PixelHeight+3)))
	}
	pageSize := columns * rowsAvailable
	start := 0
	if view.Selected >= 0 {
		start = view.Selected / pageSize * pageSize
	}
	end := min(start+pageSize, count)
	rows := (end - start + columns - 1) / columns
	t := newTile(view.Width, rows*(PixelHeight+3)-1)
	left := (view.Width - (columns*(PixelWidth+3) - 1)) / 2
	for index := start; index < end; index++ {
		position := index - start
		color, visible := highlight(view, index)
		putTile(t, frameCard(card(index), color, visible), left+(position%columns)*(PixelWidth+3), (position/columns)*(PixelHeight+3))
	}
	return r.render(t)
}

func (r *Renderer) stack(view View, card tile) string {
	color, visible := highlight(view, 0)
	framed := frameCard(card, color, visible)
	t := newTile(view.Width, len(framed))
	putTile(t, framed, (view.Width-len(framed[0]))/2, 0)
	return r.render(t)
}
