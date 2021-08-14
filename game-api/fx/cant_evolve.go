package fx

import "github.com/jyotiskaghosh/ganjifa/game-api/match"

// CantEvolve ...
func CantEvolve(card *match.Card, ctx *match.Context) {
	if event, ok := ctx.Event().(*match.PlayCardEvent); ok && event.TargetID == card.ID() {
		ctx.InterruptFlow()
	}
}
