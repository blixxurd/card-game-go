package holdem

import (
	"fmt"

	"github.com/blixxurd/card-game-go/internal/cardgame"
	"github.com/blixxurd/card-game-go/internal/cardgame/card"
	"github.com/blixxurd/card-game-go/internal/games/poker/pokerhand"
)

// MARK: Types
type PlayerHand struct {
	Player     int
	HoleCards  []card.Card
	HandResult pokerhand.HandResult
}

type HoldemGame struct {
	Game           *cardgame.Game
	CommunityCards []card.Card
	PlayerHands    []PlayerHand
	NumPlayers     int
}

// MARK: Functions

/**
 * Creates a new HoldemGame with the specified number of players.
 */
func NewGame(numPlayers int) *HoldemGame {
	return &HoldemGame{
		Game:           cardgame.NewGame(numPlayers),
		CommunityCards: make([]card.Card, 0, 5),
		PlayerHands:    make([]PlayerHand, numPlayers),
		NumPlayers:     numPlayers,
	}
}

/**
 * Deals two cards to each player.
 */
func (g *HoldemGame) DealHoleCards() error {
	for i := 0; i < 2; i++ {
		for handIndex := 0; handIndex < g.NumPlayers; handIndex++ {
			err := g.Game.Deal(handIndex)
			if err != nil {
				return fmt.Errorf("error dealing to hand %d: %v", handIndex, err)
			}
		}
	}
	return nil
}

/**
 * Deals five community cards.
 */
func (g *HoldemGame) DealCommunityCards() error {
	for i := 0; i < 5; i++ {
		card, err := g.Game.Deck.Draw()
		if err != nil {
			return fmt.Errorf("error dealing community card: %v", err)
		}
		g.CommunityCards = append(g.CommunityCards, card)
	}
	return nil
}

/**
 * Evaluates the best hand for each player.
 */
func (g *HoldemGame) EvaluateHands() error {
	for i, hand := range g.Game.Hands {
		allCards := append(hand, g.CommunityCards...)
		result, err := pokerhand.EvaluateBestHand(allCards)
		if err != nil {
			return fmt.Errorf("error evaluating hand for player %d: %v", i+1, err)
		}

		g.PlayerHands[i] = PlayerHand{
			Player:     i + 1,
			HoleCards:  hand,
			HandResult: result,
		}
	}
	return nil
}

/**
 * Determines the winner of the game.
 */
func (g *HoldemGame) DetermineWinner() PlayerHand {
	winner := g.PlayerHands[0]
	for i := 1; i < len(g.PlayerHands); i++ {
		comparison := pokerhand.CompareHands(g.PlayerHands[i].HandResult, winner.HandResult)
		if comparison > 0 {
			winner = g.PlayerHands[i]
		} else if comparison == 0 {
			winner = g.breakTie(winner, g.PlayerHands[i])
		}
	}
	return winner
}

/**
 * Breaks a tie between two hands by comparing the hole cards.
 * The function compares the hole cards from each hand in order
 * and returns the hand with the higher card.
 */
func (g *HoldemGame) breakTie(hand1, hand2 PlayerHand) PlayerHand {
	for i := 0; i < len(hand1.HoleCards) && i < len(hand2.HoleCards); i++ {
		value1 := g.getComparisonValue(hand1.HoleCards[i])
		value2 := g.getComparisonValue(hand2.HoleCards[i])
		if value1 > value2 {
			return hand1
		} else if value2 > value1 {
			return hand2
		}
	}
	return hand1
}

/**
 * Returns the comparison value of a card.
 * The comparison value is the card value, with the exception
 * of the Ace, which is assigned a value of 14 for comparison purposes.
 */
func (g *HoldemGame) getComparisonValue(card card.Card) int {
	if card.Value == 1 { // Ace
		return 14
	}
	return card.Value
}

/**
 * Prints the current state of the game.
 */
func (g *HoldemGame) PrintGameState() {
	fmt.Println("Community cards:")
	fmt.Printf("%v\n", g.CommunityCards)

	for i, hand := range g.PlayerHands {
		fmt.Printf("\n\nPlayer %d:\n", i+1)
		fmt.Printf("Hole cards: %s, %s\n", hand.HoleCards[0], hand.HoleCards[1])
		fmt.Printf("Best hand: %s\n", hand.HandResult.Name)
		fmt.Printf("%v\n", hand.HandResult.Cards)
	}
}

/**
 * Runs a simulation of a Texas Hold'em game with the specified number of players.
 */
func PlayHoldem(numPlayers int) {
	game := NewGame(numPlayers)

	err := game.DealHoleCards()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = game.DealCommunityCards()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = game.EvaluateHands()
	if err != nil {
		fmt.Println(err)
		return
	}

	game.PrintGameState()

	winner := game.DetermineWinner()
	fmt.Printf("\nWinner: Player %d with %s\n", winner.Player, winner.HandResult.Name)

	valid, invalidHands := game.Game.VerifyHands()
	if valid {
		fmt.Println("\nAll hands are valid!")
	} else {
		fmt.Printf("\nInvalid hands found: %v\n", invalidHands)
	}
}
