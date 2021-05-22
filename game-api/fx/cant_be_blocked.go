package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// CantBeBlocked ...
func CantBeBlocked(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.Block); ok && event.Attacker == card {

		ctx.Match.WarnPlayer(ctx.Match.Opponent(card.Player), fmt.Sprintf("%s can't be blocked", card.Name))
		ctx.InterruptFlow()
	}
}
