package table

import (
	"testing"

	"github.com/blixxurd/card-game-go/pkg/cardgame/hand"
	"github.com/blixxurd/card-game-go/pkg/cardgame/player"
)

// MockPlayer implements the player.Player interface for testing
type MockPlayer struct {
	id     string
	name   string
	hand   hand.Hand
	score  int
	active bool
	data   map[string]interface{}
}

func NewMockPlayer(id, name string) *MockPlayer {
	return &MockPlayer{
		id:     id,
		name:   name,
		hand:   hand.NewStandardHand(),
		score:  0,
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
	return p.hand
}

func (p *MockPlayer) SetHand(h hand.Hand) {
	p.hand = h
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

func (p *MockPlayer) IsActive() bool {
	return p.active
}

func (p *MockPlayer) SetActive(active bool) {
	p.active = active
}

func (p *MockPlayer) Data() map[string]interface{} {
	// Create a copy to avoid concurrent modification issues
	result := make(map[string]interface{}, len(p.data))
	for k, v := range p.data {
		result[k] = v
	}
	return result
}

func (p *MockPlayer) SetData(key string, value interface{}) {
	p.data[key] = value
}

func (p *MockPlayer) Clone() player.Player {
	clone := &MockPlayer{
		id:     p.id,
		name:   p.name,
		score:  p.score,
		active: p.active,
		data:   make(map[string]interface{}, len(p.data)),
	}

	if p.hand != nil {
		clone.hand = p.hand.Clone()
	}

	// Copy the data map
	for k, v := range p.data {
		clone.data[k] = v
	}

	return clone
}

// TestTableCreation tests the creation of a new table
func TestTableCreation(t *testing.T) {
	table := NewTable(9, 5, 10, 100, 1000, NoLimit)

	if len(table.Seats) != 9 {
		t.Errorf("Expected 9 seats, got %d", len(table.Seats))
	}

	if table.SmallBlindAmount != 5 {
		t.Errorf("Expected small blind of 5, got %d", table.SmallBlindAmount)
	}

	if table.BigBlindAmount != 10 {
		t.Errorf("Expected big blind of 10, got %d", table.BigBlindAmount)
	}

	if table.MinBuyIn != 100 {
		t.Errorf("Expected min buy-in of 100, got %d", table.MinBuyIn)
	}

	if table.MaxBuyIn != 1000 {
		t.Errorf("Expected max buy-in of 1000, got %d", table.MaxBuyIn)
	}

	if table.BettingRules.GetStructureType() != NoLimit {
		t.Errorf("Expected NoLimit betting structure, got %v", table.BettingRules.GetStructureType())
	}
}

// TestAddPlayer tests adding a player to the table
func TestAddPlayer(t *testing.T) {
	table := NewTable(9, 5, 10, 100, 1000, NoLimit)
	player1 := NewMockPlayer("p1", "Player 1")

	// Add player to valid seat
	err := table.AddPlayer(player1, 0, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	// Check if seat is occupied
	seat, err := table.GetSeat(0)
	if err != nil {
		t.Errorf("Failed to get seat: %v", err)
	}
	if !seat.IsOccupied {
		t.Errorf("Expected seat to be occupied")
	}
	if seat.Player.ID() != "p1" {
		t.Errorf("Expected player ID p1, got %s", seat.Player.ID())
	}
	if seat.ChipStack != 200 {
		t.Errorf("Expected chip stack of 200, got %d", seat.ChipStack)
	}

	// Try to add player to invalid seat
	err = table.AddPlayer(NewMockPlayer("p2", "Player 2"), 10, 200)
	if err == nil {
		t.Errorf("Expected error for invalid seat position")
	}

	// Try to add player to occupied seat
	err = table.AddPlayer(NewMockPlayer("p3", "Player 3"), 0, 200)
	if err == nil {
		t.Errorf("Expected error for occupied seat")
	}

	// Try to add player with insufficient chips
	err = table.AddPlayer(NewMockPlayer("p4", "Player 4"), 1, 50)
	if err == nil {
		t.Errorf("Expected error for insufficient chips")
	}

	// Try to add player with too many chips
	err = table.AddPlayer(NewMockPlayer("p5", "Player 5"), 1, 1500)
	if err == nil {
		t.Errorf("Expected error for too many chips")
	}
}

// TestRemovePlayer tests removing a player from the table
func TestRemovePlayer(t *testing.T) {
	table := NewTable(9, 5, 10, 100, 1000, NoLimit)
	player1 := NewMockPlayer("p1", "Player 1")

	// Add player
	err := table.AddPlayer(player1, 0, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	// Remove player
	err = table.RemovePlayer("p1")
	if err != nil {
		t.Errorf("Failed to remove player: %v", err)
	}

	// Check if seat is vacant
	seat, err := table.GetSeat(0)
	if err != nil {
		t.Errorf("Failed to get seat: %v", err)
	}
	if seat.IsOccupied {
		t.Errorf("Expected seat to be vacant")
	}

	// Try to remove non-existent player
	err = table.RemovePlayer("p2")
	if err == nil {
		t.Errorf("Expected error for non-existent player")
	}
}

// TestPlaceBet tests placing bets
func TestPlaceBet(t *testing.T) {
	table := NewTable(9, 5, 10, 100, 1000, NoLimit)
	player1 := NewMockPlayer("p1", "Player 1")
	player2 := NewMockPlayer("p2", "Player 2")

	// Add players
	err := table.AddPlayer(player1, 0, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}
	err = table.AddPlayer(player2, 1, 100)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	// Make sure the seats are active
	seat0, _ := table.GetSeat(0)
	seat0.IsActive = true
	seat1, _ := table.GetSeat(1)
	seat1.IsActive = true

	// Place bet
	amount, allIn, err := table.PlaceBet(0, 50)
	if err != nil {
		t.Errorf("Failed to place bet: %v", err)
	}
	if amount != 50 {
		t.Errorf("Expected bet amount of 50, got %d", amount)
	}
	if allIn {
		t.Errorf("Expected not all-in")
	}

	// Check current bet
	if table.CurrentBet != 50 {
		t.Errorf("Expected current bet of 50, got %d", table.CurrentBet)
	}

	// Check player's chip stack
	seat, _ := table.GetSeat(0)
	if seat.ChipStack != 150 {
		t.Errorf("Expected chip stack of 150, got %d", seat.ChipStack)
	}

	// Place all-in bet
	amount, allIn, err = table.PlaceBet(1, 150)
	if err != nil {
		t.Errorf("Failed to place bet: %v", err)
	}
	if amount != 100 {
		t.Errorf("Expected bet amount of 100, got %d", amount)
	}
	if !allIn {
		t.Errorf("Expected all-in")
	}

	// Check current bet
	if table.CurrentBet != 100 {
		t.Errorf("Expected current bet of 100, got %d", table.CurrentBet)
	}

	// Check player's chip stack
	seat, _ = table.GetSeat(1)
	if seat.ChipStack != 0 {
		t.Errorf("Expected chip stack of 0, got %d", seat.ChipStack)
	}
	if !seat.IsAllIn {
		t.Errorf("Expected player to be all-in")
	}

	// Check pot
	if table.GetTotalPot() != 150 {
		t.Errorf("Expected total pot of 150, got %d", table.GetTotalPot())
	}
}

// TestFoldPlayer tests folding
func TestFoldPlayer(t *testing.T) {
	table := NewTable(9, 5, 10, 100, 1000, NoLimit)
	player1 := NewMockPlayer("p1", "Player 1")

	// Add player
	err := table.AddPlayer(player1, 0, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	// Make sure the seat is active
	seat, _ := table.GetSeat(0)
	seat.IsActive = true

	// Fold player
	err = table.FoldPlayer(0)
	if err != nil {
		t.Errorf("Failed to fold player: %v", err)
	}

	// Check if player is folded
	if !seat.HasFolded {
		t.Errorf("Expected player to be folded")
	}

	// Try to fold already folded player
	err = table.FoldPlayer(0)
	if err == nil {
		t.Errorf("Expected error for folding already folded player")
	}
}

// TestCollectBlinds tests collecting blinds
func TestCollectBlinds(t *testing.T) {
	table := NewTable(9, 5, 10, 100, 1000, NoLimit)
	player1 := NewMockPlayer("p1", "Player 1")
	player2 := NewMockPlayer("p2", "Player 2")
	player3 := NewMockPlayer("p3", "Player 3")

	// Add players
	err := table.AddPlayer(player1, 0, 200) // Dealer
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}
	err = table.AddPlayer(player2, 1, 200) // Small blind
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}
	err = table.AddPlayer(player3, 2, 200) // Big blind
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	// Make sure the seats are active
	for i := 0; i < 3; i++ {
		seat, _ := table.GetSeat(i)
		seat.IsActive = true
	}

	// Collect blinds
	err = table.CollectBlinds()
	if err != nil {
		t.Errorf("Failed to collect blinds: %v", err)
	}

	// Check small blind
	seat, _ := table.GetSeat(1)
	if seat.ChipStack != 195 {
		t.Errorf("Expected small blind chip stack of 195, got %d", seat.ChipStack)
	}
	if seat.CurrentBet != 5 {
		t.Errorf("Expected small blind current bet of 5, got %d", seat.CurrentBet)
	}

	// Check big blind
	seat, _ = table.GetSeat(2)
	if seat.ChipStack != 190 {
		t.Errorf("Expected big blind chip stack of 190, got %d", seat.ChipStack)
	}
	if seat.CurrentBet != 10 {
		t.Errorf("Expected big blind current bet of 10, got %d", seat.CurrentBet)
	}

	// Check current bet
	if table.CurrentBet != 10 {
		t.Errorf("Expected current bet of 10, got %d", table.CurrentBet)
	}

	// Check pot
	if table.GetTotalPot() != 15 {
		t.Errorf("Expected total pot of 15, got %d", table.GetTotalPot())
	}
}

// TestAdvanceDealerButton tests advancing the dealer button
func TestAdvanceDealerButton(t *testing.T) {
	table := NewTable(9, 5, 10, 100, 1000, NoLimit)
	player1 := NewMockPlayer("p1", "Player 1")
	player2 := NewMockPlayer("p2", "Player 2")
	player3 := NewMockPlayer("p3", "Player 3")

	// Add players
	err := table.AddPlayer(player1, 0, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}
	err = table.AddPlayer(player2, 2, 200) // Skip seat 1
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}
	err = table.AddPlayer(player3, 4, 200) // Skip seat 3
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	// Make sure the seats are active
	seat0, _ := table.GetSeat(0)
	seat0.IsActive = true
	seat2, _ := table.GetSeat(2)
	seat2.IsActive = true
	seat4, _ := table.GetSeat(4)
	seat4.IsActive = true

	// Initial positions
	if table.DealerPosition != 0 {
		t.Errorf("Expected dealer position 0, got %d", table.DealerPosition)
	}

	// Force the initial blind positions to match our expectations
	table.SmallBlindPosition = 2
	table.BigBlindPosition = 4

	// Advance dealer button
	table.AdvanceDealerButton()

	// Check new positions
	if table.DealerPosition != 2 {
		t.Errorf("Expected dealer position 2, got %d", table.DealerPosition)
	}
	if table.SmallBlindPosition != 4 {
		t.Errorf("Expected small blind position 4, got %d", table.SmallBlindPosition)
	}
	if table.BigBlindPosition != 0 {
		t.Errorf("Expected big blind position 0, got %d", table.BigBlindPosition)
	}

	// Advance dealer button again
	table.AdvanceDealerButton()

	// Check new positions
	if table.DealerPosition != 4 {
		t.Errorf("Expected dealer position 4, got %d", table.DealerPosition)
	}
	if table.SmallBlindPosition != 0 {
		t.Errorf("Expected small blind position 0, got %d", table.SmallBlindPosition)
	}
	if table.BigBlindPosition != 2 {
		t.Errorf("Expected big blind position 2, got %d", table.BigBlindPosition)
	}
}

// TestStartNewHand tests starting a new hand
func TestStartNewHand(t *testing.T) {
	table := NewTable(9, 5, 10, 100, 1000, NoLimit)
	player1 := NewMockPlayer("p1", "Player 1")
	player2 := NewMockPlayer("p2", "Player 2")

	// Add players
	err := table.AddPlayer(player1, 0, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}
	err = table.AddPlayer(player2, 1, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	// Make sure the seats are active
	seat0, _ := table.GetSeat(0)
	seat0.IsActive = true
	seat1, _ := table.GetSeat(1)
	seat1.IsActive = true

	// Place bets and fold
	_, _, err = table.PlaceBet(0, 50)
	if err != nil {
		t.Errorf("Failed to place bet: %v", err)
	}
	err = table.FoldPlayer(1)
	if err != nil {
		t.Errorf("Failed to fold player: %v", err)
	}

	// Start new hand
	table.StartNewHand()

	// Check if seats are reset
	seat, _ := table.GetSeat(0)
	if seat.HasFolded {
		t.Errorf("Expected player not folded")
	}
	if seat.CurrentBet != 0 {
		t.Errorf("Expected current bet of 0, got %d", seat.CurrentBet)
	}
	if seat.TotalBet != 0 {
		t.Errorf("Expected total bet of 0, got %d", seat.TotalBet)
	}

	seat, _ = table.GetSeat(1)
	if seat.HasFolded {
		t.Errorf("Expected player not folded")
	}

	// Check if pot is reset
	if table.GetTotalPot() != 0 {
		t.Errorf("Expected total pot of 0, got %d", table.GetTotalPot())
	}

	// Check if betting state is reset
	if table.CurrentBet != 0 {
		t.Errorf("Expected current bet of 0, got %d", table.CurrentBet)
	}
	if table.LastRaiseAmount != 0 {
		t.Errorf("Expected last raise amount of 0, got %d", table.LastRaiseAmount)
	}
	if table.CurrentBettingRound != 0 {
		t.Errorf("Expected current betting round of 0, got %d", table.CurrentBettingRound)
	}

	// Check if dealer button advanced
	if table.DealerPosition != 1 {
		t.Errorf("Expected dealer position 1, got %d", table.DealerPosition)
	}
}

// TestStartNewBettingRound tests starting a new betting round
func TestStartNewBettingRound(t *testing.T) {
	table := NewTable(9, 5, 10, 100, 1000, NoLimit)
	player1 := NewMockPlayer("p1", "Player 1")
	player2 := NewMockPlayer("p2", "Player 2")

	// Add players
	err := table.AddPlayer(player1, 0, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}
	err = table.AddPlayer(player2, 1, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	// Make sure the seats are active
	seat0, _ := table.GetSeat(0)
	seat0.IsActive = true
	seat1, _ := table.GetSeat(1)
	seat1.IsActive = true

	// Place bets
	_, _, err = table.PlaceBet(0, 50)
	if err != nil {
		t.Errorf("Failed to place bet: %v", err)
	}
	_, _, err = table.PlaceBet(1, 100)
	if err != nil {
		t.Errorf("Failed to place bet: %v", err)
	}
	_, _, err = table.PlaceBet(0, 100)
	if err != nil {
		t.Errorf("Failed to place bet: %v", err)
	}

	// Start new betting round
	table.StartNewBettingRound()

	// Check if current bets are reset
	seat, _ := table.GetSeat(0)
	if seat.CurrentBet != 0 {
		t.Errorf("Expected current bet of 0, got %d", seat.CurrentBet)
	}
	if seat.TotalBet != 150 {
		t.Errorf("Expected total bet of 150, got %d", seat.TotalBet)
	}

	seat, _ = table.GetSeat(1)
	if seat.CurrentBet != 0 {
		t.Errorf("Expected current bet of 0, got %d", seat.CurrentBet)
	}
	if seat.TotalBet != 100 {
		t.Errorf("Expected total bet of 100, got %d", seat.TotalBet)
	}

	// Check if betting state is reset
	if table.CurrentBet != 0 {
		t.Errorf("Expected current bet of 0, got %d", table.CurrentBet)
	}
	if table.LastRaiseAmount != 0 {
		t.Errorf("Expected last raise amount of 0, got %d", table.LastRaiseAmount)
	}
	if table.CurrentBettingRound != 1 {
		t.Errorf("Expected current betting round of 1, got %d", table.CurrentBettingRound)
	}
}

// TestAwardPotsToWinners tests awarding pots to winners
func TestAwardPotsToWinners(t *testing.T) {
	table := NewTable(9, 5, 10, 100, 1000, NoLimit)
	player1 := NewMockPlayer("p1", "Player 1")
	player2 := NewMockPlayer("p2", "Player 2")
	player3 := NewMockPlayer("p3", "Player 3")

	// Add players
	err := table.AddPlayer(player1, 0, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}
	err = table.AddPlayer(player2, 1, 100)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}
	err = table.AddPlayer(player3, 2, 300)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	// Make sure the seats are active
	for i := 0; i < 3; i++ {
		seat, _ := table.GetSeat(i)
		seat.IsActive = true
	}

	// Place bets
	_, _, err = table.PlaceBet(0, 50)
	if err != nil {
		t.Errorf("Failed to place bet: %v", err)
	}
	_, allIn, err := table.PlaceBet(1, 100)
	if err != nil {
		t.Errorf("Failed to place bet: %v", err)
	}
	if !allIn {
		t.Errorf("Expected player 2 to be all-in")
	}
	_, _, err = table.PlaceBet(2, 100)
	if err != nil {
		t.Errorf("Failed to place bet: %v", err)
	}
	_, _, err = table.PlaceBet(0, 100)
	if err != nil {
		t.Errorf("Failed to place bet: %v", err)
	}

	// Award pots to winners
	winners := map[int]string{
		-1: "p3", // Main pot to player 3
		0:  "p1", // Side pot to player 1
	}
	table.AwardPotsToWinners(winners)

	// Check if winners received their chips
	seat, _ := table.GetSeat(0)
	if seat.ChipStack != 50 {
		t.Errorf("Expected player 1 chip stack of 50, got %d", seat.ChipStack)
	}

	seat, _ = table.GetSeat(1)
	if seat.ChipStack != 0 {
		t.Errorf("Expected player 2 chip stack of 0, got %d", seat.ChipStack)
	}

	seat, _ = table.GetSeat(2)
	if seat.ChipStack != 550 {
		t.Errorf("Expected player 3 chip stack of 550, got %d", seat.ChipStack)
	}

	// Check if pot is reset
	if table.GetTotalPot() != 0 {
		t.Errorf("Expected total pot of 0, got %d", table.GetTotalPot())
	}
}

// TestGetMinRaise tests getting the minimum raise amount
func TestGetMinRaise(t *testing.T) {
	table := NewTable(9, 5, 10, 100, 1000, NoLimit)
	player1 := NewMockPlayer("p1", "Player 1")
	player2 := NewMockPlayer("p2", "Player 2")

	// Add players
	err := table.AddPlayer(player1, 0, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}
	err = table.AddPlayer(player2, 1, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	// Make sure the seats are active
	seat0, _ := table.GetSeat(0)
	seat0.IsActive = true
	seat1, _ := table.GetSeat(1)
	seat1.IsActive = true

	// Initial min raise (no previous raise)
	minRaise := table.GetMinRaise()
	if minRaise != 0 {
		t.Errorf("Expected min raise of 0, got %d", minRaise)
	}

	// Place bet
	_, _, err = table.PlaceBet(0, 20)
	if err != nil {
		t.Errorf("Failed to place bet: %v", err)
	}

	// Min raise after first bet
	minRaise = table.GetMinRaise()
	if minRaise != 40 {
		t.Errorf("Expected min raise of 40, got %d", minRaise)
	}

	// Place raise
	_, _, err = table.PlaceBet(1, 50)
	if err != nil {
		t.Errorf("Failed to place bet: %v", err)
	}

	// Min raise after raise
	minRaise = table.GetMinRaise()
	if minRaise != 80 {
		t.Errorf("Expected min raise of 80, got %d", minRaise)
	}
}

// TestGetMaxBet tests getting the maximum bet amount
func TestGetMaxBet(t *testing.T) {
	// Test No Limit
	table := NewTable(9, 5, 10, 100, 1000, NoLimit)
	player1 := NewMockPlayer("p1", "Player 1")

	// Add player
	err := table.AddPlayer(player1, 0, 200)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	// Make sure the seat is active
	seat, _ := table.GetSeat(0)
	seat.IsActive = true

	// Max bet for No Limit
	maxBet, err := table.GetMaxBet(0)
	if err != nil {
		t.Errorf("Failed to get max bet: %v", err)
	}
	if maxBet != 200 {
		t.Errorf("Expected max bet of 200, got %d", maxBet)
	}

	// Test Pot Limit
	table = NewTable(9, 5, 10, 100, 1000, PotLimit)
	player1 = NewMockPlayer("p1", "Player 1")
	player2 := NewMockPlayer("p2", "Player 2")

	// Add players
	err = table.AddPlayer(player1, 0, 500)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}
	err = table.AddPlayer(player2, 1, 500)
	if err != nil {
		t.Errorf("Failed to add player: %v", err)
	}

	// Make sure the seats are active
	seat0, _ := table.GetSeat(0)
	seat0.IsActive = true
	seat1, _ := table.GetSeat(1)
	seat1.IsActive = true

	// Manually set up the pot for testing
	table.PotManager.MainPot.Amount = 150

	// Print debugging information
	t.Logf("Initial setup - Pot size: %d, Current bet: %d", table.GetTotalPot(), table.CurrentBet)

	// The pot is 150, current bet is 0
	// For a pot limit game, the max bet should be pot size + current bet = 150 + 0 = 150
	maxBet, err = table.GetMaxBet(0)
	if err != nil {
		t.Errorf("Failed to get max bet: %v", err)
	}
	t.Logf("MaxBet for player 0: %d", maxBet)
	if maxBet != 150 {
		t.Errorf("Expected max bet of 150, got %d", maxBet)
	}

	// Now let's set the current bet and check the max bet for the other player
	table.CurrentBet = 100
	table.PotManager.MainPot.Amount = 250 // Pot is now 250

	// Print debugging information
	t.Logf("After setup - Pot size: %d, Current bet: %d", table.GetTotalPot(), table.CurrentBet)

	// The pot is 250, current bet is 100
	// For a pot limit game, the max bet should be pot size + current bet = 250 + 100 = 350
	maxBet, err = table.GetMaxBet(1)
	if err != nil {
		t.Errorf("Failed to get max bet: %v", err)
	}
	t.Logf("MaxBet for player 1: %d", maxBet)
	if maxBet != 350 {
		t.Errorf("Expected max bet of 350, got %d", maxBet)
	}
}
