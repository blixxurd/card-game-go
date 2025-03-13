package table

import (
	"github.com/blixxurd/card-game-go/pkg/cardgame/player"
)

// Pot represents a pot of chips in a poker game.
type Pot struct {
	// Amount is the total number of chips in the pot
	Amount int

	// EligiblePlayers is a map of player IDs to a boolean indicating whether
	// they are eligible to win this pot
	EligiblePlayers map[string]bool
}

// NewPot creates a new pot with the specified amount and eligible players.
func NewPot(amount int, eligiblePlayers map[string]bool) *Pot {
	// Create a copy of the eligible players map
	eligibleCopy := make(map[string]bool)
	for id, eligible := range eligiblePlayers {
		eligibleCopy[id] = eligible
	}

	return &Pot{
		Amount:          amount,
		EligiblePlayers: eligibleCopy,
	}
}

// NewMainPot creates a new main pot with all players eligible.
func NewMainPot() *Pot {
	return &Pot{
		Amount:          0,
		EligiblePlayers: make(map[string]bool),
	}
}

// AddAmount adds the specified amount to the pot.
func (p *Pot) AddAmount(amount int) {
	p.Amount += amount
}

// IsPlayerEligible returns whether the specified player is eligible to win this pot.
func (p *Pot) IsPlayerEligible(playerID string) bool {
	eligible, exists := p.EligiblePlayers[playerID]
	return exists && eligible
}

// SetPlayerEligibility sets whether the specified player is eligible to win this pot.
func (p *Pot) SetPlayerEligibility(playerID string, eligible bool) {
	p.EligiblePlayers[playerID] = eligible
}

// GetEligiblePlayers returns a slice of player IDs that are eligible to win this pot.
func (p *Pot) GetEligiblePlayers() []string {
	eligibleIDs := make([]string, 0, len(p.EligiblePlayers))
	for id, eligible := range p.EligiblePlayers {
		if eligible {
			eligibleIDs = append(eligibleIDs, id)
		}
	}
	return eligibleIDs
}

// PotManager manages the main pot and side pots in a poker game.
type PotManager struct {
	// MainPot is the main pot that all players contribute to
	MainPot *Pot

	// SidePots are additional pots created when players go all-in
	SidePots []*Pot
}

// NewPotManager creates a new pot manager.
func NewPotManager() *PotManager {
	return &PotManager{
		MainPot:  NewMainPot(),
		SidePots: make([]*Pot, 0),
	}
}

// AddPlayerToPots adds a player to all pots they are eligible for.
func (pm *PotManager) AddPlayerToPots(playerID string) {
	pm.MainPot.SetPlayerEligibility(playerID, true)
	for _, pot := range pm.SidePots {
		pot.SetPlayerEligibility(playerID, true)
	}
}

// RemovePlayerFromPots removes a player from all pots.
func (pm *PotManager) RemovePlayerFromPots(playerID string) {
	pm.MainPot.SetPlayerEligibility(playerID, false)
	for _, pot := range pm.SidePots {
		pot.SetPlayerEligibility(playerID, false)
	}
}

// GetTotalPotAmount returns the total amount of chips in all pots.
func (pm *PotManager) GetTotalPotAmount() int {
	total := pm.MainPot.Amount
	for _, pot := range pm.SidePots {
		total += pot.Amount
	}
	return total
}

// CreateSidePot creates a side pot when a player goes all-in.
// allInPlayerID is the ID of the player who went all-in.
// allInAmount is the amount the player went all-in for.
// seats is the list of seats at the table.
func (pm *PotManager) CreateSidePot(allInPlayerID string, allInAmount int, seats []*Seat) {
	// Create a map of eligible players for the side pot
	// Initially, all players who have contributed to the main pot are eligible
	eligiblePlayers := make(map[string]bool)
	for playerID, eligible := range pm.MainPot.EligiblePlayers {
		if eligible && playerID != allInPlayerID {
			eligiblePlayers[playerID] = true
		}
	}

	// Calculate the excess amount that goes into the side pot
	excessAmount := 0
	for _, seat := range seats {
		if seat.IsOccupied && seat.IsActive && !seat.HasFolded {
			// Skip the all-in player
			if seat.Player.ID() == allInPlayerID {
				continue
			}

			// Calculate how much this player has contributed in excess of the all-in amount
			if seat.CurrentBet > allInAmount {
				excess := seat.CurrentBet - allInAmount
				excessAmount += excess
				// Reduce the player's current bet to match the all-in amount
				seat.CurrentBet = allInAmount
			}
		}
	}

	// Only create a side pot if there's an excess amount
	if excessAmount > 0 {
		// Create the side pot with the excess amount
		sidePot := NewPot(excessAmount, eligiblePlayers)

		// The all-in player is not eligible for the side pot
		sidePot.SetPlayerEligibility(allInPlayerID, false)

		// Add the side pot to the list of side pots
		pm.SidePots = append(pm.SidePots, sidePot)

		// Update the main pot to reflect that the all-in player is still eligible for it
		pm.MainPot.SetPlayerEligibility(allInPlayerID, true)

		// Reduce the main pot by the excess amount that was moved to the side pot
		// This is important to ensure the total pot amount remains correct
		pm.MainPot.Amount -= excessAmount
	}
}

// CreateSidePotForAllIn creates a side pot when a player goes all-in.
// This is an enhanced version that handles multiple all-ins and complex scenarios.
// allInPlayerID is the ID of the player who went all-in.
// allInAmount is the amount the player went all-in for.
// seats is the list of seats at the table.
func (pm *PotManager) CreateSidePotForAllIn(allInPlayerID string, allInAmount int, seats []*Seat) {
	// First, check if there are existing side pots
	// If there are, we need to handle this differently
	if len(pm.SidePots) > 0 {
		// Find the smallest existing side pot that the all-in player is eligible for
		smallestEligiblePotIdx := -1
		for i, pot := range pm.SidePots {
			if pot.IsPlayerEligible(allInPlayerID) {
				smallestEligiblePotIdx = i
				break
			}
		}

		// If the all-in player is eligible for any side pot,
		// we need to create a new side pot from that one
		if smallestEligiblePotIdx >= 0 {
			// Get the smallest eligible pot
			smallestEligiblePot := pm.SidePots[smallestEligiblePotIdx]

			// Create a map of eligible players for the new side pot
			eligiblePlayers := make(map[string]bool)
			for playerID, eligible := range smallestEligiblePot.EligiblePlayers {
				if eligible && playerID != allInPlayerID {
					eligiblePlayers[playerID] = true
				}
			}

			// Calculate the excess amount that goes into the new side pot
			excessAmount := 0
			for _, seat := range seats {
				if seat.IsOccupied && seat.IsActive && !seat.HasFolded {
					// Skip the all-in player
					if seat.Player.ID() == allInPlayerID {
						continue
					}

					// Calculate how much this player has contributed in excess of the all-in amount
					if seat.CurrentBet > allInAmount {
						excess := seat.CurrentBet - allInAmount
						excessAmount += excess
						// Reduce the player's current bet to match the all-in amount
						seat.CurrentBet = allInAmount
					}
				}
			}

			// Only create a new side pot if there's an excess amount
			if excessAmount > 0 {
				// Create the new side pot with the excess amount
				newSidePot := NewPot(excessAmount, eligiblePlayers)

				// The all-in player is not eligible for the new side pot
				newSidePot.SetPlayerEligibility(allInPlayerID, false)

				// Add the new side pot to the list of side pots
				pm.SidePots = append(pm.SidePots, newSidePot)

				// Update the existing side pot to reflect that the all-in player is still eligible for it
				smallestEligiblePot.SetPlayerEligibility(allInPlayerID, true)

				// Reduce the eligible pot by the excess amount that was moved to the new side pot
				smallestEligiblePot.Amount -= excessAmount
			}
		} else {
			// If the all-in player is not eligible for any existing side pot,
			// create a new side pot from the main pot
			pm.CreateSidePot(allInPlayerID, allInAmount, seats)
		}
	} else {
		// If there are no existing side pots, create a new one from the main pot
		pm.CreateSidePot(allInPlayerID, allInAmount, seats)
	}
}

// Reset resets the pot manager for a new hand.
func (pm *PotManager) Reset() {
	pm.MainPot = NewMainPot()
	pm.SidePots = make([]*Pot, 0)
}

// AwardPotToWinner awards a pot to the specified winner.
func (pm *PotManager) AwardPotToWinner(pot *Pot, winner player.Player, seats []*Seat) {
	// Find the winner's seat
	var winnerSeat *Seat
	for _, seat := range seats {
		if seat.IsOccupied && seat.Player.ID() == winner.ID() {
			winnerSeat = seat
			break
		}
	}

	// If we found the winner's seat, add the pot amount to their chip stack
	if winnerSeat != nil {
		winnerSeat.AddChips(pot.Amount)
	}
}
