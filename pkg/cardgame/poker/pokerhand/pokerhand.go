// Package pokerhand provides functionality for evaluating poker hands.
package pokerhand

import (
	"fmt"
	"sort"

	"github.com/blixxurd/card-game-go/pkg/cardgame/card"
)

// Define card ranks
const (
	Two   = 2
	Three = 3
	Four  = 4
	Five  = 5
	Six   = 6
	Seven = 7
	Eight = 8
	Nine  = 9
	Ten   = 10
	Jack  = 11
	Queen = 12
	King  = 13
	Ace   = 14 // Ace is high in poker
)

// HandRank represents the rank of a poker hand
type HandRank int

const (
	// HighCard represents a high card hand
	HighCard HandRank = iota
	// OnePair represents a one pair hand
	OnePair
	// TwoPair represents a two pair hand
	TwoPair
	// ThreeOfAKind represents a three of a kind hand
	ThreeOfAKind
	// Straight represents a straight hand
	Straight
	// Flush represents a flush hand
	Flush
	// FullHouse represents a full house hand
	FullHouse
	// FourOfAKind represents a four of a kind hand
	FourOfAKind
	// StraightFlush represents a straight flush hand
	StraightFlush
	// RoyalFlush represents a royal flush hand
	RoyalFlush
)

// HandResult represents the result of a hand evaluation
type HandResult struct {
	// Rank is the rank of the hand
	Rank HandRank
	// Description is a human-readable description of the hand
	Description string
	// Cards are the cards that make up the hand
	Cards []card.Card
	// Value is a numeric value used for comparing hands of the same rank
	Value int
}

// EvaluateHand evaluates a poker hand and returns the result
func EvaluateHand(cards []card.Card) HandResult {
	if len(cards) < 5 {
		return HandResult{
			Rank:        HighCard,
			Description: "Not enough cards",
			Cards:       cards,
			Value:       0,
		}
	}

	// For Texas Hold'em, we need to find the best 5-card hand from 7 cards
	if len(cards) > 5 {
		return findBestHand(cards)
	}

	// Sort the cards by value (descending)
	sortedCards := make([]card.Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return getCardValue(sortedCards[i]) > getCardValue(sortedCards[j])
	})

	// Check for royal flush
	if isRoyalFlush(sortedCards) {
		return HandResult{
			Rank:        RoyalFlush,
			Description: "Royal Flush",
			Cards:       sortedCards[:5],
			Value:       Ace,
		}
	}

	// Check for straight flush
	if straightFlushValue := isStraightFlush(sortedCards); straightFlushValue > 0 {
		return HandResult{
			Rank:        StraightFlush,
			Description: fmt.Sprintf("Straight Flush, %s high", rankToString(straightFlushValue)),
			Cards:       sortedCards[:5],
			Value:       straightFlushValue,
		}
	}

	// Check for four of a kind
	if fourOfAKindValue, kicker := isFourOfAKind(sortedCards); fourOfAKindValue > 0 {
		return HandResult{
			Rank:        FourOfAKind,
			Description: fmt.Sprintf("Four of a Kind, %ss", rankToString(fourOfAKindValue)),
			Cards:       sortedCards[:5],
			Value:       fourOfAKindValue*100 + kicker,
		}
	}

	// Check for full house
	if threeOfAKindValue, pairValue := isFullHouse(sortedCards); threeOfAKindValue > 0 {
		return HandResult{
			Rank:        FullHouse,
			Description: fmt.Sprintf("Full House, %ss over %ss", rankToString(threeOfAKindValue), rankToString(pairValue)),
			Cards:       sortedCards[:5],
			Value:       threeOfAKindValue*100 + pairValue,
		}
	}

	// Check for flush
	if isFlush(sortedCards) {
		return HandResult{
			Rank:        Flush,
			Description: fmt.Sprintf("Flush, %s high", rankToString(getCardValue(sortedCards[0]))),
			Cards:       sortedCards[:5],
			Value:       calculateHighCardValue(sortedCards[:5]),
		}
	}

	// Check for straight
	if straightValue := isStraight(sortedCards); straightValue > 0 {
		return HandResult{
			Rank:        Straight,
			Description: fmt.Sprintf("Straight, %s high", rankToString(straightValue)),
			Cards:       sortedCards[:5],
			Value:       straightValue,
		}
	}

	// Check for three of a kind
	if threeOfAKindValue, kickers := isThreeOfAKind(sortedCards); threeOfAKindValue > 0 {
		return HandResult{
			Rank:        ThreeOfAKind,
			Description: fmt.Sprintf("Three of a Kind, %ss", rankToString(threeOfAKindValue)),
			Cards:       sortedCards[:5],
			Value:       threeOfAKindValue*10000 + kickers[0]*100 + kickers[1],
		}
	}

	// Check for two pair
	if highPairValue, lowPairValue, kicker := isTwoPair(sortedCards); highPairValue > 0 {
		return HandResult{
			Rank:        TwoPair,
			Description: fmt.Sprintf("Two Pair, %ss and %ss", rankToString(highPairValue), rankToString(lowPairValue)),
			Cards:       sortedCards[:5],
			Value:       highPairValue*10000 + lowPairValue*100 + kicker,
		}
	}

	// Check for one pair
	if pairValue, kickers := isOnePair(sortedCards); pairValue > 0 {
		return HandResult{
			Rank:        OnePair,
			Description: fmt.Sprintf("Pair of %ss", rankToString(pairValue)),
			Cards:       sortedCards[:5],
			Value:       pairValue*1000000 + kickers[0]*10000 + kickers[1]*100 + kickers[2],
		}
	}

	// High card
	return HandResult{
		Rank:        HighCard,
		Description: fmt.Sprintf("High Card %s", rankToString(getCardValue(sortedCards[0]))),
		Cards:       sortedCards[:5],
		Value:       calculateHighCardValue(sortedCards[:5]),
	}
}

// CompareHands compares two hand results and returns:
// 1 if hand1 is better than hand2
// 0 if they are equal
// -1 if hand2 is better than hand1
func CompareHands(hand1, hand2 HandResult) int {
	if hand1.Rank > hand2.Rank {
		return 1
	}
	if hand1.Rank < hand2.Rank {
		return -1
	}
	if hand1.Value > hand2.Value {
		return 1
	}
	if hand1.Value < hand2.Value {
		return -1
	}
	return 0
}

// Helper functions

// getCardValue returns the comparison value of a card
func getCardValue(c card.Card) int {
	// Try to cast to PlayingCard
	if pc, ok := c.(card.PlayingCard); ok {
		return pc.GetComparisonValue()
	}

	// Fallback to using the card's name to determine value
	name := c.Name()
	switch {
	case name == "Ace":
		return Ace
	case name == "King":
		return King
	case name == "Queen":
		return Queen
	case name == "Jack":
		return Jack
	case name == "10":
		return Ten
	case name == "9":
		return Nine
	case name == "8":
		return Eight
	case name == "7":
		return Seven
	case name == "6":
		return Six
	case name == "5":
		return Five
	case name == "4":
		return Four
	case name == "3":
		return Three
	case name == "2":
		return Two
	default:
		return 0
	}
}

// getCardSuit returns the suit of a card
func getCardSuit(c card.Card) string {
	// Try to cast to PlayingCard
	if pc, ok := c.(card.PlayingCard); ok {
		return pc.Suit.String()
	}

	// Fallback to using the card's name to determine suit
	name := c.Name()
	if len(name) > 10 {
		suitPart := name[len(name)-6:]
		switch {
		case suitPart == "Spades":
			return "Spades"
		case suitPart == "Hearts":
			return "Hearts"
		case suitPart == "Diamonds":
			return "Diamonds"
		case suitPart == "Clubs":
			return "Clubs"
		}
	}

	return ""
}

// findBestHand finds the best 5-card hand from a set of cards
func findBestHand(cards []card.Card) HandResult {
	// Sort the cards by value
	sortedCards := make([]card.Card, len(cards))
	copy(sortedCards, cards)
	sort.Slice(sortedCards, func(i, j int) bool {
		return getCardValue(sortedCards[i]) > getCardValue(sortedCards[j])
	})

	// Generate all 5-card combinations
	combinations := generateCombinations(sortedCards, 5)

	// Evaluate each combination
	var bestResult HandResult
	for i, combo := range combinations {
		result := EvaluateHand(combo)
		if i == 0 || CompareHands(result, bestResult) > 0 {
			bestResult = result
		}
	}

	return bestResult
}

// generateCombinations generates all k-sized combinations from a set of cards
func generateCombinations(cards []card.Card, k int) [][]card.Card {
	if k == 0 || len(cards) < k {
		return [][]card.Card{}
	}

	if k == 1 {
		result := make([][]card.Card, len(cards))
		for i, c := range cards {
			result[i] = []card.Card{c}
		}
		return result
	}

	result := [][]card.Card{}

	// Include the first card
	for _, combo := range generateCombinations(cards[1:], k-1) {
		newCombo := append([]card.Card{cards[0]}, combo...)
		result = append(result, newCombo)
	}

	// Exclude the first card
	result = append(result, generateCombinations(cards[1:], k)...)

	return result
}

// isRoyalFlush checks if the hand is a royal flush
func isRoyalFlush(cards []card.Card) bool {
	if len(cards) < 5 {
		return false
	}

	// Check if it's a straight flush
	straightFlushValue := isStraightFlush(cards)
	if straightFlushValue == 0 {
		return false
	}

	// Check if the high card is an ace
	return getCardValue(cards[0]) == Ace
}

// isStraightFlush checks if the hand is a straight flush and returns the high card value
func isStraightFlush(cards []card.Card) int {
	if len(cards) < 5 {
		return 0
	}

	// Check if it's a flush
	if !isFlush(cards) {
		return 0
	}

	// Check if it's a straight
	return isStraight(cards)
}

// isFourOfAKind checks if the hand has four of a kind and returns the value and kicker
func isFourOfAKind(cards []card.Card) (int, int) {
	if len(cards) < 5 {
		return 0, 0
	}

	// Count the occurrences of each rank
	rankCounts := make(map[int]int)
	for _, c := range cards {
		rankCounts[getCardValue(c)]++
	}

	// Find the rank with 4 cards
	var fourOfAKindRank, kickerRank int
	for rank, count := range rankCounts {
		if count == 4 {
			fourOfAKindRank = rank
		} else if count == 1 && (kickerRank == 0 || rank > kickerRank) {
			kickerRank = rank
		}
	}

	if fourOfAKindRank > 0 {
		return fourOfAKindRank, kickerRank
	}

	return 0, 0
}

// isFullHouse checks if the hand is a full house and returns the three of a kind and pair values
func isFullHouse(cards []card.Card) (int, int) {
	if len(cards) < 5 {
		return 0, 0
	}

	// Count the occurrences of each rank
	rankCounts := make(map[int]int)
	for _, c := range cards {
		rankCounts[getCardValue(c)]++
	}

	// Find the ranks with 3 and 2 cards
	var threeOfAKindRank, pairRank int
	for rank, count := range rankCounts {
		if count == 3 && rank > threeOfAKindRank {
			threeOfAKindRank = rank
		} else if count >= 2 && rank > pairRank {
			pairRank = rank
		}
	}

	// Special case: if we have two three of a kinds, use the higher as the three of a kind
	// and the lower as the pair
	if threeOfAKindRank > 0 && pairRank == 0 {
		for rank, count := range rankCounts {
			if count == 3 && rank != threeOfAKindRank {
				pairRank = rank
				break
			}
		}
	}

	if threeOfAKindRank > 0 && pairRank > 0 {
		return threeOfAKindRank, pairRank
	}

	return 0, 0
}

// isFlush checks if the hand is a flush
func isFlush(cards []card.Card) bool {
	if len(cards) < 5 {
		return false
	}

	// Check if all cards have the same suit
	suit := getCardSuit(cards[0])
	for i := 1; i < 5; i++ {
		if getCardSuit(cards[i]) != suit {
			return false
		}
	}

	return true
}

// isStraight checks if the hand is a straight and returns the high card value
func isStraight(cards []card.Card) int {
	if len(cards) < 5 {
		return 0
	}

	// Get the unique ranks
	ranks := make(map[int]bool)
	for _, c := range cards {
		ranks[getCardValue(c)] = true
	}

	// Check for A-5-4-3-2 straight
	if ranks[Ace] && ranks[Five] && ranks[Four] &&
		ranks[Three] && ranks[Two] {
		return Five
	}

	// Check for regular straights
	for i := 0; i <= len(cards)-5; i++ {
		if getCardValue(cards[i])-getCardValue(cards[i+4]) == 4 {
			allDifferent := true
			for j := i; j < i+4; j++ {
				if getCardValue(cards[j]) == getCardValue(cards[j+1]) {
					allDifferent = false
					break
				}
			}
			if allDifferent {
				return getCardValue(cards[i])
			}
		}
	}

	return 0
}

// isThreeOfAKind checks if the hand has three of a kind and returns the value and kickers
func isThreeOfAKind(cards []card.Card) (int, []int) {
	if len(cards) < 5 {
		return 0, nil
	}

	// Count the occurrences of each rank
	rankCounts := make(map[int]int)
	for _, c := range cards {
		rankCounts[getCardValue(c)]++
	}

	// Find the rank with 3 cards
	var threeOfAKindRank int
	for rank, count := range rankCounts {
		if count == 3 {
			threeOfAKindRank = rank
			break
		}
	}

	if threeOfAKindRank == 0 {
		return 0, nil
	}

	// Find the kickers
	kickers := make([]int, 0, 2)
	for _, c := range cards {
		if getCardValue(c) != threeOfAKindRank && len(kickers) < 2 {
			kickers = append(kickers, getCardValue(c))
		}
	}

	return threeOfAKindRank, kickers
}

// isTwoPair checks if the hand has two pairs and returns the high pair, low pair, and kicker values
func isTwoPair(cards []card.Card) (int, int, int) {
	if len(cards) < 5 {
		return 0, 0, 0
	}

	// Count the occurrences of each rank
	rankCounts := make(map[int]int)
	for _, c := range cards {
		rankCounts[getCardValue(c)]++
	}

	// Find the ranks with 2 cards
	pairs := make([]int, 0, 2)
	for rank, count := range rankCounts {
		if count == 2 {
			pairs = append(pairs, rank)
		}
	}

	if len(pairs) < 2 {
		return 0, 0, 0
	}

	// Sort the pairs by rank
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i] > pairs[j]
	})

	// Find the kicker
	var kicker int
	for _, c := range cards {
		if getCardValue(c) != pairs[0] && getCardValue(c) != pairs[1] {
			kicker = getCardValue(c)
			break
		}
	}

	return pairs[0], pairs[1], kicker
}

// isOnePair checks if the hand has one pair and returns the pair value and kickers
func isOnePair(cards []card.Card) (int, []int) {
	if len(cards) < 5 {
		return 0, nil
	}

	// Count the occurrences of each rank
	rankCounts := make(map[int]int)
	for _, c := range cards {
		rankCounts[getCardValue(c)]++
	}

	// Find the rank with 2 cards
	var pairRank int
	for rank, count := range rankCounts {
		if count == 2 {
			pairRank = rank
			break
		}
	}

	if pairRank == 0 {
		return 0, nil
	}

	// Find the kickers
	kickers := make([]int, 0, 3)
	for _, c := range cards {
		if getCardValue(c) != pairRank && len(kickers) < 3 {
			kickers = append(kickers, getCardValue(c))
		}
	}

	return pairRank, kickers
}

// calculateHighCardValue calculates a value for high card hands
func calculateHighCardValue(cards []card.Card) int {
	if len(cards) < 5 {
		return 0
	}

	value := 0
	multiplier := 100000

	for i := 0; i < 5; i++ {
		value += getCardValue(cards[i]) * multiplier
		multiplier /= 100
	}

	return value
}

// rankToString converts a rank value to a string
func rankToString(rank int) string {
	switch rank {
	case Ace:
		return "Ace"
	case King:
		return "King"
	case Queen:
		return "Queen"
	case Jack:
		return "Jack"
	case Ten:
		return "10"
	case Nine:
		return "9"
	case Eight:
		return "8"
	case Seven:
		return "7"
	case Six:
		return "6"
	case Five:
		return "5"
	case Four:
		return "4"
	case Three:
		return "3"
	case Two:
		return "2"
	default:
		return fmt.Sprintf("%d", rank)
	}
}
