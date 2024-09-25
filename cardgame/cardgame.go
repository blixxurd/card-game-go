package cardgame

import (
	"fmt"
	"math/rand"
	"time"
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

type Deck []Card

type Hand []Card

type Game struct {
	Deck          Deck
	Hands         []Hand
	ReferenceDeck Deck // Full copy of the original deck for verification
}

// MARK: Functions
func NewDeck() Deck {
	var deck Deck
	for suit := Spades; suit <= Clubs; suit++ {
		for value := 1; value <= 13; value++ {
			deck = append(deck, Card{Suit: suit, Value: value})
		}
	}
	return deck
}

/**
 * Creates a new CardGame with the specified number of hands.
 */
func NewGame(numHands int) *Game {
	deck := NewDeck()
	referenceDeck := make(Deck, len(deck))
	copy(referenceDeck, deck)

	game := &Game{
		Deck:          deck,
		Hands:         make([]Hand, numHands),
		ReferenceDeck: referenceDeck,
	}
	game.Deck.Shuffle()
	return game
}

/**
 * Shuffles the Deck.
 */
func (d Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

/**
 * Draws a card from the Deck
 * This works by removing the first card from the slice and returning it.
 */
func (d *Deck) Draw() (Card, error) {
	if len(*d) == 0 {
		return Card{}, fmt.Errorf("no cards left in the deck")
	}
	card := (*d)[0]
	*d = (*d)[1:]
	return card, nil
}

/**
 * Deals a card to the specified hand
 */
func (g *Game) Deal(handIndex int) error {
	if handIndex < 0 || handIndex >= len(g.Hands) {
		return fmt.Errorf("invalid hand index")
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
	deckCopy := make(Deck, len(g.ReferenceDeck))
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

// MARK: String Functions

/**
 * Returns a string representation of the Card.
 */
func (c Card) String() string {
	values := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	suits := []string{"♠", "♥", "♦", "♣"}
	return fmt.Sprintf("%s%s", values[c.Value-1], suits[c.Suit])
}

/**
 * Returns a string representation of the Deck.
 */
func (h Hand) String() string {
	result := ""
	for _, card := range h {
		result += card.String() + "\n"
	}
	return result
}
