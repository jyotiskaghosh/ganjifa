package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Equipment has default functionality for equipments
func Equipment(card *match.Card, ctx *match.Context) {

	if match.AmIPlayed(card, ctx) {

		ctx.Override(func() {
			match.Equipment(
				card,
				ctx,
				fmt.Sprintf("choose a creature to equip %s from your battlezone", card.Name),
				func(x *match.Card) bool {
					return x.HasCivilisation(card.Civ, ctx) && card.GetRank(ctx) <= x.GetRank(ctx)
				})
		})
	}

	if event, ok := ctx.Event.(*match.GetAttackEvent); ok && card.AttachedTo() == event.Card {

		event.Attack += card.Attack
	}

	if event, ok := ctx.Event.(*match.GetDefenceEvent); ok && card.AttachedTo() == event.Card {

		event.Defence += card.Defence
	}

	// When the equipment is played reactively
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

	// On equip
	if event, ok := ctx.Event.(*match.Equip); ok && event.ID == card.ID {

		// Do this last in case any other cards want to interrupt the flow
		ctx.Override(func() {

			ctx.Match.Chat("Server", fmt.Sprintf("%s equipped %s on %s", card.Player.Name(), card.Name, event.Creature.Name))

			// destroy existing equipment if any
			event.Creature.RemoveEquipments()

			card.AttachTo(event.Creature)
		})
	}
}
