package auth

import (
	"regexp"
	"strings"
)

var (
	rxEmail = regexp.MustCompile("^[\\w-.]+@([\\w-]+\\.)+[\\w-]{2,4}$")
)

func (user *User) Validate() (hasErr bool, errorsSignup map[string]string) {
	errorsSignup = make(map[string]string)

	if strings.TrimSpace(user.Name) == "" {
		errorsSignup["name"] = "Please enter a name"
	}

	if strings.TrimSpace(user.Email) == "" {
		errorsSignup["email"] = "Please enter an email"
	} else {
		if match := rxEmail.Match([]byte(user.Email)); match == false {
			errorsSignup["email"] = "Please enter a valid email address"
		}
	}

	if len(user.Password) < 8 {
		errorsSignup["password"] = "Password must be at least 8 characters"
	}

	if len(errorsSignup) > 0 {
		return true, errorsSignup
	} else {
		return false, nil
	}
}

func GetMapValue(errorsMap map[string]string, key string) string {
	if value, ok := errorsMap[key]; ok {
		return value
	}
	return ""
}
