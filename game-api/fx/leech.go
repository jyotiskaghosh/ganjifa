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
				card, err := ctx.Match().GetCard(event.ID)
				if err != nil {
					logrus.Debug(err)
					return
				}

				card.Player().Heal(card, ctx, card.GetDefence(ctx))
			})
		}
	}
}
