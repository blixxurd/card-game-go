package card

import (
	"fmt"
)

// MARK: Types
type Suit int

const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
)

type Card struct {
	Suit  Suit
	Value int
}

// MARK: Methods

func (c Card) String() string {
	values := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	suits := []string{"♠", "♥", "♦", "♣"}

	if !IsValidCard(c) {
		return "Invalid Card"
	}

	return fmt.Sprintf("%s%s", values[c.Value-1], suits[c.Suit])
}

/**
 * Checks if a card has valid suit and value.
 */
func IsValidCard(c Card) bool {
	return c.Suit >= Spades && c.Suit <= Clubs && c.Value >= 1 && c.Value <= 13
}
