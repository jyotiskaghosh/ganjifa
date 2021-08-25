package fx

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, ctx *match.Context) {
	opponent := ctx.Match().Opponent(card.Player())

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
			} else {
				// Do this last in case any other cards want to interrupt the flow
				ctx.ScheduleAfter(func() {
					cards := card.Player().Search(
						match.Filter(
							card.Player().CollectCards(match.BATTLEZONE),
							func(c *match.Card) bool {
								return c.HasFamily(card.Family(), ctx) && card.GetRank(ctx)-c.GetRank(ctx) == 1
							},
						),
						fmt.Sprintf("choose a creature to evolve %s", card.Name()),
						1,
						1,
						true)

					if len(cards) > 0 {
						evoCtx := match.NewContext(ctx.Match(), &match.Evolve{
							ID:     card.ID(),
							Target: cards[0],
						})
						ctx.Match().HandleFx(evoCtx)

						if evoCtx.Cancelled() {
							ctx.InterruptFlow()
						}
					} else {
						ctx.InterruptFlow()
					}
				})
			}
		}
	// When evolving
	case *match.Evolve:
		if event.ID == card.ID() {
			// Do this last in case any other cards want to interrupt the flow
			ctx.ScheduleAfter(func() {
				event.Target.EvolveTo(card)
				ctx.Match().Chat("server", fmt.Sprintf("%s evolved %s to %s", card.Player().Name(), event.Target.Name(), card.Name()))
				card.AddCondition(CantEvolve)
			})
		}
	// When attacking player
	case *match.AttackPlayer:
		if event.ID == card.ID() {
			if card.Tapped || card.Zone() != match.BATTLEZONE || card.GetAttack(ctx) <= 0 {
				ctx.InterruptFlow()
				return
			}

			ctx.ScheduleAfter(func() {
				opponent.Action = make(chan []string)
				opponent.Cancel = make(chan bool)

				defer close(opponent.Cancel)
				defer close(opponent.Action)

				text := fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name(), opponent.Name())
				cards := match.Filter(opponent.CollectCards(match.BATTLEZONE, match.TRAPZONE), func(c *match.Card) bool { return !c.Tapped })

				if len(cards) < 1 {
					opponent.Damage(card, ctx, card.GetAttack(ctx))
					return
				}

				ctx.Match().NewAction(opponent, cards, 1, 1, text, true)
				defer ctx.Match().CloseAction(opponent)

				for {
					select {
					case action := <-opponent.Action:
						{
							if len(action) < 1 || len(action) > 1 || !match.AssertCardsIn(cards, action...) {
								ctx.Match().WarnPlayer(opponent, "The cards you selected does not meet the requirements")
								continue
							}

							for _, id := range action {
								c, err := match.GetCard(
									id,
									opponent.CollectCards(match.BATTLEZONE, match.TRAPZONE))
								if err != nil {
									logrus.Debugf("Search: %s", err)
									continue
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

								// Update card
								card, err = match.GetCard(event.ID, card.Player().CollectCards(match.BATTLEZONE))
								if err != nil {
									logrus.Debug(err)
									ctx.InterruptFlow()
									return
								}

								if card.Tapped || card.Zone() != match.BATTLEZONE || card.GetAttack(ctx) <= 0 {
									ctx.InterruptFlow()
									return
								}

								text := fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name(), opponent.Name())
								cards = match.Filter(opponent.CollectCards(match.BATTLEZONE, match.TRAPZONE), func(c *match.Card) bool { return !c.Tapped })

								if len(cards) < 1 {
									opponent.Damage(card, ctx, card.GetAttack(ctx))
									return
								}

								// refreshes the action
								ctx.Match().NewAction(opponent, cards, 1, 1, text, true)
								opponent.Action = make(chan []string)
								opponent.Cancel = make(chan bool)
							}
						}
					case cancel := <-opponent.Cancel:
						if cancel {
							opponent.Damage(card, ctx, card.GetAttack(ctx))
							return
						}
					}
				}
			})
		}
	// When attacking creature
	case *match.AttackCreature:
		if event.ID == card.ID() {
			target, err := match.GetCard(
				event.TargetID,
				match.Filter(opponent.CollectCards(match.BATTLEZONE), func(c *match.Card) bool { return c.Tapped }),
			)
			if err != nil {
				logrus.Debug(err)
				ctx.InterruptFlow()
				return
			}

			if card.Tapped ||
				card.Zone() != match.BATTLEZONE ||
				card.GetAttack(ctx) <= 0 ||
				target.Zone() != match.BATTLEZONE {
				ctx.InterruptFlow()
				return
			}

			ctx.ScheduleAfter(func() {
				opponent.Action = make(chan []string)
				opponent.Cancel = make(chan bool)

				defer close(opponent.Cancel)
				defer close(opponent.Action)

				text := fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name(), target.Name())
				cards := match.Filter(opponent.CollectCards(match.BATTLEZONE, match.TRAPZONE), func(c *match.Card) bool { return !c.Tapped })

				if len(cards) < 1 {
					ctx.Match().Battle(card, target, false)
					return
				}

				ctx.Match().NewAction(opponent, cards, 1, 1, text, true)
				defer ctx.Match().CloseAction(opponent)

				for {
					select {
					case action := <-opponent.Action:
						{
							if len(action) < 1 || len(action) > 1 || !match.AssertCardsIn(cards, action...) {
								ctx.Match().WarnPlayer(opponent, "The cards you selected does not meet the requirements")
								continue
							}

							for _, id := range action {
								c, err := match.GetCard(
									id,
									opponent.CollectCards(match.BATTLEZONE, match.TRAPZONE))
								if err != nil {
									logrus.Debugf("Search: %s", err)
									continue
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

								// Update card
								card, err = match.GetCard(event.ID, card.Player().CollectCards(match.BATTLEZONE))
								if err != nil {
									logrus.Debug(err)
									ctx.InterruptFlow()
									return
								}

								// Update target
								target, err := match.GetCard(
									event.TargetID,
									match.Filter(opponent.CollectCards(match.BATTLEZONE), func(c *match.Card) bool { return c.Tapped }),
								)
								if err != nil {
									logrus.Debug(err)
									ctx.InterruptFlow()
									return
								}

								if card.Tapped ||
									card.Zone() != match.BATTLEZONE ||
									card.GetAttack(ctx) <= 0 ||
									target.Zone() != match.BATTLEZONE {
									ctx.InterruptFlow()
									return
								}

								text = fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name(), target.Name())
								cards = match.Filter(opponent.CollectCards(match.BATTLEZONE, match.TRAPZONE), func(c *match.Card) bool { return !c.Tapped })

								if len(cards) < 1 {
									ctx.Match().Battle(card, target, false)
									return
								}

								// refreshes the action
								ctx.Match().NewAction(opponent, cards, 1, 1, text, true)
								opponent.Action = make(chan []string)
								opponent.Cancel = make(chan bool)
							}
						}
					case cancel := <-opponent.Cancel:
						if cancel {
							ctx.Match().Battle(card, target, false)
							return
						}
					}
				}
			})
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
