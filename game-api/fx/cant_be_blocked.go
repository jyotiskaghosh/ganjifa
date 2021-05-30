package fx

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// CantBeBlocked ...
func CantBeBlocked(card *match.Card, ctx *match.Context) {
	if event, ok := ctx.Event().(*match.BlockEvent); ok && event.Attacker == card {
		ctx.InterruptFlow()
	}
}
