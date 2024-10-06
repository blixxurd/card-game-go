package pokerhand

import (
	"fmt"
	"log"
	"sort"

	"github.com/blixxurd/card-game-go/internal/cardgame/card"
)

type HandRank int

const (
	HighCard HandRank = iota
	Pair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

type HandResult struct {
	Rank      HandRank
	Name      string
	Cards     []card.Card
	HighCards []int
}

/**
 * Takes a slice of cards and returns the best possible hand
 * that can be made with those cards. The function generates all possible
 * combinations of 5 cards from the input slice and evaluates each combination
 * to determine the best hand. The function returns a HandResult struct
 * containing the rank of the best hand, the name of the hand, the cards that
 * make up the hand, and the high cards used to break ties.
 */
func EvaluateBestHand(cards []card.Card) (HandResult, error) {
	if len(cards) < 5 {
		return HandResult{}, fmt.Errorf("not enough cards to evaluate hand")
	}

	combinations := generateCombinations(cards, 5)

	// Assume all hands have a high card as the best hand to start
	var bestHand HandResult = HandResult{Rank: HighCard, Name: "High Card", Cards: combinations[0], HighCards: getHighCards(combinations[0], 5)}

	for _, combo := range combinations {
		result := evaluateHand(combo)
		comparison := CompareHands(result, bestHand)
		// If there is a clear winner, update the best hand
		if comparison > 0 {
			bestHand = result
		}
	}

	log.Printf("Best hand: %v", bestHand)

	return bestHand, nil
}

/**
 * Takes a slice of cards and evaluates the best possible hand
 * that can be made with those cards. The function checks for the highest
 * possible hand starting from the best possible hand and working down to the
 * lowest possible hand. The function returns a HandResult struct containing
 * the rank of the best hand, the name of the hand, the cards that make up the
 * hand, and the high cards used to break ties.
 */
func evaluateHand(hand []card.Card) HandResult {
	sortedHand := make([]card.Card, len(hand))
	copy(sortedHand, hand)
	sort.Slice(sortedHand, func(i, j int) bool {
		return getComparisonValue(sortedHand[i]) > getComparisonValue(sortedHand[j])
	})

	isFlush := checkFlush(sortedHand)
	isStraight, highCard := checkStraight(sortedHand)

	if isFlush && isStraight {
		if highCard == 14 { // Ace high
			return HandResult{Rank: RoyalFlush, Name: "Royal Flush", Cards: sortedHand, HighCards: []int{14}}
		}
		return HandResult{Rank: StraightFlush, Name: "Straight Flush", Cards: sortedHand, HighCards: []int{highCard}}
	}

	if isFlush {
		return HandResult{Rank: Flush, Name: "Flush", Cards: sortedHand, HighCards: getHighCards(sortedHand, 5)}
	}

	if isStraight {
		return HandResult{Rank: Straight, Name: "Straight", Cards: sortedHand, HighCards: []int{highCard}}
	}

	valueCounts := countValues(sortedHand)

	if hasFourOfAKind(valueCounts) {
		return HandResult{Rank: FourOfAKind, Name: "Four of a Kind", Cards: sortedHand, HighCards: getHighCards(sortedHand, 2)}
	}

	if hasFullHouse(valueCounts) {
		return HandResult{Rank: FullHouse, Name: "Full House", Cards: sortedHand, HighCards: getHighCards(sortedHand, 2)}
	}

	if hasThreeOfAKind(valueCounts) {
		return HandResult{Rank: ThreeOfAKind, Name: "Three of a Kind", Cards: sortedHand, HighCards: getHighCards(sortedHand, 3)}
	}

	pairCount := countPairs(valueCounts)
	if pairCount == 2 {
		return HandResult{Rank: TwoPair, Name: "Two Pair", Cards: sortedHand, HighCards: getHighCards(sortedHand, 3)}
	}

	if pairCount == 1 {
		return HandResult{Rank: Pair, Name: "Pair", Cards: sortedHand, HighCards: getHighCards(sortedHand, 4)}
	}

	highCard = getComparisonValue(sortedHand[0])
	return HandResult{Rank: HighCard, Name: fmt.Sprintf("High Card %s", cardValueToString(highCard)), Cards: sortedHand, HighCards: getHighCards(sortedHand, 5)}
}

/**
 * Returns the comparison value of a card. The comparison value is the
 * card's value, except for an Ace, which is assigned a value of 14.
 */
func getComparisonValue(card card.Card) int {
	if card.Value == 1 { // Ace
		return 14
	}
	return card.Value
}

/**
 * Checks if all cards in a hand have the same suit.
 */
func checkFlush(hand []card.Card) bool {
	suit := hand[0].Suit
	for _, card := range hand[1:] {
		if card.Suit != suit {
			return false
		}
	}
	return true
}

/**
 * Checks if a hand is a straight. A straight is a hand where all cards
 * have consecutive values. The function also checks for an Ace-low straight
 * ie. A, 2, 3, 4, 5.
 */
func checkStraight(hand []card.Card) (bool, int) {
	values := make([]int, len(hand))
	for i, card := range hand {
		values[i] = getComparisonValue(card)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values)))

	// Check for regular straight
	isStraight := true
	for i := 1; i < len(values); i++ {
		if values[i-1] != values[i]+1 {
			isStraight = false
			break
		}
	}
	if isStraight {
		return true, values[0]
	}

	// Check for Ace-low straight (A, 2, 3, 4, 5)
	if values[0] == 14 && values[1] == 5 && values[2] == 4 && values[3] == 3 && values[4] == 2 {
		return true, 5 // Return 5 as the high card for A-2-3-4-5 straight
	}

	return false, 0
}

/**
 * Counts the number of cards with the same value in a hand.
 */
func countValues(hand []card.Card) map[int]int {
	counts := make(map[int]int)
	for _, card := range hand {
		counts[getComparisonValue(card)]++
	}
	return counts
}

/**
 * Checks if a hand has four of a kind.
 */
func hasFourOfAKind(counts map[int]int) bool {
	for _, count := range counts {
		if count == 4 {
			return true
		}
	}
	return false
}

/**
 * Checks if a hand has a full house.
 */
func hasFullHouse(counts map[int]int) bool {
	hasThree, hasTwo := false, false
	for _, count := range counts {
		if count == 3 {
			hasThree = true
		} else if count == 2 {
			hasTwo = true
		}
	}
	return hasThree && hasTwo
}

/**
 * Checks if a hand has three of a kind.
 */
func hasThreeOfAKind(counts map[int]int) bool {
	for _, count := range counts {
		if count == 3 {
			return true
		}
	}
	return false
}

/**
 * Counts the number of pairs in a hand.
 */
func countPairs(counts map[int]int) int {
	pairs := 0
	for _, count := range counts {
		if count == 2 {
			pairs++
		}
	}
	return pairs
}

/**
 * Returns the high cards in a hand. The number of high cards returned
 * is specified by the count parameter.
 */
func getHighCards(hand []card.Card, count int) []int {
	highCards := make([]int, 0, count)
	for i := 0; i < count && i < len(hand); i++ {
		highCards = append(highCards, getComparisonValue(hand[i]))
	}
	return highCards
}

/**
 * Converts a card value to a string representation.
 */
func cardValueToString(value int) string {
	switch value {
	case 14:
		return "A"
	case 13:
		return "K"
	case 12:
		return "Q"
	case 11:
		return "J"
	default:
		return fmt.Sprintf("%d", value)
	}
}

/**
 * Compares two HandResult structs and returns an integer value indicating
 * the result of the comparison. The function returns a negative value if
 * hand1 is less than hand2, a positive value if hand1 is greater than hand2,
 * and zero if the two hands are equal.
 */
func CompareHands(hand1, hand2 HandResult) int {
	// Compare ranks
	if hand1.Rank != hand2.Rank {
		return int(hand1.Rank) - int(hand2.Rank)
	}

	// Compare high cards
	for i := 0; i < len(hand1.HighCards) && i < len(hand2.HighCards); i++ {
		if hand1.HighCards[i] != hand2.HighCards[i] {
			return hand1.HighCards[i] - hand2.HighCards[i]
		}
	}

	return 0
}

/**
 * Generates all possible combinations of k cards from a slice of cards.
 * Uses a recursive backtracking algorithm to generate the combinations.
 * This method is used to generate all possible 5-card hands from a set of
 * cards.
 */
func generateCombinations(cards []card.Card, k int) [][]card.Card {
	var combos [][]card.Card
	var combo []card.Card
	var generate func(int, int)

	generate = func(start, k int) {
		if k == 0 {
			temp := make([]card.Card, len(combo))
			copy(temp, combo)
			combos = append(combos, temp)
			return
		}

		for i := start; i <= len(cards)-k; i++ {
			combo = append(combo, cards[i])
			generate(i+1, k-1)
			combo = combo[:len(combo)-1]
		}
	}

	generate(0, k)
	return combos
}
