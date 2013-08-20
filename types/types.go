// Steve Phillips / elimisteve
// 2013.04.28

package types

import (
	"labix.org/v2/mgo"
)

var (
	db *mgo.Database
	users *mgo.Collection
	subdomains *mgo.Collection
)

func SetDB(mgoDB *mgo.Database) {
	db = mgoDB
	users = db.C("users")
	subdomains = db.C("subdomains")
}
