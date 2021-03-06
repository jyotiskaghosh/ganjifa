package match

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
)

// EndTurnEvent is fired when a player attempts to end their turn
type EndTurnEvent struct{}

// PlayCardEvent is fired when the player attempts to play a card
type PlayCardEvent struct {
	ID string
}

// Evolve is fired when a creature is evolving
type Evolve struct {
	ID     string
	Target *Card
}

// Equip is fired when a creature is being equiped
type Equip struct {
	ID     string
	Target *Card
}

// SpellCast is fired when a spell is cast
type SpellCast struct {
	ID      string
	Targets []*Card
}

// CardMoved is fired from the *Player.MoveCard method after moving a card between containers
type CardMoved struct {
	ID   string
	From Container
	To   Container
}

// AttackPlayer is fired when the player attempts to use a creature to attack the opponent
type AttackPlayer struct {
	ID string
}

// AttackCreature is fired when the player attempts to use a creature to attack a creature
type AttackCreature struct {
	ID       string
	TargetID string
}

// TrapEvent is fired when you play cards from hiddenzone during opponent's attack
type TrapEvent struct {
	ID       string
	Attacker *Card
}

// BlockEvent is fired when a creature attempts to block an incoming attack
type BlockEvent struct {
	ID       string
	Attacker *Card
}

// Battle is fired when two creatures are fighting, i.e. from attacking a creature or blocking an attack
type Battle struct {
	Attacker *Card
	Defender *Card
	Blocked  bool
}

// CreatureDestroyed is fired when a creature dies in battle or is destroyed from another source, such as a spell
type CreatureDestroyed struct {
	ID     string
	Source *Card
}

// DamageEvent is fired when a player takes damage
type DamageEvent struct {
	Player *Player
	Source *Card
	Event  interface{}
	Health uint8
}

// HealEvent is fired when a player heals health
type HealEvent struct {
	Player *Player
	Source *Card
	Event  interface{}
	Health uint8
}

// GetAttackEvent is fired whenever a card's attack is to be used
type GetAttackEvent struct {
	ID     string
	Event  interface{}
	Attack uint8
}

// GetDefenceEvent is fired whenever a card's defence is to be used
type GetDefenceEvent struct {
	ID      string
	Event   interface{}
	Defence uint8
}

// GetCivilisationEvent is fired whenever a card's civ is to be used
type GetCivilisationEvent struct {
	ID    string
	Event interface{}
	Civ   map[civ.Civilisation]bool
}

// GetFamilyEvent is fired whenever a card's family is to be used
type GetFamilyEvent struct {
	ID     string
	Event  interface{}
	Family map[string]bool
}

// GetRankEvent is fired whenever a card's rank is to be used
type GetRankEvent struct {
	ID    string
	Event interface{}
	Rank  uint8
}

// GetHandlerEvent is fired whenever a card's handlers are to be used
type GetHandlerEvent struct {
	ID       string
	Event    interface{}
	Handlers []HandlerFunc
}
