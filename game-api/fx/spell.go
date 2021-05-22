package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
	"github.com/sirupsen/logrus"
)

// Spell has default functionality for spells
func Spell(card *match.Card, ctx *match.Context) {

	if match.AmIPlayed(card, ctx) {

		for _, creature := range card.Player.GetCreatures() {

			if creature.HasCivilisation(card.Civ, ctx) && card.GetRank(ctx) <= creature.GetRank(ctx) {
				ctx.ScheduleAfter(func() {
					if err := card.MoveCard(match.GRAVEYARD); err != nil {
						logrus.Debug(err)
					}
				})
				return
			}
		}

		ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot be cast, conditions unsatisfied", card.Name))
		ctx.InterruptFlow()
	}

	// When the spell is played reactively
	if event, ok := ctx.Event.(*match.React); ok && event.ID == card.ID {

		// Do this last in case any other cards want to interrupt the flow
		ctx.Override(func() {
			playCtx := match.NewContext(ctx.Match, &match.PlayCardEvent{
				ID: card.ID,
			})

			ctx.Match.HandleFx(playCtx)

			if ctx.Cancelled() {
				ctx.InterruptFlow()
				return
			}
		})
	}
}
