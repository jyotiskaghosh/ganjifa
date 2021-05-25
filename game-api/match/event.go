package match

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
)

// EndTurnEvent is fired when a player attempts to end their turn
type EndTurnEvent struct {
}

// PlayCardEvent is fired when the player attempts to play a card
type PlayCardEvent struct {
	ID string
}

// CardMoved is fired from the *Player.MoveCard method after moving a card between containers
type CardMoved struct {
	ID   string
	From Container
	To   Container
}

// Evolve is fired when a creature is evolving
type Evolve struct {
	ID       string
	Creature *Card
}

// Equip is fired when a creature is being equiped
type Equip struct {
	ID       string
	Creature *Card
}

// AttackPlayer is fired when the player attempts to use a creature to attack the player
type AttackPlayer struct {
	ID string
}

// AttackCreature is fired when the player attempts to use a creature to attack a creature
type AttackCreature struct {
	ID       string
	TargetID string
}

// React is fired when you play cards from hiddenzone during opponent's attack
type React struct {
	ID    string
	Event interface{}
}

// Block is fired when a creature attempts
type Block struct {
	Attacker *Card
	Blocker  *Card
}

// Battle is fired when two creatures are fighting, i.e. from attacking a creature or blocking an attack
type Battle struct {
	Attacker *Card
	Defender *Card
	Blocked  bool
}

// CreatureDestroyed is fired when a creature dies in battle or is destroyed from another source, such as a spell
type CreatureDestroyed struct {
	Card   *Card
	Source *Card
}

// DamageEvent is fired when a player takes damage
type DamageEvent struct {
	Player *Player
	Source *Card
	Event  interface{}
	Health int
}

// HealEvent is fired when a player heals health
type HealEvent struct {
	Player *Player
	Source *Card
	Event  interface{}
	Health int
}

// GetAttackEvent is fired whenever a cards attack is to be used
type GetAttackEvent struct {
	Card   *Card
	Event  interface{}
	Attack int
}

// GetDefenceEvent is fired whenever a cards defence is to be used
type GetDefenceEvent struct {
	Card    *Card
	Event   interface{}
	Defence int
}

// GetCivilisationEvent is fired whenever a cards civ is to be used
type GetCivilisationEvent struct {
	Card  *Card
	Event interface{}
	Civ   map[civ.Civilisation]bool
}

// GetFamilyEvent is fired whenever a cards family is to be used
type GetFamilyEvent struct {
	Card   *Card
	Event  interface{}
	Family map[string]bool
}

// GetRankEvent is fired whenever a cards rank is to be used
type GetRankEvent struct {
	Card  *Card
	Event interface{}
	Rank  int
}
