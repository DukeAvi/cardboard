package cards

import (
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// PixelWidth and PixelHeight are terminal cells, not image pixels.
const PixelWidth, PixelHeight = 26, 20

// Theme is the eight-color print palette. Use six-digit hex colors.
type Theme struct {
	Felt, Paper, Edge, Ink, Red, Gold, Skin, Light string
}

var Vintage = Theme{
	Felt: "#142F33", Paper: "#F9E7D2", Edge: "#D9BEA0",
	Ink: "#19474A", Red: "#913F49", Gold: "#B58B51",
	Skin: "#E8B79A", Light: "#FFF3DF",
}

const (
	felt = iota
	paper
	edge
	ink
	red
	gold
	skin
	light
)

type pixelCanvas [PixelHeight * 2][PixelWidth]int
type cell struct {
	char   string
	fg, bg int
	clear  bool
}
type tile [][]cell

// Renderer makes half-block pixel art: each terminal cell holds two pixels.
// It owns its Lip Gloss renderer and never changes global color settings.
type Renderer struct {
	Theme  Theme
	styles [8][8]lipgloss.Style
}

func NewRenderer(theme Theme) *Renderer {
	r := &Renderer{Theme: theme}
	l := lipgloss.NewRenderer(io.Discard)
	l.SetColorProfile(termenv.TrueColor)
	colors := []string{theme.Felt, theme.Paper, theme.Edge, theme.Ink, theme.Red, theme.Gold, theme.Skin, theme.Light}
	for fg := range colors {
		for bg := range colors {
			r.styles[fg][bg] = l.NewStyle().Foreground(lipgloss.Color(colors[fg])).Background(lipgloss.Color(colors[bg]))
		}
	}
	return r
}

func newCardCanvas() pixelCanvas {
	var p pixelCanvas
	for y := range p {
		for x := range p[y] {
			p[y][x] = paper
			if x == 0 || x == PixelWidth-1 || y == 0 || y == PixelHeight*2-1 {
				p[y][x] = edge
			}
			if (x == 0 || x == PixelWidth-1) && (y < 2 || y >= PixelHeight*2-2) {
				p[y][x] = -1
			}
		}
	}
	return p
}

func (p *pixelCanvas) box(x, y, w, h, color int) {
	for dx := 0; dx < w; dx++ {
		p[y][x+dx], p[y+h-1][x+dx] = color, color
	}
	for dy := 0; dy < h; dy++ {
		p[y+dy][x], p[y+dy][x+w-1] = color, color
	}
}

var suitPixels = [4][]string{
	{"..#..", ".###.", "#####", "##.##", "..#.."},
	{".#.#.", "#####", "#####", ".###.", "..#.."},
	{"..#..", ".###.", "#####", ".###.", "..#.."},
	{"..#..", ".###.", "#.#.#", "#####", "..#.."},
}

var acePixels = [4][]string{
	{
		".....#.....", "....###....", "...#####...", "..#######..",
		".#########.", "###########", "###########", "#####.#####",
		".###.#.###.", ".....#.....", "....###....", "...#####...",
	},
	{
		"..##...##..", ".####.####.", "###########", "###########",
		"###########", ".#########.", "..#######..", "...#####...",
		"....###....", ".....#.....", "...........", "...........",
	},
	{
		".....#.....", "....###....", "...#####...", "..#######..",
		".#########.", "###########", ".#########.", "..#######..",
		"...#####...", "....###....", ".....#.....", "...........",
	},
	{
		"....###....", "...#####...", "...#####...", "....###....",
		".##..#..##.", "####.#.####", "###########", "####.#.####",
		".##..#..##.", ".....#.....", "....###....", "...#####...",
	},
}

func (p *pixelCanvas) pip(s Suit, cx, cy, scale, color int, flipped bool) {
	for y, row := range suitPixels[s] {
		for x, mark := range row {
			if mark != '#' {
				continue
			}
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					px, py := x*scale+dx, y*scale+dy
					if flipped {
						px, py = 5*scale-1-px, 5*scale-1-py
					}
					p[cy-5*scale/2+py][cx-5*scale/2+px] = color
				}
			}
		}
	}
}

// These are color-indexed print sprites, not text faces. The lower half is
// rotated 180 degrees to make a traditional double-ended court card.
var courts = map[int][]string{
	11: {
		".........GG.....", ".....GGGGG......", "....IIIIIII.....", "...IAGAGAGAI....",
		"....IPPPPI......", "....IPIPPI...G..", "....IPPPPI...G..", ".....IPPI....G..",
		"...GGIIIIGG..G..", "..IAGGLGGAAI.G..", ".IAIAGLGGAIAIG..", "IAIAIAGLGGAIAIG.",
		"IGAGAIAGLGGAIGI.", "IAGAGAIAGLGGAII.", "IGAGAGAIAGLGGAI.", "IIAIAIAIAAGLGGI.",
	},
	12: {
		"....G..G..G.....", "....GGLGGLG.....", "....IGGGGGI.....", "...IIPPPPII.....",
		"...IIPIPPII.....", "...IIPPPPII..A..", "...IIPLPPII.AGA.", "...IIIPPIII..G..",
		"....GGIIGG...G..", "...IAGLLGAI..G..", "..IAGAGGAGAI.G..", ".IAGAIAGIAGAIG..",
		"IAGAIAGGAIAGAIG.", "IIGGAGLLGAGGIII.", "IAGAIAGGAIAGAAI.", "IIAGAGLLGAGAAII.",
	},
	13: {
		"....G.G.G.G.....", "....GLGLGLG.....", "....GGGGGGG.....", "....IPPPPI...G..",
		"....IPIPPI...G..", "....IPPPPI..GGG.", "....IIIPII...G..", ".....IIIII...G..",
		"...GGIIIIGG..G..", "..IAGGLLGGAI.G..", ".IAIAGGLGGAIIG..", "IAIAIAGGLGGAIIG.",
		"IGAGAIAGGLGGAII.", "IAGAGAIAGGLGGAI.", "IGAGAGAIAGGLGGI.", "IIAIAIAIAAGGLGI.",
	},
}

func (p *pixelCanvas) court(rank, primary, accent int) {
	p.box(4, 4, 18, 32, primary)
	palette := map[rune]int{'I': primary, 'A': accent, 'G': gold, 'P': skin, 'L': light}
	for y, row := range courts[rank] {
		for x, mark := range row {
			if color, ok := palette[mark]; ok {
				p[4+y][5+x] = color
				p[35-y][20-x] = color
			}
		}
	}
}

func pixelsToTile(p pixelCanvas) tile {
	t := make(tile, PixelHeight)
	for y := range t {
		t[y] = make([]cell, PixelWidth)
		for x := range t[y] {
			a, b := p[y*2][x], p[y*2+1][x]
			if a == -1 && b == -1 {
				t[y][x].clear = true
				continue
			}
			if a < 0 {
				a = felt
			}
			if b < 0 {
				b = felt
			}
			t[y][x] = cell{char: "▀", fg: a, bg: b}
			if a == b {
				t[y][x].char = " "
			}
		}
	}
	return t
}

func label(t tile, x, y int, text string, color int) {
	for _, ch := range text {
		t[y][x] = cell{char: string(ch), fg: color, bg: paper}
		x++
	}
}

func (r *Renderer) face(c Card) tile {
	p := newCardCanvas()
	if c.Rank < 1 || c.Rank > 13 || c.Suit > Clubs {
		t := pixelsToTile(p)
		label(t, 11, 9, "?", red)
		return t
	}
	primary, accent := ink, red
	if c.Red() {
		primary, accent = red, ink
	}
	switch {
	case c.Rank == 1:
		for y, row := range acePixels[c.Suit] {
			for x, mark := range row {
				if mark == '#' {
					p[13+y][7+x] = primary
				}
			}
		}
	case c.Rank >= 11:
		p.court(c.Rank, primary, accent)
	default:
		for _, spot := range layouts[c.Rank] {
			x := 12
			if spot.x == 3 {
				x = 7
			}
			if spot.x == 13 {
				x = 18
			}
			y := 7 + spot.y*3
			p.pip(c.Suit, x, y, 1, primary, spot.y > 4)
		}
	}
	t := pixelsToTile(p)
	suit := []string{"♠", "♥", "♦", "♣"}[c.Suit]
	label(t, 2, 1, labels[c.Rank], primary)
	label(t, 2, 2, suit, primary)
	label(t, PixelWidth-3, PixelHeight-3, suit, primary)
	label(t, PixelWidth-2-len(labels[c.Rank]), PixelHeight-2, labels[c.Rank], primary)
	return t
}

func (r *Renderer) back() tile {
	p := newCardCanvas()
	for y := 2; y < 38; y++ {
		for x := 2; x < 24; x++ {
			p[y][x] = ink
		}
	}
	for _, b := range [][4]int{{2, 2, 22, 36}, {4, 4, 18, 32}, {6, 6, 14, 28}} {
		p.box(b[0], b[1], b[2], b[3], edge)
	}
	for y := 3; y < 37; y++ {
		for x := 3; x < 23; x++ {
			if (x == 3 || x == 22 || y == 3 || y == 36) && (x+y)%2 == 0 {
				p[y][x] = gold
			}
		}
	}
	for _, cy := range []int{12, 27} {
		for y := -5; y <= 5; y++ {
			for x := -5; x <= 5; x++ {
				d := abs(x) + abs(y)
				if d == 5 || d == 3 {
					p[cy+y][12+x] = edge
				}
			}
		}
		p[cy][12] = gold
	}
	shape := []string{
		"##.........##",
		".##.......##.",
		"..##.....##..",
		"...#######...",
		"....#####....",
	}
	for y, row := range shape {
		for x, mark := range row {
			if mark == '#' {
				c := gold
				if y == 0 || row[x] == '#' && shape[y-1][x] == '.' {
					c = light
				}
				p[17+y][6+x] = c
			}
		}
	}
	return pixelsToTile(p)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (r *Renderer) joker(number int) tile {
	p := newCardCanvas()
	primary, accent := ink, red
	if number == 2 {
		primary, accent = red, ink
	}
	p.court(11, primary, accent)
	for _, x := range []int{8, 12, 16} {
		p[4][x], p[35][25-x] = gold, gold
		p[5][x], p[34][25-x] = accent, accent
	}
	t := pixelsToTile(p)
	label(t, 2, 1, "J", primary)
	label(t, 2, 2, "★", primary)
	label(t, 23, 17, "★", primary)
	label(t, 23, 18, "J", primary)
	return t
}

// render groups adjacent equal-color cells to avoid per-pixel ANSI resets.
func (r *Renderer) render(t tile) string {
	var out strings.Builder
	for y, row := range t {
		if y > 0 {
			out.WriteByte('\n')
		}
		for x := 0; x < len(row); {
			c := row[x]
			if c.clear {
				c = cell{char: " ", fg: felt, bg: felt}
			}
			var run strings.Builder
			end := x
			for end < len(row) {
				n := row[end]
				if n.clear {
					n = cell{char: " ", fg: felt, bg: felt}
				}
				if n.fg != c.fg || n.bg != c.bg {
					break
				}
				run.WriteString(n.char)
				end++
			}
			out.WriteString(r.styles[c.fg][c.bg].Render(run.String()))
			x = end
		}
	}
	return out.String()
}

func (r *Renderer) Card(c Card) string      { return r.render(r.face(c)) }
func (r *Renderer) Back() string            { return r.render(r.back()) }
func (r *Renderer) Joker(number int) string { return r.render(r.joker(number)) }

// Hand lays out an overlapping hand in exactly width cells. selected is raised
// and drawn last; pass -1 for no selection. Narrow terminals show a focused
// window of cards, keeping at least three columns of each index visible.
func (r *Renderer) Hand(hand []Card, selected, width int) string {
	if width < PixelWidth || len(hand) == 0 {
		return ""
	}
	if selected < -1 || selected >= len(hand) {
		selected = -1
	}
	maxCards := 1 + (width-PixelWidth)/3
	if selected >= 0 {
		// Reserve a complete selected card and a complete final card.
		// This keeps every remaining card's rank visible beside selection.
		maxCards = 1
		if width >= 2*PixelWidth {
			maxCards = 2 + (width-2*PixelWidth)/3
		}
	}
	if len(hand) > maxCards {
		start := 0
		if selected >= 0 {
			start = min(max(0, selected-maxCards/2), len(hand)-maxCards)
		}
		hand = hand[start : start+maxCards]
		if selected >= 0 {
			selected -= start
		}
	}
	step := 8
	if len(hand) > 1 {
		step = min(step, (width-PixelWidth)/(len(hand)-1))
	}
	expand := selected >= 0 && selected < len(hand)-1
	if expand && len(hand) > 2 {
		step = min(step, (width-2*PixelWidth)/(len(hand)-2))
	}
	used := PixelWidth + (len(hand)-1)*step
	if expand {
		used += PixelWidth - step
	}
	left := (width - used) / 2
	t := make(tile, PixelHeight+5)
	for y := range t {
		t[y] = make([]cell, width)
		for x := range t[y] {
			t[y][x] = cell{char: " ", fg: felt, bg: felt}
		}
	}
	draw := func(i int) {
		x := left + i*step
		if expand && i > selected {
			x += PixelWidth - step
		}
		y := 2 + min(2, abs(2*i-(len(hand)-1))/6)
		if i == selected {
			y = 0
		}
		card := r.face(hand[i])
		for dy := 1; dy <= PixelHeight; dy++ {
			for dx := 1; dx <= PixelWidth; dx++ {
				if x+dx < width && y+dy < len(t) {
					t[y+dy][x+dx] = cell{char: " ", fg: ink, bg: ink}
				}
			}
		}
		for dy, row := range card {
			for dx, c := range row {
				if !c.clear {
					t[y+dy][x+dx] = c
				}
			}
		}
	}
	for i := range hand {
		if i != selected {
			draw(i)
		}
	}
	if selected >= 0 {
		draw(selected)
	}
	return r.render(t)
}
