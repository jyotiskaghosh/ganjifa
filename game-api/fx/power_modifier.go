package fx

import "github.com/jyotiskaghosh/ganjifa/game-api/match"

// AttackModifier ...
func AttackModifier(card *match.Card, ctx *match.Context, n uint8) {
	if event, ok := ctx.Event().(*match.GetAttackEvent); ok && event.ID == card.ID() {
		event.Attack += n
	}
}

// DefenceModifier ...
func DefenceModifier(card *match.Card, ctx *match.Context, n uint8) {
	if event, ok := ctx.Event().(*match.GetDefenceEvent); ok && event.ID == card.ID() {
		event.Defence += n
	}
}
