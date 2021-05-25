package match

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
	"github.com/jyotiskaghosh/ganjifa/game-api/family"

	"github.com/sirupsen/logrus"
)

// Search prompts the user to select n cards from the specified container
func (p *Player) Search(cards []*Card, text string, min int, max int, cancellable bool) []*Card {

	return p.Filter(cards, text, min, max, cancellable, func(c *Card) bool { return true })
}

// Filter prompts the user to select n cards from the specified container that matches the given filter
func (p *Player) Filter(cards []*Card, text string, min int, max int, cancellable bool, filter func(*Card) bool) []*Card {

	result := make([]*Card, 0)
	filtered := make([]*Card, 0)

	for _, mCard := range cards {
		if filter(mCard) {
			filtered = append(filtered, mCard)
		}
	}

	if len(filtered) < 1 {
		return result
	}

	p.match.NewAction(p, filtered, min, max, text, cancellable)
	defer p.match.CloseAction(p)

	for {
		select {
		case action := <-p.Action:
			{
				if len(action) < min || len(action) > max || !AssertCardsIn(filtered, action) {
					p.match.WarnPlayer(p, "The cards you selected does not meet the requirements")
					continue
				}

				for _, id := range action {

					c, err := p.GetCard(id)
					if err != nil {
						logrus.Debug(err)
						return result
					}
					result = append(result, c)
				}

				return result
			}

		case cancel := <-p.Cancel:
			if cancellable && cancel {
				return result
			}
		}
	}
}

// AssertCardsIn returns true or false based on if the specified card ids are present in the source []*Card
func AssertCardsIn(src []*Card, test []string) bool {

	for _, toTest := range test {

		ok := false

		for _, card := range src {
			if card.ID == toTest {
				ok = true
			}
		}

		if !ok {
			return false
		}
	}

	return true
}

// AmIPlayed returns true or false based on if the card is played
func AmIPlayed(card *Card, ctx *Context) bool {

	event, ok := ctx.Event.(*PlayCardEvent)

	return ok && event.ID == card.ID
}

// Evolution handles evolution
func Evolution(card *Card, ctx *Context, text string, filter func(*Card) bool) {

	creatures := card.player.Filter(
		card.player.battlezone,
		text,
		1,
		1,
		true,
		filter,
	)

	if len(creatures) < 1 {
		ctx.InterruptFlow()
		return
	}

	evoCtx := NewContext(ctx.Match, &Evolve{
		ID:       card.ID,
		Creature: creatures[0],
	})

	ctx.Match.HandleFx(evoCtx)

	if evoCtx.Cancelled() {
		ctx.InterruptFlow()
	}
}

// Equipment handles equiping
func Equipment(card *Card, ctx *Context, text string, filter func(*Card) bool) {

	creatures := card.player.Filter(
		card.player.battlezone,
		text,
		1,
		1,
		true,
		filter,
	)

	if len(creatures) < 1 {
		ctx.InterruptFlow()
		return
	}

	equipCtx := NewContext(ctx.Match, &Equip{
		ID:       card.ID,
		Creature: creatures[0],
	})

	ctx.Match.HandleFx(equipCtx)

	if equipCtx.Cancelled() {
		ctx.InterruptFlow()
	}
}

// RemoveEquipments ...
func (c *Card) RemoveEquipments() {
	for _, c = range c.Attachments() {
		if c.family == family.Equipment {
			if err := c.MoveCard(GRAVEYARD); err != nil {
				logrus.Debug(err)
				return
			}
		}
	}
}

// GetCreatures ...
func (p *Player) GetCreatures() []*Card {

	creatures := make([]*Card, 0)

	if cards, err := p.Container(BATTLEZONE); err == nil {
		creatures = append(creatures, cards...)
	}

	return creatures
}

// HasFamily if card has given family
func (c *Card) HasFamily(family string, ctx *Context) bool {

	for f := range c.GetFamily(ctx) {
		if family == f {
			return true
		}
	}
	return false
}

// HasCivilisation if card has given civilisation
func (c *Card) HasCivilisation(civilisation civ.Civilisation, ctx *Context) bool {

	for civ := range c.GetCivilisation(ctx) {
		if civilisation == civ {
			return true
		}
	}
	return false
}
