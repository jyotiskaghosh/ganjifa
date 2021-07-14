package fx

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
	"github.com/sirupsen/logrus"
)

// Spell has default functionality for spells
func Spell(card *match.Card, ctx *match.Context) {
	switch event := ctx.Event().(type) {
	// On card played
	case *match.PlayCardEvent:
		if event.ID == card.ID() {
			for _, creature := range card.Player().GetCreatures() {
				if creature.HasCivilisation(card.Civ(), ctx) && card.GetRank(ctx) <= creature.GetRank(ctx) && !card.Tapped {
					ctx.ScheduleAfter(func() {
						if err := card.MoveCard(match.GRAVEYARD); err != nil {
							logrus.Debug(err)
						}
					})
					return
				}
			}
			ctx.InterruptFlow()
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
