package fx

import "github.com/jyotiskaghosh/ganjifa/game-api/match"

// CantEvolve ...
func CantEvolve(card *match.Card, ctx *match.Context) {
	if event, ok := ctx.Event().(*match.EvolveEvent); ok && event.Creature == card {
		ctx.InterruptFlow()
	}
}
