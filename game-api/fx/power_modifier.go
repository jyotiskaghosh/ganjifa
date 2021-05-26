package fx

import "github.com/jyotiskaghosh/ganjifa/game-api/match"

// AttackModifier ...
func AttackModifier(card *match.Card, ctx *match.Context, n int) {
	if event, ok := ctx.Event().(*match.GetAttackEvent); ok && event.Card == card {
		event.Attack += n
	}
}

// DefenceModifier ...
func DefenceModifier(card *match.Card, ctx *match.Context, n int) {
	if event, ok := ctx.Event().(*match.GetDefenceEvent); ok && event.Card == card {
		event.Defence += n
	}
}
