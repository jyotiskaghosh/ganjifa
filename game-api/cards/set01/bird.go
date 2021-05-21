package set01

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
	"github.com/jyotiskaghosh/ganjifa/game-api/family"
	"github.com/jyotiskaghosh/ganjifa/game-api/fx"
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Kukkutah ...
func Kukkutah() *match.Card {

	c := &match.Card{
		Name:    "Kukutah",
		Rank:    0,
		Civ:     civ.VAYU,
		Family:  family.Bird,
		Attack:  200,
		Defence: 200,
	}

	c.Use(fx.Creature)

	return c
}

// Cataka ...
func Cataka() *match.Card {

	c := &match.Card{
		Name:    "Cataka",
		Rank:    0,
		Civ:     civ.VAYU,
		Family:  family.Bird,
		Attack:  100,
		Defence: 100,
	}

	c.Use(fx.Creature, fx.CantBeAttacked)

	return c
}

// Syena ...
func Syena() *match.Card {

	c := &match.Card{
		Name:    "Syena",
		Rank:    1,
		Civ:     civ.VAYU,
		Family:  family.Bird,
		Attack:  300,
		Defence: 300,
	}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.GetAttackEvent); ok && event.Card == card {

			if _, ok := event.Event.(*match.AttackPlayer); ok {
				event.Attack += 400
			}
		}
	})

	return c
}

// Atayi ...
func Atayi() *match.Card {

	c := &match.Card{
		Name:    "Atayi",
		Rank:    1,
		Civ:     civ.VAYU,
		Family:  family.Bird,
		Attack:  400,
		Defence: 200,
	}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.DamageEvent); ok && card.Zone == match.BATTLEZONE && event.Player != card.Player {
			c.AddCondition(fx.CantBeBlocked)
		}
	})

	return c
}
