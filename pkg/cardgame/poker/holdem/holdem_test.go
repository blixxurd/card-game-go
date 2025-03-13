package holdem

import (
	"testing"

	"github.com/blixxurd/card-game-go/pkg/cardgame/hand"
	"github.com/blixxurd/card-game-go/pkg/cardgame/player"
	"github.com/blixxurd/card-game-go/pkg/cardgame/table"
)

// MockPlayer implements the player.Player interface for testing
type MockPlayer struct {
	id     string
	name   string
	active bool
	data   map[string]interface{}
}

func NewMockPlayer(id, name string) *MockPlayer {
	return &MockPlayer{
		id:     id,
		name:   name,
		active: true,
		data:   make(map[string]interface{}),
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

func (p *MockPlayer) Hand() hand.Hand {
	return nil
}

func (p *MockPlayer) SetHand(h hand.Hand) {
}

func (p *MockPlayer) Score() int {
	return 0
}

func (p *MockPlayer) SetScore(score int) {
}

func (p *MockPlayer) AddToScore(points int) {
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

func (p *MockPlayer) Clone() player.Player {
	clone := &MockPlayer{
		id:     p.id,
		name:   p.name,
		active: p.active,
		data:   make(map[string]interface{}),
	}

	for k, v := range p.data {
		clone.data[k] = v
	}

	return clone
}

// TestHoldemGameCreation tests the creation of a new HoldemGame
func TestHoldemGameCreation(t *testing.T) {
	game := NewHoldemGame("test", "Test Game", 5, 10)

	if game.ID() != "test" {
		t.Errorf("Expected game ID 'test', got '%s'", game.ID())
	}

	if game.Name() != "Test Game" {
		t.Errorf("Expected game name 'Test Game', got '%s'", game.Name())
	}

	if game.smallBlind != 5 {
		t.Errorf("Expected small blind 5, got %d", game.smallBlind)
	}

	if game.bigBlind != 10 {
		t.Errorf("Expected big blind 10, got %d", game.bigBlind)
	}

	if game.table == nil {
		t.Errorf("Expected table to be initialized")
	}

	if len(game.table.Seats) != 9 {
		t.Errorf("Expected 9 seats, got %d", len(game.table.Seats))
	}
}

// TestAddAndRemovePlayer tests adding and removing players from the game
func TestAddAndRemovePlayer(t *testing.T) {
	game := NewHoldemGame("test", "Test Game", 5, 10)
	player1 := NewMockPlayer("p1", "Player 1")
	player2 := NewMockPlayer("p2", "Player 2")

	// Add players
	err := game.AddPlayer(player1)
	if err != nil {
		t.Errorf("Failed to add player1: %v", err)
	}

	err = game.AddPlayer(player2)
	if err != nil {
		t.Errorf("Failed to add player2: %v", err)
	}

	// Check if players were added to the game
	if len(game.players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(game.players))
	}

	// Check if players were added to the table
	occupiedSeats := 0
	for _, seat := range game.table.Seats {
		if seat.IsOccupied {
			occupiedSeats++
		}
	}

	if occupiedSeats != 2 {
		t.Errorf("Expected 2 occupied seats, got %d", occupiedSeats)
	}

	// Remove a player
	err = game.RemovePlayer("p1")
	if err != nil {
		t.Errorf("Failed to remove player1: %v", err)
	}

	// Check if player was removed from the game
	if len(game.players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(game.players))
	}

	// Check if player was removed from the table
	occupiedSeats = 0
	for _, seat := range game.table.Seats {
		if seat.IsOccupied {
			occupiedSeats++
		}
	}

	if occupiedSeats != 1 {
		t.Errorf("Expected 1 occupied seat, got %d", occupiedSeats)
	}
}

// TestGameStart tests starting a game
func TestGameStart(t *testing.T) {
	game := NewHoldemGame("test", "Test Game", 5, 10)
	player1 := NewMockPlayer("p1", "Player 1")
	player2 := NewMockPlayer("p2", "Player 2")

	// Add players
	err := game.AddPlayer(player1)
	if err != nil {
		t.Errorf("Failed to add player1: %v", err)
	}

	err = game.AddPlayer(player2)
	if err != nil {
		t.Errorf("Failed to add player2: %v", err)
	}

	// Start the game
	err = game.Start()
	if err != nil {
		t.Errorf("Failed to start game: %v", err)
	}

	// Check if the game state is PreFlop
	if game.state != StatePreFlop {
		t.Errorf("Expected game state PreFlop, got %s", game.state)
	}

	// Check if blinds were collected
	if game.pot != 15 { // 5 (small blind) + 10 (big blind)
		t.Errorf("Expected pot 15, got %d", game.pot)
	}

	// Check if the table pot matches
	if game.table.GetTotalPot() != 15 {
		t.Errorf("Expected table pot 15, got %d", game.table.GetTotalPot())
	}
}

// TestAllInAndSidePots tests the handling of all-in scenarios and side pots
func TestAllInAndSidePots(t *testing.T) {
	// Create a table directly instead of using the game
	pokerTable := table.NewTable(9, 10, 20, 100, 1000, table.NoLimit)

	// Create players
	player1 := NewMockPlayer("p1", "Player 1")
	player2 := NewMockPlayer("p2", "Player 2")
	player3 := NewMockPlayer("p3", "Player 3")

	// Add players to the table with specific chip stacks
	err := pokerTable.AddPlayer(player1, 0, 200) // Player 1 has 200 chips
	if err != nil {
		t.Fatalf("Failed to add player1: %v", err)
	}

	err = pokerTable.AddPlayer(player2, 1, 100) // Player 2 has 100 chips (will go all-in)
	if err != nil {
		t.Fatalf("Failed to add player2: %v", err)
	}

	err = pokerTable.AddPlayer(player3, 2, 300) // Player 3 has 300 chips
	if err != nil {
		t.Fatalf("Failed to add player3: %v", err)
	}

	// Activate all seats
	for i := 0; i < 3; i++ {
		seat, _ := pokerTable.GetSeat(i)
		seat.IsActive = true
	}

	// Record initial chip stacks
	seat0, _ := pokerTable.GetSeat(0)
	seat1, _ := pokerTable.GetSeat(1)
	seat2, _ := pokerTable.GetSeat(2)

	initialP1Chips := seat0.ChipStack
	initialP2Chips := seat1.ChipStack
	initialP3Chips := seat2.ChipStack

	t.Logf("Initial chip stacks - P1: %d, P2: %d, P3: %d", initialP1Chips, initialP2Chips, initialP3Chips)

	// Player 1 bets 50
	amount1, _, err := pokerTable.PlaceBet(0, 50)
	if err != nil {
		t.Fatalf("Failed to place bet for player1: %v", err)
	}
	t.Logf("Player 1 bet %d", amount1)
	t.Logf("Current bet: %d, Last raise: %d", pokerTable.CurrentBet, pokerTable.LastRaiseAmount)
	t.Logf("Main pot: %d, Side pots: %d", pokerTable.PotManager.MainPot.Amount, len(pokerTable.PotManager.SidePots))
	t.Logf("Player 1 current bet: %d, chip stack: %d", seat0.CurrentBet, seat0.ChipStack)
	t.Logf("Player 2 current bet: %d, chip stack: %d", seat1.CurrentBet, seat1.ChipStack)
	t.Logf("Player 3 current bet: %d, chip stack: %d", seat2.CurrentBet, seat2.ChipStack)

	// Player 2 goes all-in with 100
	amount2, isAllIn, err := pokerTable.PlaceBet(1, 100)
	if err != nil {
		t.Fatalf("Failed to place bet for player2: %v", err)
	}
	t.Logf("Player 2 bet %d, all-in: %v", amount2, isAllIn)
	t.Logf("Current bet: %d, Last raise: %d", pokerTable.CurrentBet, pokerTable.LastRaiseAmount)
	t.Logf("Main pot: %d, Side pots: %d", pokerTable.PotManager.MainPot.Amount, len(pokerTable.PotManager.SidePots))
	t.Logf("Player 1 current bet: %d, chip stack: %d", seat0.CurrentBet, seat0.ChipStack)
	t.Logf("Player 2 current bet: %d, chip stack: %d", seat1.CurrentBet, seat1.ChipStack)
	t.Logf("Player 3 current bet: %d, chip stack: %d", seat2.CurrentBet, seat2.ChipStack)

	// Verify player 2 is all-in
	if !isAllIn {
		t.Errorf("Expected player2 to be all-in")
	}

	// Player 3 raises to 200
	amount3, _, err := pokerTable.PlaceBet(2, 200)
	if err != nil {
		t.Fatalf("Failed to place bet for player3: %v", err)
	}
	t.Logf("Player 3 bet %d", amount3)
	t.Logf("Current bet: %d, Last raise: %d", pokerTable.CurrentBet, pokerTable.LastRaiseAmount)
	t.Logf("Main pot: %d, Side pots: %d", pokerTable.PotManager.MainPot.Amount, len(pokerTable.PotManager.SidePots))
	t.Logf("Player 1 current bet: %d, chip stack: %d", seat0.CurrentBet, seat0.ChipStack)
	t.Logf("Player 2 current bet: %d, chip stack: %d", seat1.CurrentBet, seat1.ChipStack)
	t.Logf("Player 3 current bet: %d, chip stack: %d", seat2.CurrentBet, seat2.ChipStack)

	// Player 1 calls the raise
	amount4, _, err := pokerTable.PlaceBet(0, 150) // Additional 150 to match the 200 total
	if err != nil {
		t.Fatalf("Failed to place bet for player1: %v", err)
	}
	t.Logf("Player 1 called with %d", amount4)
	t.Logf("Current bet: %d, Last raise: %d", pokerTable.CurrentBet, pokerTable.LastRaiseAmount)
	t.Logf("Main pot: %d, Side pots: %d", pokerTable.PotManager.MainPot.Amount, len(pokerTable.PotManager.SidePots))
	t.Logf("Player 1 current bet: %d, chip stack: %d", seat0.CurrentBet, seat0.ChipStack)
	t.Logf("Player 2 current bet: %d, chip stack: %d", seat1.CurrentBet, seat1.ChipStack)
	t.Logf("Player 3 current bet: %d, chip stack: %d", seat2.CurrentBet, seat2.ChipStack)

	// Start a new betting round to ensure all bets are processed
	pokerTable.StartNewBettingRound()
	t.Logf("After starting new betting round:")
	t.Logf("Current bet: %d, Last raise: %d", pokerTable.CurrentBet, pokerTable.LastRaiseAmount)
	t.Logf("Player 1 current bet: %d, chip stack: %d", seat0.CurrentBet, seat0.ChipStack)
	t.Logf("Player 2 current bet: %d, chip stack: %d", seat1.CurrentBet, seat1.ChipStack)
	t.Logf("Player 3 current bet: %d, chip stack: %d", seat2.CurrentBet, seat2.ChipStack)

	// Verify side pot was created
	if len(pokerTable.PotManager.SidePots) == 0 {
		t.Errorf("Expected at least one side pot to be created")
	} else {
		t.Logf("Number of side pots: %d", len(pokerTable.PotManager.SidePots))
	}

	// Log the current pot amounts
	t.Logf("Main pot: %d", pokerTable.PotManager.MainPot.Amount)
	for i, pot := range pokerTable.PotManager.SidePots {
		t.Logf("Side pot %d: %d", i, pot.Amount)

		// Log eligible players for this side pot
		eligiblePlayers := pot.GetEligiblePlayers()
		t.Logf("Side pot %d eligible players: %v", i, eligiblePlayers)
	}

	// Verify pot amounts - the exact amounts depend on how the side pots are created
	totalPot := pokerTable.GetTotalPot()
	expectedTotalPot := amount1 + amount2 + amount3 + amount4
	if totalPot != expectedTotalPot {
		t.Errorf("Expected total pot to be %d, got %d", expectedTotalPot, totalPot)
	}

	// Verify player eligibility for main pot
	mainPotEligiblePlayers := pokerTable.PotManager.MainPot.GetEligiblePlayers()
	t.Logf("Main pot eligible players: %v", mainPotEligiblePlayers)

	if !pokerTable.PotManager.MainPot.IsPlayerEligible("p2") {
		t.Errorf("Expected player 2 to be eligible for main pot")
	}

	// Verify player eligibility for side pot (if any)
	if len(pokerTable.PotManager.SidePots) > 0 && pokerTable.PotManager.SidePots[0].IsPlayerEligible("p2") {
		t.Errorf("Expected player 2 to NOT be eligible for side pot")
	}

	// Create a map of winners for testing pot distribution
	winners := make(map[int]string)
	winners[-1] = "p3" // Main pot goes to player 3
	for i := range pokerTable.PotManager.SidePots {
		winners[i] = "p3" // Side pots go to player 3
	}

	// Award pots to winners
	pokerTable.AwardPotsToWinners(winners)

	// Verify chip stacks after pot distribution
	finalP1Chips := seat0.ChipStack
	finalP2Chips := seat1.ChipStack
	finalP3Chips := seat2.ChipStack

	t.Logf("Final chip stacks - P1: %d, P2: %d, P3: %d", finalP1Chips, finalP2Chips, finalP3Chips)

	// Player 1 should have lost their bet
	expectedP1Chips := initialP1Chips - (amount1 + amount4)
	if finalP1Chips != expectedP1Chips {
		t.Errorf("Expected player 1 to have %d chips, got %d", expectedP1Chips, finalP1Chips)
	}

	// Player 2 should have lost their all-in amount
	expectedP2Chips := initialP2Chips - amount2
	if finalP2Chips != expectedP2Chips {
		t.Errorf("Expected player 2 to have %d chips, got %d", expectedP2Chips, finalP2Chips)
	}

	// Player 3 should have won all pots
	expectedP3Chips := initialP3Chips - amount3 + totalPot
	if finalP3Chips != expectedP3Chips {
		t.Errorf("Expected player 3 to have %d chips, got %d", expectedP3Chips, finalP3Chips)
	}

	// Verify the pots are reset after awarding
	if pokerTable.GetTotalPot() != 0 {
		t.Errorf("Expected total pot to be 0 after awarding, got %d", pokerTable.GetTotalPot())
	}
}
