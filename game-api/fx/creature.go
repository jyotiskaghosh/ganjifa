package fx

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, ctx *match.Context) {

	// Untap the card
	if _, ok := ctx.Event.(*match.UntapStep); ok && card.Player.IsPlayerTurn() {

		card.Tapped = false
		card.ClearConditions()
	}

	if match.AmIPlayed(card, ctx) {

		if card.GetRank(ctx) > 0 {

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {
				match.Evolution(
					card,
					ctx,
					fmt.Sprintf("choose a %s creature to evolve %s from your battlezone", card.Family, card.Name),
					func(x *match.Card) bool {
						dif := card.GetRank(ctx) - x.GetRank(ctx)
						return x.HasFamily(card.Family, ctx) && dif >= 0 && dif <= 1
					})
			})
		} else {

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {

				if err := card.MoveCard(match.BATTLEZONE); err != nil {
					logrus.Debug(err)
					return
				}
				ctx.Match.Chat("Server", fmt.Sprintf("%s summoned %s to the battle zone", card.Player.Name(), card.Name))

				card.AddCondition(CantEvolve)
			})
		}
	}

	// On evolve
	if event, ok := ctx.Event.(*match.Evolve); ok && event.ID == card.ID {

		// Do this last in case any other cards want to interrupt the flow
		ctx.Override(func() {

			if err := card.MoveCard(match.BATTLEZONE); err != nil {
				logrus.Debug(err)
				return
			}

			card.Tapped = event.Creature.Tapped

			card.Attach(event.Creature)

			ctx.Match.Chat("Server", fmt.Sprintf("%s summoned %s to the battle zone", card.Player.Name(), card.Name))

			card.AddCondition(event.Creature.Conditions()...)
			card.AddCondition(CantEvolve)
		})
	}

	// Attack the player
	if event, ok := ctx.Event.(*match.AttackPlayer); ok && event.ID == card.ID {

		if card.Tapped || card.Zone != match.BATTLEZONE {
			ctx.InterruptFlow()
			return
		}

		// Do this last in case any other cards want to interrupt the flow
		ctx.Override(func() {

			opponent := ctx.Match.Opponent(card.Player)

			react(card, ctx, opponent, func() {
				card.Tapped = true
				opponent.Damage(card, ctx, card.GetAttack(ctx))
			})
		})
	}

	// Attack a creature
	if event, ok := ctx.Event.(*match.AttackCreature); ok && event.ID == card.ID {

		if card.Tapped || card.Zone != match.BATTLEZONE {
			ctx.InterruptFlow()
			return
		}

		opponent := ctx.Match.Opponent(card.Player)

		defender, err := opponent.GetCard(event.TargetID)
		if err != nil {
			ctx.Match.WarnPlayer(card.Player, "creature to attack was not found")
			ctx.InterruptFlow()
			return
		}

		// Do this last in case any other cards want to interrupt the flow
		ctx.Override(func() {

			react(card, ctx, opponent, func() { ctx.Match.Battle(card, defender, false) })
		})
	}

	// When blocking
	if event, ok := ctx.Event.(*match.BlockEvent); ok && event.Blocker == card {

		if card.Tapped || card.Zone != match.BATTLEZONE {
			ctx.InterruptFlow()
			return
		}

		ctx.Override(func() {
			event.Blocker.Tapped = true
			ctx.Match.Battle(event.Attacker, card, true)
		})
	}

	// When destroyed
	if event, ok := ctx.Event.(*match.CreatureDestroyed); ok && event.Card == card {

		ctx.Override(func() {
			if err := card.MoveCard(match.GRAVEYARD); err != nil {
				logrus.Debug(err)
			}
		})
	}
}

func react(card *match.Card, ctx *match.Context, opponent *match.Player, f func()) {

	spellzone, err := ctx.Match.Opponent(card.Player).ContainerRef(match.SPELLZONE)
	if err != nil {
		ctx.InterruptFlow()
		logrus.Debug(err)
		return
	}

	blockers := make([]*match.Card, 0)

	for _, c := range opponent.GetCreatures() {
		if !c.Tapped {
			blockers = append(blockers, c)
		}
	}

	ctx.Match.NewAction(opponent, append(blockers, *spellzone...), 1, 1, fmt.Sprintf("%s is attacking player, you may play a set down card or block with a creature", card.Name), true)
	defer ctx.Match.CloseAction(opponent)

	for {
		select {
		case action := <-opponent.Action:
			{
				if len(action) < 1 || len(action) > 1 || !match.AssertCardsIn(append(blockers, *spellzone...), action) {
					ctx.Match.WarnPlayer(opponent, "The cards you selected does not meet the requirements")
					continue
				}

				for _, id := range action {

					c, err := opponent.GetCard(id)
					if err != nil {
						logrus.Debug(err)
					}

					if c.Zone == match.SPELLZONE {

						ctx.Match.React(c.ID, ctx.Event)

						ctx.Match.ResolveEvent(ctx)

						if ctx.Cancelled() {
							return
						}

						ctx.Match.NewAction(opponent, append(blockers, *spellzone...), 1, 1, fmt.Sprintf("%s is attacking player, you may play a set down card or block with a creature", card.Name), true)

					} else {

						ctx.Match.Block(card, c)

						ctx.Match.ResolveEvent(ctx)

						if ctx.Cancelled() {
							return
						}
					}
				}
			}

		case cancel := <-opponent.Cancel:
			if cancel {
				f()
				return
			}
		}
	}
}
