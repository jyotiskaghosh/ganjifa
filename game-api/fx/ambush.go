package fx

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Ambush ...
func Ambush(card *match.Card, ctx *match.Context) {
	if event, ok := ctx.Event().(*match.TrapEvent); ok && event.ID == card.ID() {
		// Do this last in case any other cards want to interrupt the flow
		ctx.ScheduleAfter(func() {
			playCtx := match.NewContext(ctx.Match(), &match.PlayCardEvent{
				ID: card.ID(),
			})
			ctx.Match().HandleFx(playCtx)

			if playCtx.Cancelled() {
				ctx.InterruptFlow()
				return
			}

			if card.Tapped || card.Zone() != match.BATTLEZONE {
				return
			}

			ctx.Match().Battle(card, event.Attacker, false)
			card.Tapped = true
		})
	}
}
