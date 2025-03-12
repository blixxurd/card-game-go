package main

import (
	"fmt"

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
	// Create a new Texas Hold'em game
	game := holdem.NewHoldemGame("game1", "Texas Hold'em", 5, 10)

	// Initialize the game
	err := game.Initialize(map[string]interface{}{})
	if err != nil {
		fmt.Printf("Error initializing game: %v\n", err)
		return
	}

	// Add players
	for i := 1; i <= 4; i++ {
		player := &SimplePlayer{
			id:     fmt.Sprintf("player%d", i),
			name:   fmt.Sprintf("Player %d", i),
			active: true,
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
	fmt.Println("Game started!")
	fmt.Printf("Players: %d\n", len(game.Players()))
	fmt.Printf("Current player: %s\n", game.CurrentPlayer().Name())
	fmt.Printf("Game state: %s\n", game.State())

	// Simulate some actions
	simulateGame(game)
}

func simulateGame(game *holdem.HoldemGame) {
	// Get the players' hands
	playerHands := make(map[string][]card.Card)
	for _, p := range game.Players() {
		playerHand := p.Hand()
		if playerHand != nil {
			playerHands[p.ID()] = playerHand.Cards()
		}
	}

	// Print the players' hands
	fmt.Println("\nPlayer hands:")
	for _, p := range game.Players() {
		cards := playerHands[p.ID()]
		fmt.Printf("%s: %v\n", p.Name(), cards)
	}

	// Simulate the flop, turn, and river
	fmt.Println("\nSimulating the game...")

	// Simulate all players checking
	for i := 0; i < len(game.Players()); i++ {
		currentPlayer := game.CurrentPlayer()
		fmt.Printf("%s checks\n", currentPlayer.Name())

		action := holdem.NewHoldemAction(holdem.ActionCheck, currentPlayer.ID(), 0)
		err := game.ProcessAction(action)
		if err != nil {
			fmt.Printf("Error processing action: %v\n", err)
			return
		}
	}

	// Print the flop
	fmt.Println("\nFlop:")
	printCommunityCards(game)

	// Simulate all players checking again
	for i := 0; i < len(game.Players()); i++ {
		currentPlayer := game.CurrentPlayer()
		fmt.Printf("%s checks\n", currentPlayer.Name())

		action := holdem.NewHoldemAction(holdem.ActionCheck, currentPlayer.ID(), 0)
		err := game.ProcessAction(action)
		if err != nil {
			fmt.Printf("Error processing action: %v\n", err)
			return
		}
	}

	// Print the turn
	fmt.Println("\nTurn:")
	printCommunityCards(game)

	// Simulate all players checking again
	for i := 0; i < len(game.Players()); i++ {
		currentPlayer := game.CurrentPlayer()
		fmt.Printf("%s checks\n", currentPlayer.Name())

		action := holdem.NewHoldemAction(holdem.ActionCheck, currentPlayer.ID(), 0)
		err := game.ProcessAction(action)
		if err != nil {
			fmt.Printf("Error processing action: %v\n", err)
			return
		}
	}

	// Print the river
	fmt.Println("\nRiver:")
	printCommunityCards(game)

	// Evaluate hands
	fmt.Println("\nEvaluating hands...")
	evaluateHands(game, playerHands)
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

	for _, c := range communityCards {
		fmt.Printf("%s ", c.String())
	}
	fmt.Println()
}

func evaluateHands(game *holdem.HoldemGame, playerHands map[string][]card.Card) {
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

	// Evaluate each player's hand
	results := make(map[string]pokerhand.HandResult)
	for playerID, cards := range playerHands {
		// Combine player cards and community cards
		allCards := append(cards, communityCards...)

		// Evaluate the hand
		result := pokerhand.EvaluateHand(allCards)
		results[playerID] = result

		// Print the result
		player := findPlayer(game, playerID)
		fmt.Printf("%s: %s\n", player.Name(), result.Description)
	}

	// Determine the winner
	determineWinner(game, results)
}

func findPlayer(game *holdem.HoldemGame, playerID string) player.Player {
	for _, p := range game.Players() {
		if p.ID() == playerID {
			return p
		}
	}
	return nil
}

func determineWinner(game *holdem.HoldemGame, results map[string]pokerhand.HandResult) {
	// Find the player with the best hand
	var bestPlayerID string
	var bestResult pokerhand.HandResult

	for playerID, result := range results {
		if bestPlayerID == "" || pokerhand.CompareHands(result, bestResult) > 0 {
			bestPlayerID = playerID
			bestResult = result
		}
	}

	// Print the winner
	winner := findPlayer(game, bestPlayerID)
	fmt.Printf("\nWinner: %s with %s\n", winner.Name(), bestResult.Description)
}
