package fx

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Spell has default functionality for spells
func Spell(card *match.Card, ctx *match.Context) {
	switch event := ctx.Event().(type) {
	// On card played
	case *match.PlayCardEvent:
		if event.ID == card.ID() {
			card.SpellCast(ctx, func() []*match.Card { return []*match.Card{} })
		}
	// When the spell is played reactively
	case *match.TrapEvent:
		if event.ID == card.ID() {
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
			})
		}
	}
}
