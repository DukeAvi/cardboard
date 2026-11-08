package cardboard

import "fmt"

type TurnOrder struct {
	playerCount int
	current     int
	direction   int
}

func NewTurnOrder(playerCount int) (*TurnOrder, error) {
	if playerCount < 1 {
		return nil, fmt.Errorf("player count must be positive")
	}

	return &TurnOrder{
		playerCount: playerCount,
		current:     0,
		direction:   1,
	}, nil
}

func (t *TurnOrder) Current() int {
	if t == nil {
		return -1
	}
	return t.current
}

func (t *TurnOrder) Advance(steps int) (int, error) {
	if t == nil {
		return -1, fmt.Errorf("turn order must not be nil")
	}

	if steps < 1 {
		return -1, fmt.Errorf("steps must be positive")
	}

	steps %= t.playerCount
	t.current = (t.current + t.direction*steps) % t.playerCount

	if t.current < 0 {
		t.current += t.playerCount
	}

	return t.current, nil
}

func (t *TurnOrder) Reverse() {
	if t == nil {
		return
	}

	t.direction *= -1
}

func (t *TurnOrder) Reset(startPlayer int) error {
	if t == nil {
		return fmt.Errorf("turn order must not be nil")
	}

	if startPlayer < 0 || startPlayer >= t.playerCount {
		return fmt.Errorf("invalid starting player")
	}

	t.current = startPlayer
	t.direction = 1
	return nil
}
