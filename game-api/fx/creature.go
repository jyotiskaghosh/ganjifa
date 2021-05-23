package fx

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Creature has default behaviours for creatures
func Creature(card *match.Card, ctx *match.Context) {

	switch event := ctx.Event.(type) {

	// Untap the card
	case *match.UntapStep:
		if card.Player.IsPlayerTurn() {

			card.Tap(false)
			card.ClearConditions()
		}

	// On playing the card
	case *match.PlayCardEvent:
		if event.ID == card.ID {

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
	case *match.Evolve:
		if event.ID == card.ID {

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {

				if err := card.MoveCard(match.BATTLEZONE); err != nil {
					logrus.Debug(err)
					return
				}

				ctx.Match.Chat("Server", fmt.Sprintf("%s summoned %s to the battle zone", card.Player.Name(), card.Name))

				card.Tap(event.Creature.Tapped())

				event.Creature.AttachTo(card)

				card.AddCondition(event.Creature.Conditions()...)
				card.AddCondition(CantEvolve)

				// This is done to maintain a singular identity for a creature
				card.ID, event.Creature.ID = event.Creature.ID, card.ID
			})
		}

	// Attack the player
	case *match.AttackPlayer:
		if event.ID == card.ID {

			if card.Tapped() || card.Zone != match.BATTLEZONE {
				ctx.InterruptFlow()
				return
			}

			if card.GetAttack(ctx) <= 0 {
				ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack", card.Name))
				ctx.InterruptFlow()
				return
			}

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {

				opponent := ctx.Match.Opponent(card.Player)

				spellzone, err := ctx.Match.Opponent(card.Player).ContainerRef(match.SPELLZONE)
				if err != nil {
					ctx.InterruptFlow()
					logrus.Debug(err)
					return
				}

				blockers := make([]*match.Card, 0)
				for _, c := range opponent.GetCreatures() {
					if !c.Tapped() {
						blockers = append(blockers, c)
					}
				}

				ctx.Match.NewAction(opponent, append(blockers, *spellzone...), 1, 1, fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name, opponent.Name()), true)
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

									// Playing set down card
									ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.React{
										ID:    c.ID,
										Event: ctx.Event,
									}))
									ctx.Match.BroadcastState()

									// Check if attack can still go through
									ctx.Match.ResolveEvent(ctx)
									if ctx.Cancelled() {
										return
									}

									ctx.Match.NewAction(opponent, append(blockers, *spellzone...), 1, 1, fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name, opponent.Name()), true)

								} else {

									// Blocking attack
									ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.Block{
										Attacker: card,
										Blocker:  c,
									}))
									ctx.Match.BroadcastState()

									// Check if attack can still go through
									ctx.Match.ResolveEvent(ctx)
									if ctx.Cancelled() {
										return
									}
								}
							}
						}

					case cancel := <-opponent.Cancel:
						if cancel {

							card.Tap(true)
							opponent.Damage(card, ctx, card.GetAttack(ctx))
							return
						}
					}
				}
			})
		}

	// Attack a creature
	case *match.AttackCreature:
		if event.ID == card.ID {

			if card.Tapped() || card.Zone != match.BATTLEZONE {
				ctx.InterruptFlow()
				return
			}

			if card.GetAttack(ctx) <= 0 {
				ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack", card.Name))
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

			if defender.Zone != match.BATTLEZONE {
				ctx.Match.WarnPlayer(card.Player, "creature to attack is not in battlezone")
				ctx.InterruptFlow()
				return
			}

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {

				spellzone, err := ctx.Match.Opponent(card.Player).ContainerRef(match.SPELLZONE)
				if err != nil {
					ctx.InterruptFlow()
					logrus.Debug(err)
					return
				}

				blockers := make([]*match.Card, 0)
				for _, c := range opponent.GetCreatures() {
					if !c.Tapped() {
						blockers = append(blockers, c)
					}
				}

				ctx.Match.NewAction(opponent, append(blockers, *spellzone...), 1, 1, fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name, defender.Name), true)
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

									// Playing set down card
									ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.React{
										ID:    c.ID,
										Event: ctx.Event,
									}))
									ctx.Match.BroadcastState()

									// Check if attack can still go through
									ctx.Match.ResolveEvent(ctx)
									if ctx.Cancelled() {
										return
									}

									// This is done in cases the target evolves
									defender, err = opponent.GetCard(event.TargetID)
									if err != nil {
										ctx.Match.WarnPlayer(card.Player, "creature to attack was not found")
										ctx.InterruptFlow()
										return
									}

									ctx.Match.NewAction(opponent, append(blockers, *spellzone...), 1, 1, fmt.Sprintf("%s is attacking %s, you may play a set down card or block with a creature", card.Name, defender.Name), true)

								} else {

									// Blocking attack
									ctx.Match.HandleFx(match.NewContext(ctx.Match, &match.Block{
										Attacker: card,
										Blocker:  c,
									}))
									ctx.Match.BroadcastState()

									// Check if attack can still go through
									ctx.Match.ResolveEvent(ctx)
									if ctx.Cancelled() {
										return
									}
								}
							}
						}

					case cancel := <-opponent.Cancel:
						if cancel {

							ctx.Match.Battle(card, defender, false)
							return
						}
					}
				}
			})
		}

	// When blocking
	case *match.Block:
		if event.Blocker == card {

			if card.Tapped() || card.Zone != match.BATTLEZONE {
				ctx.InterruptFlow()
				return
			}

			if card.GetDefence(ctx) <= 0 {
				ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't block", card.Name))
				ctx.InterruptFlow()
				return
			}

			// Do this last in case any other cards want to interrupt the flow
			ctx.Override(func() {
				event.Blocker.Tap(true)
				ctx.Match.Battle(event.Attacker, card, true)
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
