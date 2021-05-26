package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Equipment has default functionality for equipments
func Equipment(card *match.Card, ctx *match.Context) {

	switch event := ctx.Event().(type) {

	// On card played
	case *match.PlayCardEvent:
		if event.ID == card.ID() {

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {
				match.Equipment(
					card,
					ctx,
					fmt.Sprintf("choose a creature to equip %s from your battlezone", card.Name()),
					func(x *match.Card) bool {
						return x.HasCivilisation(card.Civ(), ctx) && card.GetRank(ctx) <= x.GetRank(ctx)
					})
			})
		}

	// When calculating attack
	case *match.GetAttackEvent:
		if card.AttachedTo() == event.Card {

			event.Attack += card.Attack()
		}

	// When calculating defence
	case *match.GetDefenceEvent:
		if card.AttachedTo() == event.Card {

			event.Defence += card.Defence()
		}

	// When the equipment is played reactively
	case *match.React:
		if event.ID == card.ID() {

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {
				playCtx := match.NewContext(ctx.Match(), &match.PlayCardEvent{
					ID: card.ID(),
				})

				ctx.Match().HandleFx(playCtx)

				if ctx.Cancelled() {
					ctx.InterruptFlow()
					return
				}
			})
		}

	// On equip
	case *match.Equip:
		if event.ID == card.ID() {

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {

				// Destroy existing equipment if any
				event.Creature.RemoveEquipments()

				card.AttachTo(event.Creature)

				ctx.Match().Chat("Server", fmt.Sprintf("%s equipped %s on %s", card.Player().Name(), card.Name(), event.Creature.Name()))
			})
		}
	}
}
