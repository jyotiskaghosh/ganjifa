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
	cb := match.CardBuilder{
		Name:    "Salavrka",
		Rank:    0,
		Civ:     civ.PRITHVI,
		Family:  family.Beast,
		Attack:  1,
		Defence: 1,
		Effects: []match.HandlerFunc{
			fx.Creature,
			func(card *match.Card, ctx *match.Context) {
				if cards, err := card.Player().Container(match.BATTLEZONE); err == nil {
					for _, c := range cards {
						if c != card && c.Family() == family.Beast {
							fx.AttackModifier(card, ctx, 1)
							fx.DefenceModifier(card, ctx, 1)
						}
					}
				}
			},
		},
	}

	return cb.Build()
}

// Vanara ...
func Vanara() *match.Card {
	cb := match.CardBuilder{
		Name:    "Vanara",
		Rank:    0,
		Civ:     civ.PRITHVI,
		Family:  family.Beast,
		Attack:  2,
		Defence: 2,
		Effects: []match.HandlerFunc{
			fx.Creature,
		},
	}

	return cb.Build()
}

// Krostr ...
func Krostr() *match.Card {
	cb := match.CardBuilder{
		Name:    "Krostr",
		Rank:    1,
		Civ:     civ.PRITHVI,
		Family:  family.Beast,
		Attack:  2,
		Defence: 2,
		Effects: []match.HandlerFunc{
			fx.Creature,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.ScheduleAfter(func() {
						cards, err := card.Player().Container(match.DECK)
						if err != nil {
							ctx.InterruptFlow()
							logrus.Debug(err)
							return
						}

						cards = card.Player().Search(
							match.Filter(cards, func(x *match.Card) bool { return x.Family() == family.Beast }),
							fmt.Sprintf("Select 1 %s", family.Beast),
							1,
							1,
							true)

						for _, c := range cards {
							if err := c.MoveCard(match.HAND); err != nil {
								logrus.Debug(err)
								return
							}

							ctx.Match().Chat("Server", fmt.Sprintf("%s was moved from %s's deck to their hand", c.Name(), card.Player().Name()))
						}

						card.Player().ShuffleDeck()
					})
				}
			},
		},
	}

	return cb.Build()
}

// Dvipin ...
func Dvipin() *match.Card {
	cb := match.CardBuilder{
		Name:    "Dvipin",
		Rank:    1,
		Civ:     civ.PRITHVI,
		Family:  family.Beast,
		Attack:  4,
		Defence: 2,
		Effects: []match.HandlerFunc{
			fx.Creature,
			fx.Ambush,
		},
	}

	return cb.Build()
}

// Simha ...
func Simha() *match.Card {
	cb := match.CardBuilder{
		Name:    "Simha",
		Rank:    2,
		Civ:     civ.PRITHVI,
		Family:  family.Beast,
		Attack:  9,
		Defence: 7,
		Effects: []match.HandlerFunc{
			fx.Creature,
		},
	}

	return cb.Build()
}
