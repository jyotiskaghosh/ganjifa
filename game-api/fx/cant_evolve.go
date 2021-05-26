package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// CantEvolve ...
func CantEvolve(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event().(*match.EvolveEvent); ok && event.Creature == card {

		ctx.Match().WarnPlayer(card.Player(), fmt.Sprintf("%s can't evolve", card.Name()))
		ctx.InterruptFlow()
	}
}
