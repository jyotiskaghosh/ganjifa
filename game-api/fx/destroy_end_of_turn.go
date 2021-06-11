package fx

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// DestroyEndOfTurn ...
func DestroyEndOfTurn(card *match.Card, ctx *match.Context) {
	if _, ok := ctx.Event().(*match.EndStep); ok && card.Zone() == match.BATTLEZONE {
		ctx.Match().Destroy(card, nil)
	}
}
