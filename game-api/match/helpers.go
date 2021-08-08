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
	cards := card.player.SearchAction(
		Filter(card.Attachments(), func(x *Card) bool { return x.Family() != family.Equipment }),
		"Select a card",
		1,
		1,
		false)

	if len(cards) < 1 {
		card.player.match.Destroy(card, src)
		return
	}

	if err := cards[0].MoveCard(BATTLEZONE); err != nil {
		logrus.Debugf("Couldn't devolve: %s", err)
		return
	}

	cards[0].Tapped = card.Tapped

	for _, card := range card.Attachments() {
		card.AttachTo(cards[0])
	}

	// This is done to maintain a single identity for a creature
	card.id, cards[0].id = cards[0].id, card.id
	card.conditions, cards[0].conditions = cards[0].conditions, card.conditions

	card.player.match.Destroy(card, src)
}

// Filter filters a slice of cards according to the filter func
func Filter(cards []*Card, filter func(*Card) bool) []*Card {
	filtered := make([]*Card, 0)

	for _, mCard := range cards {
		if filter(mCard) {
			filtered = append(filtered, mCard)
		}
	}

	return filtered
}
