package api

import (
	"LinkingAPI/share_my_feed/database"
	"LinkingAPI/share_my_feed/database/databaseUtil"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type listHandler struct {
	ListsFunctions map[string]func(w http.ResponseWriter, r *http.Request)
}

func (handler *listHandler) initialize() {
	handler.ListsFunctions = map[string]func(w http.ResponseWriter, r *http.Request){}
	handler.ListsFunctions["lists/Make-list"] = makeList
	handler.ListsFunctions["lists/Get-lists"] = getLists
	fmt.Println(handler)
}

func listWrapper(w http.ResponseWriter, r *http.Request, endpoint func(http.ResponseWriter, *http.Request)) {
	endpoint(w, r)
}

//func makeList(w http.ResponseWriter, r *http.Request) {
//	bodyBytes, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//	}
//
//	var list database.LinksList
//	err = json.Unmarshal(bodyBytes, &list)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//	}
//
//	fmt.Println(list)
//	//var user thing
//	//err = json.Unmarshal(bodyBytes, &thi)
//	//fmt.Println(thi.DiscordUsername)
//}

func getLists(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("I got till here.")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var list database.LinksList
	err = json.Unmarshal(bodyBytes, &list)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	lists, err := database.GetLists(list.UserDataId)

	databaseUtil.PrettyPrint(lists)

	fmt.Println("WE got the lists")
	data, err := json.Marshal(lists)

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)

	//var user thing
	//err = json.Unmarshal(bodyBytes, &thi)
	//fmt.Println(thi.DiscordUsername)
}

func makeList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Making a new list")
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var list database.LinksList
	err = json.Unmarshal(bodyBytes, &list)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	databaseUtil.PrettyPrint(list)
}
