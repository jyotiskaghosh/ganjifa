package fx

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// CantBeAttacked ...
func CantBeAttacked(card *match.Card, ctx *match.Context) {
	if event, ok := ctx.Event().(*match.AttackEvent); ok && event.TargetID == card.ID() {
		ctx.InterruptFlow()
	}
}
