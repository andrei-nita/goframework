package auth

import (
	"encoding/json"
	"time"
)

type Session struct {
	Uuid      string
	Email     string
	CreatedAt time.Time
}

// Check if session is valid in the database
func (session *Session) Check() (bool, error) {
	return sessionsHM.Exists(session.Uuid)
}

// Delete session from database
func (session *Session) Delete() error {
	return sessionsHM.Del(session.Uuid)
}

// GetUser gets the user from the session
func (session *Session) GetUser() (user User, err error) {
	str, err := usersHM.Get(session.Email, "")
	err = json.Unmarshal([]byte(str), &user)
	return
}

// SessionDeleteAll deletes all sessions from database
func SessionDeleteAll() error {
	return sessionsHM.Clear()
}
