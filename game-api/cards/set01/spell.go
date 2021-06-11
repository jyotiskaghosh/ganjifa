package set01

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
	"github.com/jyotiskaghosh/ganjifa/game-api/family"
	"github.com/jyotiskaghosh/ganjifa/game-api/fx"
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// EnergySurge ...
func EnergySurge() *match.Card {

	cb := match.CardBuilder{
		Name:   "Energy Surge",
		Rank:   0,
		Civ:    civ.AGNI,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if event, ok := ctx.Event().(*match.PlayCardEvent); ok && event.ID == card.ID() {

					if len(event.Targets) <= 0 {
						ctx.InterruptFlow()
						return
					}

					target, err := card.Player().GetCard(event.Targets[0])
					if err != nil {
						ctx.InterruptFlow()
						logrus.Debug(err)
						return
					}

					ctx.Override(func() {
						target.AddCondition(func(card *match.Card, ctx *match.Context) {
							fx.AttackModifier(card, ctx, 400)
						})
					})
				}
			},
		},
	}

	return cb.Build()
}

// Fireball ...
func Fireball() *match.Card {

	cb := match.CardBuilder{
		Name:   "Fireball",
		Rank:   0,
		Civ:    civ.AGNI,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {
						if event, ok := ctx.Event().(*match.PlayCardEvent); ok && event.ID == card.ID() {

							if len(event.Targets) <= 0 {
								ctx.InterruptFlow()
								return
							}

							target, err := card.Player().GetCard(event.Targets[0])
							if err != nil {
								ctx.InterruptFlow()
								logrus.Debug(err)
								return
							}

							if target.GetDefence(ctx) <= 200 {
								ctx.Match().Destroy(target, card)
							} else {
								ctx.InterruptFlow()
							}
						}
					})
				}
			},
		},
	}

	return cb.Build()
}

// RainOfArrows ...
func RainOfArrows() *match.Card {

	cb := match.CardBuilder{
		Name:   "Rain Of Arrows",
		Rank:   0,
		Civ:    civ.AGNI,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {
						for _, c := range append(card.Player().GetCreatures(), ctx.Match().Opponent(card.Player()).GetCreatures()...) {
							if c.GetDefence(ctx) <= 100 {
								ctx.Match().Destroy(c, card)
							}
						}
					})
				}
			},
		},
	}

	return cb.Build()
}

// MagmaGeyser ...
func MagmaGeyser() *match.Card {

	cb := match.CardBuilder{
		Name:   "MagmaGeyser",
		Rank:   1,
		Civ:    civ.AGNI,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {
						if event, ok := ctx.Event().(*match.PlayCardEvent); ok && event.ID == card.ID() {

							if len(event.Targets) <= 0 {
								ctx.InterruptFlow()
								return
							}

							target, err := card.Player().GetCard(event.Targets[0])
							if err != nil {
								ctx.InterruptFlow()
								logrus.Debug(err)
								return
							}

							if target.GetDefence(ctx) <= 400 {
								ctx.Match().Destroy(target, card)
							} else {
								ctx.InterruptFlow()
							}
						}
					})
				}
			},
		},
	}

	return cb.Build()
}

// Degenerate ...
func Degenerate() *match.Card {

	cb := match.CardBuilder{
		Name:   "Degenerate",
		Rank:   0,
		Civ:    civ.PRITHVI,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {
						if event, ok := ctx.Event().(*match.PlayCardEvent); ok && event.ID == card.ID() {

							if len(event.Targets) <= 0 {
								ctx.InterruptFlow()
								return
							}

							target, err := card.Player().GetCard(event.Targets[0])
							if err != nil {
								ctx.InterruptFlow()
								logrus.Debug(err)
								return
							}

							match.Devolve(target, card)
						}
					})
				}
			},
		},
	}

	return cb.Build()
}

// LeechLife ...
func LeechLife() *match.Card {

	cb := match.CardBuilder{
		Name:   "Leech Life",
		Rank:   0,
		Civ:    civ.PRITHVI,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {
						if event, ok := ctx.Event().(*match.PlayCardEvent); ok && event.ID == card.ID() {

							if len(event.Targets) <= 0 {
								ctx.InterruptFlow()
								return
							}

							target, err := card.Player().GetCard(event.Targets[0])
							if err != nil {
								ctx.InterruptFlow()
								logrus.Debug(err)
								return
							}

							target.AddCondition(func(card *match.Card, ctx *match.Context) {
								fx.AttackModifier(card, ctx, 200)
							})
							target.AddCondition(fx.Leech)
						}
					})
				}
			},
		},
	}

	return cb.Build()
}

// RapidEvolution ...
func RapidEvolution() *match.Card {

	cb := match.CardBuilder{
		Name:   "Rapid Evolution",
		Rank:   0,
		Civ:    civ.PRITHVI,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {
						if event, ok := ctx.Event().(*match.PlayCardEvent); ok && event.ID == card.ID() {

							if len(event.Targets) <= 0 {
								ctx.InterruptFlow()
								return
							}

							target, err := card.Player().GetCard(event.Targets[0])
							if err != nil {
								ctx.InterruptFlow()
								logrus.Debug(err)
								return
							}

							target.Tap(true)
							target.RemoveCondition(fx.CantEvolve)
						}
					})
				}
			},
		},
	}

	return cb.Build()
}

// AirMail ...
func AirMail() *match.Card {

	cb := match.CardBuilder{
		Name:   "Air Mail",
		Rank:   0,
		Civ:    civ.VAYU,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {

						cards, err := card.Player().Container(match.DECK)
						if err != nil {
							ctx.InterruptFlow()
							logrus.Debug(err)
							return
						}

						cards = card.Player().Search(
							cards,
							"Select a card",
							1,
							1,
							false)

						for _, c := range cards {
							if err := c.MoveCard(match.HAND); err != nil {
								logrus.Debug(err)
								return
							}
							ctx.Match().Chat("Server", fmt.Sprintf("%s was moved from %s's deck to their hand", c.Name(), c.Player().Name()))
						}

						card.Player().ShuffleDeck()
					})
				}
			},
		},
	}

	return cb.Build()
}

// Whirlwind ...
func Whirlwind() *match.Card {

	cb := match.CardBuilder{
		Name:   "Whirlwind",
		Rank:   0,
		Civ:    civ.VAYU,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {

						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))

						cards, err := ctx.Match().Opponent(card.Player()).Container(match.TRAPZONE)
						if err != nil {
							ctx.InterruptFlow()
							logrus.Debug(err)
							return
						}

						cards = card.Player().Search(
							cards,
							"Select a card",
							1,
							1,
							false)

						for _, c := range cards {
							if err := c.MoveCard(match.GRAVEYARD); err != nil {
								logrus.Debug(err)
							}
							ctx.Match().Chat("Server", fmt.Sprintf("%s was moved to %s %s by %s", c.Name(), c.Player().Name(), match.GRAVEYARD, card.Name()))
						}
					})
				}
			},
		},
	}

	return cb.Build()
}

// Tailwind ...
func Tailwind() *match.Card {

	cb := match.CardBuilder{
		Name:   "Tailwind",
		Rank:   0,
		Civ:    civ.VAYU,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {
						if event, ok := ctx.Event().(*match.PlayCardEvent); ok && event.ID == card.ID() {

							if len(event.Targets) <= 0 {
								ctx.InterruptFlow()
								return
							}

							target, err := card.Player().GetCard(event.Targets[0])
							if err != nil {
								ctx.InterruptFlow()
								logrus.Debug(err)
								return
							}

							target.AddCondition(fx.CantBeBlocked)
						}
					})
				}
			},
		},
	}

	return cb.Build()
}

// Tornado ...
func Tornado() *match.Card {

	cb := match.CardBuilder{
		Name:   "Tornado",
		Rank:   1,
		Civ:    civ.VAYU,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {
						if event, ok := ctx.Event().(*match.PlayCardEvent); ok && event.ID == card.ID() {

							if len(event.Targets) <= 0 {
								ctx.InterruptFlow()
								return
							}

							target, err := card.Player().GetCard(event.Targets[0])
							if err != nil {
								ctx.InterruptFlow()
								logrus.Debug(err)
								return
							}

							if err := target.MoveCard(match.DECK); err != nil {
								logrus.Debug(err)
							}
						}
					})
				}
			},
		},
	}

	return cb.Build()
}

// FrostBreath ...
func FrostBreath() *match.Card {

	cb := match.CardBuilder{
		Name:   "FrostBreath",
		Rank:   0,
		Civ:    civ.APAS,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {
						if event, ok := ctx.Event().(*match.PlayCardEvent); ok && event.ID == card.ID() {

							if len(event.Targets) <= 0 {
								ctx.InterruptFlow()
								return
							}

							target, err := card.Player().GetCard(event.Targets[0])
							if err != nil {
								ctx.InterruptFlow()
								logrus.Debug(err)
								return
							}

							target.Tap(true)
						}
					})
				}
			},
		},
	}

	return cb.Build()
}

// TidalWave ...
func TidalWave() *match.Card {

	cb := match.CardBuilder{
		Name:   "Tidal Wave",
		Rank:   1,
		Civ:    civ.APAS,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {
						card.Player().DrawCards(2)
					})
				}
			},
		},
	}

	return cb.Build()
}

// Amrita ...
func Amrita() *match.Card {

	cb := match.CardBuilder{
		Name:   "Amrita",
		Rank:   1,
		Civ:    civ.APAS,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {
						card.Player().Heal(card, ctx, 800)
					})
				}
			},
		},
	}

	return cb.Build()
}

// Blizzard ...
func Blizzard() *match.Card {

	cb := match.CardBuilder{
		Name:   "Blizzard",
		Rank:   2,
		Civ:    civ.APAS,
		Family: family.Spell,
		Handlers: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {
						for _, c := range ctx.Match().Opponent(card.Player()).GetCreatures() {
							c.Tap(true)
						}
					})
				}
			},
		},
	}

	return cb.Build()
}
