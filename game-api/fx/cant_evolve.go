package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// CantEvolve ...
func CantEvolve(card *match.Card, ctx *match.Context) {
	if event, ok := ctx.Event().(*match.Evolve); ok && event.Target == card {
		ctx.Match().WarnPlayer(card.Player(), fmt.Sprintf("Can't evolve %s", card.Name()))
		ctx.InterruptFlow()
	}
}
