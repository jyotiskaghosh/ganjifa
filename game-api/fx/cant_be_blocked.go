package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// CantBeBlocked ...
func CantBeBlocked(card *match.Card, ctx *match.Context) {
	if event, ok := ctx.Event().(*match.BlockEvent); ok && event.Attacker == card {
		ctx.Match().WarnPlayer(ctx.Match().Opponent(card.Player()), fmt.Sprintf("Can't block %s", card.Name()))
		ctx.InterruptFlow()
	}
}
