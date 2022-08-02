package api

import (
	"LinkingAPI/share_my_feed/database"
	"fmt"
	"io/ioutil"
	"net/http"
)

type userHandler struct {
	UsersFunctions map[string]func(w http.ResponseWriter, r *http.Request)
}

var userCommands []command
var newDiscordUserCommand command

func (handler *userHandler) initialize() {
	handler.UsersFunctions = map[string]func(w http.ResponseWriter, r *http.Request){}

	newDiscordUserCommand = command{
		Address:         "users/New-discord",
		AccountSpecific: false,
		Testing:         false,
		Function:        newDiscordUser,
	}

	userCommands = []command{
		newDiscordUserCommand,
	}

	for _, v := range userCommands {
		handler.UsersFunctions[v.Address] = v.Function
		allCommands["/"+v.Address] = v
	}
}

func userWrapper(w http.ResponseWriter, r *http.Request, endpoint func(http.ResponseWriter, *http.Request)) {
	endpoint(w, r)
}

func newDiscordUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got till here")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = database.NewDiscordUser(bodyBytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
