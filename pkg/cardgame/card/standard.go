// Package card defines the core interfaces and implementations for cards in various card games.
package card

import (
	"fmt"
	"strconv"
)

// Suit represents the suit of a standard playing card
type Suit int

const (
	// Spades suit
	Spades Suit = iota
	// Hearts suit
	Hearts
	// Diamonds suit
	Diamonds
	// Clubs suit
	Clubs
)

// String returns a string representation of the suit
func (s Suit) String() string {
	suits := []string{"Spades", "Hearts", "Diamonds", "Clubs"}
	if s < Spades || s > Clubs {
		return "Unknown"
	}
	return suits[s]
}

// Symbol returns the symbol representation of the suit
func (s Suit) Symbol() string {
	symbols := []string{"♠", "♥", "♦", "♣"}
	if s < Spades || s > Clubs {
		return "?"
	}
	return symbols[s]
}

// PlayingCard represents a standard playing card
type PlayingCard struct {
	Suit  Suit
	Value int
}

// ID returns a unique identifier for the card
func (c PlayingCard) ID() string {
	return fmt.Sprintf("%d-%d", c.Suit, c.Value)
}

// Name returns a human-readable name for the card
func (c PlayingCard) Name() string {
	if !c.IsValid() {
		return "Invalid Card"
	}

	valueName := ""
	switch c.Value {
	case 1:
		valueName = "Ace"
	case 11:
		valueName = "Jack"
	case 12:
		valueName = "Queen"
	case 13:
		valueName = "King"
	default:
		valueName = strconv.Itoa(c.Value)
	}

	return valueName + " of " + c.Suit.String()
}

// String returns a string representation of the card
func (c PlayingCard) String() string {
	values := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	if !c.IsValid() {
		return "Invalid Card"
	}

	return fmt.Sprintf("%s%s", values[c.Value-1], c.Suit.Symbol())
}

// Equal checks if two cards are equivalent
func (c PlayingCard) Equal(other Card) bool {
	otherCard, ok := other.(PlayingCard)
	if !ok {
		return false
	}
	return c.Suit == otherCard.Suit && c.Value == otherCard.Value
}

// Clone creates a deep copy of the card
func (c PlayingCard) Clone() Card {
	return PlayingCard{
		Suit:  c.Suit,
		Value: c.Value,
	}
}

// IsValid checks if a card has valid suit and value
func (c PlayingCard) IsValid() bool {
	return c.Suit >= Spades && c.Suit <= Clubs && c.Value >= 1 && c.Value <= 13
}

// GetComparisonValue returns the comparison value of a card
// The comparison value is the card value, with the exception
// of the Ace, which is assigned a value of 14 for comparison purposes
func (c PlayingCard) GetComparisonValue() int {
	if c.Value == 1 { // Ace
		return 14
	}
	return c.Value
}

// PlayingCardFactory creates standard playing cards
type PlayingCardFactory struct{}

// CreateCard creates a new standard playing card
func (f *PlayingCardFactory) CreateCard(params map[string]interface{}) (Card, error) {
	suitVal, ok := params["suit"].(int)
	if !ok {
		return nil, fmt.Errorf("invalid suit parameter")
	}

	valueVal, ok := params["value"].(int)
	if !ok {
		return nil, fmt.Errorf("invalid value parameter")
	}

	card := PlayingCard{
		Suit:  Suit(suitVal),
		Value: valueVal,
	}

	if !card.IsValid() {
		return nil, fmt.Errorf("invalid card parameters")
	}

	return card, nil
}

// Type returns the type of cards this factory creates
func (f *PlayingCardFactory) Type() CardType {
	return StandardCard
}

// NewPlayingCard creates a new standard playing card
func NewPlayingCard(suit Suit, value int) (PlayingCard, error) {
	card := PlayingCard{
		Suit:  suit,
		Value: value,
	}

	if !card.IsValid() {
		return PlayingCard{}, fmt.Errorf("invalid card parameters: suit=%d, value=%d", suit, value)
	}

	return card, nil
}
