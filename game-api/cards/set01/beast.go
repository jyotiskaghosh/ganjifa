package set01

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
	"github.com/jyotiskaghosh/ganjifa/game-api/family"
	"github.com/jyotiskaghosh/ganjifa/game-api/fx"
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
	"github.com/sirupsen/logrus"
)

// Salavrka ...
func Salavrka() *match.Card {

	c := &match.Card{
		Name:   "Salavrka",
		Rank:   0,
		Civ:    civ.PRITHVI,
		Family: family.Beast,
	}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if creatures, err := card.Player.Container(match.BATTLEZONE); err == nil {
			for _, c := range creatures {
				if c.Family == family.Beast {
					fx.AttackModifier(card, ctx, 100)
					fx.DefenceModifier(card, ctx, 100)
				}
			}
		}
	})

	return c
}

// Vanara ...
func Vanara() *match.Card {

	c := &match.Card{
		Name:    "Vanara",
		Rank:    0,
		Civ:     civ.PRITHVI,
		Family:  family.Beast,
		Attack:  200,
		Defence: 200,
	}

	c.Use(fx.Creature)

	return c
}

// Krostr ...
func Krostr() *match.Card {

	c := &match.Card{
		Name:    "Krostr",
		Rank:    1,
		Civ:     civ.PRITHVI,
		Family:  family.Beast,
		Attack:  200,
		Defence: 200,
	}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.ScheduleAfter(func() {
				cards, err := card.Player.Container(match.DECK)
				if err != nil {
					ctx.InterruptFlow()
					logrus.Debug(err)
					return
				}

				creatures := card.Player.Filter(
					cards,
					fmt.Sprintf("Select 1 %s", family.Beast),
					1,
					1,
					true,
					func(x *match.Card) bool { return x.Family == family.Beast },
				)

				for _, c := range creatures {
					if err := c.MoveCard(match.HAND); err != nil {
						logrus.Debug(err)
						return
					}
					ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's deck to their hand", c.Name, card.Player.Name()))
				}

				card.Player.ShuffleDeck()
			})
		}
	})

	return c
}

// Dvipin ...
func Dvipin() *match.Card {

	c := &match.Card{
		Name:    "Dvipin",
		Rank:    1,
		Civ:     civ.PRITHVI,
		Family:  family.Beast,
		Attack:  300,
		Defence: 100,
	}

	c.Use(fx.Creature, fx.Ambush)

	return c
}

// Simha ...
func Simha() *match.Card {

	c := &match.Card{
		Name:    "Simha",
		Rank:    2,
		Civ:     civ.PRITHVI,
		Family:  family.Beast,
		Attack:  900,
		Defence: 700,
	}

	c.Use(fx.Creature)

	return c
}
