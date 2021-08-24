package fx

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
	"github.com/sirupsen/logrus"
)

// Leech ...
func Leech(card *match.Card, ctx *match.Context) {
	switch event := ctx.Event().(type) {
	case *match.DamageEvent:
		if event.Source == card {
			ctx.ScheduleAfter(func() {
				card.Player().Heal(card, ctx, event.Health)
			})
		}
	case *match.CreatureDestroyed:
		if event.Source == card {
			ctx.ScheduleAfter(func() {
				c, err := match.GetCard(
					event.ID,
					ctx.Match().Opponent(card.Player()).CollectCards(match.GRAVEYARD))
				if err != nil {
					logrus.Debug(err)
					return
				}

				card.Player().Heal(card, ctx, c.GetDefence(ctx))
			})
		}
	}
}
