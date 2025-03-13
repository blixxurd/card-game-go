// Package deck defines interfaces and implementations for card decks in various card games.
package deck

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/blixxurd/card-game-go/pkg/cardgame/card"
)

// PlayingDeck represents a standard 52-card deck
type PlayingDeck struct {
	cards []card.Card
}

// NewPlayingDeck creates a new standard 52-card deck
func NewPlayingDeck() (*PlayingDeck, error) {
	deck := &PlayingDeck{
		cards: make([]card.Card, 0, 52),
	}

	// Create all 52 cards
	for suit := card.Spades; suit <= card.Clubs; suit++ {
		for value := 1; value <= 13; value++ {
			c, err := card.NewPlayingCard(suit, value)
			if err != nil {
				return nil, fmt.Errorf("failed to create card: %w", err)
			}
			deck.cards = append(deck.cards, c)
		}
	}

	return deck, nil
}

// Shuffle randomizes the order of cards in the deck using secure random number generation
func (d *PlayingDeck) Shuffle() error {
	for i := range d.cards {
		// Generate cryptographically secure random number
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(d.cards)-i)))
		if err != nil {
			return fmt.Errorf("failed to generate random number: %w", err)
		}
		j := int(n.Int64()) + i
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
	return nil
}

// Draw removes and returns the top card from the deck
func (d *PlayingDeck) Draw() (card.Card, error) {
	if len(d.cards) == 0 {
		return nil, errors.New("no cards left in the deck")
	}

	card := d.cards[0]
	d.cards = d.cards[1:]
	return card, nil
}

// DrawN removes and returns the top n cards from the deck
func (d *PlayingDeck) DrawN(n int) ([]card.Card, error) {
	if len(d.cards) < n {
		return nil, fmt.Errorf("not enough cards in the deck, requested %d but only have %d", n, len(d.cards))
	}

	cards := make([]card.Card, n)
	copy(cards, d.cards[:n])
	d.cards = d.cards[n:]
	return cards, nil
}

// AddCard adds a card to the deck
func (d *PlayingDeck) AddCard(c card.Card) error {
	// Verify it's a playing card
	_, ok := c.(card.PlayingCard)
	if !ok {
		return errors.New("card is not a standard playing card")
	}

	d.cards = append(d.cards, c)
	return nil
}

// AddCards adds multiple cards to the deck
func (d *PlayingDeck) AddCards(cards []card.Card) error {
	for _, c := range cards {
		if err := d.AddCard(c); err != nil {
			return err
		}
	}
	return nil
}

// RemoveCard removes a specific card from the deck
func (d *PlayingDeck) RemoveCard(c card.Card) (bool, error) {
	for i, deckCard := range d.cards {
		if deckCard.Equal(c) {
			d.cards = append(d.cards[:i], d.cards[i+1:]...)
			return true, nil
		}
	}
	return false, nil
}

// Size returns the number of cards in the deck
func (d *PlayingDeck) Size() int {
	return len(d.cards)
}

// IsEmpty returns true if the deck has no cards
func (d *PlayingDeck) IsEmpty() bool {
	return len(d.cards) == 0
}

// Reset returns the deck to its initial state
func (d *PlayingDeck) Reset() error {
	newDeck, err := NewPlayingDeck()
	if err != nil {
		return fmt.Errorf("failed to create new deck: %w", err)
	}

	d.cards = newDeck.cards
	return nil
}

// Cards returns all cards in the deck (without removing them)
func (d *PlayingDeck) Cards() []card.Card {
	result := make([]card.Card, len(d.cards))
	copy(result, d.cards)
	return result
}

// Clone creates a deep copy of the deck
func (d *PlayingDeck) Clone() Deck {
	clone := &PlayingDeck{
		cards: make([]card.Card, len(d.cards)),
	}

	for i, c := range d.cards {
		clone.cards[i] = c.Clone()
	}

	return clone
}

// PlayingDeckFactory creates standard 52-card decks
type PlayingDeckFactory struct{}

// CreateDeck creates a new standard 52-card deck
func (f *PlayingDeckFactory) CreateDeck(params map[string]interface{}) (Deck, error) {
	deck, err := NewPlayingDeck()
	if err != nil {
		return nil, fmt.Errorf("failed to create deck: %w", err)
	}

	// Check if we should shuffle the deck
	if shuffle, ok := params["shuffle"].(bool); ok && shuffle {
		if err := deck.Shuffle(); err != nil {
			return nil, fmt.Errorf("failed to shuffle deck: %w", err)
		}
	}

	return deck, nil
}

// Type returns the type of decks this factory creates
func (f *PlayingDeckFactory) Type() DeckType {
	return StandardDeck
}
