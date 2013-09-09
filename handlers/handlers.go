// Steve Phillips / elimisteve
// 2013.04.28

package handlers

import (
	"../types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to deluged!")
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Figure out the auth dance involving PayPalCallback

	// TODO: Read as a stream up to a certain number of bytes to
	// prevent huge-upload-induced DDOS attacks
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer r.Body.Close()

	// Unmarshal JSON body to User
	tempUser := types.User{}
	if err := json.Unmarshal(body, &tempUser); err != nil {
		e := fmt.Errorf("Error parsing JSON: %v", err)
		http.Error(w, e.Error(), 400)
		return
	}

	// Create and save new User in DB. Ignore other JSON fields from
	// user
	u := types.NewUser(tempUser.Username)
	if err = u.Save(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Return new User to user as JSON. Re-use `body` for efficiency
	if err = json.Unmarshal(body, u); err != nil {
		log.Printf("PostUser: Couldn't unmarshal '%s': %v\n", body, err)
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"username": "%s", "token": "%s"}`,
		u.Username, u.Token)
}

func PostSubdomain(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	// Unmarshal JSON body to map
	m := map[string]string{}
	if err := json.Unmarshal(body, &m); err != nil {
		e := fmt.Errorf("Error parsing JSON: %v", err)
		http.Error(w, e.Error(), 500)
		return
	}

	username, okU := m["username"]
	token, okT := m["token"]
	subdomain, okS := m["subdomain"]

	if !okU || !okT || !okS || username == "" || token == "" || subdomain == "" {
		str := "'username', 'token', and 'subdomain' must be provided"
		log.Printf("PostSubdomain: %s -- %v", str, err)
		http.Error(w, str, 400)
		return
	}

	user, err := types.QueryUser(username)
	if err != nil {
		log.Printf("PostSubdomain: QueryUser: %v\n", err)
		http.Error(w, "Error finding user", 500)
		return
	}

	if user.Token != token {
		log.Printf("PostSubdomain: Invalid token\n")
		http.Error(w, "Invalid token", 400)
		return
	}

	// TODO: Ensure that the desired subdomain isn't banned

	// Nothing else can go wrong... right???

	// Update subdomain
	user.Subdomain = subdomain
	if err = user.Update(); err != nil {
		log.Printf("Error after user.Update() -- %v\n", err)
		http.Error(w, "Error updating subdomain", 500)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"username": "%s", "subdomain": "%s"}`, username, subdomain)
}

// func PayPalCallback(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("NOT IMPLEMENTED\n"))
// }
