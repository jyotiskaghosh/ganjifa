package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
	"github.com/sirupsen/logrus"
)

// Ambush ...
func Ambush(card *match.Card, ctx *match.Context) {

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

		ctx.ScheduleAfter(func() {

			if card.Tapped() || card.Zone() != match.BATTLEZONE {
				return
			}

			defer card.Tap(true)

			ctx.Match.NewAction(card.Player(), nil, 0, 0, fmt.Sprintf("Should %s ambush?", card.Name()), true)
			defer ctx.Match.CloseAction(card.Player())

			for {
				select {
				case <-card.Player().Action:
					{
						switch event := event.Event.(type) {

						case *match.AttackPlayer:
							{
								c, err := ctx.Match.Opponent(card.Player()).GetCard(event.ID)
								if err != nil {
									logrus.Debug(err)
									return
								}

								ctx.Match.Battle(card, c, false)
							}

						case *match.AttackCreature:
							{
								c, err := ctx.Match.Opponent(card.Player()).GetCard(event.ID)
								if err != nil {
									logrus.Debug(err)
									return
								}

								ctx.Match.Battle(card, c, false)
							}
						}

						return
					}

				case cancel := <-card.Player().Cancel:
					if cancel {
						return
					}
				}
			}
		})
	}
}
