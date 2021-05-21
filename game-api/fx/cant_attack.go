package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// CantAttackPlayers prevents a card from attacking players
func CantAttackPlayers(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.AttackPlayer); ok && event.ID == card.ID {

		ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack players", card.Name))
		ctx.InterruptFlow()
	}
}

// CantAttackCreatures prevents a card from attacking players
func CantAttackCreatures(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.AttackCreature); ok && event.ID == card.ID {

		ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't attack creatures", card.Name))
		ctx.InterruptFlow()
	}
}
