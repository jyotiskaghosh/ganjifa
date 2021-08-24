package fx

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Equipment has default functionality for equipments
func Equipment(card *match.Card, ctx *match.Context) {
	switch event := ctx.Event().(type) {
	// On card played
	case *match.PlayCardEvent:
		if event.ID == card.ID() {
			// Do this last in case any other cards want to interrupt the flow
			ctx.ScheduleAfter(func() {
				cards := card.Player().Search(
					match.Filter(
						card.Player().CollectCards(match.BATTLEZONE),
						func(c *match.Card) bool {
							return c.HasCivilisation(card.Civ(), ctx) && card.GetRank(ctx) <= c.GetRank(ctx)
						},
					),
					fmt.Sprintf("choose a creature to equip %s", card.Name()),
					1,
					1,
					true)

				for _, c := range cards {
					ctx.Match().Equip(card.ID(), c)
				}
			})
		}
	// On equip
	case *match.Equip:
		if event.ID == card.ID() {
			// Do this last in case any other cards want to interrupt the flow
			ctx.ScheduleAfter(func() {
				target, err := match.GetCard(event.Target.ID(), card.Player().CollectCards(match.BATTLEZONE))
				if err != nil {
					logrus.Debug(err)
					ctx.InterruptFlow()
					return
				}

				target.RemoveEquipments()
				card.AttachTo(target)
				ctx.Match().Chat("server", fmt.Sprintf("%s equipped %s to %s", card.Player().Name(), card.Name(), target.Name()))
			})
		}
	// When the equipment is played reactively
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
	// When calculating attack
	case *match.GetAttackEvent:
		if card.AttachedTo() != nil && event.ID == card.AttachedTo().ID() {
			event.Attack += card.Attack()
		}
	// When calculating defence
	case *match.GetDefenceEvent:
		if card.AttachedTo() != nil && event.ID == card.AttachedTo().ID() {
			event.Defence += card.Defence()
		}
	}
}
