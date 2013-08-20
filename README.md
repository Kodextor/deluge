# deluge


## API

### Create New User

    $ curl -X POST -d '{"username": "some_new_user"}' http://localhost:9090/users/new

Sample response:

    {"username": "some_new_user", "token": "0ce7d54d-33ac-4911-5105-9d6b6a9b8230"}

The returned token is what should be used by this user's Genesis
instance for future API calls.


### Associate User with New Subdomain

    $ curl -X POST -d \
    '{"username":"some_new_user","subdomain":"mysubdomain","token":"0ce7d54d-33ac-4911-5105-9d6b6a9b8230"}' \
    http://localhost:9090/subdomains/new



## Tips

### MongoDB

Type `mongo deluged` at the command line to enter the MongoDB shell


#### View All Users

    db.users.find()

#### Delete a User (TEMPORARY; will be added to API)

    db.users.remove({"username": "some_new_user"})



## TODO

* Create TCP proxy to route traffic to user-hosted askOS nodes

  * Use `types.DNSMapper` or similar

* Associate currently-recorded user-specific subdomains with domains

  * (Or don't, and have mysubdomain.*.* all route to the user's home machine)

* Add auth to new user creation

* Make it clear to users that a dot (`.`) cannot be included in their subdomain

  * This is for SSL cert-related reasons

* Decide whether to allow for multiple Genesis instances that are each routed to

* Add to TODO :-)
