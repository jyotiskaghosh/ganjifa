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
		Attack: 300,
		Handlers: []match.HandlerFunc{
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
		Attack: 600,
		Handlers: []match.HandlerFunc{
			fx.Equipment,
		},
	}

	return cb.Build()
}

// VampireFangs ...
func VampireFangs() *match.Card {

	cb := match.CardBuilder{
		Name:   "Vampire Fangs",
		Rank:   0,
		Civ:    civ.PRITHVI,
		Family: family.Equipment,
		Attack: 100,
		Handlers: []match.HandlerFunc{
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
		Name:    "Wind Cloak",
		Rank:    0,
		Civ:     civ.VAYU,
		Family:  family.Equipment,
		Defence: 100,
		Handlers: []match.HandlerFunc{
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
		Attack: 200,
		Handlers: []match.HandlerFunc{
			fx.Equipment,
			func(card *match.Card, ctx *match.Context) {

				if event, ok := ctx.Event().(*match.GetAttackEvent); ok && event.Card == card.AttachedTo() {

					if _, ok := event.Event.(*match.AttackPlayer); ok {
						event.Attack += 200
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
		Defence: 300,
		Handlers: []match.HandlerFunc{
			fx.Equipment,
		},
	}

	return cb.Build()
}
