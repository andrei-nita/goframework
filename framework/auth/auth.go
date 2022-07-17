package auth

import (
	"github.com/xyproto/simplebolt"
	"os"
)

var (
	usersDB *simplebolt.Database
	usersHM *simplebolt.HashMap

	sessionsDB *simplebolt.Database
	sessionsHM *simplebolt.HashMap
)

func OpenAuth() (err error) {
	// create directory if not exists
	if _, err = os.Stat("db"); os.IsNotExist(err) {
		err = os.Mkdir("db", 0755)
		if err != nil {
			return err
		}
	}

	if err = openUsersDB(); err != nil {
		return err
	}
	if err = openSessionsDB(); err != nil {
		return err
	}

	usersHM, err = simplebolt.NewHashMap(usersDB, "users")
	if err != nil {
		return err
	}

	sessionsHM, err = simplebolt.NewHashMap(sessionsDB, "sessions")
	if err != nil {
		return err
	}

	return err
}

func CloseAuth() {
	closeUsersDB()
	closeSessionsDB()
}

func openUsersDB() (err error) {
	usersDB, err = simplebolt.New("db/users.db")
	if err != nil {
		return err
	}
	return err
}

func closeUsersDB() {
	usersDB.Close()
}

func openSessionsDB() (err error) {
	sessionsDB, err = simplebolt.New("db/sessions.db")
	if err != nil {
		return err
	}
	return err
}

func closeSessionsDB() {
	sessionsDB.Close()
}
