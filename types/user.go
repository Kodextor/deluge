// Steve Phillips / elimisteve
// 2013.04.28

package types

import (
	"labix.org/v2/mgo/bson"
	"fmt"
	"log"
	"time"
)

// TODO: Add Password, Email, etc later if this is going into
// production
type User struct {
	Username   string    `json:"username" bson:"username"`
	Token      string    `json:"token" bson:"token"`
	Subdomain  string    `json:"subdomain" bson:"subdomain"`
	CreatedAt  time.Time `json:"-" bson:"-"`
	ModifiedAt time.Time `json:"-" bson:"-"`
}

// String returns the user's username
func (user *User) String() string {
	return user.Username
}

// NewUser creates a new user with the given username and a fresh
// timestamp
func NewUser(username string) *User {
	now := time.Now()
	return &User{
		Username:   username,
		Token:      NewToken(),
		CreatedAt:  now,
		ModifiedAt: now,
	}
}

func QueryUser(username string) (*User, error) {
	u := User{}
	err := users.
		Find(bson.M{"username": username}).
		Sort("-createdat").
		One(&u)
	return &u, err
}

// func (user *User) GetAddresses(n int) (addrs []*net.TCPAddr, err error) {
// 	if user == nil {
// 		err = fmt.Errorf("Can't get subdomain for nil User")
// 		return
// 	}
// 	log.Printf("Trying to get %d messages from user %s\n", n, user.Username)
// 	subs = make([]Subdomain, n)
// 	err = subdomains.
// 		Find(bson.M{"user": user}).
// 		Sort("-createdat").
// 		Limit(n).
// 		All(&subs)

// 	return
// }

// Save inserts a new user into MongoDB
func (user *User) Save() error {
	if err := users.Insert(user); err != nil {
		return fmt.Errorf("Error creating new user: %v", err)
	}
	return nil
}

func (user *User) Update() error {
	info, err := users.Upsert(bson.M{"username": user.Username}, user)
	if err != nil {
		return fmt.Errorf("Error creating new user: %v", err)
	}
	log.Printf("info == %+v\n", info)
	return nil
}
