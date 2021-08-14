package fx

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, ctx *match.Context) {
	canAttack := func() bool { return !card.Tapped && card.Zone() == match.BATTLEZONE && card.GetAttack(ctx) > 0 }

	switch event := ctx.Event().(type) {
	// Untap the card
	case *match.UntapStep:
		if card.Player().IsPlayerTurn() {
			card.ClearConditions()
			card.Tapped = false
		}
	// On playing the card
	case *match.PlayCardEvent:
		if event.ID == card.ID() {
			if card.GetRank(ctx) == 0 {
				// Do this last in case any other cards want to interrupt the flow
				ctx.ScheduleAfter(func() {
					if err := card.MoveCard(match.BATTLEZONE); err != nil {
						logrus.Debug(err)
						return
					}
					ctx.Match().Chat("server", fmt.Sprintf("%s summoned creature %s", card.Player().Name(), card.Name()))

					card.AddCondition(CantEvolve)
				})
			} else if card.GetRank(ctx) > 0 {
				target, err := card.Player().GetCard(event.TargetID)
				if err != nil {
					ctx.InterruptFlow()
					logrus.Debug(err)
					return
				}

				dif := card.GetRank(ctx) - target.GetRank(ctx)

				if target.HasFamily(card.Family(), ctx) && dif >= 0 && dif <= 1 {
					ctx.InterruptFlow()
					return
				}

				// Do this last in case any other cards want to interrupt the flow
				ctx.ScheduleAfter(func() {
					target.EvolveTo(card)
					ctx.Match().Chat("server", fmt.Sprintf("%s evolved %s to %s", card.Player().Name(), target.Name(), card.Name()))

					card.AddCondition(CantEvolve)
				})
			} else {
				ctx.InterruptFlow()
			}
		}
	// When Attacking
	case *match.AttackEvent:
		if event.ID == card.ID() {
			if !canAttack() {
				ctx.InterruptFlow()
				return
			}

			opponent := ctx.Match().Opponent(card.Player())

			attackfx := func(text string, fx func()) {
				// Do this last in case any other cards want to interrupt the flow
				ctx.ScheduleAfter(func() {
					trapzone, err := opponent.Container(match.TRAPZONE)
					if err != nil {
						logrus.Debug(err)
						return
					}

					ctx.Match().MessagePlayer(opponent, text)

					opponent.Action(
						append(opponent.GetCreatures(), trapzone...),
						1,
						1,
						true,
						func(action []string) {
							for _, id := range action {
								c, err := opponent.GetCard(id)
								if err != nil {
									logrus.Debugf("Search: %s", err)
									return
								}

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

								if !canAttack() {
									ctx.InterruptFlow()
									return
								}
							}
						},
						func() {
							// Update card
							card, err = card.Player().GetCard(event.ID)
							if err != nil {
								logrus.Debug(err)
								return
							}

							fx()
						},
					)
				})
			}

			if event.TargetID == "" {
				attackfx(
					fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name(), opponent.Name()),
					func() {
						opponent.Damage(card, ctx, card.GetAttack(ctx))
					})
			} else {
				target, err := opponent.GetCard(event.TargetID)
				if err != nil {
					logrus.Debug(err)
					return
				}

				if target.Zone() != match.BATTLEZONE {
					ctx.InterruptFlow()
					return
				}

				attackfx(
					fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name(), target.Name()),
					func() {
						// Update target
						target, err = opponent.GetCard(event.TargetID)
						if err != nil {
							logrus.Debug(err)
							return
						}

						ctx.Match().Battle(card, target, false)
					})
			}
		}
	// When blocking
	case *match.BlockEvent:
		if event.ID == card.ID() {
			if card.Tapped || card.Zone() != match.BATTLEZONE || card.GetDefence(ctx) <= 0 {
				ctx.InterruptFlow()
				return
			}

			// Do this last in case any other cards want to interrupt the flow
			ctx.ScheduleAfter(func() {
				card.Tapped = true
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
			ctx.ScheduleAfter(func() {
				if err := card.MoveCard(match.GRAVEYARD); err != nil {
					logrus.Debug(err)
				}
			})
		}
	}
}
