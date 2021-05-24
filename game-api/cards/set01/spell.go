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

	c := &match.Card{
		Name:   "Energy Surge",
		Rank:   0,
		Civ:    civ.AGNI,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				creatures := card.Player.Search(
					card.Player.GetCreatures(),
					"Select 1 of your creatures",
					1,
					1,
					false)

				for _, c := range creatures {
					c.AddCondition(func(card *match.Card, ctx *match.Context) {
						fx.AttackModifier(card, ctx, 400)
						ctx.Match.Chat("Server", fmt.Sprintf("%s used spell %s on %s", card.Player.Name(), card.Name, c.Name))
					})
				}
			})
		}
	})

	return c
}

// Fireball ...
func Fireball() *match.Card {

	c := &match.Card{
		Name:   "Fireball",
		Rank:   0,
		Civ:    civ.AGNI,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				creatures := card.Player.Filter(
					ctx.Match.Opponent(card.Player).GetCreatures(),
					"Select 1 of your opponent's creature with defense 200 or lesser",
					1,
					1,
					false,
					func(x *match.Card) bool { return x.GetDefence(ctx) <= 200 })

				for _, c := range creatures {
					ctx.Match.Destroy(c, card, fmt.Sprintf("%s was destroyed by %s", c.Name, card.Name))
				}
			})
		}
	})

	return c
}

// RainOfArrows ...
func RainOfArrows() *match.Card {

	c := &match.Card{
		Name:   "Rain Of Arrows",
		Rank:   0,
		Civ:    civ.AGNI,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				for _, c := range append(card.Player.GetCreatures(), ctx.Match.Opponent(card.Player).GetCreatures()...) {

					if c.GetDefence(ctx) <= 100 {
						ctx.Match.Destroy(c, card, fmt.Sprintf("%s was destroyed by %s", c.Name, card.Name))
					}
				}
			})
		}
	})

	return c
}

// MagmaGeyser ...
func MagmaGeyser() *match.Card {

	c := &match.Card{
		Name:   "MagmaGeyser",
		Rank:   1,
		Civ:    civ.AGNI,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				creatures := card.Player.Filter(
					ctx.Match.Opponent(card.Player).GetCreatures(),
					"Select 1 of your opponent's creature with defense 400 or lesser",
					1,
					1,
					false,
					func(x *match.Card) bool { return x.GetDefence(ctx) <= 400 })

				for _, c := range creatures {
					ctx.Match.Destroy(c, card, fmt.Sprintf("%s was destroyed by %s", c.Name, card.Name))
				}
			})
		}
	})

	return c
}

// Degenerate ...
func Degenerate() *match.Card {

	c := &match.Card{
		Name:   "Degenerate",
		Rank:   0,
		Civ:    civ.PRITHVI,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				creatures := card.Player.Search(
					ctx.Match.Opponent(card.Player).GetCreatures(),
					"Select 1 of your opponent's creature",
					1,
					1,
					false)

				for _, c := range creatures {

					creatures = ctx.Match.Opponent(card.Player).Filter(
						c.Attachments(),
						"Select a creature",
						1,
						1,
						false,
						func(x *match.Card) bool { return x.Family != family.Equipment })

					if len(creatures) < 1 {
						ctx.Match.Destroy(c, card, fmt.Sprintf("%s was destroyed by %s", c.Name, card.Name))
						return
					}

					if err := creatures[0].MoveCard(match.BATTLEZONE); err != nil {
						logrus.Debug(err)
					}

					c.AttachTo(creatures[0])

					// This is done to maintain a single identity for a creature
					card.ID, creatures[0].ID = creatures[0].ID, card.ID

					ctx.Match.Destroy(c, card, fmt.Sprintf("%s was destroyed by %s", c.Name, card.Name))
				}
			})
		}
	})

	return c
}

// LeechLife ...
func LeechLife() *match.Card {

	c := &match.Card{
		Name:   "Leech Life",
		Rank:   0,
		Civ:    civ.PRITHVI,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				creatures := card.Player.Search(
					card.Player.GetCreatures(),
					"Select 1 of your creature",
					1,
					1,
					false)

				for _, c := range creatures {
					c.AddCondition(func(card *match.Card, ctx *match.Context) {
						fx.AttackModifier(card, ctx, 200)
					})
					c.AddCondition(fx.Leech)
					ctx.Match.Chat("Server", fmt.Sprintf("%s used spell %s on %s", card.Player.Name(), card.Name, c.Name))
				}
			})
		}
	})

	return c
}

// RapidEvolution ...
func RapidEvolution() *match.Card {

	c := &match.Card{
		Name:   "Rapid Evolution",
		Rank:   0,
		Civ:    civ.PRITHVI,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				creatures := card.Player.Filter(
					card.Player.GetCreatures(),
					"Select 1 of your creature",
					1,
					1,
					false,
					func(x *match.Card) bool { return !x.Tapped() && x.HasCondition(fx.CantEvolve) })

				for _, c := range creatures {
					c.RemoveCondition(fx.CantEvolve)
					c.Tap(true)
					ctx.Match.Chat("Server", fmt.Sprintf("%s used spell %s on %s", card.Player.Name(), card.Name, c.Name))
				}
			})
		}
	})

	return c
}

// AirMail ...
func AirMail() *match.Card {

	c := &match.Card{
		Name:   "Air Mail",
		Rank:   0,
		Civ:    civ.VAYU,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {
		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				cards, err := card.Player.Container(match.DECK)
				if err != nil {
					ctx.InterruptFlow()
					logrus.Debug(err)
					return
				}

				cards = card.Player.Search(
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
					ctx.Match.Chat("Server", fmt.Sprintf("%s was moved from %s's deck to their hand", c.Name, c.Player.Name()))
				}

				card.Player.ShuffleDeck()
			})
		}
	})

	return c
}

// Whirlwind ...
func Whirlwind() *match.Card {

	c := &match.Card{
		Name:   "Whirlwind",
		Rank:   0,
		Civ:    civ.VAYU,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				cards, err := ctx.Match.Opponent(card.Player).Container(match.HIDDENZONE)
				if err != nil {
					ctx.InterruptFlow()
					logrus.Debug(err)
					return
				}

				cards = card.Player.Search(
					cards,
					"Select a card",
					1,
					1,
					false)

				for _, c := range cards {
					if err := c.MoveCard(match.GRAVEYARD); err != nil {
						logrus.Debug(err)
					}
					ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s %s by %s", c.Name, c.Player.Name(), match.GRAVEYARD, card.Name))
				}
			})
		}
	})

	return c
}

// Tailwind ...
func Tailwind() *match.Card {

	c := &match.Card{
		Name:   "Tailwind",
		Rank:   0,
		Civ:    civ.VAYU,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				creatures := card.Player.Search(
					card.Player.GetCreatures(),
					"Select 1 of your creature",
					1,
					1,
					false)

				for _, c := range creatures {
					c.AddCondition(fx.CantBeBlocked)
					ctx.Match.Chat("Server", fmt.Sprintf("%s used spell %s on %s", card.Player.Name(), card.Name, c.Name))
				}
			})
		}
	})

	return c
}

// Tornado ...
func Tornado() *match.Card {

	c := &match.Card{
		Name:   "Tornado",
		Rank:   1,
		Civ:    civ.VAYU,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				creatures := card.Player.Search(
					ctx.Match.Opponent(card.Player).GetCreatures(),
					"Select 1 of your opponent's creature",
					1,
					1,
					false)

				for _, c := range creatures {
					if err := c.MoveCard(match.DECK); err != nil {
						logrus.Debug(err)
					}
					ctx.Match.Chat("Server", fmt.Sprintf("%s was moved to %s %s by %s", c.Name, c.Player.Name(), match.DECK, card.Name))
				}
			})
		}
	})

	return c
}

// FrostBreath ...
func FrostBreath() *match.Card {

	c := &match.Card{
		Name:   "FrostBreath",
		Rank:   0,
		Civ:    civ.APAS,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				creatures := card.Player.Filter(
					ctx.Match.Opponent(card.Player).GetCreatures(),
					"Select 1 of your opponent's creature",
					1,
					1,
					false,
					func(x *match.Card) bool { return !x.Tapped() },
				)

				for _, c := range creatures {
					c.Tap(true)
					ctx.Match.Chat("Server", fmt.Sprintf("%s used spell %s on %s", card.Player.Name(), card.Name, c.Name))
				}
			})
		}
	})

	return c
}

// TidalWave ...
func TidalWave() *match.Card {

	c := &match.Card{
		Name:   "Tidal Wave",
		Rank:   1,
		Civ:    civ.APAS,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {
		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {
				card.Player.DrawCards(2)
				ctx.Match.Chat("Server", fmt.Sprintf("%s used spell %s", card.Player.Name(), card.Name))
			})
		}
	})

	return c
}

// Amrita ...
func Amrita() *match.Card {

	c := &match.Card{
		Name:   "Amrita",
		Rank:   1,
		Civ:    civ.APAS,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {
		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {
				card.Player.Heal(card, ctx, 800)
			})
		}
	})

	return c
}

// Blizzard ...
func Blizzard() *match.Card {

	c := &match.Card{
		Name:   "Blizzard",
		Rank:   2,
		Civ:    civ.APAS,
		Family: family.Spell,
	}

	c.Use(fx.Spell, func(card *match.Card, ctx *match.Context) {

		if match.AmIPlayed(card, ctx) {

			ctx.Override(func() {

				for _, c := range ctx.Match.Opponent(card.Player).GetCreatures() {
					c.Tap(true)
					ctx.Match.Chat("Server", fmt.Sprintf("%s used spell %s", card.Player.Name(), card.Name))
				}
			})
		}
	})

	return c
}
