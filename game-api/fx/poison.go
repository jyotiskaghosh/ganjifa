package fx

import (
	"fmt"

	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Poisonous ...
func Poisonous(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event().(*match.Battle); ok && event.Defender == card {
		ctx.ScheduleAfter(func() {
			event.Attacker.AddCondition(DestroyEndOfTurn)
			ctx.Match().Chat("server", fmt.Sprintf("%s is poisoned by %s", event.Attacker.Name(), card.Name()))
		})
	}
}

// Venomous ...
func Venomous(card *match.Card, ctx *match.Context) {

	if event, ok := ctx.Event().(*match.Battle); ok && event.Attacker == card {
		ctx.ScheduleAfter(func() {
			event.Defender.AddCondition(DestroyEndOfTurn)
			ctx.Match().Chat("server", fmt.Sprintf("%s is poisoned by %s", event.Attacker.Name(), card.Name()))
		})
	}
}
