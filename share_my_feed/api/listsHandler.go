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

var listCommands []command
var getListsCommand command
var makeListCommand command

func (handler *listHandler) initialize() {
	handler.ListsFunctions = map[string]func(w http.ResponseWriter, r *http.Request){}

	getListsCommand = command{
		Address:         "lists/Get-lists",
		AccountSpecific: true,
		Testing:         false,
		Function:        getLists,
	}
	makeListCommand = command{
		Address:         "lists/Make-list",
		AccountSpecific: true,
		Testing:         false,
		Function:        makeList,
	}

	listCommands = []command{
		getListsCommand,
		makeListCommand,
	}

	for _, v := range listCommands {
		handler.ListsFunctions[v.Address] = v.Function
		allCommands["/"+v.Address] = v
	}

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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	fmt.Println("We got the lists.")
	databaseUtil.PrettyPrint(lists)
	//fmt.Println(lists)
	data2 := map[string]interface{}{
		"lists": lists}

	data, err := json.Marshal(data2)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	fmt.Println(string(data[:]))

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	//var user thing
	//err = json.Unmarshal(bodyBytes, &thi)
	//fmt.Println(thi.DiscordUsername)
}

func makeList(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var list database.LinksList
	err = json.Unmarshal(bodyBytes, &list)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	//databaseUtil.PrettyPrint(list)

	err = database.AddList(list, list.UserDataId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusAccepted)
}
