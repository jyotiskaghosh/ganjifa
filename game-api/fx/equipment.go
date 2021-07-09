package fx

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
	"github.com/sirupsen/logrus"
)

// Equipment has default functionality for equipments
func Equipment(card *match.Card, ctx *match.Context) {

	switch event := ctx.Event().(type) {

	// On card played
	case *match.PlayCardEvent:
		if event.ID == card.ID() {

			if len(event.Targets) <= 0 {
				ctx.InterruptFlow()
				return
			}

			target, err := card.Player().GetCard(event.Targets[0])
			if err != nil {
				ctx.InterruptFlow()
				logrus.Debug(err)
				return
			}

			if !target.HasCivilisation(card.Civ(), ctx) && card.GetRank(ctx) > target.GetRank(ctx) {
				ctx.InterruptFlow()
				return
			}

			// Do this last in case any other cards want to interrupt the flow
			ctx.ScheduleAfter(func() {
				target.RemoveEquipments()
				card.AttachTo(target)
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
