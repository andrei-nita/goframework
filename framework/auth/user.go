package auth

import (
	"encoding/json"
	"time"
)

type User struct {
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

// CreateSession creates a new session for an existing user
func (user *User) CreateSession() (session *Session, err error) {
	session = &Session{}
	session.Email = user.Email
	session.Uuid = user.Uuid
	session.CreatedAt = time.Now()
	bytes, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}
	err = sessionsHM.Set(user.Uuid, "", string(bytes))
	return
}

// GetSession gets the session for an existing user
func (user *User) GetSession() (session *Session, err error) {
	str, err := sessionsHM.Get(user.Uuid, "")
	err = json.Unmarshal([]byte(str), &session)
	return
}

// Create a new user, save user info into the database
func (user *User) Create() error {
	user.Uuid = createUUID()
	user.Password = HashPsw(user.Password)
	user.CreatedAt = time.Now()
	bytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return usersHM.Set(user.Email, "", string(bytes))
}

// Exists check if the user is present in the database
func (user *User) Exists() (bool, error) {
	return usersHM.Exists(user.Email)
}

// Delete user from database
func (user *User) Delete() error {
	return usersHM.Del(user.Email)
}

// Update user information in the database
func (user *User) Update() error {
	return user.Create()
}

// UserDeleteAll deletes all users from database
func UserDeleteAll() error {
	return usersHM.Clear()
}

// Users get all users in the database and returns it
func Users() (users []User, err error) {
	ids, err := usersHM.All()
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		bytes, err := usersHM.Get(id, "")
		if err != nil {
			return nil, err
		}
		var user User
		err = json.Unmarshal([]byte(bytes), &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return
}

// UserByEmail get a single user given the email
func UserByEmail(email string) (user *User, err error) {
	bytes, err := usersHM.Get(email, "")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(bytes), &user)
	return
}
