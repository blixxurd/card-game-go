// Package holdem provides a Texas Hold'em poker game implementation.
package holdem

import (
	"errors"

	"github.com/blixxurd/card-game-go/pkg/cardgame/card"
	"github.com/blixxurd/card-game-go/pkg/cardgame/deck"
	"github.com/blixxurd/card-game-go/pkg/cardgame/game"
	"github.com/blixxurd/card-game-go/pkg/cardgame/player"
	"github.com/blixxurd/card-game-go/pkg/cardgame/poker"
	"github.com/blixxurd/card-game-go/pkg/cardgame/table"
)

// PokerActionAdapter adapts a poker.PokerAction to a HoldemAction
type PokerActionAdapter struct {
	pokerAction poker.PokerAction
}

// NewPokerActionAdapter creates a new PokerActionAdapter
func NewPokerActionAdapter(pokerAction poker.PokerAction) *HoldemAction {
	return &HoldemAction{
		actionType: ActionType(pokerAction.Type()),
		playerID:   pokerAction.PlayerID(),
		amount:     pokerAction.Amount(),
	}
}

// HoldemActionAdapter adapts a HoldemAction to a poker.PokerAction
type HoldemActionAdapter struct {
	holdemAction *HoldemAction
}

// NewHoldemActionAdapter creates a new HoldemActionAdapter
func NewHoldemActionAdapter(holdemAction *HoldemAction) *HoldemActionAdapter {
	return &HoldemActionAdapter{
		holdemAction: holdemAction,
	}
}

// Type returns the type of the action as a string
func (a *HoldemActionAdapter) Type() string {
	return a.holdemAction.Type()
}

// PlayerID returns the ID of the player who performed the action
func (a *HoldemActionAdapter) PlayerID() string {
	return a.holdemAction.PlayerID()
}

// Params returns the parameters for the action
func (a *HoldemActionAdapter) Params() map[string]interface{} {
	return a.holdemAction.Params()
}

// Amount returns the amount of the bet or raise
func (a *HoldemActionAdapter) Amount() int {
	return a.holdemAction.Params()["amount"].(int)
}

// ActionType returns the type of poker action
func (a *HoldemActionAdapter) ActionType() poker.PokerActionType {
	return poker.PokerActionType(a.holdemAction.Type())
}

// HoldemGameAdapter adapts a HoldemGame to implement the PokerGame interface
type HoldemGameAdapter struct {
	holdemGame *HoldemGame
}

// NewHoldemGameAdapter creates a new HoldemGameAdapter
func NewHoldemGameAdapter(holdemGame *HoldemGame) *HoldemGameAdapter {
	return &HoldemGameAdapter{
		holdemGame: holdemGame,
	}
}

// ID returns the ID of the game
func (g *HoldemGameAdapter) ID() string {
	return g.holdemGame.ID()
}

// Name returns the name of the game
func (g *HoldemGameAdapter) Name() string {
	return g.holdemGame.Name()
}

// Initialize sets up the game with the provided configuration
func (g *HoldemGameAdapter) Initialize(config map[string]interface{}) error {
	// HoldemGame doesn't have an Initialize method, so we'll just return nil
	return nil
}

// Start begins the game
func (g *HoldemGameAdapter) Start() error {
	return g.holdemGame.Start()
}

// AddPlayer adds a player to the game
func (g *HoldemGameAdapter) AddPlayer(p player.Player) error {
	return g.holdemGame.AddPlayer(p)
}

// RemovePlayer removes a player from the game
func (g *HoldemGameAdapter) RemovePlayer(playerID string) error {
	return g.holdemGame.RemovePlayer(playerID)
}

// Players returns all players in the game
func (g *HoldemGameAdapter) Players() []player.Player {
	return g.holdemGame.players
}

// CurrentPlayer returns the player whose turn it is
func (g *HoldemGameAdapter) CurrentPlayer() player.Player {
	if g.holdemGame.currentPlayerIdx < 0 || g.holdemGame.currentPlayerIdx >= len(g.holdemGame.players) {
		return nil
	}
	return g.holdemGame.players[g.holdemGame.currentPlayerIdx]
}

// NextTurn advances to the next player's turn
func (g *HoldemGameAdapter) NextTurn() error {
	return g.holdemGame.NextTurn()
}

// State returns the current state of the game
func (g *HoldemGameAdapter) State() game.GameState {
	switch g.holdemGame.state {
	case StateInitial:
		return game.GameStateInitialized
	case StateComplete:
		return game.GameStateCompleted
	default:
		return game.GameStateInProgress
	}
}

// Deck returns the main deck used in the game
func (g *HoldemGameAdapter) Deck() deck.Deck {
	return g.holdemGame.mainDeck
}

// ProcessAction processes a player action
func (g *HoldemGameAdapter) ProcessAction(action game.Action) error {
	// Convert the generic action to a HoldemAction
	holdemAction := &HoldemAction{
		actionType: ActionType(action.Type()),
		playerID:   action.PlayerID(),
		amount:     action.Params()["amount"].(int),
	}

	return g.holdemGame.ProcessAction(holdemAction)
}

// IsValidAction checks if an action is valid in the current game state
func (g *HoldemGameAdapter) IsValidAction(action game.Action) bool {
	// Convert the generic action to a HoldemAction
	holdemAction := &HoldemAction{
		actionType: ActionType(action.Type()),
		playerID:   action.PlayerID(),
		amount:     action.Params()["amount"].(int),
	}

	return g.holdemGame.IsValidAction(holdemAction)
}

// AllowedActions returns the actions that are currently allowed
func (g *HoldemGameAdapter) AllowedActions(playerID string) []string {
	allowedActions := g.holdemGame.AllowedActions(playerID)
	return allowedActions
}

// Winners returns the winners of the game (if any)
func (g *HoldemGameAdapter) Winners() []player.Player {
	return g.holdemGame.winners
}

// Reset resets the game to its initial state
func (g *HoldemGameAdapter) Reset() error {
	return g.holdemGame.Reset()
}

// Data returns custom game data as a map
func (g *HoldemGameAdapter) Data() map[string]interface{} {
	return g.holdemGame.customData
}

// SetData sets custom game data
func (g *HoldemGameAdapter) SetData(key string, value interface{}) {
	g.holdemGame.customData[key] = value
}

// Table returns the table used in the game
func (g *HoldemGameAdapter) Table() *table.Table {
	return g.holdemGame.table
}

// CommunityCards returns the community cards
func (g *HoldemGameAdapter) CommunityCards() []card.Card {
	return g.holdemGame.communityCards
}

// ActivePlayers returns the players who are still in the hand
func (g *HoldemGameAdapter) ActivePlayers() []player.Player {
	activePlayers := make([]player.Player, 0)
	for _, p := range g.holdemGame.players {
		if g.holdemGame.activePlayers[p.ID()] {
			activePlayers = append(activePlayers, p)
		}
	}
	return activePlayers
}

// Pot returns the current pot size
func (g *HoldemGameAdapter) Pot() int {
	return g.holdemGame.pot
}

// CurrentBet returns the current bet amount
func (g *HoldemGameAdapter) CurrentBet() int {
	return g.holdemGame.table.CurrentBet
}

// MinRaise returns the minimum raise amount
func (g *HoldemGameAdapter) MinRaise() int {
	return g.holdemGame.table.GetMinRaise()
}

// MaxBet returns the maximum bet amount for a player
func (g *HoldemGameAdapter) MaxBet(playerID string) (int, error) {
	seat, err := g.holdemGame.table.GetPlayerSeat(playerID)
	if err != nil {
		return 0, err
	}

	maxBet, err := g.holdemGame.table.GetMaxBet(seat.Position)
	if err != nil {
		return 0, err
	}

	return maxBet, nil
}

// ProcessPokerAction processes a poker-specific action
func (g *HoldemGameAdapter) ProcessPokerAction(action poker.PokerAction) error {
	// Convert the PokerAction to a HoldemAction
	holdemAction := NewPokerActionAdapter(action)

	return g.holdemGame.ProcessAction(holdemAction)
}

// IsValidPokerAction checks if a poker action is valid
func (g *HoldemGameAdapter) IsValidPokerAction(action poker.PokerAction) bool {
	// Convert the PokerAction to a HoldemAction
	holdemAction := NewPokerActionAdapter(action)

	return g.holdemGame.IsValidAction(holdemAction)
}

// AllowedPokerActions returns the poker actions that are currently allowed for a player
func (g *HoldemGameAdapter) AllowedPokerActions(playerID string) []poker.PokerActionType {
	allowedActions := g.holdemGame.AllowedActions(playerID)
	pokerActions := make([]poker.PokerActionType, len(allowedActions))

	for i, action := range allowedActions {
		pokerActions[i] = poker.PokerActionType(action)
	}

	return pokerActions
}

// BettingRound returns the current betting round
func (g *HoldemGameAdapter) BettingRound() int {
	return g.holdemGame.table.CurrentBettingRound
}

// AdvanceBettingRound advances to the next betting round
func (g *HoldemGameAdapter) AdvanceBettingRound() error {
	// This is a simplified implementation
	g.holdemGame.table.StartNewBettingRound()
	return nil
}

// DealCommunityCards deals community cards according to the current game state
func (g *HoldemGameAdapter) DealCommunityCards() error {
	// This is a simplified implementation
	// In a real implementation, we would need to check the current state
	// and deal the appropriate number of cards
	return errors.New("not implemented")
}

// EvaluateHands evaluates the hands of all active players
func (g *HoldemGameAdapter) EvaluateHands() error {
	// This is a simplified implementation
	// In a real implementation, we would need to evaluate all hands
	// and determine the winners
	return errors.New("not implemented")
}

// AwardPots awards the pots to the winners
func (g *HoldemGameAdapter) AwardPots() error {
	// This is a simplified implementation
	// In a real implementation, we would need to award the pots
	// to the winners based on hand rankings
	return errors.New("not implemented")
}
