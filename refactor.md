# Card Game Framework Refactoring Plan
Creating a flexible architecture that would allow you to easily add other card games (like Uno, Texas Hold'em, or Pokémon) while maintaining generic base classes
This document outlines the plan for refactoring our card game codebase into a flexible framework that can support multiple card games (Texas Hold'em, Uno, Pokémon TCG, etc.) while maintaining a strong core foundation.

## Phase 1: Core Infrastructure ✅ COMPLETED

### 1.1 Define Core Interfaces ✅ COMPLETED
- [x] Create `pkg/cardgame/card` package with generic Card interface
- [x] Create `pkg/cardgame/deck` package with Deck interface and operations
- [x] Create `pkg/cardgame/hand` package with Hand interface
- [x] Create `pkg/cardgame/player` package with Player interface
- [x] Create `pkg/cardgame/game` package with Game interface and lifecycle methods

### 1.2 Implement Standard Playing Card Model ✅ COMPLETED
- [x] Implement StandardCard struct in `pkg/cardgame/card`
- [x] Implement StandardDeck in `pkg/cardgame/deck`
- [x] Implement secure shuffling using crypto/rand
- [x] Create utility functions for common card operations

### 1.3 Create Testing Framework (PENDING)
- [ ] Write unit tests for Card implementations
- [ ] Write unit tests for Deck operations
- [ ] Create test fixtures and helpers
- [ ] Implement property-based testing for shuffle operations

## Phase 2: Game State Management ✅ COMPLETED

### 2.1 Design State Management System ✅ COMPLETED
- [x] Create game state management in `pkg/cardgame/poker/holdem`
- [x] Implement state transition system
- [x] Create Action interface for player actions
- [x] Implement event system for game events

### 2.2 Implement Rule Engine ✅ COMPLETED
- [x] Implement game rules within specific game implementations
- [x] Create HandEvaluator interface for evaluating hands
- [x] Implement validation system for player actions

### 2.3 Create Configuration System (MODIFIED APPROACH)
- [x] Implement game configuration through constructor parameters
- [x] Create default configurations for different games
- [ ] Consider adding more robust configuration options in the future

## Phase 3: Refactor Texas Hold'em ✅ COMPLETED

### 3.1 Adapt Existing Hold'em Implementation ✅ COMPLETED
- [x] Move poker-specific card logic to `pkg/cardgame/poker`
- [x] Implement poker hand evaluation in `pkg/cardgame/poker/pokerhand`
- [x] Create Texas Hold'em implementation in `pkg/cardgame/poker/holdem`

### 3.2 Implement Hold'em Game States ✅ COMPLETED
- [x] Implement game state management in `pkg/cardgame/poker/holdem`
- [x] Implement dealing logic for hole cards and community cards
- [x] Implement betting rounds
- [x] Implement showdown and hand comparison

### 3.3 Create Hold'em Rules ✅ COMPLETED
- [x] Implement Texas Hold'em specific rules
- [x] Create betting rules
- [x] Implement hand ranking rules
- [x] Create dealer button and player turn management

## Phase 4: Implement Second Game (REVISED PRIORITY)

### 4.1 Design Additional Card Games
- [ ] Consider implementing Blackjack as the next game
- [ ] Create `pkg/cardgame/blackjack` package
- [ ] Implement BlackjackGame that implements Game interface
- [ ] Implement dealer logic and game rules

### 4.2 Implement Game States
- [ ] Create states for game flow
- [ ] Implement card dealing rules
- [ ] Implement player actions (hit, stand, etc.)
- [ ] Create win condition checks

### 4.3 Test Cross-Game Compatibility
- [ ] Ensure core interfaces work with multiple games
- [ ] Verify state management system flexibility
- [ ] Test player interfaces across different games
- [ ] Validate event system across games

## Phase 5: Testing and Documentation (REVISED PRIORITY)

### 5.1 Comprehensive Testing
- [ ] Add unit tests for all core components
- [ ] Create integration tests for game flows
- [ ] Test edge cases and error handling
- [ ] Implement benchmarks for performance-critical operations

### 5.2 Improve Documentation
- [ ] Add godoc comments to all exported types and functions
- [ ] Create usage examples for each game type
- [ ] Document interface contracts and expectations
- [ ] Create architecture diagrams

### 5.3 Refine API Design
- [ ] Review public APIs for consistency
- [ ] Ensure proper error handling throughout
- [ ] Consider backward compatibility for future changes
- [ ] Optimize for developer experience

## Phase 6: User Interface (FUTURE ENHANCEMENT)

### 6.1 Design Game Factory System
- [ ] Create factory pattern for game creation
- [ ] Implement game registry for available games
- [ ] Create unified game launcher

### 6.2 Implement CLI
- [ ] Create `cmd` package with CLI commands
- [ ] Implement game selection and configuration
- [ ] Create interactive game play mode
- [ ] Add help and documentation commands

### 6.3 Separate UI Concerns
- [ ] Create clean separation between game logic and presentation
- [ ] Implement observer pattern for UI updates
- [ ] Create text-based renderer for CLI
- [ ] Design interface for potential future GUI

## Future Enhancements

### Potential Future Features
- [ ] Network play support
- [ ] AI players
- [ ] Persistence layer for saving games
- [ ] Statistics and analytics
- [ ] Tournament mode
- [ ] Additional card games (Solitaire, Bridge, etc.)

## Notes

- ✅ Successfully migrated from internal package to pkg package structure
- ✅ Implemented clean interfaces with proper separation of concerns
- ✅ Created a working Texas Hold'em implementation
- Focus on adding comprehensive tests before implementing new games
- Consider implementing Blackjack next as it shares many components with poker
- Use dependency injection for flexibility
- Document design decisions and trade-offs