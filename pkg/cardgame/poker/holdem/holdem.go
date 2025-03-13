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
	"github.com/blixxurd/card-game-go/pkg/cardgame/table"
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
	activePlayers    map[string]bool // Track which players are still in the hand
	smallBlind       int
	bigBlind         int
	pot              int
	winners          []player.Player
	mutex            sync.RWMutex
	customData       map[string]interface{}

	// New field for the Table entity
	table *table.Table
}

// NewHoldemGame creates a new HoldemGame
func NewHoldemGame(id, name string, smallBlind, bigBlind int) *HoldemGame {
	// Create a new table with 9 seats (standard poker table)
	pokerTable := table.NewTable(9, smallBlind, bigBlind, 100, 10000, table.NoLimit)

	return &HoldemGame{
		id:               id,
		name:             name,
		state:            StateInitial,
		players:          make([]player.Player, 0),
		currentPlayerIdx: 0,
		mainDeck:         deck.NewStandardDeck(),
		communityCards:   make([]card.Card, 0),
		playerHands:      make(map[string]hand.Hand),
		activePlayers:    make(map[string]bool),
		smallBlind:       smallBlind,
		bigBlind:         bigBlind,
		pot:              0,
		winners:          make([]player.Player, 0),
		customData:       make(map[string]interface{}),
		table:            pokerTable,
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

	// Initialize all players as active
	g.activePlayers = make(map[string]bool)
	for _, p := range g.players {
		g.activePlayers[p.ID()] = true
	}

	// Start a new hand on the table
	g.table.StartNewHand()

	// Activate all occupied seats
	for _, seat := range g.table.Seats {
		if seat.IsOccupied {
			seat.IsActive = true
		}
	}

	// Collect blinds
	err := g.table.CollectBlinds()
	if err != nil {
		return fmt.Errorf("failed to collect blinds: %w", err)
	}

	// Update the pot from the table (for backward compatibility)
	g.pot = g.table.GetTotalPot()

	// Deal hole cards to each player
	g.dealHoleCards()

	// Set the game state to pre-flop
	g.state = StatePreFlop

	return nil
}

// AddPlayer adds a player to the game with a default chip stack of 1000
func (g *HoldemGame) AddPlayer(p player.Player) error {
	return g.AddPlayerWithChips(p, 1000)
}

// AddPlayerWithChips adds a player to the game with a specified chip stack
func (g *HoldemGame) AddPlayerWithChips(p player.Player, chipStack int) error {
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

	// Add player to the game's player list (for backward compatibility)
	g.players = append(g.players, p)

	// Find an empty seat at the table
	for i := 0; i < len(g.table.Seats); i++ {
		seat, err := g.table.GetSeat(i)
		if err != nil {
			continue
		}

		if !seat.IsOccupied {
			// Add player to the table with the specified chip stack
			err := g.table.AddPlayer(p, i, chipStack)
			if err != nil {
				// If there's an error adding to the table, remove from players list
				g.players = g.players[:len(g.players)-1]
				return fmt.Errorf("failed to add player to table: %w", err)
			}
			break
		}
	}

	return nil
}

// RemovePlayer removes a player from the game
func (g *HoldemGame) RemovePlayer(playerID string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.state != StateInitial {
		return errors.New("cannot remove player after game has started")
	}

	// Remove player from the game's player list
	playerIndex := -1
	for i, p := range g.players {
		if p.ID() == playerID {
			playerIndex = i
			break
		}
	}

	if playerIndex >= 0 {
		g.players = append(g.players[:playerIndex], g.players[playerIndex+1:]...)
	} else {
		return errors.New("player not found")
	}

	// Remove player from the table
	err := g.table.RemovePlayer(playerID)
	if err != nil {
		// If there's an error removing from the table, but we've already removed from players list,
		// we should log this but not fail the operation
		fmt.Printf("Warning: Failed to remove player from table: %v\n", err)
	}

	return nil
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

	// Skip players who have folded
	for !g.activePlayers[g.players[g.currentPlayerIdx].ID()] && g.state != StateComplete {
		g.currentPlayerIdx = (g.currentPlayerIdx + 1) % len(g.players)
	}

	// Check if we've completed a round (all players have acted)
	if g.currentPlayerIdx == 0 {
		// Advance the game state
		switch g.state {
		case StatePreFlop:
			g.DealFlop()
			g.state = StateFlop
		case StateFlop:
			g.DealTurn()
			g.state = StateTurn
		case StateTurn:
			g.DealRiver()
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

	// Find the player's seat at the table
	var playerSeat *table.Seat
	var playerSeatPosition int
	var err error

	for i, seat := range g.table.Seats {
		if seat.IsOccupied && seat.Player.ID() == holdemAction.playerID {
			playerSeat = seat
			playerSeatPosition = i
			break
		}
	}

	if playerSeat == nil {
		return errors.New("player not found at table")
	}

	// Get the current bet amount that the player needs to call
	currentBet := g.table.CurrentBet

	// Get the player's current bet in this round
	playerCurrentBet := playerSeat.CurrentBet

	// Calculate how much more the player needs to call
	callAmount := currentBet - playerCurrentBet

	switch holdemAction.actionType {
	case ActionCheck:
		// Check is only valid if there's no bet to call
		if callAmount > 0 {
			return errors.New("cannot check when there's a bet to call")
		}
		// No action needed for check
	case ActionCall:
		// Call matches the current bet
		betAmount := callAmount

		// Use the table to place the bet
		amount, isAllIn, err := g.table.PlaceBet(playerSeatPosition, betAmount)
		if err != nil {
			return fmt.Errorf("failed to place bet: %w", err)
		}

		// Update the pot from the table (for backward compatibility)
		g.pot = g.table.GetTotalPot()

		// If the player went all-in with a call, it might create a side pot
		if isAllIn && amount < betAmount {
			// The player couldn't match the full bet amount, so they're all-in for less
			// This is handled by the Table's PlaceBet method, which creates a side pot
			g.SetData("last_action", fmt.Sprintf("Player %s called %d and is all-in", holdemAction.playerID, amount))
		} else {
			g.SetData("last_action", fmt.Sprintf("Player %s called %d", holdemAction.playerID, amount))
		}
	case ActionBet:
		// Bet is only valid if there's no bet to call
		if callAmount > 0 {
			return errors.New("cannot bet when there's a bet to call")
		}

		// Validate minimum bet
		minBet := g.table.GetMinRaise()
		if holdemAction.amount < minBet {
			return fmt.Errorf("bet must be at least %d", minBet)
		}

		// Use the table to place the bet
		amount, isAllIn, err := g.table.PlaceBet(playerSeatPosition, holdemAction.amount)
		if err != nil {
			return fmt.Errorf("failed to place bet: %w", err)
		}

		// Update the pot from the table (for backward compatibility)
		g.pot = g.table.GetTotalPot()

		if isAllIn {
			g.SetData("last_action", fmt.Sprintf("Player %s bet %d and is all-in", holdemAction.playerID, amount))
		} else {
			g.SetData("last_action", fmt.Sprintf("Player %s bet %d", holdemAction.playerID, amount))
		}
	case ActionRaise:
		// Raise is only valid if there's a bet to call
		if callAmount <= 0 {
			return errors.New("cannot raise when there's no bet to call")
		}

		// Calculate the total amount (call + raise)
		totalAmount := callAmount + holdemAction.amount

		// Validate minimum raise
		minRaise := g.table.GetMinRaise()
		if holdemAction.amount < minRaise {
			return fmt.Errorf("raise must be at least %d", minRaise)
		}

		// Use the table to place the bet
		amount, isAllIn, err := g.table.PlaceBet(playerSeatPosition, totalAmount)
		if err != nil {
			return fmt.Errorf("failed to place bet: %w", err)
		}

		// Update the pot from the table (for backward compatibility)
		g.pot = g.table.GetTotalPot()

		if isAllIn {
			g.SetData("last_action", fmt.Sprintf("Player %s raised to %d and is all-in", holdemAction.playerID, amount))
		} else {
			g.SetData("last_action", fmt.Sprintf("Player %s raised to %d", holdemAction.playerID, amount))
		}
	case ActionFold:
		// Use the table to fold the player
		err = g.table.FoldPlayer(playerSeatPosition)
		if err != nil {
			return fmt.Errorf("failed to fold player: %w", err)
		}

		// Mark the player as inactive (for backward compatibility)
		g.activePlayers[holdemAction.playerID] = false

		g.SetData("last_action", fmt.Sprintf("Player %s folded", holdemAction.playerID))

		// Check if there's only one active player left
		activeCount := 0
		var lastActivePlayer player.Player

		for _, p := range g.players {
			if g.activePlayers[p.ID()] {
				activeCount++
				lastActivePlayer = p
			}
		}

		// If only one player remains, they win automatically
		if activeCount == 1 {
			g.winners = []player.Player{lastActivePlayer}
			g.state = StateComplete

			// Award the pot to the winner
			winners := map[int]string{
				-1: lastActivePlayer.ID(), // Main pot goes to the last active player
			}

			// Award side pots if they exist
			for i := range g.table.PotManager.SidePots {
				// Check if the winner is eligible for this side pot
				if g.table.PotManager.SidePots[i].IsPlayerEligible(lastActivePlayer.ID()) {
					winners[i] = lastActivePlayer.ID()
				}
			}

			g.table.AwardPotsToWinners(winners)

			g.SetData("winner", lastActivePlayer.ID())
			g.SetData("win_type", "last_player_standing")

			return nil
		}

		// If no players remain, end the game with no winner
		if activeCount == 0 {
			g.winners = []player.Player{}
			g.state = StateComplete
			return nil
		}
	default:
		return errors.New("invalid action type")
	}

	// Move to the next player
	g.currentPlayerIdx = (g.currentPlayerIdx + 1) % len(g.players)

	// Skip players who have folded or are all-in
	for (!g.activePlayers[g.players[g.currentPlayerIdx].ID()] || g.isPlayerAllIn(g.players[g.currentPlayerIdx].ID())) && g.state != StateComplete {
		g.currentPlayerIdx = (g.currentPlayerIdx + 1) % len(g.players)

		// If we've gone all the way around and everyone is either folded or all-in,
		// we need to deal the remaining community cards and go to showdown
		if g.currentPlayerIdx == 0 {
			// Deal remaining community cards
			switch g.state {
			case StatePreFlop:
				g.DealFlop()
				g.DealTurn()
				g.DealRiver()
				g.evaluateHands()
				g.state = StateShowdown
				g.awardPotsToWinners()
				g.state = StateComplete
				return nil
			case StateFlop:
				g.DealTurn()
				g.DealRiver()
				g.evaluateHands()
				g.state = StateShowdown
				g.awardPotsToWinners()
				g.state = StateComplete
				return nil
			case StateTurn:
				g.DealRiver()
				g.evaluateHands()
				g.state = StateShowdown
				g.awardPotsToWinners()
				g.state = StateComplete
				return nil
			case StateRiver:
				g.evaluateHands()
				g.state = StateShowdown
				g.awardPotsToWinners()
				g.state = StateComplete
				return nil
			}
		}
	}

	// Check if we've completed a round (all players have acted)
	// This happens when we return to the first player who still needs to act
	// or when everyone has either folded or gone all-in
	if g.hasRoundCompleted() {
		// Start a new betting round on the table
		g.table.StartNewBettingRound()

		// Advance the game state
		switch g.state {
		case StatePreFlop:
			g.DealFlop()
			g.state = StateFlop
		case StateFlop:
			g.DealTurn()
			g.state = StateTurn
		case StateTurn:
			g.DealRiver()
			g.state = StateRiver
		case StateRiver:
			g.evaluateHands()
			g.state = StateShowdown
			g.awardPotsToWinners()
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
	g.activePlayers = make(map[string]bool)
	g.pot = 0
	g.winners = make([]player.Player, 0)

	// Reset the table for a new hand
	// We don't use StartNewHand here because that would advance the dealer button
	// and we want to keep the same dealer position for the next hand
	for _, seat := range g.table.Seats {
		if seat.IsOccupied {
			seat.Reset()
		}
	}

	// Reset the pot manager
	g.table.PotManager.Reset()

	// Reset betting state
	g.table.CurrentBettingRound = 0
	g.table.CurrentBet = 0
	g.table.LastRaiseAmount = 0

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
func (g *HoldemGame) DealFlop() {
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
func (g *HoldemGame) DealTurn() {
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
func (g *HoldemGame) DealRiver() {
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
		// Skip players who have folded
		if !g.activePlayers[p.ID()] {
			continue
		}

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

// isPlayerAllIn checks if a player is all-in
func (g *HoldemGame) isPlayerAllIn(playerID string) bool {
	for _, seat := range g.table.Seats {
		if seat.IsOccupied && seat.Player.ID() == playerID {
			return seat.IsAllIn
		}
	}
	return false
}

// hasRoundCompleted checks if the current betting round has completed
func (g *HoldemGame) hasRoundCompleted() bool {
	// If we're back to the first player who still needs to act, the round is complete
	if g.currentPlayerIdx == 0 {
		return true
	}

	// Count active players who aren't all-in
	activePlayersNotAllIn := 0
	for _, p := range g.players {
		if g.activePlayers[p.ID()] && !g.isPlayerAllIn(p.ID()) {
			activePlayersNotAllIn++
		}
	}

	// If there are fewer than 2 active players who aren't all-in, the round is complete
	if activePlayersNotAllIn < 2 {
		return true
	}

	// Check if all active players have bet the same amount or are all-in
	currentBet := g.table.CurrentBet
	for _, seat := range g.table.Seats {
		if seat.IsOccupied && seat.IsActive && !seat.HasFolded && !seat.IsAllIn {
			if seat.CurrentBet != currentBet {
				return false
			}
		}
	}

	return true
}

// awardPotsToWinners awards the pots to the winners based on hand evaluation
func (g *HoldemGame) awardPotsToWinners() {
	if len(g.winners) == 0 {
		return
	}

	// Create a map to track which pots each winner is eligible for
	winners := make(map[int][]string)

	// Award main pot
	// If there are multiple winners, they split the main pot
	mainPotWinners := make([]string, 0, len(g.winners))
	for _, winner := range g.winners {
		mainPotWinners = append(mainPotWinners, winner.ID())
	}
	winners[-1] = mainPotWinners

	if len(mainPotWinners) > 1 {
		g.SetData("split_pot", true)
	}

	// Award side pots
	// For each side pot, determine which winners are eligible
	for i, sidePot := range g.table.PotManager.SidePots {
		sidePotWinners := make([]string, 0)
		for _, winner := range g.winners {
			if sidePot.IsPlayerEligible(winner.ID()) {
				sidePotWinners = append(sidePotWinners, winner.ID())
			}
		}
		if len(sidePotWinners) > 0 {
			winners[i] = sidePotWinners
		}
	}

	// Award the pots to the winners
	g.table.AwardPotsToMultipleWinners(winners)

	// Store the winners in game data
	winnerIDs := make([]string, 0, len(g.winners))
	for _, w := range g.winners {
		winnerIDs = append(winnerIDs, w.ID())
	}
	g.SetData("winners", winnerIDs)
}
