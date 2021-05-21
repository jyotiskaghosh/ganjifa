package match

// BeginTurnStep ...
// Resolve any summoning sickness from creatures in the battle zone.
type BeginTurnStep struct{}

// UntapStep ...
// Your creatures in the battle zone and cards in your mana zone are untapped. This is forced.
type UntapStep struct{}

// StartOfTurnStep ...
// Any abilities that trigger at "the start of your turn" are resolved now.
type StartOfTurnStep struct{}

// DrawStep ...
// You draw a card. This is forced. If you have no deck, you lose when you draw the last card.
type DrawStep struct{}

// EndStep ...
// The turn finishes after you have no more creatures to attack with.
type EndStep struct{}
