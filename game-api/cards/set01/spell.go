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
		Effects: []match.HandlerFunc{
			func(card *match.Card, ctx *match.Context) {
				switch event := ctx.Event().(type) {
				case *match.PlayCardEvent:
					if event.ID == card.ID() {
						card.SpellCast(
							ctx,
							func() []*match.Card {
								return card.Player().Search(
									card.Player().CollectCards(match.BATTLEZONE),
									"Select 1 of your creatures",
									1,
									1,
									false)
							},
						)
					}
				case *match.SpellCast:
					if event.ID == card.ID() {
						ctx.ScheduleAfter(func() {
							for _, c := range event.Targets {
								c.AddCondition(func(card *match.Card, ctx *match.Context) {
									fx.AttackModifier(card, ctx, 4)
								})
							}
						})
					}
				default:
					fx.Spell(card, ctx)
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
		Effects: []match.HandlerFunc{
			func(card *match.Card, ctx *match.Context) {
				switch event := ctx.Event().(type) {
				case *match.PlayCardEvent:
					if event.ID == card.ID() {
						card.SpellCast(
							ctx,
							func() []*match.Card {
								return card.Player().Search(
									ctx.Match().Opponent(card.Player()).CollectCards(match.BATTLEZONE),
									"Select 1 of your opponents creatures",
									1,
									1,
									false)
							},
						)
					}
				case *match.SpellCast:
					if event.ID == card.ID() {
						ctx.ScheduleAfter(func() {
							for _, c := range event.Targets {
								if c.GetDefence(ctx) <= 2 {
									ctx.Match().Destroy(c, card)
								} else {
									ctx.InterruptFlow()
								}
							}
						})
					}
				default:
					fx.Spell(card, ctx)
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
		Effects: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if event, ok := ctx.Event().(*match.SpellCast); ok && event.ID == card.ID() {
					ctx.ScheduleAfter(func() {
						for _, c := range append(
							card.Player().CollectCards(match.BATTLEZONE),
							ctx.Match().Opponent(card.Player()).CollectCards(match.BATTLEZONE)...) {
							if c.GetDefence(ctx) <= 1 {
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
		Effects: []match.HandlerFunc{
			func(card *match.Card, ctx *match.Context) {
				switch event := ctx.Event().(type) {
				case *match.PlayCardEvent:
					if event.ID == card.ID() {
						card.SpellCast(
							ctx,
							func() []*match.Card {
								return card.Player().Search(
									ctx.Match().Opponent(card.Player()).CollectCards(match.BATTLEZONE),
									"Select 1 of your opponents creatures",
									1,
									1,
									false)
							},
						)
					}
				case *match.SpellCast:
					if event.ID == card.ID() {
						ctx.ScheduleAfter(func() {
							for _, c := range event.Targets {
								if c.GetDefence(ctx) <= 4 {
									ctx.Match().Destroy(c, card)
								} else {
									ctx.InterruptFlow()
								}
							}
						})
					}
				default:
					fx.Spell(card, ctx)
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
		Effects: []match.HandlerFunc{
			func(card *match.Card, ctx *match.Context) {
				switch event := ctx.Event().(type) {
				case *match.PlayCardEvent:
					if event.ID == card.ID() {
						card.SpellCast(
							ctx,
							func() []*match.Card {
								return card.Player().Search(
									append(
										card.Player().CollectCards(match.BATTLEZONE),
										ctx.Match().Opponent(card.Player()).CollectCards(match.BATTLEZONE)...,
									),
									"Select a creature",
									1,
									1,
									false)
							},
						)
					}
				case *match.SpellCast:
					if event.ID == card.ID() {
						ctx.ScheduleAfter(func() {
							for _, c := range event.Targets {
								match.Devolve(c, card)
							}
						})
					}
				default:
					fx.Spell(card, ctx)
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
		Effects: []match.HandlerFunc{
			func(card *match.Card, ctx *match.Context) {
				switch event := ctx.Event().(type) {
				case *match.PlayCardEvent:
					if event.ID == card.ID() {
						card.SpellCast(
							ctx,
							func() []*match.Card {
								return card.Player().Search(
									card.Player().CollectCards(match.BATTLEZONE),
									"Select 1 of your creatures",
									1,
									1,
									false)
							},
						)
					}
				case *match.SpellCast:
					if event.ID == card.ID() {
						ctx.ScheduleAfter(func() {
							for _, c := range event.Targets {
								c.AddCondition(func(card *match.Card, ctx *match.Context) {
									fx.AttackModifier(card, ctx, 2)
								})
								c.AddCondition(fx.Leech)
							}
						})
					}
				default:
					fx.Spell(card, ctx)
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
		Effects: []match.HandlerFunc{
			func(card *match.Card, ctx *match.Context) {
				switch event := ctx.Event().(type) {
				case *match.PlayCardEvent:
					if event.ID == card.ID() {
						card.SpellCast(
							ctx,
							func() []*match.Card {
								return card.Player().Search(
									append(
										card.Player().CollectCards(match.BATTLEZONE),
										ctx.Match().Opponent(card.Player()).CollectCards(match.BATTLEZONE)...,
									),
									"Select a creature",
									1,
									1,
									false)
							},
						)
					}
				case *match.SpellCast:
					if event.ID == card.ID() {
						ctx.ScheduleAfter(func() {
							for _, c := range event.Targets {
								c.Tapped = true
								c.RemoveCondition(fx.CantEvolve)
							}
						})
					}
				default:
					fx.Spell(card, ctx)
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
		Effects: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if event, ok := ctx.Event().(*match.SpellCast); ok && event.ID == card.ID() {
					ctx.ScheduleAfter(func() {
						cards := card.Player().Search(
							card.Player().CollectCards(match.DECK),
							"Select a card",
							1,
							1,
							false,
						)

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
		Effects: []match.HandlerFunc{
			func(card *match.Card, ctx *match.Context) {
				switch event := ctx.Event().(type) {
				case *match.PlayCardEvent:
					if event.ID == card.ID() {
						card.SpellCast(
							ctx,
							func() []*match.Card {
								return card.Player().Search(
									ctx.Match().Opponent(card.Player()).CollectCards(match.TRAPZONE),
									"Select 1 of your opponents traps",
									1,
									1,
									false)
							},
						)
					}
				case *match.SpellCast:
					if event.ID == card.ID() {
						ctx.ScheduleAfter(func() {
							for _, c := range event.Targets {
								if err := c.MoveCard(match.GRAVEYARD); err != nil {
									logrus.Debug(err)
								}
								ctx.Match().Chat("Server", fmt.Sprintf("%s was moved to %s %s by %s", c.Name(), c.Player().Name(), match.GRAVEYARD, card.Name()))
							}
						})
					}
				default:
					fx.Spell(card, ctx)
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
		Effects: []match.HandlerFunc{
			func(card *match.Card, ctx *match.Context) {
				switch event := ctx.Event().(type) {
				case *match.PlayCardEvent:
					if event.ID == card.ID() {
						card.SpellCast(
							ctx,
							func() []*match.Card {
								return card.Player().Search(
									card.Player().CollectCards(match.BATTLEZONE),
									"Select a creature",
									1,
									1,
									false)
							},
						)
					}
				case *match.SpellCast:
					if event.ID == card.ID() {
						ctx.ScheduleAfter(func() {
							for _, c := range event.Targets {
								c.AddCondition(fx.CantBeBlocked)
							}
						})
					}
				default:
					fx.Spell(card, ctx)
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
		Effects: []match.HandlerFunc{
			func(card *match.Card, ctx *match.Context) {
				switch event := ctx.Event().(type) {
				case *match.PlayCardEvent:
					if event.ID == card.ID() {
						card.SpellCast(
							ctx,
							func() []*match.Card {
								return card.Player().Search(
									ctx.Match().Opponent(card.Player()).CollectCards(match.BATTLEZONE),
									"Select 1 of your opponents creatures",
									1,
									1,
									false)
							},
						)
					}
				case *match.SpellCast:
					if event.ID == card.ID() {
						ctx.ScheduleAfter(func() {
							for _, c := range event.Targets {
								if err := c.MoveCard(match.DECK); err != nil {
									logrus.Debug(err)
								}
							}
						})
					}
				default:
					fx.Spell(card, ctx)
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
		Effects: []match.HandlerFunc{
			func(card *match.Card, ctx *match.Context) {
				switch event := ctx.Event().(type) {
				case *match.PlayCardEvent:
					if event.ID == card.ID() {
						card.SpellCast(
							ctx,
							func() []*match.Card {
								return card.Player().Search(
									ctx.Match().Opponent(card.Player()).CollectCards(match.BATTLEZONE),
									"Select 1 of your opponents creatures",
									1,
									1,
									false)
							},
						)
					}
				case *match.SpellCast:
					if event.ID == card.ID() {
						ctx.ScheduleAfter(func() {
							for _, c := range event.Targets {
								c.Tapped = true
							}
						})
					}
				default:
					fx.Spell(card, ctx)
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
		Effects: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if event, ok := ctx.Event().(*match.SpellCast); ok && event.ID == card.ID() {
					ctx.ScheduleAfter(func() {
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
		Effects: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if event, ok := ctx.Event().(*match.SpellCast); ok && event.ID == card.ID() {
					ctx.ScheduleAfter(func() {
						card.Player().Heal(card, ctx, 8)
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
		Effects: []match.HandlerFunc{
			fx.Spell,
			func(card *match.Card, ctx *match.Context) {
				if event, ok := ctx.Event().(*match.SpellCast); ok && event.ID == card.ID() {
					ctx.ScheduleAfter(func() {
						for _, c := range ctx.Match().Opponent(card.Player()).CollectCards(match.BATTLEZONE) {
							c.Tapped = true
						}
					})
				}
			},
		},
	}

	return cb.Build()
}
