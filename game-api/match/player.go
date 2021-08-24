package match

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/sirupsen/logrus"
)

// Container for cards
type Container string

// containers
const (
	DECK       Container = "deck"
	HAND       Container = "hand"
	GRAVEYARD  Container = "graveyard"
	BATTLEZONE Container = "battlezone"
	TRAPZONE   Container = "trapzone"
	SOUL       Container = "soul"
)

// LIFE is the default life of player
const LIFE int = 20

// Writer is the interface for the output
type Writer interface {
	Write(interface{})
}

// Player holds the player data
type Player struct {
	name   string
	writer Writer

	deck       []*Card
	hand       []*Card
	graveyard  []*Card
	battlezone []*Card
	trapzone   []*Card
	soul       []*Card
	life       int

	ready bool

	turn   bool
	turnNo uint8

	Action chan []string
	Cancel chan bool

	match *Match

	wait  bool
	mutex *sync.Mutex
}

// newPlayer returns a new player
func newPlayer(name string, writer Writer, match *Match, turn bool) *Player {
	return &Player{
		name:       name,
		writer:     writer,
		deck:       make([]*Card, 0),
		hand:       make([]*Card, 0),
		graveyard:  make([]*Card, 0),
		battlezone: make([]*Card, 0),
		trapzone:   make([]*Card, 0),
		soul:       make([]*Card, 0),
		life:       LIFE,
		turn:       turn,
		match:      match,
		mutex:      &sync.Mutex{},
	}
}

// Name ...
func (p *Player) Name() string {
	return p.name
}

// Write data
func (p *Player) Write(msg interface{}) {
	p.writer.Write(msg)
}

// IsPlayerTurn is it the Player's turnNo
func (p *Player) IsPlayerTurn() bool {
	return p.turn
}

// Turn returns turn no.
func (p *Player) Turn() uint8 {
	return p.turnNo
}

// waiting assigns bool value to m.wait
func (p *Player) waiting(b bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.wait = b
}

// containerRef returns a pointer to one of the player's card zones based on the specified string
func (p *Player) containerRef(c Container) (*[]*Card, error) {
	switch c {
	case DECK:
		return &p.deck, nil
	case HAND:
		return &p.hand, nil
	case GRAVEYARD:
		return &p.graveyard, nil
	case BATTLEZONE:
		return &p.battlezone, nil
	case TRAPZONE:
		return &p.trapzone, nil
	case SOUL:
		return &p.soul, nil
	default:
		return nil, errors.New("Invalid container")
	}
}

// Container returns a copy of one of the player's card zones based on the specified string
func (p *Player) Container(c Container) ([]*Card, error) {
	switch c {
	case DECK:
		return p.deck, nil
	case HAND:
		return p.hand, nil
	case GRAVEYARD:
		return p.graveyard, nil
	case BATTLEZONE:
		return p.battlezone, nil
	case TRAPZONE:
		return p.trapzone, nil
	case SOUL:
		return p.soul, nil
	default:
		return nil, errors.New("Invalid container")
	}
}

// MapContainers performs the given Action on all cards in the specified containers
func (p *Player) MapContainers(fn func(*Card), containers ...Container) {
	for _, card := range p.CollectCards(containers...) {
		fn(card)
	}
}

// createDeck initializes a new deck from a list of card ids
func (p *Player) createDeck(cards []int) error {
	deck := make([]*Card, 0)

	for _, card := range cards {
		c, err := NewCard(p, card)
		if err != nil {
			logrus.Warnf("Failed to create card with id %d", card)
			continue
		}

		deck = append(deck, c)
	}

	if len(deck) != 40 {
		return errors.New("deck must have exactly 40 cards")
	}

	count := make(map[int]int)

	for _, card := range deck {
		count[card.cardID]++
		if count[card.cardID] > 4 {
			return errors.New("deck must have only 4 copies of a card")
		}
	}

	p.deck = deck
	p.ready = true

	return nil
}

// ShuffleDeck randomizes the order of cards in the players deck
func (p *Player) ShuffleDeck() {
	rand.Shuffle(len(p.deck), func(i, j int) { p.deck[i], p.deck[j] = p.deck[j], p.deck[i] })
	p.match.Chat("Server", fmt.Sprintf("%s's deck was shuffled", p.Name()))
}

// PeekDeck returns references to the next n cards in the deck
func (p *Player) PeekDeck(n int) []*Card {
	result := make([]*Card, 0)

	if len(p.deck) < n {
		n = len(p.deck)
	}

	for i := 0; i < n; i++ {
		result = append(result, p.deck[i])
	}

	return result
}

// DrawCards moves n cards from the players deck to their hand
func (p *Player) DrawCards(n int) {
	if len(p.deck) < n {
		n = len(p.deck)
	}

	for i := 0; i < n; i++ {
		if err := p.deck[i].MoveCard(HAND); err != nil {
			logrus.Debugf("Couldn't draw card: %s", err)
			return
		}
	}

	if n > 1 {
		p.match.Chat("Server", fmt.Sprintf("%s drew %v cards", p.Name(), n))
	} else {
		p.match.Chat("Server", fmt.Sprintf("%s drew %v card", p.Name(), n))
	}

	if len(p.deck) == 0 {
		p.match.End(p.match.Opponent(p), fmt.Sprintf("%s has no cards left in his deck", p.Name()))
	}
}

// HasCard checks if the given containers have the card
func (p *Player) HasCard(id string, containers ...Container) bool {
	for _, c := range p.CollectCards(containers...) {
		if c.id == id {
			return true
		}
	}
	return false
}

// CollectCards collects cards from the specified containers
func (p *Player) CollectCards(containers ...Container) []*Card {
	cards := make([]*Card, 0)

	for _, container := range containers {
		c, err := p.Container(container)
		if err != nil {
			logrus.Debug(err)
		}
		cards = append(cards, c...)
	}

	return cards
}

// Damage reduces life of player
func (p *Player) Damage(source *Card, ctx *Context, health uint8) {
	ctx = NewContext(ctx.match,
		&DamageEvent{
			Player: p,
			Source: source,
			Event:  ctx.event,
			Health: health,
		})

	ctx.ScheduleAfter(func() {
		if event, ok := ctx.event.(*DamageEvent); ok {
			p.life -= int(event.Health)
			ctx.match.Chat("Server", fmt.Sprintf("%s did %d damage to %s", source.name, health, p.Name()))
		}
	})

	ctx.match.HandleFx(ctx)

	if p.life <= 0 {
		ctx.match.End(p.match.Opponent(p), fmt.Sprintf("%s has no life left", p.Name()))
	}
}

// Heal reduces life of player
func (p *Player) Heal(source *Card, ctx *Context, health uint8) {
	ctx = NewContext(ctx.match,
		&HealEvent{
			Player: p,
			Source: source,
			Health: health,
		})

	ctx.ScheduleAfter(func() {
		if event, ok := ctx.event.(*HealEvent); ok {
			p.life += int(event.Health)
			ctx.match.Chat("Server", fmt.Sprintf("%s healed %d life for %s", source.name, health, p.Name()))
		}
	})

	ctx.match.HandleFx(ctx)
}

// Search prompts the user to select n cards from a slice of cards
func (p *Player) Search(cards []*Card, text string, min int, max int, cancellable bool) []*Card {
	p.Action = make(chan []string)
	p.Cancel = make(chan bool)

	defer close(p.Cancel)
	defer close(p.Action)

	p.waiting(false)

	result := make([]*Card, 0)

	p.match.NewAction(p, cards, min, max, text, cancellable)
	defer p.match.CloseAction(p)

	for {
		select {
		case Action := <-p.Action:
			{
				if len(Action) < min || len(Action) > max || !AssertCardsIn(cards, Action...) {
					p.match.WarnPlayer(p, "The cards you selected does not meet the requirements")
					continue
				}
				for _, id := range Action {
					c, err := GetCard(id, cards)
					if err != nil {
						logrus.Debugf("Search: %s", err)
						return result
					}

					result = append(result, c)
				}
				return result
			}
		case Cancel := <-p.Cancel:
			if cancellable && Cancel {
				return result
			}
		}
	}
}

// denormalized returns a server.PlayerState
func (p *Player) denormalized() PlayerState {
	return PlayerState{
		Life:       p.life,
		Deck:       len(p.deck),
		Hand:       denormalizeCards(p.hand),
		Graveyard:  denormalizeCards(p.graveyard),
		Battlezone: denormalizeCards(p.battlezone),
		Trapzone:   denormalizeCards(p.trapzone),
	}
}
