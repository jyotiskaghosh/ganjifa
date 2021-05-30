package match

import "github.com/jyotiskaghosh/ganjifa/game-api/civ"

// Message is the default message struct
type Message struct {
	Header string `json:"header"`
}

// ChatMessage stores information about a chat message
type ChatMessage struct {
	Header  string `json:"header"`
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

// CreateDeck sends an array of card id's for creating a deck
type CreateDeck struct {
	Cards []int `json:"cards"`
}

// CardState stores information about the state of a card
type CardState struct {
	ID            string           `json:"id"`
	UID           int              `json:"uid"`
	Name          string           `json:"name"`
	Civ           civ.Civilisation `json:"civilization"`
	Tapped        bool             `json:"tapped"`
	FaceDown      bool             `json:"faceDown"`
	AttachedCards []CardState      `json:"attachedCards"`
}

// PlayerState stores information about the state of the current player
type PlayerState struct {
	Life       int         `json:"life"`
	Deck       int         `json:"deck"`
	Hand       []CardState `json:"hand"`
	Graveyard  []CardState `json:"graveyard"`
	Battlezone []CardState `json:"battlezone"`
	Trapzone   []CardState `json:"trapzone"`
}

// State stores information about the current state of the match in the eyes of a given player
type State struct {
	MyTurn   bool        `json:"myTurn"`
	Me       PlayerState `json:"me"`
	Opponent PlayerState `json:"opponent"`
}

// StateMessage is the message that should be sent to the client for state updates
type StateMessage struct {
	Header string `json:"header"`
	State  State  `json:"state"`
}

// ActionMessage is used to prompt the user to make a selection of the specified cards
type ActionMessage struct {
	Header        string      `json:"header"`
	Cards         []CardState `json:"cards"`
	Text          string      `json:"text"`
	MinSelections int         `json:"minSelections"`
	MaxSelections int         `json:"maxSelections"`
	Cancellable   bool        `json:"cancellable"`
}

// ShowCardsMessage is used to show the user n cards without an action to perform
type ShowCardsMessage struct {
	Header  string   `json:"header"`
	Message string   `json:"message"`
	Cards   []string `json:"cards"`
}
