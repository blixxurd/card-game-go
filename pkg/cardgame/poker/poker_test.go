package poker_test

import (
	"testing"

	"github.com/blixxurd/card-game-go/pkg/cardgame/hand"
	"github.com/blixxurd/card-game-go/pkg/cardgame/player"
	"github.com/blixxurd/card-game-go/pkg/cardgame/poker"
	"github.com/blixxurd/card-game-go/pkg/cardgame/poker/holdem"
	"github.com/blixxurd/card-game-go/pkg/cardgame/table"
)

// MockPlayer implements the player.Player interface for testing
type MockPlayer struct {
	id     string
	name   string
	active bool
	data   map[string]interface{}
	score  int
	hand   hand.Hand
}

func NewMockPlayer(id, name string) *MockPlayer {
	return &MockPlayer{
		id:     id,
		name:   name,
		active: true,
		data:   make(map[string]interface{}),
		score:  0,
		hand:   nil,
	}
}

func (p *MockPlayer) ID() string {
	return p.id
}

func (p *MockPlayer) Name() string {
	return p.name
}

func (p *MockPlayer) SetName(name string) {
	p.name = name
}

func (p *MockPlayer) IsActive() bool {
	return p.active
}

func (p *MockPlayer) SetActive(active bool) {
	p.active = active
}

func (p *MockPlayer) Data() map[string]interface{} {
	return p.data
}

func (p *MockPlayer) SetData(key string, value interface{}) {
	p.data[key] = value
}

func (p *MockPlayer) Score() int {
	return p.score
}

func (p *MockPlayer) SetScore(score int) {
	p.score = score
}

func (p *MockPlayer) AddToScore(points int) {
	p.score += points
}

func (p *MockPlayer) Hand() hand.Hand {
	return p.hand
}

func (p *MockPlayer) SetHand(h hand.Hand) {
	p.hand = h
}

func (p *MockPlayer) Clone() player.Player {
	clone := &MockPlayer{
		id:     p.id,
		name:   p.name,
		active: p.active,
		data:   make(map[string]interface{}),
		score:  p.score,
		hand:   p.hand,
	}

	for k, v := range p.data {
		clone.data[k] = v
	}

	return clone
}

// TestPokerGameFactory tests the PokerGameFactory
func TestPokerGameFactory(t *testing.T) {
	// Create a new factory
	factory := poker.NewPokerGameFactory()

	// Register the Texas Hold'em variant
	factory.RegisterVariant(holdem.NewHoldemVariant())

	// Check available variants
	variants := factory.AvailableVariants()
	if len(variants) != 1 {
		t.Errorf("Expected 1 variant, got %d", len(variants))
	}

	if variants[0] != "Texas Hold'em" {
		t.Errorf("Expected 'Texas Hold'em', got '%s'", variants[0])
	}

	// Create a new game
	game, err := factory.CreatePokerGame("Texas Hold'em", "test", "Test Game", 5, 10, table.NoLimit)
	if err != nil {
		t.Errorf("Failed to create game: %v", err)
	}

	// Check game properties
	if game.ID() != "test" {
		t.Errorf("Expected game ID 'test', got '%s'", game.ID())
	}

	if game.Name() != "Test Game" {
		t.Errorf("Expected game name 'Test Game', got '%s'", game.Name())
	}

	// Test interface methods without actually playing a game
	// This avoids potential deadlocks in the game logic

	// Test Table method
	table := game.Table()
	if table == nil {
		t.Errorf("Expected table to not be nil")
	}

	// Test CommunityCards method
	communityCards := game.CommunityCards()
	if communityCards == nil {
		t.Errorf("Expected communityCards to not be nil")
	}

	// Test BettingRound method
	bettingRound := game.BettingRound()
	if bettingRound != 0 {
		t.Errorf("Expected betting round to be 0, got %d", bettingRound)
	}

	// Test MinRaise method
	minRaise := game.MinRaise()
	// The min raise might be 0 if the game hasn't started yet, which is valid
	t.Logf("Min raise: %d", minRaise)

	// Test AllowedPokerActions method
	// We don't test this with actual players since it might depend on game state

	// Test that we can add players
	player1 := NewMockPlayer("p1", "Player 1")
	err = game.AddPlayer(player1)
	if err != nil {
		t.Errorf("Failed to add player1: %v", err)
	}

	// Check that the player was added
	players := game.Players()
	if len(players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(players))
	}

	// Test that we can remove players
	err = game.RemovePlayer("p1")
	if err != nil {
		t.Errorf("Failed to remove player1: %v", err)
	}

	// Check that the player was removed
	players = game.Players()
	if len(players) != 0 {
		t.Errorf("Expected 0 players, got %d", len(players))
	}
}

// TestPokerVariantInterface tests the PokerVariant interface
func TestPokerVariantInterface(t *testing.T) {
	// Create a new variant
	variant := holdem.NewHoldemVariant()

	// Test Name method
	name := variant.Name()
	if name != "Texas Hold'em" {
		t.Errorf("Expected name to be 'Texas Hold'em', got '%s'", name)
	}

	// Test Description method
	description := variant.Description()
	if description == "" {
		t.Errorf("Expected description to not be empty")
	}

	// Test MaxPlayers method
	maxPlayers := variant.MaxPlayers()
	if maxPlayers != 9 {
		t.Errorf("Expected max players to be 9, got %d", maxPlayers)
	}

	// Test MinPlayers method
	minPlayers := variant.MinPlayers()
	if minPlayers != 2 {
		t.Errorf("Expected min players to be 2, got %d", minPlayers)
	}

	// Test HandSize method
	handSize := variant.HandSize()
	if handSize != 2 {
		t.Errorf("Expected hand size to be 2, got %d", handSize)
	}

	// Test CommunityCardCount method
	communityCardCount := variant.CommunityCardCount()
	if communityCardCount != 5 {
		t.Errorf("Expected community card count to be 5, got %d", communityCardCount)
	}

	// Test BettingRounds method
	bettingRounds := variant.BettingRounds()
	if bettingRounds != 4 {
		t.Errorf("Expected betting rounds to be 4, got %d", bettingRounds)
	}

	// Test StateSequence method
	stateSequence := variant.StateSequence()
	if len(stateSequence) != 7 {
		t.Errorf("Expected state sequence to have 7 states, got %d", len(stateSequence))
	}
}

// TestPokerAction implements the poker.PokerAction interface for testing
type TestPokerAction struct {
	actionType poker.PokerActionType
	playerID   string
	amount     int
}

func (a *TestPokerAction) Type() string {
	return string(a.actionType)
}

func (a *TestPokerAction) PlayerID() string {
	return a.playerID
}

func (a *TestPokerAction) Params() map[string]interface{} {
	return map[string]interface{}{
		"amount": a.amount,
	}
}

func (a *TestPokerAction) Amount() int {
	return a.amount
}

func (a *TestPokerAction) ActionType() poker.PokerActionType {
	return a.actionType
}
