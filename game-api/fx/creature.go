package fx

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, ctx *match.Context) {

	switch event := ctx.Event().(type) {

	// Untap the card
	case *match.UntapStep:
		if card.Player().IsPlayerTurn() {
			card.ClearConditions()
			card.Tap(false)
		}

	// On playing the card
	case *match.PlayCardEvent:
		if event.ID == card.ID() {

			if card.GetRank(ctx) > 0 {

				// Do this last in case any other cards want to interrupt the flow
				ctx.Override(func() {
					match.Evolve(
						card,
						ctx,
						fmt.Sprintf("choose a %s creature to evolve %s from your battlezone", card.Family(), card.Name()),
						func(x *match.Card) bool {
							dif := card.GetRank(ctx) - x.GetRank(ctx)
							return x.HasFamily(card.Family(), ctx) && dif >= 0 && dif <= 1
						})
				})
			} else {

				// Do this last in case any other cards want to interrupt the flow
				ctx.Override(func() {

					if err := card.MoveCard(match.BATTLEZONE); err != nil {
						logrus.Debug(err)
						return
					}
					card.AddCondition(CantEvolve)

					ctx.Match().Chat("Server", fmt.Sprintf("%s summoned %s to the battle zone", card.Player().Name(), card.Name()))
				})
			}
		}

	// On evolve
	case *match.EvolveEvent:
		if event.ID == card.ID() {

			if event.Creature.Zone() != match.BATTLEZONE {
				ctx.InterruptFlow()
				return
			}

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {
				event.Creature.EvolveTo(card)
				card.AddCondition(CantEvolve)
			})
		}

	// Attack the player
	case *match.AttackPlayer:
		if event.ID == card.ID() {

			if card.Tapped() || card.Zone() != match.BATTLEZONE || card.GetAttack(ctx) <= 0 {
				ctx.InterruptFlow()
				return
			}

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {

				opponent := ctx.Match().Opponent(card.Player())

				trapzone, err := ctx.Match().Opponent(card.Player()).Container(match.TRAPZONE)
				if err != nil {
					logrus.Debug(err)
					return
				}

				defer func() {
					card.Tap(true)
				}()

				cards := opponent.Filter(
					append(opponent.GetCreatures(), trapzone...),
					fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name(), opponent.Name()),
					1,
					1,
					true,
					func(x *match.Card) bool {
						if x.Zone() == match.BATTLEZONE {
							return match.CanPerformEvent(match.NewContext(ctx.Match(), &match.BlockEvent{
								ID:       x.ID(),
								Attacker: card,
							}))
						}
						return true
					})

				for _, c := range cards {
					if c.Zone() == match.TRAPZONE {

						ctx.Match().HandleFx(match.NewContext(ctx.Match(), &match.TrapEvent{
							ID:       c.ID(),
							Attacker: card,
						}))

					} else {

						// Blocking attack
						blockCtx := match.NewContext(ctx.Match(), &match.BlockEvent{
							ID:       c.ID(),
							Attacker: card,
						})
						ctx.Match().HandleFx(blockCtx)

						if !blockCtx.Cancelled() {
							ctx.InterruptFlow()
							return
						}
					}
				}

				// Update card
				card, err = card.Player().GetCard(event.ID)
				if err != nil {
					logrus.Debug(err)
					return
				}

				// Can attack?
				if match.CanPerformEvent(ctx) {
					opponent.Damage(card, ctx, card.GetAttack(ctx))
				}
			})
		}

	// Attack a creature
	case *match.AttackCreature:
		if event.ID == card.ID() {

			if card.Tapped() || card.Zone() != match.BATTLEZONE || card.GetAttack(ctx) <= 0 {
				ctx.InterruptFlow()
				return
			}

			opponent := ctx.Match().Opponent(card.Player())

			target, err := opponent.GetCard(event.TargetID)
			if err != nil {
				logrus.Debug(err)
				return
			}

			if target.Zone() != match.BATTLEZONE {
				ctx.InterruptFlow()
				return
			}

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {

				trapzone, err := ctx.Match().Opponent(card.Player()).Container(match.TRAPZONE)
				if err != nil {
					logrus.Debug(err)
					return
				}

				defer func() {
					card.Tap(true)
				}()

				cards := opponent.Filter(
					append(opponent.GetCreatures(), trapzone...),
					fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name(), target.Name()),
					1,
					1,
					true,
					func(x *match.Card) bool {
						if x.Zone() == match.BATTLEZONE {
							return match.CanPerformEvent(match.NewContext(ctx.Match(), &match.BlockEvent{
								ID:       x.ID(),
								Attacker: card,
							}))
						}
						return true
					})

				for _, c := range cards {
					if c.Zone() == match.TRAPZONE {

						ctx.Match().HandleFx(match.NewContext(ctx.Match(), &match.TrapEvent{
							ID:       c.ID(),
							Attacker: card,
						}))

					} else {

						// Blocking attack
						blockCtx := match.NewContext(ctx.Match(), &match.BlockEvent{
							ID:       c.ID(),
							Attacker: card,
						})
						ctx.Match().HandleFx(blockCtx)

						if !blockCtx.Cancelled() {
							ctx.InterruptFlow()
							return
						}
					}
				}

				// Update cards
				card, err = card.Player().GetCard(event.ID)
				if err != nil {
					logrus.Debug(err)
					return
				}
				target, err = opponent.GetCard(event.TargetID)
				if err != nil {
					logrus.Debug(err)
					return
				}

				// Can attack?
				if match.CanPerformEvent(ctx) {
					ctx.Match().Battle(card, target, false)
				}
			})
		}

	// When blocking
	case *match.BlockEvent:
		if event.ID == card.ID() {

			if card.Tapped() || card.Zone() != match.BATTLEZONE || card.GetDefence(ctx) <= 0 {
				ctx.InterruptFlow()
				return
			}

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {
				card.Tap(true)
				ctx.Match().Battle(event.Attacker, card, true)
			})
		}

	// When destroyed
	case *match.CreatureDestroyed:
		if event.ID == card.ID() {

			if card.Zone() != match.BATTLEZONE {
				ctx.InterruptFlow()
				return
			}

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {
				if err := card.MoveCard(match.GRAVEYARD); err != nil {
					logrus.Debug(err)
				}
			})
		}
	}
}
