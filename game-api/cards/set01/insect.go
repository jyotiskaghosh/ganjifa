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

	cb := match.CardBuilder{
		Name:    "Pipilika",
		Rank:    0,
		Civ:     civ.PRITHVI,
		Family:  family.Insect,
		Attack:  100,
		Defence: 100,
		Handlers: []match.HandlerFunc{
			fx.Creature,
			func(card *match.Card, ctx *match.Context) {
				if event, ok := ctx.Event().(*match.GetAttackEvent); ok && event.ID == card.ID() {
					ctx.ScheduleAfter(func() {
						event.Attack *= 2
					})
				}
			},
		},
	}

	return cb.Build()
}

// Masaka ...
func Masaka() *match.Card {

	cb := match.CardBuilder{
		Name:    "Masaka",
		Rank:    0,
		Civ:     civ.PRITHVI,
		Family:  family.Insect,
		Attack:  100,
		Defence: 100,
		Handlers: []match.HandlerFunc{
			fx.Creature,
			fx.Leech,
		},
	}

	return cb.Build()
}

// MahisiPipilika ...
func MahisiPipilika() *match.Card {

	cb := match.CardBuilder{
		Name:    "Mahisi Pipilika",
		Rank:    1,
		Civ:     civ.PRITHVI,
		Family:  family.Insect,
		Defence: 200,
		Handlers: []match.HandlerFunc{
			fx.Creature,
			func(card *match.Card, ctx *match.Context) {

				if card.Zone() != match.BATTLEZONE {
					return
				}

				if _, ok := ctx.Event().(*match.BeginTurnStep); ok && card.Player().IsPlayerTurn() {

					cards, err := card.Player().Container(match.DECK)
					if err != nil {
						ctx.InterruptFlow()
						logrus.Debug(err)
						return
					}

					cards = card.Player().Filter(
						cards,
						fmt.Sprintf("Select 1 %s", family.Insect),
						1,
						1,
						true,
						func(x *match.Card) bool { return x.Family() == family.Insect },
					)

					for _, c := range cards {
						if err := c.MoveCard(match.HAND); err != nil {
							logrus.Debug(err)
							return
						}
						ctx.Match().Chat("Server", fmt.Sprintf("%s was moved from %s's deck to their hand", c.Name(), card.Player().Name()))
					}

					card.Player().ShuffleDeck()
				}

				if event, ok := ctx.Event().(*match.GetAttackEvent); ok && event.ID != card.ID() {

					card, err := card.Player().GetCard(event.ID)
					if err != nil {
						logrus.Debug(err)
						return
					}

					if card.HasFamily(family.Insect, ctx) {
						event.Attack += 100
					}
				}

				if event, ok := ctx.Event().(*match.GetDefenceEvent); ok && event.ID != card.ID() {

					card, err := card.Player().GetCard(event.ID)
					if err != nil {
						logrus.Debug(err)
						return
					}

					if card.HasFamily(family.Insect, ctx) {
						event.Defence += 100
					}
				}
			},
		},
	}

	return cb.Build()
}
