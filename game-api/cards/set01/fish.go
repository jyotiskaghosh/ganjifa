package set01

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
	"github.com/jyotiskaghosh/ganjifa/game-api/family"
	"github.com/jyotiskaghosh/ganjifa/game-api/fx"
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Matsyaka ...
func Matsyaka() *match.Card {

	c := &match.Card{
		Name:    "Matsyaka",
		Rank:    0,
		Civ:     civ.APAS,
		Family:  family.Fish,
		Attack:  100,
		Defence: 100,
	}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.GetAttackEvent); ok && event.Card == card.AttachedTo() {
			event.Attack += 100
		}

		if event, ok := ctx.Event.(*match.GetDefenceEvent); ok && event.Card == card.AttachedTo() {
			event.Defence += 100
		}
	})

	return c
}

// DeadlyZebrafish ...
func DeadlyZebrafish() *match.Card {

	c := &match.Card{
		Name:    "Deadly Zebrafish",
		Rank:    0,
		Civ:     civ.APAS,
		Family:  family.Fish,
		Attack:  100,
		Defence: 100,
	}

	c.Use(fx.Creature, fx.Poisonous)

	return c
}

// TorpedoingBarracuda ...
func TorpedoingBarracuda() *match.Card {

	c := &match.Card{
		Name:    "Torpedoing Barracuda",
		Rank:    1,
		Civ:     civ.APAS,
		Family:  family.Fish,
		Attack:  300,
		Defence: 100,
	}

	c.Use(fx.Creature, fx.CantBeBlocked)

	return c
}