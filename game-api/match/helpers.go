package match

import (
	"fmt"

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
			if card.id == toTest {
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

	event, ok := ctx.event.(*PlayCardEvent)

	return ok && event.ID == card.id
}

// Evolution handles evolution
func Evolution(card *Card, ctx *Context, text string, filter func(*Card) bool) {

	cards := card.player.Filter(
		card.player.battlezone,
		text,
		1,
		1,
		true,
		filter,
	)

	if len(cards) < 1 {
		ctx.InterruptFlow()
		return
	}

	evoCtx := NewContext(ctx.match, &EvolveEvent{
		ID:       card.id,
		Creature: cards[0],
	})

	ctx.match.HandleFx(evoCtx)

	if evoCtx.Cancelled() {
		ctx.InterruptFlow()
	}
}

// Equipment handles equiping
func Equipment(card *Card, ctx *Context, text string, filter func(*Card) bool) {

	cards := card.player.Filter(
		card.player.battlezone,
		text,
		1,
		1,
		true,
		filter,
	)

	if len(cards) < 1 {
		ctx.InterruptFlow()
		return
	}

	equipCtx := NewContext(ctx.match, &Equip{
		ID:       card.id,
		Creature: cards[0],
	})

	ctx.match.HandleFx(equipCtx)

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

// Evolve ...
func Evolve(card *Card, target *Card) {

	if err := card.MoveCard(BATTLEZONE); err != nil {
		logrus.Debug(err)
		return
	}

	card.tapped = target.tapped

	target.AttachTo(card)

	// This is done to maintain a single identity for a creature
	card.id, target.id = target.id, card.id
	card.conditions, target.conditions = target.conditions, card.conditions
}

// Devolve ...
func Devolve(card *Card, source *Card) {

	cards := card.player.Filter(
		card.Attachments(),
		"Select a card",
		1,
		1,
		false,
		func(x *Card) bool { return x.Family() != family.Equipment })

	if len(cards) < 1 {
		card.player.match.Destroy(card, source, fmt.Sprintf("%s was destroyed by %s", card.name, source.name))
		return
	}

	if err := cards[0].MoveCard(BATTLEZONE); err != nil {
		logrus.Debug(err)
		return
	}

	cards[0].tapped = card.tapped

	for _, c := range card.Attachments() {
		c.AttachTo(cards[0])
	}

	// This is done to maintain a single identity for a creature
	card.id, cards[0].id = cards[0].id, card.id
	card.conditions, cards[0].conditions = cards[0].conditions, card.conditions

	card.player.match.Destroy(card, source, fmt.Sprintf("%s was destroyed by %s", card.name, source.name))
}
