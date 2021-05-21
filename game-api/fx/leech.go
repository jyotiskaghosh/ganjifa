package fx

import "github.com/jyotiskaghosh/ganjifa/game-api/match"

// Leech ...
func Leech(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.DamageEvent); ok && event.Source == card {

		ctx.ScheduleAfter(func() {
			card.Player.Heal(card, ctx, event.Health)
		})
	}

	if event, ok := ctx.Event.(*match.CreatureDestroyed); ok && event.Source == card {

		ctx.ScheduleAfter(func() {
			card.Player.Heal(card, ctx, event.Card.GetDefence(ctx))
		})
	}
}
