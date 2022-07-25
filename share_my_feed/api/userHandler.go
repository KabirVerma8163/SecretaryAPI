package api

import (
	"LinkingAPI/share_my_feed/database"
	"io/ioutil"
	"net/http"
)

type userHandler struct {
	UsersFunctions map[string]func(w http.ResponseWriter, r *http.Request)
}

func (handler *userHandler) initialize() {
	handler.UsersFunctions = map[string]func(w http.ResponseWriter, r *http.Request){}
	handler.UsersFunctions["/users/New-discord"] = getLists
}

func userWrapper(w http.ResponseWriter, r *http.Request, endpoint func(http.ResponseWriter, *http.Request)) {
	endpoint(w, r)
}

func newDiscordUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = database.NewDiscordUser(bodyBytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}
