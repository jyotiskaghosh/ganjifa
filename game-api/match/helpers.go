package match

import (
	"github.com/jyotiskaghosh/ganjifa/game-api/family"

	"github.com/sirupsen/logrus"
)

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

// Devolve ...
func Devolve(card *Card, src *Card) {

	cards := card.player.Filter(
		card.Attachments(),
		"Select a card",
		1,
		1,
		false,
		func(x *Card) bool { return x.Family() != family.Equipment })

	if len(cards) < 1 {
		card.player.match.Destroy(card, src)
		return
	}

	if err := cards[0].MoveCard(BATTLEZONE); err != nil {
		logrus.Debug(err)
		return
	}

	cards[0].tapped = card.tapped

	for _, card := range card.Attachments() {
		card.AttachTo(cards[0])
	}

	// This is done to maintain a single identity for a creature
	card.id, cards[0].id = cards[0].id, card.id
	card.conditions, cards[0].conditions = cards[0].conditions, card.conditions

	card.player.match.Destroy(card, src)
}

// CanPerformEvent ...
func CanPerformEvent(ctx *Context) bool {
	for _, c := range ctx.match.CollectCards() {
		for _, h := range c.GetHandlers(ctx) {
			h(c, ctx)
		}
	}
	return !ctx.cancel
}
