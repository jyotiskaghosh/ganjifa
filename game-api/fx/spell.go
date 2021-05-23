package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
	"github.com/sirupsen/logrus"
)

// Spell has default functionality for spells
func Spell(card *match.Card, ctx *match.Context) {

	switch event := ctx.Event.(type) {

	// On card played
	case *match.PlayCardEvent:
		if event.ID == card.ID {

			for _, creature := range card.Player.GetCreatures() {

				if creature.HasCivilisation(card.Civ, ctx) && card.GetRank(ctx) <= creature.GetRank(ctx) {
					ctx.ScheduleAfter(func() {
						if err := card.MoveCard(match.GRAVEYARD); err != nil {
							logrus.Debug(err)
						}
						ctx.Match.Chat("Server", fmt.Sprintf("%s moved to %s's %s", card.Name, card.Player.Name(), match.GRAVEYARD))
					})
					return
				}
			}

			ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s cannot be cast; conditions unsatisfied", card.Name))
			ctx.InterruptFlow()
		}

	// When the spell is played reactively
	case *match.React:
		if event.ID == card.ID {

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
}
