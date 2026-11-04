package cardboard

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
	theme  Theme
	styles [8][8]lipgloss.Style
}

func NewRenderer(theme Theme) *Renderer {
	r := &Renderer{theme: theme}
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

func (r *Renderer) Theme() Theme { return r.theme }

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

type point struct{ x, y int }

var pipLayouts = map[Rank][]point{
	Two:   {{8, 0}, {8, 8}},
	Three: {{8, 0}, {8, 4}, {8, 8}},
	Four:  {{3, 0}, {13, 0}, {3, 8}, {13, 8}},
	Five:  {{3, 0}, {13, 0}, {8, 4}, {3, 8}, {13, 8}},
	Six:   {{3, 0}, {13, 0}, {3, 4}, {13, 4}, {3, 8}, {13, 8}},
	Seven: {{3, 0}, {13, 0}, {8, 2}, {3, 4}, {13, 4}, {3, 8}, {13, 8}},
	Eight: {{3, 0}, {13, 0}, {8, 2}, {3, 4}, {13, 4}, {8, 6}, {3, 8}, {13, 8}},
	Nine:  {{3, 0}, {13, 0}, {3, 2}, {13, 2}, {8, 4}, {3, 6}, {13, 6}, {3, 8}, {13, 8}},
	Ten:   {{3, 0}, {13, 0}, {8, 1}, {3, 2}, {13, 2}, {3, 6}, {13, 6}, {8, 7}, {3, 8}, {13, 8}},
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
var courts = map[Rank][]string{
	Jack: {
		".........GG.....", ".....GGGGG......", "....IIIIIII.....", "...IAGAGAGAI....",
		"....IPPPPI......", "....IPIPPI...G..", "....IPPPPI...G..", ".....IPPI....G..",
		"...GGIIIIGG..G..", "..IAGGLGGAAI.G..", ".IAIAGLGGAIAIG..", "IAIAIAGLGGAIAIG.",
		"IGAGAIAGLGGAIGI.", "IAGAGAIAGLGGAII.", "IGAGAGAIAGLGGAI.", "IIAIAIAIAAGLGGI.",
	},
	Queen: {
		"....G..G..G.....", "....GGLGGLG.....", "....IGGGGGI.....", "...IIPPPPII.....",
		"...IIPIPPII.....", "...IIPPPPII..A..", "...IIPLPPII.AGA.", "...IIIPPIII..G..",
		"....GGIIGG...G..", "...IAGLLGAI..G..", "..IAGAGGAGAI.G..", ".IAGAIAGIAGAIG..",
		"IAGAIAGGAIAGAIG.", "IIGGAGLLGAGGIII.", "IAGAIAGGAIAGAAI.", "IIAGAGLLGAGAAII.",
	},
	King: {
		"....G.G.G.G.....", "....GLGLGLG.....", "....GGGGGGG.....", "....IPPPPI...G..",
		"....IPIPPI...G..", "....IPPPPI..GGG.", "....IIIPII...G..", ".....IIIII...G..",
		"...GGIIIIGG..G..", "..IAGGLLGGAI.G..", ".IAIAGGLGGAIIG..", "IAIAIAGGLGGAIIG.",
		"IGAGAIAGGLGGAII.", "IAGAGAIAGGLGGAI.", "IGAGAGAIAGGLGGI.", "IIAIAIAIAAGGLGI.",
	},
}

func (p *pixelCanvas) court(rank Rank, primary, accent int) {
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

func (r *Renderer) pixelsToTile(p pixelCanvas) tile {
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
	if c.Suit == Joker && (c.Rank == Ace || c.Rank == Two) {
		return r.joker(int(c.Rank))
	}
	p := newCardCanvas()
	if !c.Valid() {
		t := r.pixelsToTile(p)
		label(t, PixelWidth/2, PixelHeight/2, "?", red)
		return t
	}
	primary, accent := ink, red
	if c.Red() {
		primary, accent = red, ink
	}
	switch {
	case c.Rank == Ace:
		for y, row := range acePixels[c.Suit] {
			for x, mark := range row {
				if mark == '#' {
					p[13+y][7+x] = primary
				}
			}
		}
	case c.Rank >= Jack:
		p.court(c.Rank, primary, accent)
	default:
		for _, spot := range pipLayouts[c.Rank] {
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
	t := r.pixelsToTile(p)
	suit := []string{"♠", "♥", "♦", "♣"}[c.Suit]
	label(t, 2, 1, rankLabels[c.Rank], primary)
	label(t, 2, 2, suit, primary)
	label(t, PixelWidth-3, PixelHeight-3, suit, primary)
	label(t, PixelWidth-2-len(rankLabels[c.Rank]), PixelHeight-2, rankLabels[c.Rank], primary)
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
	return r.pixelsToTile(p)
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
	p.court(Jack, primary, accent)
	for _, x := range []int{8, 12, 16} {
		p[4][x], p[35][25-x] = gold, gold
		p[5][x], p[34][25-x] = accent, accent
	}
	t := r.pixelsToTile(p)
	label(t, 2, 1, "J", primary)
	label(t, 2, 2, "★", primary)
	label(t, PixelWidth-3, PixelHeight-3, "★", primary)
	label(t, PixelWidth-3, PixelHeight-2, "J", primary)
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

func (r *Renderer) Card(c Card) string { return r.render(r.face(c)) }
func (r *Renderer) Back() string       { return r.render(r.back()) }
