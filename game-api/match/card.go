package match

import (
	"reflect"
	"sync"

	"github.com/jyotiskaghosh/ganjifa/game-api/civ"

	"github.com/sirupsen/logrus"
	"github.com/ventu-io/go-shortid"
)

// Card stores card data
type Card struct {
	ID string

	cardID   int
	name     string
	rank     int
	civ      civ.Civilisation
	family   string
	attack   int
	defence  int
	handlers []HandlerFunc

	zone       Container
	tapped     bool
	attachedTo *Card
	conditions []HandlerFunc // Temporary handlers

	player *Player
	mutex  *sync.Mutex
}

// CardBuilder is a builder for Card
type CardBuilder struct {
	Name     string
	Rank     int
	Civ      civ.Civilisation
	Family   string
	Attack   int
	Defence  int
	Handlers []HandlerFunc
}

// Build constructs a card with the values of CardBuilder
func (cb *CardBuilder) Build() *Card {
	return &Card{
		name:     cb.Name,
		rank:     cb.Rank,
		civ:      cb.Civ,
		family:   cb.Family,
		attack:   cb.Attack,
		defence:  cb.Defence,
		handlers: cb.Handlers,
	}
}

// NewCard returns a new, initialized card
func NewCard(p *Player, cardID int) (*Card, error) {

	c, err := CardCtor(cardID)
	if err != nil {
		logrus.Warn(err)
		return nil, err
	}

	id, err := shortid.Generate()
	if err != nil {
		logrus.Debug(err)
		return nil, err
	}

	c.ID = id
	c.cardID = cardID
	c.player = p
	c.zone = DECK
	c.mutex = &sync.Mutex{}

	return c, nil
}

// Name ...
func (c *Card) Name() string {
	return c.name
}

// Rank ...
func (c *Card) Rank() int {
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
func (c *Card) Attack() int {
	return c.attack
}

// Defence ...
func (c *Card) Defence() int {
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

// Tap taps or untaps based on the bool value passed
func (c *Card) Tap(tap bool) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.tapped = tap
}

// Tapped returns if card is tapped
func (c *Card) Tapped() bool {
	return c.tapped
}

// Conditions returns all conditions
func (c *Card) Conditions() []HandlerFunc {
	return c.conditions
}

// AddCondition adds temporary handler functions
func (c *Card) AddCondition(handlers ...HandlerFunc) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.conditions = append(c.conditions, handlers...)
}

// HasCondition returns true or false based on if an handler is added to the cards list of conditions
func (c *Card) HasCondition(handler HandlerFunc) bool {

	for _, condition := range c.conditions {
		if reflect.ValueOf(condition).Pointer() == reflect.ValueOf(handler).Pointer() {
			return true
		}
	}

	return false
}

// RemoveCondition removes all instances of the given handler from the cards conditions
func (c *Card) RemoveCondition(handler HandlerFunc) {

	c.mutex.Lock()
	defer c.mutex.Unlock()

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

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.conditions = make([]HandlerFunc, 0)
}

// AttachedTo returns the card that this card is attached to
func (c *Card) AttachedTo() *Card {
	return c.attachedTo
}

// Attachments returns a copy of the card's attached cards
func (c *Card) Attachments() []*Card {

	result := make([]*Card, 0)

	cards, err := c.player.Container(SOUL)
	if err != nil {
		logrus.Debug(err)
		return result
	}

	for _, card := range cards {
		if card.attachedTo == c {
			result = append(result, card)
		}
	}

	return result
}

// AttachTo attaches c to card
func (c *Card) AttachTo(card *Card) {

	if card.zone != BATTLEZONE || c == card {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.attachedTo = card
	if err := c.MoveCard(SOUL); err != nil {
		logrus.Debug(err)
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
		logrus.Debug(err)
		return err
	}

	to, err := c.player.containerRef(destination)
	if err != nil {
		logrus.Debug(err)
		return err
	}

	c.player.mutex.Lock()
	defer c.player.mutex.Unlock()

	temp := make([]*Card, 0)

	for _, card := range *from {
		if card.ID != c.ID {
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

	c.player.match.HandleFx(NewContext(c.player.match, &CardMoved{
		ID:   c.ID,
		From: f,
		To:   destination,
	}))

	for _, card := range c.Attachments() {

		card.attachedTo = c.attachedTo

		card.player.match.HandleFx(NewContext(card.player.match, &CardMoved{
			ID:   card.ID,
			From: card.zone,
			To:   destination,
		}))
	}

	return nil
}

// GetRank returns the rank of a given card
func (c *Card) GetRank(ctx *Context) int {

	e := &GetRankEvent{
		Card:  c,
		Event: ctx.Event,
		Rank:  c.rank,
	}
	ctx.Match.HandleFx(NewContext(ctx.Match, e))

	return e.Rank
}

// GetCivilisation returns the Civ of a given card
func (c *Card) GetCivilisation(ctx *Context) map[civ.Civilisation]bool {

	e := &GetCivilisationEvent{
		Card:  c,
		Event: ctx.Event,
		Civ:   map[civ.Civilisation]bool{c.civ: true},
	}
	ctx.Match.HandleFx(NewContext(ctx.Match, e))

	return e.Civ
}

// GetFamily returns the family of a given card
func (c *Card) GetFamily(ctx *Context) map[string]bool {

	e := &GetFamilyEvent{
		Card:   c,
		Event:  ctx.Event,
		Family: map[string]bool{c.family: true},
	}
	ctx.Match.HandleFx(NewContext(ctx.Match, e))

	return e.Family
}

// GetAttack returns the attack of a given card
func (c *Card) GetAttack(ctx *Context) int {

	e := &GetAttackEvent{
		Card:   c,
		Event:  ctx.Event,
		Attack: c.attack,
	}
	ctx.Match.HandleFx(NewContext(ctx.Match, e))

	return e.Attack
}

// GetDefence returns the defence of a given card
func (c *Card) GetDefence(ctx *Context) int {

	e := &GetDefenceEvent{
		Card:    c,
		Event:   ctx.Event,
		Defence: c.defence,
	}
	ctx.Match.HandleFx(NewContext(ctx.Match, e))

	return e.Defence
}
