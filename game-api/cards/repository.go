package cards

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/cards/set01"
	"github.com/jyotiskaghosh/ganjifa/game-api/match"
)

// Sets is a map of pointers to the available card sets
var Sets = map[string]*map[int]match.CardConstructor{
	"set-01": &Set01,
}

var ctors = []match.CardConstructor{
	set01.AirMail,
	set01.Amrita,
	set01.Astrakara,
	set01.Atayi,
	set01.Ayudhabhrt,
	set01.Blizzard,
	set01.Cataka,
	set01.Churika,
	set01.DeadlyZebrafish,
	set01.Degenerate,
	set01.Dvipin,
	set01.EnergySurge,
	set01.Fireball,
	set01.FrostBreath,
	set01.Khadga,
	set01.Krostr,
	set01.Kukkutah,
	set01.LeechLife,
	set01.MagmaGeyser,
	set01.MahisiPipilika,
	set01.Masaka,
	set01.Matsyaka,
	set01.Pipilika,
	set01.RainOfArrows,
	set01.RapidEvolution,
	set01.Sainika,
	set01.Salavrka,
	set01.Sastravikrayin,
	set01.ScopeLens,
	set01.ShellArmor,
	set01.Simha,
	set01.Syena,
	set01.Tailwind,
	set01.TidalWave,
	set01.Tornado,
	set01.TorpedoingBarracuda,
	set01.Vanara,
	set01.VampireFangs,
	set01.Whirlwind,
	set01.WindCloak,
}

// Set01 is a map with all the card id's in the game and corresponding CardConstructor for set01
var Set01 = func() map[int]match.CardConstructor {
	m := make(map[int]match.CardConstructor, 0)
	for i, ctor := range ctors {
		m[i] = ctor
	}
	return m
}()
