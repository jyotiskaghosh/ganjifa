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

			card.Tap(false)
			card.ClearConditions()
		}

	// On playing the card
	case *match.PlayCardEvent:
		if event.ID == card.ID() {

			if card.GetRank(ctx) > 0 {

				// Do this last in case any other cards want to interrupt the flow
				ctx.Override(func() {
					match.Evolution(
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
					ctx.Match().Chat("Server", fmt.Sprintf("%s summoned %s to the battle zone", card.Player().Name(), card.Name()))

					card.AddCondition(CantEvolve)
				})
			}
		}

	// On evolve
	case *match.EvolveEvent:
		if event.ID == card.ID() {

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {

				match.Evolve(card, event.Creature)

				card.AddCondition(CantEvolve)

				ctx.Match().Chat("Server", fmt.Sprintf("%s summoned %s to the battle zone", card.Player().Name(), card.Name()))
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

				blockers := make([]*match.Card, 0)
				for _, c := range opponent.GetCreatures() {
					if !c.Tapped() {
						blockers = append(blockers, c)
					}
				}

				hiddenzone, err := ctx.Match().Opponent(card.Player()).Container(match.HIDDENZONE)
				if err != nil {
					logrus.Debug(err)
					return
				}

				defer func() {
					card.Tap(true)
				}()

				ctx.Match().NewAction(opponent, append(blockers, hiddenzone...), 1, 1, fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name(), opponent.Name()), true)
				defer ctx.Match().CloseAction(opponent)

				for {
					select {
					case action := <-opponent.Action:
						{
							if len(action) < 1 || len(action) > 1 || !match.AssertCardsIn(append(blockers, hiddenzone...), action) {
								ctx.Match().WarnPlayer(opponent, "The cards you selected does not meet the requirements")
								continue
							}

							for _, id := range action {

								c, err := opponent.GetCard(id)
								if err != nil {
									logrus.Debug(err)
									return
								}

								if c.Zone() == match.HIDDENZONE {

									// Playing set down card
									ctx.Match().HandleFx(match.NewContext(ctx.Match(), &match.React{
										ID:    c.ID(),
										Event: ctx.Event(),
									}))
									ctx.Match().BroadcastState()

									// Check if attack can still go through
									ctx.Match().ResolveEvent(ctx)
									if ctx.Cancelled() {
										return
									}

									// Update card
									card, err = card.Player().GetCard(event.ID)
									if err != nil {
										logrus.Debug(err)
										return
									}

									blockers = make([]*match.Card, 0)
									for _, c := range opponent.GetCreatures() {
										if !c.Tapped() {
											blockers = append(blockers, c)
										}
									}

									hiddenzone, err := ctx.Match().Opponent(card.Player()).Container(match.HIDDENZONE)
									if err != nil {
										logrus.Debug(err)
										return
									}

									ctx.Match().NewAction(opponent, append(blockers, hiddenzone...), 1, 1, fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name(), opponent.Name()), true)

								} else {

									// Blocking attack
									blockCtx := match.NewContext(ctx.Match(), &match.Block{
										Attacker: card,
										Blocker:  c,
									})

									ctx.Match().HandleFx(blockCtx)
									if blockCtx.Cancelled() {
										continue
									}

									return
								}
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
				ctx.InterruptFlow()
				logrus.Debug(err)
				return
			}

			if target.Zone() != match.BATTLEZONE {
				ctx.InterruptFlow()
				return
			}

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {

				blockers := make([]*match.Card, 0)
				for _, c := range opponent.GetCreatures() {
					if !c.Tapped() {
						blockers = append(blockers, c)
					}
				}

				hiddenzone, err := ctx.Match().Opponent(card.Player()).Container(match.HIDDENZONE)
				if err != nil {
					logrus.Debug(err)
					return
				}

				defer func() {
					card.Tap(true)
				}()

				ctx.Match().NewAction(opponent, append(blockers, hiddenzone...), 1, 1, fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name(), target.Name()), true)
				defer ctx.Match().CloseAction(opponent)

				for {
					select {
					case action := <-opponent.Action:
						{
							if len(action) < 1 || len(action) > 1 || !match.AssertCardsIn(append(blockers, hiddenzone...), action) {
								ctx.Match().WarnPlayer(opponent, "The cards you selected does not meet the requirements")
								continue
							}

							for _, id := range action {

								c, err := opponent.GetCard(id)
								if err != nil {
									logrus.Debug(err)
									return
								}

								if c.Zone() == match.HIDDENZONE {

									// Playing set down card
									ctx.Match().HandleFx(match.NewContext(ctx.Match(), &match.React{
										ID:    c.ID(),
										Event: ctx.Event(),
									}))
									ctx.Match().BroadcastState()

									// Check if attack can still go through
									ctx.Match().ResolveEvent(ctx)
									if ctx.Cancelled() {
										return
									}

									// Update card
									card, err = card.Player().GetCard(event.ID)
									if err != nil {
										logrus.Debug(err)
										return
									}

									// Update target
									target, err = opponent.GetCard(event.TargetID)
									if err != nil {
										logrus.Debug(err)
										return
									}

									blockers = make([]*match.Card, 0)
									for _, c := range opponent.GetCreatures() {
										if !c.Tapped() {
											blockers = append(blockers, c)
										}
									}

									hiddenzone, err := ctx.Match().Opponent(card.Player()).Container(match.HIDDENZONE)
									if err != nil {
										logrus.Debug(err)
										return
									}

									ctx.Match().NewAction(opponent, append(blockers, hiddenzone...), 1, 1, fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name(), target.Name()), true)

								} else {

									// Blocking attack
									blockCtx := match.NewContext(ctx.Match(), &match.Block{
										Attacker: card,
										Blocker:  c,
									})

									ctx.Match().HandleFx(blockCtx)
									if blockCtx.Cancelled() {
										continue
									}

									return
								}
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
	case *match.Block:
		if event.Blocker == card {

			if card.Tapped() || card.Zone() != match.BATTLEZONE || card.GetDefence(ctx) <= 0 {
				ctx.InterruptFlow()
				return
			}

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {

				defer event.Blocker.Tap(true)

				ctx.Match().Battle(event.Attacker, card, true)
			})
		}

	// When destroyed
	case *match.CreatureDestroyed:
		if event.Card == card {

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {
				if err := card.MoveCard(match.GRAVEYARD); err != nil {
					logrus.Debug(err)
				}
			})
		}
	}
}
