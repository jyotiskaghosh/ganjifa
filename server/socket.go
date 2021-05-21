package server

import (
	"sync"
	"time"

	"github.com/jyotiskaghosh/ganjifa/db"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var sockets = make(map[*Socket]Hub)
var socketsMutex = sync.Mutex{}

// Sockets returns a list of the current sockets
func Sockets() []*Socket {
	result := make([]*Socket, 0)
	socketsMutex.Lock()
	defer socketsMutex.Unlock()
	for s := range sockets {
		result = append(result, s)
	}
	return result
}

// SocketsInHub ...
func SocketsInHub(hubID string) []*Socket {
	result := make([]*Socket, 0)
	socketsMutex.Lock()
	defer socketsMutex.Unlock()
	for s, hub := range sockets {
		if hub.Name() == hubID {
			result = append(result, s)
		}
	}
	return result
}

// Socket links a ws connection to a user id and handles safe reading and writing of data
type Socket struct {
	conn   *websocket.Conn
	User   db.User
	hub    Hub
	ready  bool
	mutex  *sync.Mutex
	closed bool
	lost   bool
}

// NewSocket creates and returns a new Socket instance
func NewSocket(c *websocket.Conn, hub Hub) *Socket {

	s := &Socket{
		conn:   c,
		hub:    hub,
		ready:  false,
		mutex:  &sync.Mutex{},
		closed: false,
		lost:   false,
	}

	socketsMutex.Lock()
	sockets[s] = hub
	socketsMutex.Unlock()

	logrus.Debugf("Opened a connection")

	return s

}

// Ready returns true or false based on if the socket is ready or not
func (s *Socket) Ready() bool {
	return s.ready
}

// Listen sets up reader and writer for the socket
func (s *Socket) Listen() {

	s.conn.SetReadLimit(maxMessageSize)
	s.conn.SetReadDeadline(time.Now().Add(pongWait))
	s.conn.SetPongHandler(func(string) error { s.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	defer s.Close()

	go s.handlePing()

	for {

		_, message, err := s.conn.ReadMessage()

		if err != nil {
			return
		}

		if !s.ready {

			// Look for authorization token as the first message
			u, err := db.GetUserForToken(string(message))

			if err != nil {
				continue
			}

			s.User = u
			s.ready = true

			s.Write(Message{Header: "hello"})

			continue

		}

		go s.hub.Parse(s, message)

	}

}

func (s *Socket) handlePing() {

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("recovered from handlePing: %v", r)
		}
	}()

	for {

		if s.closed || s.lost {
			return
		}

		select {
		case <-ticker.C:
			s.mutex.Lock()
			s.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := s.conn.WriteMessage(websocket.PingMessage, nil)
			s.mutex.Unlock()
			if err != nil {
				if !s.closed && !s.lost {
					s.conn.Close()
				}
				return
			}
		}
	}

}

// Write sends a struct v to the client
func (s *Socket) Write(v interface{}) {

	if s.closed || s.lost {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from panic in socket Write. %v", r)
			return
		}
	}()

	s.mutex.Lock()
	s.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err := s.conn.WriteJSON(v); err != nil {
		logrus.Debug(err)
	}
	s.mutex.Unlock()

}

// Close closes the client connection
func (s *Socket) Close() {

	defer func() {
		if r := recover(); r != nil {
			logrus.Warnf("Recovered from socket close. %v", r)
			return
		}
	}()

	if s.closed {
		return
	}

	s.closed = true

	socketsMutex.Lock()

	delete(sockets, s)

	socketsMutex.Unlock()

	s.hub.OnSocketClose(s)

	if s.conn != nil {
		s.conn.Close()
	}

	logrus.Debug("Closed a connection")

}

// GetUserList returns a list of users currently online
func GetUserList() UserListMessage {

	usersMap := make(map[string]UserMessage)

	socketsMutex.Lock()
	defer socketsMutex.Unlock()

	for s, h := range sockets {

		userEntry := UserMessage{
			Username:    s.User.Username,
			Color:       s.User.Color,
			Hub:         h.Name(),
			Permissions: s.User.Permissions,
		}

		if _, ok := usersMap[s.User.Username]; ok {

			// Replace if this socket is in a match because the client shows
			// an icon for if the player is in a match or just the lobby
			if userEntry.Hub == "match" {
				usersMap[s.User.Username] = userEntry
			}

		} else {
			usersMap[s.User.Username] = userEntry
		}

	}

	users := make([]UserMessage, 0)

	for _, user := range usersMap {
		users = append(users, user)
	}

	return UserListMessage{
		Header: "users",
		Users:  users,
	}

}
