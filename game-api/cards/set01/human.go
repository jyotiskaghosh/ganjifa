package set01

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
	"github.com/jyotiskaghosh/ganjifa/game-api/family"
	"github.com/jyotiskaghosh/ganjifa/game-api/fx"
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
	"github.com/sirupsen/logrus"
)

// Ayudhabhrt ...
func Ayudhabhrt() *match.Card {

	c := &match.Card{
		Name:    "Ayudhabhrt",
		Rank:    0,
		Civ:     civ.AGNI,
		Family:  family.Human,
		Attack:  200,
		Defence: 100,
	}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		for _, c := range card.Attachments() {
			if c.Family == family.Equipment {
				fx.AttackModifier(card, ctx, 100)
				fx.DefenceModifier(card, ctx, 100)
				return
			}
		}
	})

	return c
}

// Sainika ...
func Sainika() *match.Card {

	c := &match.Card{
		Name:    "Sainika",
		Rank:    0,
		Civ:     civ.AGNI,
		Family:  family.Human,
		Attack:  200,
		Defence: 200,
	}

	c.Use(fx.Creature)

	return c
}

// Sastravikrayin ...
func Sastravikrayin() *match.Card {

	c := &match.Card{
		Name:    "Sastravikrayin",
		Rank:    1,
		Civ:     civ.AGNI,
		Family:  family.Human,
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
					fmt.Sprintf("Select 1 %s", family.Equipment),
					1,
					1,
					true,
					func(x *match.Card) bool { return x.Family == family.Equipment },
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

// Astrakara ...
func Astrakara() *match.Card {

	c := &match.Card{
		Name:    "Astrakara",
		Rank:    1,
		Civ:     civ.AGNI,
		Family:  family.Human,
		Attack:  200,
		Defence: 200,
	}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if event, ok := ctx.Event.(*match.GetRankEvent); ok && event.Card.Player == card.Player {

			if event.Card.Family == family.Equipment {
				if event.Rank > 0 {
					event.Rank--
				}
			}
		}
	})

	return c
}