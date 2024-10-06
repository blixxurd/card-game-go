package cardgame

import (
	"errors"

	"github.com/blixxurd/card-game-go/internal/cardgame/card"
	"github.com/blixxurd/card-game-go/internal/cardgame/deck"
)

// MARK: Types
type Hand []card.Card

type Game struct {
	Deck          deck.Deck
	Hands         []Hand
	ReferenceDeck deck.Deck // Full copy of the original deck for verification
}

// MARK: Functions

/**
 * Creates a new CardGame with the specified number of hands.
 */
func NewGame(numHands int) *Game {
	newDeck := deck.NewDeck()
	referenceDeck := make(deck.Deck, len(newDeck))
	copy(referenceDeck, newDeck)

	game := &Game{
		Deck:          newDeck,
		Hands:         make([]Hand, numHands),
		ReferenceDeck: referenceDeck,
	}
	game.Deck.Shuffle()
	return game
}

/**
 * Deals a card to the specified hand
 */
func (g *Game) Deal(handIndex int) error {
	if handIndex < 0 || handIndex >= len(g.Hands) {
		return errors.New("no cards left in the deck")
	}
	card, err := g.Deck.Draw()
	if err != nil {
		return err
	}
	g.Hands[handIndex] = append(g.Hands[handIndex], card)
	return nil
}

/**
 * Verifies that all cards in the hands are present in the reference deck.
 * Returns true if all hands are valid, along with a slice of indices of invalid hands.
 */
func (g *Game) VerifyHands() (bool, []int) {
	deckCopy := make(deck.Deck, len(g.ReferenceDeck))
	copy(deckCopy, g.ReferenceDeck)

	invalidHandIndices := []int{}

	for i, hand := range g.Hands {
		for _, card := range hand {
			found := false
			for j, refCard := range deckCopy {
				if card == refCard {
					// Remove the card from deckCopy to ensure it's not counted twice
					deckCopy = append(deckCopy[:j], deckCopy[j+1:]...)
					found = true
					break
				}
			}
			if !found {
				invalidHandIndices = append(invalidHandIndices, i)
				break
			}
		}
	}

	return len(invalidHandIndices) == 0, invalidHandIndices
}

// MARK: Methods

/**
 * Returns a string representation of the Hand.
 */
func (h Hand) String() string {
	result := ""
	for _, card := range h {
		result += card.String() + " "
	}
	return result
}
