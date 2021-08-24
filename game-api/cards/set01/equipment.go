package set01

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
	"github.com/jyotiskaghosh/ganjifa/game-api/family"
	"github.com/jyotiskaghosh/ganjifa/game-api/fx"
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Churika ...
func Churika() *match.Card {
	cb := match.CardBuilder{
		Name:   "Churika",
		Rank:   0,
		Civ:    civ.AGNI,
		Family: family.Equipment,
		Attack: 2,
		Effects: []match.HandlerFunc{
			fx.Equipment,
		},
	}

	return cb.Build()
}

// Khadga ...
func Khadga() *match.Card {
	cb := match.CardBuilder{
		Name:   "Khadga",
		Rank:   1,
		Civ:    civ.AGNI,
		Family: family.Equipment,
		Attack: 4,
		Effects: []match.HandlerFunc{
			fx.Equipment,
		},
	}

	return cb.Build()
}

// VampireFangs ...
func VampireFangs() *match.Card {
	cb := match.CardBuilder{
		Name:   "Vampire Fangs",
		Rank:   1,
		Civ:    civ.PRITHVI,
		Family: family.Equipment,
		Attack: 2,
		Effects: []match.HandlerFunc{
			fx.Equipment,
			func(card *match.Card, ctx *match.Context) {
				if card.AttachedTo() != nil {
					fx.Leech(card.AttachedTo(), ctx)
				}
			},
		},
	}

	return cb.Build()
}

// WindCloak ...
func WindCloak() *match.Card {
	cb := match.CardBuilder{
		Name:   "Wind Cloak",
		Rank:   0,
		Civ:    civ.VAYU,
		Family: family.Equipment,
		Effects: []match.HandlerFunc{
			fx.Equipment,
			func(card *match.Card, ctx *match.Context) {
				if card.AttachedTo() != nil {
					fx.CantBeAttacked(card.AttachedTo(), ctx)
				}
			},
		},
	}

	return cb.Build()
}

// ScopeLens ...
func ScopeLens() *match.Card {
	cb := match.CardBuilder{
		Name:   "Scope Lens",
		Rank:   0,
		Civ:    civ.VAYU,
		Family: family.Equipment,
		Attack: 1,
		Effects: []match.HandlerFunc{
			fx.Equipment,
			func(card *match.Card, ctx *match.Context) {
				if event, ok := ctx.Event().(*match.GetAttackEvent); ok &&
					card.AttachedTo() != nil && event.ID == card.AttachedTo().ID() {
					if _, ok := event.Event.(*match.AttackPlayer); ok {
						event.Attack += 2
					}
				}
			},
		},
	}

	return cb.Build()
}

// ShellArmor ...
func ShellArmor() *match.Card {
	cb := match.CardBuilder{
		Name:    "Shell Armor",
		Rank:    0,
		Civ:     civ.APAS,
		Family:  family.Equipment,
		Defence: 2,
		Effects: []match.HandlerFunc{
			fx.Equipment,
		},
	}

	return cb.Build()
}
