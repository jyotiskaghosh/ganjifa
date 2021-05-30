package set01

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
	"github.com/jyotiskaghosh/ganjifa/game-api/family"
	"github.com/jyotiskaghosh/ganjifa/game-api/fx"
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Kukkutah ...
func Kukkutah() *match.Card {

	cb := match.CardBuilder{
		Name:    "Kukutah",
		Rank:    0,
		Civ:     civ.VAYU,
		Family:  family.Bird,
		Attack:  200,
		Defence: 200,
		Handlers: []match.HandlerFunc{
			fx.Creature,
		},
	}

	return cb.Build()
}

// Cataka ...
func Cataka() *match.Card {

	cb := match.CardBuilder{
		Name:    "Cataka",
		Rank:    0,
		Civ:     civ.VAYU,
		Family:  family.Bird,
		Attack:  100,
		Defence: 100,
		Handlers: []match.HandlerFunc{
			fx.Creature,
			fx.CantBeAttacked,
		},
	}

	return cb.Build()
}

// Syena ...
func Syena() *match.Card {

	cb := match.CardBuilder{
		Name:    "Syena",
		Rank:    1,
		Civ:     civ.VAYU,
		Family:  family.Bird,
		Attack:  300,
		Defence: 300,
		Handlers: []match.HandlerFunc{
			fx.Creature,
			func(card *match.Card, ctx *match.Context) {
				if event, ok := ctx.Event().(*match.GetAttackEvent); ok && event.ID == card.ID() {
					if _, ok := event.Event.(*match.AttackPlayer); ok {
						event.Attack += 400
					}
				}
			},
		},
	}

	return cb.Build()
}

// Atayi ...
func Atayi() *match.Card {

	cb := match.CardBuilder{
		Name:    "Atayi",
		Rank:    1,
		Civ:     civ.VAYU,
		Family:  family.Bird,
		Attack:  400,
		Defence: 200,
		Handlers: []match.HandlerFunc{
			fx.Creature,
			func(card *match.Card, ctx *match.Context) {
				if event, ok := ctx.Event().(*match.DamageEvent); ok && card.Zone() == match.BATTLEZONE && event.Player != card.Player() {
					card.AddCondition(fx.CantBeBlocked)
				}
			},
		},
	}

	return cb.Build()
}
