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
	cb := match.CardBuilder{
		Name:    "Ayudhabhrt",
		Rank:    0,
		Civ:     civ.AGNI,
		Family:  family.Human,
		Attack:  200,
		Defence: 100,
		Handlers: []match.HandlerFunc{
			fx.Creature,
			func(card *match.Card, ctx *match.Context) {
				for _, c := range card.Attachments() {
					if c.Family() == family.Equipment {
						fx.AttackModifier(card, ctx, 100)
						fx.DefenceModifier(card, ctx, 100)
					}
				}
			},
		},
	}

	return cb.Build()
}

// Sainika ...
func Sainika() *match.Card {
	cb := match.CardBuilder{
		Name:    "Sainika",
		Rank:    0,
		Civ:     civ.AGNI,
		Family:  family.Human,
		Attack:  200,
		Defence: 200,
		Handlers: []match.HandlerFunc{
			fx.Creature,
		},
	}

	return cb.Build()
}

// Sastravikrayin ...
func Sastravikrayin() *match.Card {
	cb := match.CardBuilder{
		Name:    "Sastravikrayin",
		Rank:    1,
		Civ:     civ.AGNI,
		Family:  family.Human,
		Attack:  200,
		Defence: 200,
		Handlers: []match.HandlerFunc{
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

						cards = card.Player().SearchAction(
							match.Filter(cards, func(x *match.Card) bool { return x.Family() == family.Equipment }),
							fmt.Sprintf("Select 1 %s", family.Equipment),
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

// Astrakara ...
func Astrakara() *match.Card {
	cb := match.CardBuilder{
		Name:    "Astrakara",
		Rank:    1,
		Civ:     civ.AGNI,
		Family:  family.Human,
		Attack:  200,
		Defence: 200,
		Handlers: []match.HandlerFunc{
			fx.Creature,
			func(card *match.Card, ctx *match.Context) {
				if event, ok := ctx.Event().(*match.GetRankEvent); ok && card.Zone() == match.BATTLEZONE {
					card, err := card.Player().GetCard(event.ID)
					if err != nil {
						logrus.Debug(err)
						return
					}

					if card.Family() == family.Equipment && event.Rank > 0 {
						event.Rank--
					}
				}
			},
		},
	}

	return cb.Build()
}
