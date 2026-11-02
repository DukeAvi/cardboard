package cards

type CardSize int

const (
	Large CardSize = iota
	Medium
	Small
)

func (s CardSize) Dimensions() (int, int) {
	switch s {
	case Small:
		return 12, 9
	case Medium:
		return 18, 14
	default:
		return PixelWidth, PixelHeight
	}
}

func (s CardSize) String() string {
	switch s {
	case Small:
		return "small"
	case Medium:
		return "medium"
	default:
		return "large"
	}
}
