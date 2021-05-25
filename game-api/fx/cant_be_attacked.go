package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// CantBeAttacked ...
func CantBeAttacked(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.AttackCreature); ok && event.TargetID == card.ID {

		ctx.Match.WarnPlayer(ctx.Match.Opponent(card.Player()), fmt.Sprintf("%s can't be attacked", card.Name()))
		ctx.InterruptFlow()
	}
}
