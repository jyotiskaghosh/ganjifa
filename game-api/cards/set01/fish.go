package set01

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
	"github.com/jyotiskaghosh/ganjifa/game-api/family"
	"github.com/jyotiskaghosh/ganjifa/game-api/fx"
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Matsyaka ...
func Matsyaka() *match.Card {
	cb := match.CardBuilder{
		Name:    "Matsyaka",
		Rank:    0,
		Civ:     civ.APAS,
		Family:  family.Fish,
		Attack:  1,
		Defence: 1,
		Effects: []match.HandlerFunc{
			fx.Creature,
			func(card *match.Card, ctx *match.Context) {
				if event, ok := ctx.Event().(*match.GetAttackEvent); ok {
					if card.AttachedTo() != nil && event.ID == card.AttachedTo().ID() {
						event.Attack++
					}
				}

				if event, ok := ctx.Event().(*match.GetDefenceEvent); ok {
					if card.AttachedTo() != nil && event.ID == card.AttachedTo().ID() {
						event.Defence++
					}
				}
			},
		},
	}

	return cb.Build()
}

// DeadlyZebrafish ...
func DeadlyZebrafish() *match.Card {
	cb := match.CardBuilder{
		Name:    "Deadly Zebrafish",
		Rank:    0,
		Civ:     civ.APAS,
		Family:  family.Fish,
		Attack:  1,
		Defence: 1,
		Effects: []match.HandlerFunc{
			fx.Creature,
			fx.Poisonous,
		},
	}

	return cb.Build()
}

// TorpedoingBarracuda ...
func TorpedoingBarracuda() *match.Card {
	cb := match.CardBuilder{
		Name:    "Torpedoing Barracuda",
		Rank:    1,
		Civ:     civ.APAS,
		Family:  family.Fish,
		Attack:  3,
		Defence: 1,
		Effects: []match.HandlerFunc{
			fx.Creature,
			fx.CantBeBlocked,
		},
	}

	return cb.Build()
}
