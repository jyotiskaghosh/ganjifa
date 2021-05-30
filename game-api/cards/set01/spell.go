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
				if card.AmIPlayed(ctx) {
					ctx.Override(func() {

						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))

						cards := card.Player().Search(
							card.Player().GetCreatures(),
							"Select 1 of your creatures",
							1,
							1,
							false)

						for _, c := range cards {
							c.AddCondition(func(card *match.Card, ctx *match.Context) {
								fx.AttackModifier(card, ctx, 400)
								ctx.Match().Chat("Server", fmt.Sprintf("%s used spell %s on %s", card.Player().Name(), card.Name(), c.Name()))
							})
						}
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

						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))

						cards := card.Player().Filter(
							ctx.Match().Opponent(card.Player()).GetCreatures(),
							"Select 1 of your opponent's creature with defense 200 or lesser",
							1,
							1,
							false,
							func(x *match.Card) bool { return x.GetDefence(ctx) <= 200 })

						for _, c := range cards {
							ctx.Match().Destroy(c, card, fmt.Sprintf("%s was destroyed by %s", c.Name(), card.Name()))
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

						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))

						for _, c := range append(card.Player().GetCreatures(), ctx.Match().Opponent(card.Player()).GetCreatures()...) {
							if c.GetDefence(ctx) <= 100 {
								ctx.Match().Destroy(c, card, fmt.Sprintf("%s was destroyed by %s", c.Name(), card.Name()))
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

						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))

						cards := card.Player().Filter(
							ctx.Match().Opponent(card.Player()).GetCreatures(),
							"Select 1 of your opponent's creature with defense 400 or lesser",
							1,
							1,
							false,
							func(x *match.Card) bool { return x.GetDefence(ctx) <= 400 })

						for _, c := range cards {
							ctx.Match().Destroy(c, card, fmt.Sprintf("%s was destroyed by %s", c.Name(), card.Name()))
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

						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))

						cards := card.Player().Search(
							ctx.Match().Opponent(card.Player()).GetCreatures(),
							"Select 1 of your opponent's creature",
							1,
							1,
							false)

						for _, c := range cards {
							match.Devolve(c, card)
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

						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))

						cards := card.Player().Search(
							card.Player().GetCreatures(),
							"Select 1 of your creature",
							1,
							1,
							false)

						for _, c := range cards {
							c.AddCondition(func(card *match.Card, ctx *match.Context) {
								fx.AttackModifier(card, ctx, 200)
							})
							c.AddCondition(fx.Leech)
							ctx.Match().Chat("Server", fmt.Sprintf("%s used spell %s on %s", card.Player().Name(), card.Name(), c.Name()))
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

						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))

						cards := card.Player().Filter(
							card.Player().GetCreatures(),
							"Select 1 of your creature",
							1,
							1,
							false,
							func(x *match.Card) bool { return !x.Tapped() && x.HasHandler(fx.CantEvolve, ctx) })

						for _, c := range cards {
							c.RemoveCondition(fx.CantEvolve)
							c.Tap(true)
							ctx.Match().Chat("Server", fmt.Sprintf("%s used spell %s on %s", card.Player().Name(), card.Name(), c.Name()))
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

						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))

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

						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))

						cards := card.Player().Search(
							card.Player().GetCreatures(),
							"Select 1 of your creature",
							1,
							1,
							false)

						for _, c := range cards {
							c.AddCondition(fx.CantBeBlocked)
							ctx.Match().Chat("Server", fmt.Sprintf("%s used spell %s on %s", card.Player().Name(), card.Name(), c.Name()))
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

						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))

						cards := card.Player().Search(
							ctx.Match().Opponent(card.Player()).GetCreatures(),
							"Select 1 of your opponent's creature",
							1,
							1,
							false)

						for _, c := range cards {
							if err := c.MoveCard(match.DECK); err != nil {
								logrus.Debug(err)
							}
							ctx.Match().Chat("Server", fmt.Sprintf("%s was moved to %s %s by %s", c.Name(), c.Player().Name(), match.DECK, card.Name()))
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

						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))

						cards := card.Player().Filter(
							ctx.Match().Opponent(card.Player()).GetCreatures(),
							"Select 1 of your opponent's creature",
							1,
							1,
							false,
							func(x *match.Card) bool { return !x.Tapped() },
						)

						for _, c := range cards {
							c.Tap(true)
							ctx.Match().Chat("Server", fmt.Sprintf("%s used spell %s on %s", card.Player().Name(), card.Name(), c.Name()))
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
						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))
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
						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))
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
						ctx.Match().Chat("Server", fmt.Sprintf("%s played spell %s", card.Player().Name(), card.Name()))
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
