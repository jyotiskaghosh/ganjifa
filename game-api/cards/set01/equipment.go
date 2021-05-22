package set01

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
	"github.com/jyotiskaghosh/ganjifa/game-api/family"
	"github.com/jyotiskaghosh/ganjifa/game-api/fx"
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Churika ...
func Churika() *match.Card {

	c := &match.Card{
		Name:   "Churika",
		Rank:   0,
		Civ:    civ.AGNI,
		Family: family.Equipment,
		Attack: 300,
	}

	c.Use(fx.Equipment)

	return c
}

// Khadga ...
func Khadga() *match.Card {

	c := &match.Card{
		Name:   "Khadga",
		Rank:   1,
		Civ:    civ.AGNI,
		Family: family.Equipment,
		Attack: 600,
	}

	c.Use(fx.Equipment)

	return c
}

// VampireFangs ...
func VampireFangs() *match.Card {

	c := &match.Card{
		Name:   "Vampire Fangs",
		Rank:   0,
		Civ:    civ.PRITHVI,
		Family: family.Equipment,
		Attack: 100,
	}

	c.Use(fx.Equipment, func(card *match.Card, ctx *match.Context) {

		if card.AttachedTo() != nil {
			fx.Leech(card.AttachedTo(), ctx)
		}
	})

	return c
}

// WindCloak ...
func WindCloak() *match.Card {

	c := &match.Card{
		Name:    "Wind Cloak",
		Rank:    0,
		Civ:     civ.VAYU,
		Family:  family.Equipment,
		Defence: 100,
	}

	c.Use(fx.Equipment, func(card *match.Card, ctx *match.Context) {

		if card.AttachedTo() != nil {
			fx.CantBeAttacked(card.AttachedTo(), ctx)
		}
	})

	return c
}

// ScopeLens ...
func ScopeLens() *match.Card {

	c := &match.Card{
		Name:   "Scope Lens",
		Rank:   0,
		Civ:    civ.VAYU,
		Family: family.Equipment,
		Attack: 200,
	}

	c.Use(fx.Equipment, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.GetAttackEvent); ok && event.Card == card.AttachedTo() {

			if _, ok := event.Event.(*match.AttackPlayer); ok {
				event.Attack += 200
			}
		}
	})

	return c
}

// ShellArmor ...
func ShellArmor() *match.Card {

	c := &match.Card{
		Name:    "Shell Armor",
		Rank:    0,
		Civ:     civ.APAS,
		Family:  family.Equipment,
		Defence: 300,
	}

	c.Use(fx.Equipment)

	return c
}
