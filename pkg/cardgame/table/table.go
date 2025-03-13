package table

import (
	"errors"
	"sync"

	"github.com/blixxurd/card-game-go/pkg/cardgame/player"
)

// BettingStructure represents the type of betting structure used in a poker game.
type BettingStructure int

const (
	// NoLimit allows players to bet any amount up to their chip stack.
	NoLimit BettingStructure = iota

	// PotLimit limits the maximum bet to the current pot size.
	PotLimit

	// FixedLimit uses fixed bet sizes.
	FixedLimit
)

// Table represents a poker table with seats, chips, and betting.
type Table struct {
	// Seats are the positions at the table where players can sit
	Seats []*Seat

	// PotManager manages the main pot and side pots
	PotManager *PotManager

	// DealerPosition is the position of the dealer button (0-indexed)
	DealerPosition int

	// SmallBlindPosition is the position of the small blind (0-indexed)
	SmallBlindPosition int

	// BigBlindPosition is the position of the big blind (0-indexed)
	BigBlindPosition int

	// SmallBlindAmount is the amount of the small blind
	SmallBlindAmount int

	// BigBlindAmount is the amount of the big blind
	BigBlindAmount int

	// MinBuyIn is the minimum amount of chips a player can bring to the table
	MinBuyIn int

	// MaxBuyIn is the maximum amount of chips a player can bring to the table
	MaxBuyIn int

	// BettingRules defines the betting structure and rules
	BettingRules BettingRules

	// CurrentBettingRound is the current betting round (0-indexed)
	CurrentBettingRound int

	// CurrentBet is the current bet amount that players must match
	CurrentBet int

	// LastRaiseAmount is the amount of the last raise
	LastRaiseAmount int

	// mutex protects the table from concurrent access
	mutex sync.RWMutex
}

// NewTable creates a new poker table with the specified number of seats and betting structure.
func NewTable(numSeats int, smallBlind, bigBlind, minBuyIn, maxBuyIn int, bettingStructure BettingStructure) *Table {
	// Create seats
	seats := make([]*Seat, numSeats)
	for i := 0; i < numSeats; i++ {
		seats[i] = NewSeat(i)
	}

	return &Table{
		Seats:               seats,
		PotManager:          NewPotManager(),
		DealerPosition:      0,
		SmallBlindPosition:  1,
		BigBlindPosition:    2,
		SmallBlindAmount:    smallBlind,
		BigBlindAmount:      bigBlind,
		MinBuyIn:            minBuyIn,
		MaxBuyIn:            maxBuyIn,
		BettingRules:        NewBettingRules(bettingStructure, smallBlind, bigBlind*2, 4),
		CurrentBettingRound: 0,
		CurrentBet:          0,
		LastRaiseAmount:     0,
	}
}

// AddPlayer adds a player to the specified seat with the specified chip stack.
func (t *Table) AddPlayer(player player.Player, seatPosition int, chipStack int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Validate seat position
	if seatPosition < 0 || seatPosition >= len(t.Seats) {
		return errors.New("invalid seat position")
	}

	// Check if seat is already occupied
	if t.Seats[seatPosition].IsOccupied {
		return errors.New("seat already occupied")
	}

	// Validate chip stack
	if chipStack < t.MinBuyIn {
		return errors.New("chip stack below minimum buy-in")
	}
	if chipStack > t.MaxBuyIn {
		return errors.New("chip stack above maximum buy-in")
	}

	// Add player to seat
	t.Seats[seatPosition].Occupy(player, chipStack)

	// Add player to pot eligibility
	t.PotManager.AddPlayerToPots(player.ID())

	return nil
}

// RemovePlayer removes a player from the table.
func (t *Table) RemovePlayer(playerID string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Find the player's seat
	for _, seat := range t.Seats {
		if seat.IsOccupied && seat.Player.ID() == playerID {
			// Remove player from pot eligibility
			t.PotManager.RemovePlayerFromPots(playerID)

			// Remove player from seat
			seat.Vacate()
			return nil
		}
	}

	return errors.New("player not found at table")
}

// GetSeat returns the seat at the specified position.
func (t *Table) GetSeat(position int) (*Seat, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if position < 0 || position >= len(t.Seats) {
		return nil, errors.New("invalid seat position")
	}

	return t.Seats[position], nil
}

// GetPlayerSeat returns the seat occupied by the specified player.
func (t *Table) GetPlayerSeat(playerID string) (*Seat, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for _, seat := range t.Seats {
		if seat.IsOccupied && seat.Player.ID() == playerID {
			return seat, nil
		}
	}

	return nil, errors.New("player not found at table")
}

// GetOccupiedSeats returns all occupied seats at the table.
func (t *Table) GetOccupiedSeats() []*Seat {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	occupiedSeats := make([]*Seat, 0)
	for _, seat := range t.Seats {
		if seat.IsOccupied {
			occupiedSeats = append(occupiedSeats, seat)
		}
	}

	return occupiedSeats
}

// GetActiveSeats returns all active seats at the table (occupied and not folded).
func (t *Table) GetActiveSeats() []*Seat {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	activeSeats := make([]*Seat, 0)
	for _, seat := range t.Seats {
		if seat.IsOccupied && seat.IsActive && !seat.HasFolded {
			activeSeats = append(activeSeats, seat)
		}
	}

	return activeSeats
}

// GetActivePlayers returns all active players at the table.
func (t *Table) GetActivePlayers() []player.Player {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	activePlayers := make([]player.Player, 0)
	for _, seat := range t.Seats {
		if seat.IsOccupied && seat.IsActive && !seat.HasFolded {
			activePlayers = append(activePlayers, seat.Player)
		}
	}

	return activePlayers
}

// GetTotalPot returns the total amount in all pots.
func (t *Table) GetTotalPot() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.PotManager.GetTotalPotAmount()
}

// PlaceBet places a bet for the player at the specified seat position.
// Returns the actual amount bet and whether the player went all-in.
func (t *Table) PlaceBet(seatPosition int, amount int) (int, bool, error) {
	// Get the total pot size before locking the mutex
	potSize := t.GetTotalPot()

	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Validate seat position
	if seatPosition < 0 || seatPosition >= len(t.Seats) {
		return 0, false, errors.New("invalid seat position")
	}

	seat := t.Seats[seatPosition]

	// Check if player can bet
	if !seat.CanBet() {
		return 0, false, errors.New("player cannot bet")
	}

	// Validate bet according to betting rules
	validatedAmount, err := t.BettingRules.ValidateBet(
		t.CurrentBet,
		t.LastRaiseAmount,
		seat.ChipStack,
		amount,
		potSize, // Use the pot size we got before locking
	)
	if err != nil {
		return 0, false, err
	}

	// Place the bet
	actualAmount, isAllIn := seat.PlaceBet(validatedAmount)

	// Update the current bet if this is a raise
	if seat.CurrentBet > t.CurrentBet {
		t.LastRaiseAmount = seat.CurrentBet - t.CurrentBet
		t.CurrentBet = seat.CurrentBet
	}

	// Add the bet to the main pot
	t.PotManager.MainPot.AddAmount(actualAmount)

	// Make the player eligible for the main pot
	t.PotManager.MainPot.SetPlayerEligibility(seat.Player.ID(), true)

	// If player went all-in, create a side pot
	if isAllIn {
		// Use the enhanced side pot creation method
		t.PotManager.CreateSidePotForAllIn(seat.Player.ID(), seat.CurrentBet, t.Seats)
	}

	return actualAmount, isAllIn, nil
}

// FoldPlayer marks the player at the specified seat position as folded.
func (t *Table) FoldPlayer(seatPosition int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Validate seat position
	if seatPosition < 0 || seatPosition >= len(t.Seats) {
		return errors.New("invalid seat position")
	}

	seat := t.Seats[seatPosition]

	// Check if player can fold
	if !seat.IsOccupied || !seat.IsActive || seat.HasFolded {
		return errors.New("player cannot fold")
	}

	// Fold the player
	seat.Fold()

	return nil
}

// CollectBlinds collects the small and big blinds from the appropriate players.
func (t *Table) CollectBlinds() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Get the small blind and big blind seats
	smallBlindSeat := t.Seats[t.SmallBlindPosition]
	bigBlindSeat := t.Seats[t.BigBlindPosition]

	// Check if the seats are occupied and active
	if !smallBlindSeat.IsOccupied || !smallBlindSeat.IsActive {
		return errors.New("small blind seat not occupied or not active")
	}
	if !bigBlindSeat.IsOccupied || !bigBlindSeat.IsActive {
		return errors.New("big blind seat not occupied or not active")
	}

	// Collect small blind
	smallBlindAmount, smallBlindAllIn := smallBlindSeat.PlaceBet(t.SmallBlindAmount)
	t.PotManager.MainPot.AddAmount(smallBlindAmount)

	// Make the small blind player eligible for the main pot
	t.PotManager.MainPot.SetPlayerEligibility(smallBlindSeat.Player.ID(), true)

	// Collect big blind
	bigBlindAmount, bigBlindAllIn := bigBlindSeat.PlaceBet(t.BigBlindAmount)
	t.PotManager.MainPot.AddAmount(bigBlindAmount)

	// Make the big blind player eligible for the main pot
	t.PotManager.MainPot.SetPlayerEligibility(bigBlindSeat.Player.ID(), true)

	// Set the current bet to the big blind amount
	t.CurrentBet = bigBlindAmount

	// Handle all-ins if necessary
	if smallBlindAllIn || bigBlindAllIn {
		// Create side pots if needed
		if smallBlindAllIn {
			// Use the enhanced side pot creation method
			t.PotManager.CreateSidePotForAllIn(smallBlindSeat.Player.ID(), smallBlindSeat.TotalBet, t.Seats)
		}

		if bigBlindAllIn {
			// Use the enhanced side pot creation method
			t.PotManager.CreateSidePotForAllIn(bigBlindSeat.Player.ID(), bigBlindSeat.TotalBet, t.Seats)
		}
	}

	return nil
}

// AdvanceDealerButton advances the dealer button to the next occupied seat.
func (t *Table) AdvanceDealerButton() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Find the next occupied seat
	originalPosition := t.DealerPosition
	for {
		t.DealerPosition = (t.DealerPosition + 1) % len(t.Seats)
		if t.Seats[t.DealerPosition].IsOccupied || t.DealerPosition == originalPosition {
			break
		}
	}

	// Update small and big blind positions
	t.updateBlindPositions()
}

// updateBlindPositions updates the small and big blind positions based on the dealer position.
func (t *Table) updateBlindPositions() {
	// Small blind is the next occupied seat after the dealer
	t.SmallBlindPosition = t.DealerPosition
	for {
		t.SmallBlindPosition = (t.SmallBlindPosition + 1) % len(t.Seats)
		if t.Seats[t.SmallBlindPosition].IsOccupied {
			break
		}
		// Avoid infinite loop if there's only one player
		if t.SmallBlindPosition == t.DealerPosition {
			break
		}
	}

	// Big blind is the next occupied seat after the small blind
	t.BigBlindPosition = t.SmallBlindPosition
	for {
		t.BigBlindPosition = (t.BigBlindPosition + 1) % len(t.Seats)
		if t.Seats[t.BigBlindPosition].IsOccupied {
			break
		}
		// Avoid infinite loop if there's only one or two players
		if t.BigBlindPosition == t.DealerPosition || t.BigBlindPosition == t.SmallBlindPosition {
			break
		}
	}
}

// StartNewHand prepares the table for a new hand.
func (t *Table) StartNewHand() {
	t.mutex.Lock()

	// Reset all seats for a new hand
	for _, seat := range t.Seats {
		seat.Reset()
	}

	// Reset pot manager
	t.PotManager.Reset()

	// Reset betting state
	t.CurrentBettingRound = 0
	t.CurrentBet = 0
	t.LastRaiseAmount = 0

	// If using fixed limit rules, update the current round
	if fixedRules, ok := t.BettingRules.(*FixedLimitRules); ok {
		fixedRules.SetCurrentRound(0)
	}

	// Release the mutex before advancing the dealer button
	t.mutex.Unlock()

	// Advance dealer button (this will acquire its own lock)
	t.AdvanceDealerButton()
}

// StartNewBettingRound starts a new betting round.
func (t *Table) StartNewBettingRound() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Process any all-in players and create side pots if needed
	for _, seat := range t.Seats {
		if seat.IsOccupied && seat.IsActive && !seat.HasFolded && seat.ChipStack == 0 {
			// This player is all-in, create a side pot
			t.PotManager.CreateSidePotForAllIn(seat.Player.ID(), seat.TotalBet, t.Seats)
		}
	}

	// Reset betting state for the new round
	t.CurrentBettingRound++
	t.CurrentBet = 0
	t.LastRaiseAmount = 0

	// If using fixed limit rules, update the current round
	if fixedRules, ok := t.BettingRules.(*FixedLimitRules); ok {
		fixedRules.SetCurrentRound(t.CurrentBettingRound)
	}

	// Reset current bets for all seats
	for _, seat := range t.Seats {
		seat.ResetBetForRound()
	}
}

// AwardPotsToWinners awards the pots to the winners.
// The winners map maps pot indices (-1 for main pot, 0+ for side pots) to winner player IDs.
func (t *Table) AwardPotsToWinners(winners map[int]string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Award main pot
	if winnerID, ok := winners[-1]; ok {
		for _, seat := range t.Seats {
			if seat.IsOccupied && seat.Player.ID() == winnerID {
				seat.AddChips(t.PotManager.MainPot.Amount)
				break
			}
		}
	}

	// Award side pots
	for potIdx, winnerID := range winners {
		if potIdx >= 0 && potIdx < len(t.PotManager.SidePots) {
			for _, seat := range t.Seats {
				if seat.IsOccupied && seat.Player.ID() == winnerID {
					seat.AddChips(t.PotManager.SidePots[potIdx].Amount)
					break
				}
			}
		}
	}

	// Reset the pot manager after awarding pots
	t.PotManager.Reset()
}

// AwardPotsToMultipleWinners awards the pots to multiple winners, splitting pots when necessary.
// The winners map maps pot indices (-1 for main pot, 0+ for side pots) to slices of winner player IDs.
func (t *Table) AwardPotsToMultipleWinners(winners map[int][]string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Award main pot
	if winnerIDs, ok := winners[-1]; ok && len(winnerIDs) > 0 {
		// Calculate the amount each winner gets
		amountPerWinner := t.PotManager.MainPot.Amount / len(winnerIDs)
		// Calculate the remainder (if the pot can't be split evenly)
		remainder := t.PotManager.MainPot.Amount % len(winnerIDs)

		// Award the pot to each winner
		for i, winnerID := range winnerIDs {
			for _, seat := range t.Seats {
				if seat.IsOccupied && seat.Player.ID() == winnerID {
					// Add the base amount
					seat.AddChips(amountPerWinner)

					// Add the remainder to the first winner (if any)
					// This is a common approach to handle odd chip amounts
					if i == 0 && remainder > 0 {
						seat.AddChips(remainder)
					}
					break
				}
			}
		}
	}

	// Award side pots
	for potIdx, winnerIDs := range winners {
		if potIdx >= 0 && potIdx < len(t.PotManager.SidePots) && len(winnerIDs) > 0 {
			// Calculate the amount each winner gets
			amountPerWinner := t.PotManager.SidePots[potIdx].Amount / len(winnerIDs)
			// Calculate the remainder (if the pot can't be split evenly)
			remainder := t.PotManager.SidePots[potIdx].Amount % len(winnerIDs)

			// Award the pot to each winner
			for i, winnerID := range winnerIDs {
				for _, seat := range t.Seats {
					if seat.IsOccupied && seat.Player.ID() == winnerID {
						// Add the base amount
						seat.AddChips(amountPerWinner)

						// Add the remainder to the first winner (if any)
						if i == 0 && remainder > 0 {
							seat.AddChips(remainder)
						}
						break
					}
				}
			}
		}
	}

	// Reset the pot manager after awarding pots
	t.PotManager.Reset()
}

// GetMinRaise returns the minimum raise amount according to the betting rules.
func (t *Table) GetMinRaise() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.BettingRules.GetMinRaise(t.CurrentBet, t.LastRaiseAmount)
}

// GetMaxBet returns the maximum bet amount a player can make according to the betting rules.
func (t *Table) GetMaxBet(seatPosition int) (int, error) {
	// Get the total pot size before locking the mutex
	potSize := t.GetTotalPot()

	t.mutex.RLock()
	defer t.mutex.RUnlock()

	// Validate seat position
	if seatPosition < 0 || seatPosition >= len(t.Seats) {
		return 0, errors.New("invalid seat position")
	}

	seat := t.Seats[seatPosition]

	// Check if player can bet
	if !seat.CanBet() {
		return 0, errors.New("player cannot bet")
	}

	return t.BettingRules.GetMaxBet(seat.ChipStack, t.CurrentBet, potSize), nil
}
