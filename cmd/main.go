package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/blixxurd/card-game-go/pkg/cardgame/card"
	"github.com/blixxurd/card-game-go/pkg/cardgame/hand"
	"github.com/blixxurd/card-game-go/pkg/cardgame/player"
	"github.com/blixxurd/card-game-go/pkg/cardgame/poker/holdem"
	"github.com/blixxurd/card-game-go/pkg/cardgame/poker/pokerhand"
)

// SimplePlayer represents a player in the game
type SimplePlayer struct {
	id       string
	name     string
	score    int
	active   bool
	hand     hand.Hand
	userData map[string]interface{}
	chips    int
}

// ID returns the player's ID
func (p *SimplePlayer) ID() string {
	return p.id
}

// Name returns the player's name
func (p *SimplePlayer) Name() string {
	return p.name
}

// SetName sets the player's name
func (p *SimplePlayer) SetName(name string) {
	p.name = name
}

// Hand returns the player's current hand
func (p *SimplePlayer) Hand() hand.Hand {
	return p.hand
}

// SetHand sets the player's hand
func (p *SimplePlayer) SetHand(h hand.Hand) {
	p.hand = h
}

// Score returns the player's current score
func (p *SimplePlayer) Score() int {
	return p.score
}

// SetScore sets the player's score
func (p *SimplePlayer) SetScore(score int) {
	p.score = score
}

// AddToScore adds points to the player's score
func (p *SimplePlayer) AddToScore(points int) {
	p.score += points
}

// IsActive returns whether the player is active in the current game
func (p *SimplePlayer) IsActive() bool {
	return p.active
}

// SetActive sets whether the player is active
func (p *SimplePlayer) SetActive(active bool) {
	p.active = active
}

// Data returns the player's data
func (p *SimplePlayer) Data() map[string]interface{} {
	if p.userData == nil {
		p.userData = make(map[string]interface{})
	}
	return p.userData
}

// SetData sets the player's data
func (p *SimplePlayer) SetData(key string, value interface{}) {
	if p.userData == nil {
		p.userData = make(map[string]interface{})
	}
	p.userData[key] = value
}

// Clone creates a copy of the player
func (p *SimplePlayer) Clone() player.Player {
	clonedPlayer := &SimplePlayer{
		id:       p.id,
		name:     p.name,
		score:    p.score,
		active:   p.active,
		chips:    p.chips,
		userData: make(map[string]interface{}),
	}

	// Clone the hand if it exists
	if p.hand != nil {
		// For simplicity, we're not doing a deep clone of the hand
		clonedPlayer.hand = p.hand
	}

	// Clone the user data
	for k, v := range p.userData {
		clonedPlayer.userData[k] = v
	}

	return clonedPlayer
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a new Texas Hold'em game
	game := holdem.NewHoldemGame("game1", "Texas Hold'em", 5, 10)

	// Initialize the game
	err := game.Initialize(map[string]interface{}{})
	if err != nil {
		fmt.Printf("Error initializing game: %v\n", err)
		return
	}

	// Add players with starting chips
	for i := 1; i <= 4; i++ {
		player := &SimplePlayer{
			id:     fmt.Sprintf("player%d", i),
			name:   fmt.Sprintf("Player %d", i),
			active: true,
			chips:  1000, // Starting with 1000 chips
		}
		err := game.AddPlayer(player)
		if err != nil {
			fmt.Printf("Error adding player: %v\n", err)
			return
		}
	}

	// Start the game
	err = game.Start()
	if err != nil {
		fmt.Printf("Error starting game: %v\n", err)
		return
	}

	// Print the initial state
	fmt.Println("===== TEXAS HOLD'EM POKER SIMULATION =====")
	fmt.Println("Game started!")
	fmt.Printf("Players: %d\n", len(game.Players()))
	fmt.Printf("Current player: %s\n", game.CurrentPlayer().Name())
	fmt.Printf("Game state: %s\n", game.State())

	// Simulate the game
	simulateSimplePokerGame(game)
}

func simulateSimplePokerGame(game *holdem.HoldemGame) {
	// Get the players' hands
	playerHands := make(map[string][]card.Card)
	for _, p := range game.Players() {
		playerHand := p.Hand()
		if playerHand != nil {
			playerHands[p.ID()] = playerHand.Cards()
		}
	}

	// Print the players' hands
	fmt.Println("\n===== PLAYER HANDS =====")
	for _, p := range game.Players() {
		cards := playerHands[p.ID()]
		fmt.Printf("%s: %v\n", p.Name(), cards)
	}

	// Track active players and pot
	activePlayers := make(map[string]bool)
	for _, p := range game.Players() {
		activePlayers[p.ID()] = true
	}
	potSize := 0

	// Simulate pre-flop betting (simplified)
	fmt.Println("\n===== PRE-FLOP BETTING =====")
	simulateSimpleBetting(game, &potSize, activePlayers)

	// Deal flop
	fmt.Println("\n===== FLOP =====")
	game.DealFlop()
	printCommunityCards(game)

	// Simulate flop betting (simplified)
	fmt.Println("\n===== FLOP BETTING =====")
	simulateSimpleBetting(game, &potSize, activePlayers)

	// Deal turn
	fmt.Println("\n===== TURN =====")
	game.DealTurn()
	printCommunityCards(game)

	// Simulate turn betting (simplified)
	fmt.Println("\n===== TURN BETTING =====")
	simulateSimpleBetting(game, &potSize, activePlayers)

	// Deal river
	fmt.Println("\n===== RIVER =====")
	game.DealRiver()
	printCommunityCards(game)

	// Simulate river betting (simplified)
	fmt.Println("\n===== RIVER BETTING =====")
	simulateSimpleBetting(game, &potSize, activePlayers)

	// Print final state
	fmt.Println("\n===== GAME SUMMARY =====")
	fmt.Printf("Final pot size: %d chips\n", potSize)

	// Print final community cards
	fmt.Println("\nFinal community cards:")
	printCommunityCards(game)

	// Print active players
	fmt.Println("\nPlayers still in the hand:")
	for _, p := range game.Players() {
		if activePlayers[p.ID()] {
			fmt.Printf("%s: %v\n", p.Name(), playerHands[p.ID()])
		}
	}

	// Evaluate hands
	communityCardsData, ok := game.Data()["community_cards"]
	if ok {
		communityCards, ok := communityCardsData.([]card.Card)
		if ok && len(communityCards) > 0 {
			fmt.Println("\n===== HAND EVALUATION =====")

			// Evaluate active players' hands
			results := make(map[string]pokerhand.HandResult)
			for playerID, active := range activePlayers {
				if active {
					// Get player's cards
					playerCards := playerHands[playerID]

					// Combine with community cards
					allCards := append(playerCards, communityCards...)

					// Evaluate hand
					result := pokerhand.EvaluateHand(allCards)
					results[playerID] = result

					// Print result
					player := findPlayer(game, playerID)
					fmt.Printf("%s: %s\n", player.Name(), result.Description)
				}
			}

			// Determine winner
			var bestPlayerID string
			var bestResult pokerhand.HandResult

			for playerID, result := range results {
				if bestPlayerID == "" || pokerhand.CompareHands(result, bestResult) > 0 {
					bestPlayerID = playerID
					bestResult = result
				}
			}

			// Print winner
			if bestPlayerID != "" {
				winner := findPlayer(game, bestPlayerID)
				fmt.Printf("\n🏆 WINNER: %s with %s 🏆\n", winner.Name(), bestResult.Description)
				fmt.Printf("%s wins %d chips!\n", winner.Name(), potSize)
			}
		}
	}
}

// simulateSimpleBetting simulates a simplified betting round
func simulateSimpleBetting(game *holdem.HoldemGame, potSize *int, activePlayers map[string]bool) {
	// Randomly decide which players fold
	for _, p := range game.Players() {
		if activePlayers[p.ID()] {
			// 25% chance to fold
			if rand.Float64() < 0.25 {
				fmt.Printf("%s folds\n", p.Name())
				activePlayers[p.ID()] = false
			} else {
				// Random bet between 10 and 50
				bet := 10 + rand.Intn(41)
				fmt.Printf("%s bets %d\n", p.Name(), bet)
				*potSize += bet
			}
		}
	}

	// Count active players
	activeCount := 0
	for _, active := range activePlayers {
		if active {
			activeCount++
		}
	}

	fmt.Printf("Betting round complete. %d players still active. Current pot: %d chips\n",
		activeCount, *potSize)
}

// SimpleCard is a simple implementation of the Card interface
type SimpleCard struct {
	rank string
	suit string
}

// ID returns a unique identifier for the card
func (c *SimpleCard) ID() string {
	return c.String()
}

// Name returns a human-readable name for the card
func (c *SimpleCard) Name() string {
	return c.String()
}

func (c *SimpleCard) String() string {
	return c.rank + c.suit
}

// Equal checks if two cards are equivalent
func (c *SimpleCard) Equal(other card.Card) bool {
	return c.String() == other.String()
}

func (c *SimpleCard) Value() int {
	// Convert rank to value
	switch c.rank {
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "10":
		return 10
	case "J":
		return 11
	case "Q":
		return 12
	case "K":
		return 13
	case "A":
		return 14
	default:
		return 0
	}
}

// Clone returns a copy of the card
func (c *SimpleCard) Clone() card.Card {
	return &SimpleCard{
		rank: c.rank,
		suit: c.suit,
	}
}

func printCommunityCards(game *holdem.HoldemGame) {
	// Get the community cards from the game's data
	communityCardsData, ok := game.Data()["community_cards"]
	if !ok {
		fmt.Println("No community cards available")
		return
	}

	communityCards, ok := communityCardsData.([]card.Card)
	if !ok {
		fmt.Println("Invalid community cards data")
		return
	}

	if len(communityCards) == 0 {
		fmt.Println("No community cards dealt yet")
		return
	}

	fmt.Print("Community cards: ")
	for _, c := range communityCards {
		fmt.Printf("%s ", c.String())
	}
	fmt.Println()
}

func findPlayer(game *holdem.HoldemGame, playerID string) player.Player {
	for _, p := range game.Players() {
		if p.ID() == playerID {
			return p
		}
	}
	return nil
}
