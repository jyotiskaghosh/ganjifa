package match

import (
	"fmt"
	"reflect"

	"github.com/jyotiskaghosh/ganjifa/game-api/civ"
	"github.com/jyotiskaghosh/ganjifa/game-api/family"

	"github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
)

// Card stores card data
type Card struct {
	id string

	cardID  int
	name    string
	rank    uint8
	civ     civ.Civilisation
	family  string
	attack  uint8
	defence uint8
	effects []HandlerFunc

	zone       Container
	Tapped     bool
	attachedTo *Card
	conditions []HandlerFunc // Temporary effects

	player *Player
}

// CardBuilder is a builder for Card
type CardBuilder struct {
	Name    string
	Rank    uint8
	Civ     civ.Civilisation
	Family  string
	Attack  uint8
	Defence uint8
	Effects []HandlerFunc
}

// Build constructs a card with the values of CardBuilder
func (cb *CardBuilder) Build() *Card {
	return &Card{
		name:    cb.Name,
		rank:    cb.Rank,
		civ:     cb.Civ,
		family:  cb.Family,
		attack:  cb.Attack,
		defence: cb.Defence,
		effects: cb.Effects,
	}
}

// NewCard returns a new, initialized card
func NewCard(p *Player, cardID int) (*Card, error) {
	c, err := CardCtor(cardID)
	if err != nil {
		logrus.Warnf("NewCard: %s", err)
		return nil, err
	}

	id, err := shortid.Generate()
	if err != nil {
		return nil, err
	}

	c.id = id
	c.cardID = cardID
	c.player = p
	c.zone = DECK

	return c, nil
}

// ID ...
func (c *Card) ID() string {
	return c.id
}

// CardID ...
func (c *Card) CardID() int {
	return c.cardID
}

// Name ...
func (c *Card) Name() string {
	return c.name
}

// Rank ...
func (c *Card) Rank() uint8 {
	return c.rank
}

// Civ ...
func (c *Card) Civ() civ.Civilisation {
	return c.civ
}

// Family ...
func (c *Card) Family() string {
	return c.family
}

// Attack ...
func (c *Card) Attack() uint8 {
	return c.attack
}

// Defence ...
func (c *Card) Defence() uint8 {
	return c.defence
}

// Zone ...
func (c *Card) Zone() Container {
	return c.zone
}

// Player ...
func (c *Card) Player() *Player {
	return c.player
}

// AddCondition adds temporary handler functions
func (c *Card) AddCondition(effects ...HandlerFunc) {
	c.conditions = append(c.conditions, effects...)
}

// RemoveCondition removes all instances of the given handler from the cards conditions
func (c *Card) RemoveCondition(handler HandlerFunc) {
	tmp := make([]HandlerFunc, 0)

	for _, condition := range c.conditions {
		if reflect.ValueOf(condition).Pointer() != reflect.ValueOf(handler).Pointer() {
			tmp = append(tmp, condition)
		}
	}

	c.conditions = tmp
}

// ClearConditions removes all conditions from the card
func (c *Card) ClearConditions() {
	c.conditions = make([]HandlerFunc, 0)
}

// AttachedTo returns the card that this card is attached to
func (c *Card) AttachedTo() *Card {
	return c.attachedTo
}

// Attachments returns a copy of the card's attached cards
func (c *Card) Attachments() []*Card {
	result := make([]*Card, 0)

	for _, card := range c.player.CollectCards(SOUL) {
		if card.attachedTo == c {
			result = append(result, card)
		}
	}

	return result
}

// AttachTo attaches c to card
func (c *Card) AttachTo(card *Card) {
	if card.zone != BATTLEZONE {
		logrus.Debug("Can't attach card to creature not on battlezone")
		return
	}

	c.attachedTo = card

	if err := c.MoveCard(SOUL); err != nil {
		logrus.Debugf("AttachTo: %s", err)
		return
	}
}

// Detach detaches a card
func (c *Card) Detach() {
	c.attachedTo = nil
}

// MoveCard tries to move a card to container b
func (c *Card) MoveCard(destination Container) error {
	from, err := c.player.containerRef(c.zone)
	if err != nil {
		return fmt.Errorf("Couldn't move card %s(%s) to %s", c.id, c.name, destination)
	}

	to, err := c.player.containerRef(destination)
	if err != nil {
		return fmt.Errorf("Couldn't move card %s(%s) to %s", c.id, c.name, destination)
	}

	temp := make([]*Card, 0)

	for _, card := range *from {
		if card.id != c.id {
			temp = append(temp, card)
		}
	}

	*from = temp

	*to = append(*to, c)

	f := c.zone
	c.zone = destination

	if f == SOUL && destination != SOUL {
		c.Detach()
	}

	ctx := NewContext(c.player.match, &CardMoved{
		ID:   c.id,
		From: f,
		To:   destination,
	})
	c.player.match.HandleFx(ctx)

	for _, card := range c.Attachments() {
		card.attachedTo = c.attachedTo
		c.MoveCard(destination)
	}

	return nil
}

// GetRank returns the rank of a given card
func (c *Card) GetRank(ctx *Context) uint8 {
	e := &GetRankEvent{
		ID:    c.id,
		Event: ctx.event,
		Rank:  c.rank,
	}
	ctx.match.HandleFx(NewContext(ctx.match, e))

	return e.Rank
}

// GetCivilisation returns the Civ of a given card
func (c *Card) GetCivilisation(ctx *Context) map[civ.Civilisation]bool {
	e := &GetCivilisationEvent{
		ID:    c.id,
		Event: ctx.event,
		Civ:   map[civ.Civilisation]bool{c.civ: true},
	}
	ctx.match.HandleFx(NewContext(ctx.match, e))

	return e.Civ
}

// GetFamily returns the family of a given card
func (c *Card) GetFamily(ctx *Context) map[string]bool {
	e := &GetFamilyEvent{
		ID:     c.id,
		Event:  ctx.event,
		Family: map[string]bool{c.family: true},
	}
	ctx.match.HandleFx(NewContext(ctx.match, e))

	return e.Family
}

// GetAttack returns the attack of a given card
func (c *Card) GetAttack(ctx *Context) uint8 {
	e := &GetAttackEvent{
		ID:     c.id,
		Event:  ctx.event,
		Attack: c.attack,
	}
	ctx.match.HandleFx(NewContext(ctx.match, e))

	return e.Attack
}

// GetDefence returns the defence of a given card
func (c *Card) GetDefence(ctx *Context) uint8 {
	e := &GetDefenceEvent{
		ID:      c.id,
		Event:   ctx.event,
		Defence: c.defence,
	}
	ctx.match.HandleFx(NewContext(ctx.match, e))

	return e.Defence
}

// HasHandler returns true or false based on if c has the specified handler
func (c *Card) HasHandler(handler HandlerFunc, ctx *Context) bool {
	for _, condition := range append(c.effects, c.conditions...) {
		if reflect.ValueOf(condition).Pointer() == reflect.ValueOf(handler).Pointer() {
			return true
		}
	}
	return false
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

// RemoveEquipments ...
func (c *Card) RemoveEquipments() {
	for _, c = range c.Attachments() {
		if c.family == family.Equipment {
			if err := c.MoveCard(GRAVEYARD); err != nil {
				logrus.Debugf("RemoveEquipments: %s", err)
				return
			}
		}
	}
}

// AmIPlayed returns true or false based on if the card is played
func (c *Card) AmIPlayed(ctx *Context) bool {
	event, ok := ctx.event.(*PlayCardEvent)
	return ok && event.ID == c.id
}

// EvolveTo handles evolution
func (c *Card) EvolveTo(card *Card) {
	if err := card.MoveCard(BATTLEZONE); err != nil {
		logrus.Debugf("EvolveTo: %s", err)
		return
	}

	card.Tapped = c.Tapped

	c.AttachTo(card)

	// This is done to maintain a single identity for a creature
	c.id, card.id = card.id, c.id
	c.conditions, card.conditions = card.conditions, c.conditions

	c.player.match.Chat("Server", fmt.Sprintf("%s evolved %s to %s", c.player.Name(), c.name, card.name))
}

// SpellCast handles spell cast
func (c *Card) SpellCast(ctx *Context, fx func() []*Card) {
	for _, creature := range c.player.CollectCards(BATTLEZONE) {
		if creature.HasCivilisation(c.civ, ctx) &&
			c.GetRank(ctx) <= creature.GetRank(ctx) &&
			!creature.Tapped {
			ctx.ScheduleAfter(func() {
				ctx.match.SpellCast(c.id, fx())

				if err := c.MoveCard(GRAVEYARD); err != nil {
					logrus.Debug(err)
				}
				ctx.match.Chat("server", fmt.Sprintf("%s played %s", c.player.name, c.name))
			})
			return
		}
	}
	ctx.InterruptFlow()
}

// denormalizeCard returns the CardState of a card
func (c *Card) denormalizeCard() CardState {
	return CardState{
		ID:            c.id,
		UID:           c.cardID,
		Name:          c.name,
		Civ:           c.civ,
		Tapped:        c.Tapped,
		AttachedCards: denormalizeCards(c.Attachments()),
	}
}
