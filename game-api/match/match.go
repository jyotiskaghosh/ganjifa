package match

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// Match struct
type Match struct {
	player1 *PlayerReference
	player2 *PlayerReference

	mutex *sync.Mutex

	started bool

	winner *Player
	quit   chan bool
}

// New returns a new match object
func New() *Match {
	return &Match{
		mutex: &sync.Mutex{},
		quit:  make(chan bool),
	}
}

// playerRef returns the player ref for a given player
func (m *Match) playerRef(p *Player) (*PlayerReference, error) {
	if m.player1.Player == p {
		return m.player1, nil
	}

	if m.player2.Player == p {
		return m.player2, nil
	}

	return nil, errors.New("player is not a player of this match")
}

// PlayerForWriter returns the player ref for a given output or an error if the output is not p1 or p2
func (m *Match) PlayerForWriter(w Writer) (*PlayerReference, error) {
	if m.player1.Writer == w {
		return m.player1, nil
	}

	if m.player2.Writer == w {
		return m.player2, nil
	}

	return nil, errors.New("not a player of this match")
}

// Started if the game has started
func (m *Match) Started() bool {
	return m.started
}

// Quit if the game has ended
func (m *Match) Quit() <-chan bool {
	return m.quit
}

// Winner returns true or false based on if the playerref is the winner
func (m *Match) Winner(pr *PlayerReference) bool {
	return pr.Player == m.winner
}

// AddPlayer adds a new player
func (m *Match) AddPlayer(w Writer) (*PlayerReference, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var pr *PlayerReference

	switch {
	case m.player1 == nil:
		{
			m.player1 = &PlayerReference{
				Name:   "player1",
				Player: newPlayer(m, true),
				Writer: w,
			}
		}
	case m.player2 == nil:
		{
			m.player2 = &PlayerReference{
				Name:   "player2",
				Player: newPlayer(m, false),
				Writer: w,
			}
		}
	default:
		return nil, errors.New("players at max capacity")
	}

	return pr, nil
}

// Parse processes the data provided by player
func (m *Match) Parse(pr *PlayerReference, data []byte) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from parsing message in match. %v", r)
		}
	}()

	if pr.Player.wait {
		Warn(pr, "waiting for an action to resolve")
		return
	}

	m.player1.Player.waiting(true)
	m.player2.Player.waiting(true)
	defer func() {
		m.player1.Player.waiting(false)
		m.player2.Player.waiting(false)
	}()

	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return
	}

	switch msg.Header {
	case "choose_deck":
		{
			if m.started {
				Warn(pr, "match started cannot choose deck again")
				return
			}

			msg := CreateDeck{}
			if err := json.Unmarshal(data, &msg); err != nil {
				return
			}

			if err := pr.Player.createDeck(msg.Cards); err != nil {
				Warn(pr, err.Error())
				return
			}

			m.Chat("Server", fmt.Sprintf("%s has chosen their deck", pr.Name))

			if m.player1 != nil && m.player2 != nil && m.player1.Player.ready && m.player2.Player.ready {
				m.start()
			}
		}
	case "end_turn":
		{
			if pr.Player.turn {
				m.EndTurn()
			}
		}
	case "set_card":
		{
			if !pr.Player.turn {
				return
			}

			var msg struct {
				ID string `json:"id"`
			}
			if err := json.Unmarshal(data, &msg); err != nil {
				logrus.Debug(err)
				return
			}

			c, err := pr.Player.GetCard(msg.ID)
			if err != nil {
				logrus.Debug(err)
				return
			}

			if err := c.MoveCard(TRAPZONE); err != nil {
				logrus.Debug(err)
				return
			}
		}
	case "play_card":
		{
			if !pr.Player.turn {
				return
			}

			msg := PlayCardEvent{}
			if err := json.Unmarshal(data, &msg); err != nil {
				logrus.Debug(err)
				return
			}

			if ok := pr.Player.HasCard(HAND, msg.ID); ok {
				m.PlayCard(msg.ID, msg.TargetID)
			}
		}
	case "action":
		{
			var msg struct {
				Cards []string `json:"cards"`
			}
			if err := json.Unmarshal(data, &msg); err != nil {
				Warn(pr, "Invalid selection")
				return
			}

			m := make(map[string]bool)
			for _, card := range msg.Cards {
				m[card] = true
			}

			var cards []string
			for card := range m {
				cards = append(cards, card)
			}

			if pr.Player.action != nil {
				pr.Player.action <- cards
			}
		}
	case "cancel":
		{
			if pr.Player.cancel != nil {
				pr.Player.cancel <- true
			}
		}
	case "attack":
		{
			if !pr.Player.turn {
				return
			}

			if pr.Player.turnNo == 1 && pr == m.player1 {
				Warn(pr, "player 1 can't attack on first turn")
				return
			}

			msg := AttackEvent{}
			if err := json.Unmarshal(data, &msg); err != nil {
				logrus.Debug(err)
				return
			}

			m.Attack(pr, msg.ID, msg.TargetID)
		}
	default:
		logrus.Debugf("Received message in incorrect format: %v", string(data))
	}
}

// Opponent returns the opponent of the given player
func (m *Match) Opponent(p *Player) *Player {
	if m.player1.Player == p {
		return m.player2.Player
	}
	return m.player1.Player
}

// CurrentPlayer returns the turn player
func (m *Match) CurrentPlayer() *PlayerReference {
	if m.player1.Player.turn {
		return m.player1
	}
	return m.player2
}

// Chat sends a chat message
func (m *Match) Chat(sender string, message string) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from panic during sending chat message. %v", r)
			return
		}
	}()

	msg := ChatMessage{
		Header:  "chat",
		Message: message,
		Sender:  sender,
	}

	m.player1.Write(msg)
	m.player2.Write(msg)
}

// MessagePlayer sends a chat message from the server to the specified player
func (m *Match) MessagePlayer(p *Player, message string) {
	m.WritePlayer(p, ChatMessage{
		Header:  "chat",
		Message: message,
		Sender:  "server",
	})
}

// WritePlayer sends a message to a player
func (m *Match) WritePlayer(p *Player, msg interface{}) {
	pr, err := m.playerRef(p)
	if err != nil {
		logrus.Debug(fmt.Sprintf("error while writing to player: %s", err))
		return
	}

	pr.Write(msg)
}

// Warn sends a warning to the specified player ref
func Warn(p *PlayerReference, message string) {
	p.Write(ChatMessage{
		Header:  "warn",
		Message: message,
	})
}

// WarnPlayer sends a warning to the specified player
func (m *Match) WarnPlayer(p *Player, message string) {
	m.WritePlayer(p, ChatMessage{
		Header:  "warn",
		Message: message,
	})
}

// Battle handles a battle between two creatures
func (m *Match) Battle(attacker *Card, defender *Card, blocked bool) {
	if attacker.zone != BATTLEZONE || defender.zone != BATTLEZONE {
		return
	}

	ctx := NewContext(m, &Battle{Attacker: attacker, Defender: defender, Blocked: blocked})
	ctx.ScheduleAfter(func() {
		if attacker.GetAttack(ctx) > defender.GetDefence(ctx) {
			m.Destroy(defender, attacker)
		}
	})
	m.HandleFx(ctx)
}

// Destroy sends the given card to its players graveyard
func (m *Match) Destroy(card *Card, source *Card) {
	ctx := NewContext(m, &CreatureDestroyed{ID: card.id, Source: source})
	m.HandleFx(ctx)
}

// CollectCards ...
func (m *Match) CollectCards() []*Card {
	players := make([]*PlayerReference, 0)

	// The player in which turn it is is to be handled first
	if m.player1.Player.turn {
		players = append(players, m.player1, m.player2)
	} else {
		players = append(players, m.player2, m.player1)
	}

	cards := make([]*Card, 0)

	for _, p := range players {
		cards = append(cards, p.Player.battlezone...)
		cards = append(cards, p.Player.soul...)
		cards = append(cards, p.Player.trapzone...)
		cards = append(cards, p.Player.hand...)
		cards = append(cards, p.Player.graveyard...)
		cards = append(cards, p.Player.deck...)
	}

	return cards
}

// HandleFx ...
func (m *Match) HandleFx(ctx *Context) {
	defer m.BroadcastState()

	for _, c := range m.CollectCards() {
		for _, h := range c.GetHandlers(ctx) {
			h(c, ctx)
		}
	}

	for _, h := range ctx.postFxs {
		if ctx.cancel {
			return
		}

		h()
	}
}

// BroadcastState sends the current game's state to both players, hiding the opponent's hand
func (m *Match) BroadcastState() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from panic during sending state update. %v", r)
			return
		}
	}()

	player1 := m.player1.Player.denormalized()
	player2 := m.player2.Player.denormalized()

	p1state := StateMessage{
		Header: "state_update",
		State: State{
			MyTurn:   m.player1.Player.turn,
			Me:       player1,
			Opponent: player2,
		},
	}

	p2state := StateMessage{
		Header: "state_update",
		State: State{
			MyTurn:   m.player2.Player.turn,
			Me:       player2,
			Opponent: player1,
		},
	}

	p1state.State.Opponent.Hand = hideCards(len(p1state.State.Opponent.Hand))
	p1state.State.Opponent.Trapzone = hideCards(len(p1state.State.Opponent.Trapzone))

	p2state.State.Opponent.Hand = hideCards(len(p2state.State.Opponent.Hand))
	p2state.State.Opponent.Trapzone = hideCards(len(p2state.State.Opponent.Trapzone))

	m.player1.Write(p1state)
	m.player2.Write(p2state)
}

// End ends the match
func (m *Match) End(winner *Player, reason string) {
	logrus.Debugf("Attempting to end match")

	m.Chat("server", fmt.Sprintf("%s won the match, %s", winner.Name(), reason))

	m.winner = winner

	m.quit <- true
	close(m.quit)
}

// NewAction prompts the user to make a selection of the specified []Cards
func (m *Match) NewAction(p *Player, cards []*Card, minSelections int, maxSelections int, text string, cancellable bool) {
	m.WritePlayer(p, ActionMessage{
		Header:        "action",
		Cards:         denormalizeCards(cards),
		Text:          text,
		MinSelections: minSelections,
		MaxSelections: maxSelections,
		Cancellable:   cancellable,
	})
}

// CloseAction closes the card selection popup for the given player
func (m *Match) CloseAction(p *Player) {
	m.WritePlayer(p, Message{
		Header: "close_action",
	})
}

// ShowCards shows the specified cards to the player with a message of why it is being shown
func (m *Match) ShowCards(p *Player, message string, cards []string) {
	m.WritePlayer(p, ShowCardsMessage{
		Header:  "show_cards",
		Message: message,
		Cards:   cards,
	})
}

// Highlight highlights creatures
func (m *Match) Highlight(ids ...string) {
	msg := HighlightMessage{"highlight", append(make([]string, 0), ids...)}
	m.player1.Write(msg)
	m.player2.Write(msg)
}

// changeCurrentPlayer changes the current player
func (m *Match) changeCurrentPlayer() {
	m.player1.Player.turn = !m.player1.Player.turn
	m.player2.Player.turn = !m.player2.Player.turn
}

// start starts the match
func (m *Match) start() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.started = true

	m.player1.Player.ShuffleDeck()
	m.player2.Player.ShuffleDeck()

	m.player1.Player.DrawCards(5)
	m.player2.Player.DrawCards(5)

	m.Chat("Server", "The match has begun!")

	// This is done to offset beginNewTurn which changes current player
	m.changeCurrentPlayer()

	m.beginNewTurn()
}

// beginNewTurn starts a new turn
func (m *Match) beginNewTurn() {
	m.changeCurrentPlayer()
	m.CurrentPlayer().Player.turnNo++

	m.HandleFx(NewContext(m, &BeginTurnStep{}))

	m.untapStep()
}

// untapStep ...
func (m *Match) untapStep() {
	m.HandleFx(NewContext(m, &UntapStep{}))

	m.startOfTurnStep()
}

// startOfTurnStep ...
func (m *Match) startOfTurnStep() {
	m.HandleFx(NewContext(m, &StartOfTurnStep{}))

	m.drawStep()
}

// drawStep ...
func (m *Match) drawStep() {
	m.HandleFx(NewContext(m, &DrawStep{}))

	p := m.CurrentPlayer().Player
	p.DrawCards(1)
}

// endStep ...
func (m *Match) endStep() {
	m.HandleFx(NewContext(m, &EndStep{}))

	m.beginNewTurn()
}

// EndTurn is called when the player attempts to end their turn
// If the context is not cancelled by a card, the endStep is called
func (m *Match) EndTurn() {
	ctx := NewContext(m, &EndTurnEvent{})
	m.HandleFx(ctx)

	if !ctx.cancel {
		m.endStep()
	}
}

// PlayCard is called when the player attempts to play a card
func (m *Match) PlayCard(id string, targetID string) {
	m.Highlight(targetID)
	defer m.Highlight()

	ctx := NewContext(m, &PlayCardEvent{
		ID:       id,
		TargetID: targetID,
	})
	m.HandleFx(ctx)
}

// Attack is called when the player attempts to attack an opponent's creature
func (m *Match) Attack(pr *PlayerReference, id string, targetID string) {
	m.Highlight(id, targetID)
	defer m.Highlight()

	ctx := NewContext(m, &AttackEvent{
		ID:       id,
		TargetID: targetID,
	})
	m.HandleFx(ctx)

	// Tap card after attack
	// card gets tapped even if attack fails, this is done to prevent reusing before attacking effects
	card, err := m.GetCard(id)
	if err != nil {
		logrus.Debug(err)
	} else {
		card.Tapped = true
	}

	m.BroadcastState()
}

// GetCard returns the *Card from a container from either players
func (m *Match) GetCard(id string) (*Card, error) {
	c, err := m.player1.Player.GetCard(id)
	if err != nil {
		return m.player2.Player.GetCard(id)
	}

	return c, nil
}
