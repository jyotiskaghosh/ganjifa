package set01

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
	"github.com/jyotiskaghosh/ganjifa/game-api/family"
	"github.com/jyotiskaghosh/ganjifa/game-api/fx"
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
	"github.com/sirupsen/logrus"
)

// Pipilika ...
func Pipilika() *match.Card {

	c := &match.Card{
		Name:    "Pipilika",
		Rank:    0,
		Civ:     civ.PRITHVI,
		Family:  family.Insect,
		Attack:  100,
		Defence: 100,
	}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if event, ok := ctx.Event.(*match.GetAttackEvent); ok && event.Card == card {

			ctx.ScheduleAfter(func() {
				event.Attack *= 2
			})
		}
	})

	return c
}

// Masaka ...
func Masaka() *match.Card {

	c := &match.Card{
		Name:    "Masaka",
		Rank:    0,
		Civ:     civ.PRITHVI,
		Family:  family.Insect,
		Attack:  100,
		Defence: 100,
	}

	c.Use(fx.Creature, fx.Leech)

	return c
}

// MahisiPipilika ...
func MahisiPipilika() *match.Card {

	c := &match.Card{
		Name:    "Mahisi Pipilika",
		Rank:    1,
		Civ:     civ.PRITHVI,
		Family:  family.Insect,
		Defence: 200,
	}

	c.Use(fx.Creature, func(card *match.Card, ctx *match.Context) {

		if card.Zone != match.BATTLEZONE {
			return
		}

		if _, ok := ctx.Event.(*match.BeginTurnStep); ok && card.Player.IsPlayerTurn() {

			cards, err := card.Player.Container(match.DECK)
			if err != nil {
				ctx.InterruptFlow()
				logrus.Debug(err)
				return
			}

			creatures := card.Player.Filter(
				cards,
				fmt.Sprintf("Select 1 %s", family.Insect),
				1,
				1,
				true,
				func(x *match.Card) bool { return x.Family == family.Insect },
			)

			for _, c := range creatures {
				if err := c.MoveCard(match.HAND); err != nil {
					logrus.Debug(err)
					return
				}
				ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's deck to their hand", c.Name, card.Player.Name()))
			}

			card.Player.ShuffleDeck()
		}

		if event, ok := ctx.Event.(*match.GetAttackEvent); ok && event.Card != card {

			if event.Card.Player == card.Player && event.Card.HasFamily(family.Insect, ctx) {
				event.Attack += 100
			}
		}

		if event, ok := ctx.Event.(*match.GetDefenceEvent); ok && event.Card != card {

			if event.Card.Player == card.Player && event.Card.HasFamily(family.Insect, ctx) {
				event.Defence += 100
			}
		}
	})

	return c
}