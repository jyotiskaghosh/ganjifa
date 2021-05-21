package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// CantBlock ...
func CantBlock(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event.(*match.BlockEvent); ok && event.Blocker == card {

		ctx.Match.WarnPlayer(card.Player, fmt.Sprintf("%s can't block", card.Name))
		ctx.InterruptFlow()
	}
}
