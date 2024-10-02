package cardgame

import (
	"errors"
	"math/rand"
	"time"
)

// MARK: Types
type Deck []Card

// MARK: Functions
func NewDeck() Deck {
	deck := make(Deck, 0, 52)
	for suit := Spades; suit <= Clubs; suit++ {
		for value := 1; value <= 13; value++ {
			deck = append(deck, Card{Suit: suit, Value: value})
		}
	}
	return deck
}

// MARK: Methods

/**
 * Shuffles the Deck.
 */
func (d Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

/**
 * Draws a card from the Deck
 * This works by removing the first card from the slice and returning it.
 */
func (d *Deck) Draw() (Card, error) {
	if len(*d) == 0 {
		return Card{}, errors.New("no cards left in the deck")
	}
	card := (*d)[0]
	*d = (*d)[1:]
	return card, nil
}

/**
 * Adds a card to the deck if it's valid.
 */
func (d *Deck) AddCardToDeck(c Card) error {
	if !IsValidCard(c) {
		return errors.New("invalid card")
	}
	*d = append(*d, c)
	return nil
}

/**
 * Removes a specific card from the deck.
 * Returns true if the card was found and removed, false otherwise.
 */
func (d *Deck) RemoveCard(c Card) bool {
	for i, card := range *d {
		if card == c {
			*d = append((*d)[:i], (*d)[i+1:]...)
			return true
		}
	}
	return false
}
