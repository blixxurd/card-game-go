// Package holdem provides a Texas Hold'em poker game implementation.
package holdem

import (
	"errors"
	"fmt"
	"sync"

	"github.com/blixxurd/card-game-go/pkg/cardgame/card"
	"github.com/blixxurd/card-game-go/pkg/cardgame/deck"
	"github.com/blixxurd/card-game-go/pkg/cardgame/game"
	"github.com/blixxurd/card-game-go/pkg/cardgame/hand"
	"github.com/blixxurd/card-game-go/pkg/cardgame/player"
	"github.com/blixxurd/card-game-go/pkg/cardgame/poker/pokerhand"
)

// GameState represents the state of a Texas Hold'em game
type GameState string

const (
	// StateInitial is the initial state of the game
	StateInitial GameState = "initial"
	// StatePreFlop is the pre-flop state
	StatePreFlop GameState = "pre_flop"
	// StateFlop is the flop state
	StateFlop GameState = "flop"
	// StateTurn is the turn state
	StateTurn GameState = "turn"
	// StateRiver is the river state
	StateRiver GameState = "river"
	// StateShowdown is the showdown state
	StateShowdown GameState = "showdown"
	// StateComplete is the complete state
	StateComplete GameState = "complete"
)

// ActionType represents the type of action a player can take
type ActionType string

const (
	// ActionCheck represents a check action
	ActionCheck ActionType = "check"
	// ActionBet represents a bet action
	ActionBet ActionType = "bet"
	// ActionCall represents a call action
	ActionCall ActionType = "call"
	// ActionRaise represents a raise action
	ActionRaise ActionType = "raise"
	// ActionFold represents a fold action
	ActionFold ActionType = "fold"
)

// HoldemAction represents an action in a Texas Hold'em game
type HoldemAction struct {
	actionType ActionType
	playerID   string
	amount     int
}

// NewHoldemAction creates a new HoldemAction
func NewHoldemAction(actionType ActionType, playerID string, amount int) *HoldemAction {
	return &HoldemAction{
		actionType: actionType,
		playerID:   playerID,
		amount:     amount,
	}
}

// Type returns the type of the action
func (a *HoldemAction) Type() string {
	return string(a.actionType)
}

// PlayerID returns the ID of the player who performed the action
func (a *HoldemAction) PlayerID() string {
	return a.playerID
}

// Data returns the data associated with the action
func (a *HoldemAction) Data() map[string]interface{} {
	return map[string]interface{}{
		"amount": a.amount,
	}
}

// Params returns the parameters for the action (for compatibility with game.Action)
func (a *HoldemAction) Params() map[string]interface{} {
	return map[string]interface{}{
		"amount": a.amount,
	}
}

// HoldemGame represents a Texas Hold'em game
type HoldemGame struct {
	id               string
	name             string
	state            GameState
	players          []player.Player
	currentPlayerIdx int
	mainDeck         deck.Deck
	communityCards   []card.Card
	playerHands      map[string]hand.Hand
	smallBlind       int
	bigBlind         int
	pot              int
	winners          []player.Player
	mutex            sync.RWMutex
	customData       map[string]interface{}
}

// NewHoldemGame creates a new HoldemGame
func NewHoldemGame(id, name string, smallBlind, bigBlind int) *HoldemGame {
	return &HoldemGame{
		id:               id,
		name:             name,
		state:            StateInitial,
		players:          make([]player.Player, 0),
		currentPlayerIdx: 0,
		mainDeck:         deck.NewStandardDeck(),
		communityCards:   make([]card.Card, 0),
		playerHands:      make(map[string]hand.Hand),
		smallBlind:       smallBlind,
		bigBlind:         bigBlind,
		pot:              0,
		winners:          make([]player.Player, 0),
		customData:       make(map[string]interface{}),
	}
}

// ID returns the ID of the game
func (g *HoldemGame) ID() string {
	return g.id
}

// Name returns the name of the game
func (g *HoldemGame) Name() string {
	return g.name
}

// Type returns the type of the game
func (g *HoldemGame) Type() string {
	return string("poker")
}

// Initialize initializes the game with the provided configuration
func (g *HoldemGame) Initialize(config map[string]interface{}) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.state != StateInitial {
		return errors.New("game already initialized")
	}

	// Shuffle the deck
	g.mainDeck.Shuffle()

	return nil
}

// Start starts the game
func (g *HoldemGame) Start() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.state != StateInitial {
		return errors.New("game already started")
	}

	if len(g.players) < 2 {
		return errors.New("not enough players")
	}

	// Deal hole cards to each player
	g.dealHoleCards()

	// Set the game state to pre-flop
	g.state = StatePreFlop

	return nil
}

// AddPlayer adds a player to the game
func (g *HoldemGame) AddPlayer(p player.Player) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.state != StateInitial {
		return errors.New("cannot add player after game has started")
	}

	// Check if player already exists
	for _, existingPlayer := range g.players {
		if existingPlayer.ID() == p.ID() {
			return errors.New("player already exists")
		}
	}

	g.players = append(g.players, p)
	return nil
}

// RemovePlayer removes a player from the game
func (g *HoldemGame) RemovePlayer(playerID string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.state != StateInitial {
		return errors.New("cannot remove player after game has started")
	}

	for i, p := range g.players {
		if p.ID() == playerID {
			g.players = append(g.players[:i], g.players[i+1:]...)
			return nil
		}
	}

	return errors.New("player not found")
}

// Players returns the players in the game
func (g *HoldemGame) Players() []player.Player {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return g.players
}

// CurrentPlayer returns the current player
func (g *HoldemGame) CurrentPlayer() player.Player {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if len(g.players) == 0 {
		return nil
	}

	return g.players[g.currentPlayerIdx]
}

// NextTurn advances to the next player's turn
func (g *HoldemGame) NextTurn() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.state == StateInitial || g.state == StateComplete {
		return errors.New("game not in progress")
	}

	// Move to the next player
	g.currentPlayerIdx = (g.currentPlayerIdx + 1) % len(g.players)

	// Check if we've completed a round (all players have acted)
	if g.currentPlayerIdx == 0 {
		// Advance the game state
		switch g.state {
		case StatePreFlop:
			g.dealFlop()
			g.state = StateFlop
		case StateFlop:
			g.dealTurn()
			g.state = StateTurn
		case StateTurn:
			g.dealRiver()
			g.state = StateRiver
		case StateRiver:
			g.evaluateHands()
			g.state = StateShowdown
		case StateShowdown:
			g.state = StateComplete
		}
	}

	return nil
}

// State returns the current state of the game
func (g *HoldemGame) State() string {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return string(g.state)
}

// Deck returns the main deck used in the game
func (g *HoldemGame) Deck() deck.Deck {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return g.mainDeck
}

// ProcessAction processes a player action
func (g *HoldemGame) ProcessAction(action game.Action) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.state == StateInitial || g.state == StateComplete {
		return errors.New("game not in progress")
	}

	// Check if it's the player's turn
	currentPlayer := g.players[g.currentPlayerIdx]
	if currentPlayer.ID() != action.PlayerID() {
		return errors.New("not player's turn")
	}

	// Process the action based on its type
	holdemAction, ok := action.(*HoldemAction)
	if !ok {
		return errors.New("invalid action type")
	}

	switch holdemAction.actionType {
	case ActionCheck:
		// No action needed for check
	case ActionBet:
		g.pot += holdemAction.amount
	case ActionCall:
		g.pot += holdemAction.amount
	case ActionRaise:
		g.pot += holdemAction.amount
	case ActionFold:
		// Handle fold logic
	default:
		return errors.New("invalid action type")
	}

	// Move to the next player
	g.currentPlayerIdx = (g.currentPlayerIdx + 1) % len(g.players)

	// Check if we've completed a round (all players have acted)
	if g.currentPlayerIdx == 0 {
		// Advance the game state
		switch g.state {
		case StatePreFlop:
			g.dealFlop()
			g.state = StateFlop
		case StateFlop:
			g.dealTurn()
			g.state = StateTurn
		case StateTurn:
			g.dealRiver()
			g.state = StateRiver
		case StateRiver:
			g.evaluateHands()
			g.state = StateShowdown
		case StateShowdown:
			g.state = StateComplete
		}
	}

	return nil
}

// IsValidAction checks if an action is valid in the current game state
func (g *HoldemGame) IsValidAction(action game.Action) bool {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if g.state == StateInitial || g.state == StateComplete {
		return false
	}

	// Check if it's the player's turn
	currentPlayer := g.players[g.currentPlayerIdx]
	if currentPlayer.ID() != action.PlayerID() {
		return false
	}

	// Check if the action type is valid
	holdemAction, ok := action.(*HoldemAction)
	if !ok {
		return false
	}

	switch holdemAction.actionType {
	case ActionCheck, ActionBet, ActionCall, ActionRaise, ActionFold:
		return true
	default:
		return false
	}
}

// AllowedActions returns the actions currently allowed for the player
func (g *HoldemGame) AllowedActions(playerID string) []string {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if g.state == StateInitial || g.state == StateComplete {
		return []string{}
	}

	// Check if it's the player's turn
	currentPlayer := g.players[g.currentPlayerIdx]
	if currentPlayer.ID() != playerID {
		return []string{}
	}

	// Return all possible actions for simplicity
	return []string{
		string(ActionCheck),
		string(ActionBet),
		string(ActionCall),
		string(ActionRaise),
		string(ActionFold),
	}
}

// Winners returns the winners of the game if it has completed
func (g *HoldemGame) Winners() []player.Player {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	if g.state != StateComplete {
		return []player.Player{}
	}

	return g.winners
}

// Reset resets the game to its initial state
func (g *HoldemGame) Reset() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.state = StateInitial
	g.currentPlayerIdx = 0
	g.mainDeck = deck.NewStandardDeck()
	g.communityCards = make([]card.Card, 0)
	g.playerHands = make(map[string]hand.Hand)
	g.pot = 0
	g.winners = make([]player.Player, 0)

	return nil
}

// Data returns custom game data
func (g *HoldemGame) Data() map[string]interface{} {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return g.customData
}

// SetData sets custom game data
func (g *HoldemGame) SetData(key string, value interface{}) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.customData[key] = value
}

// GetPlayerHand returns the hand of a player
func (g *HoldemGame) GetPlayerHand(playerID string) hand.Hand {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return g.playerHands[playerID]
}

// GetCommunityCards returns the community cards
func (g *HoldemGame) GetCommunityCards() []card.Card {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return g.communityCards
}

// Internal methods

// dealHoleCards deals two cards to each player
func (g *HoldemGame) dealHoleCards() {
	for _, p := range g.players {
		playerHand := hand.NewStandardHand()

		// Deal two cards to the player
		for i := 0; i < 2; i++ {
			card, err := g.mainDeck.Draw()
			if err != nil {
				// Handle error (in a real implementation)
				continue
			}
			playerHand.AddCard(card)
		}

		// Store the hand in our internal map
		g.playerHands[p.ID()] = playerHand

		// Set the hand on the player
		p.SetHand(playerHand)
	}
}

// dealFlop deals the flop (three community cards)
func (g *HoldemGame) dealFlop() {
	// Burn a card
	_, _ = g.mainDeck.Draw()

	// Deal three cards to the community
	for i := 0; i < 3; i++ {
		card, err := g.mainDeck.Draw()
		if err != nil {
			// Handle error (in a real implementation)
			continue
		}
		g.communityCards = append(g.communityCards, card)
	}

	// Store community cards in game data
	g.SetData("community_cards", g.communityCards)
}

// dealTurn deals the turn (fourth community card)
func (g *HoldemGame) dealTurn() {
	// Burn a card
	_, _ = g.mainDeck.Draw()

	// Deal one card to the community
	card, err := g.mainDeck.Draw()
	if err != nil {
		// Handle error (in a real implementation)
		return
	}
	g.communityCards = append(g.communityCards, card)

	// Store community cards in game data
	g.SetData("community_cards", g.communityCards)
}

// dealRiver deals the river (fifth community card)
func (g *HoldemGame) dealRiver() {
	// Burn a card
	_, _ = g.mainDeck.Draw()

	// Deal one card to the community
	card, err := g.mainDeck.Draw()
	if err != nil {
		// Handle error (in a real implementation)
		return
	}
	g.communityCards = append(g.communityCards, card)

	// Store community cards in game data
	g.SetData("community_cards", g.communityCards)
}

// evaluateHands evaluates the hands of all players and determines the winner
func (g *HoldemGame) evaluateHands() {
	// Convert community cards to Card interface
	communityCards := make([]card.Card, 0, len(g.communityCards))
	for _, c := range g.communityCards {
		communityCards = append(communityCards, c)
	}

	// Evaluate each player's hand
	bestRank := -1
	bestPlayers := make([]player.Player, 0)

	for _, p := range g.players {
		playerHand := g.playerHands[p.ID()]
		if playerHand == nil {
			continue
		}

		// Convert player cards to Card interface
		playerCards := make([]card.Card, 0, len(playerHand.Cards()))
		for _, c := range playerHand.Cards() {
			playerCards = append(playerCards, c)
		}

		// Combine player cards and community cards
		allCards := append(playerCards, communityCards...)

		// Evaluate the hand
		result := pokerhand.EvaluateHand(allCards)

		// Check if this is the best hand so far
		if int(result.Rank) > bestRank {
			bestRank = int(result.Rank)
			bestPlayers = []player.Player{p}
		} else if int(result.Rank) == bestRank {
			bestPlayers = append(bestPlayers, p)
		}
	}

	// Set the winners
	g.winners = bestPlayers
}

// String returns a string representation of the game
func (g *HoldemGame) String() string {
	return fmt.Sprintf("HoldemGame{id: %s, name: %s, state: %s, players: %d}", g.id, g.name, g.state, len(g.players))
}
