// Steve Phillips / elimisteve
// 2013.04.28

package types

import (
	"fmt"
	"log"
	"time"
)

var (
	// TODO: Replace with Postgres
	users = map[string]*User{}
)

// TODO: Add Password, Email, etc later if this is going into
// production
type User struct {
	Username   string    `json:"username"`
	Token      string    `json:"token"`
	Subdomain  string    `json:"subdomain"`
	CreatedAt  time.Time `json:"-"`
	ModifiedAt time.Time `json:"-"`
}

// String returns the user's username
func (user *User) String() string {
	return user.Username
}

// NewUser creates a new user with the given username and a fresh
// timestamp
func NewUser(username string) *User {
	log.Printf("Creating new user '%v'\n", username)
	now := time.Now()
	return &User{
		Username:   username,
		Token:      NewToken(),
		CreatedAt:  now,
		ModifiedAt: now,
	}
}

func QueryUser(username string) (*User, error) {
	u, ok := users[username]
	if !ok {
		return nil, fmt.Errorf("User not found")
	}
	log.Printf("User successfully queried: '%s'\n", username)
	return u, nil
}

func (user *User) Save() error {
	if user == nil {
		return fmt.Errorf("Can't save nil user!")
	}
	// TODO: Make thread-safe
	// TODO: Make sure user doesn't already exist
	log.Printf("Saving user '%s'\n", user.Username)
	users[user.Username] = user
	return nil
}

func (user *User) Update() error {
	// TODO: Truly update instead of saving/replacing
	return user.Save()
}
