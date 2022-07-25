package api

import (
	"LinkingAPI/share_my_feed/database"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var DiscordPassword string

/*
Structure
universal wrapper checks the request and makes sure it is valid for each kind of function.
<Type> = List, User, Link
<Type>Functions is the map of functions and their addresses and is added to each individual file.
*/

func Start() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	DiscordPassword = os.Getenv("DISCORD_PASSWORD")

	//http.HandleFunc("/test", test)
	//
	//http.Handle("/list", listWrapper(CreateList))
	//http.Handle("/test", UniversalWrapper(testWrapper, test))

	ListHandler := listHandler{}
	ListHandler.initialize()
	for k, v := range ListHandler.ListsFunctions {
		http.Handle(k, UniversalWrapper(listWrapper, v))
	}

	UserHandler := userHandler{}
	UserHandler.initialize()
	for k, v := range UserHandler.UsersFunctions {
		http.Handle(k, UniversalWrapper(listWrapper, v))
	}

	log.Fatalln(http.ListenAndServe(":8000", nil))

	//databaseUtil.PrettyPrint()

	return nil
}

func UniversalWrapper(wrapper func(http.ResponseWriter, *http.Request, func(http.ResponseWriter, *http.Request)), endpoint func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if typeArr, ok := r.Header["Client-Type"]; ok {
			switch typeArr[0] {
			case "Discord":
				{
					if passArr, ok := r.Header["Client-Password"]; ok {
						if passArr[0] != DiscordPassword {
							w.WriteHeader(http.StatusBadRequest)
							return
						}

						dataBytes, err := ioutil.ReadAll(r.Body)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}

						var unwrapper UnWrapper
						err = json.Unmarshal(dataBytes, &unwrapper)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}

						userDataID, err := database.GetUserDataIDWithDiscordID(unwrapper.DiscordID)
						if err != nil {
							if err.Error() == "ServerError: userDataType for given user does not exist" {
								// What happens if the user account does not exist.
								// TODO
								//responseString := []byte(`{}`)
								//_, err = w.Write(responseString)
								//if err != nil {
								//	w.WriteHeader(http.StatusInternalServerError)
								//}
							}
							w.WriteHeader(http.StatusNotFound)
							return
						}
						unwrapper.UserDataID = userDataID

						data, err := json.Marshal(unwrapper)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						body, err := ioutil.ReadAll(r.Body)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						body = append(body, data...)
						r.Body = io.NopCloser(bytes.NewReader(data))

						wrapper(w, r, endpoint)
					}
				}
			case "WebApp":
				{
				}
			default:
				{
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}

type UnWrapper struct {
	DiscordID  string             `json:"discord_id"`
	Username   string             `json:"username"`
	UserDataID primitive.ObjectID `json:"user_data_id" json:"owner_data_id"`
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Ryan sent a request")

	body := r.Body
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(bodyBytes)
}

func testWrapper(w http.ResponseWriter, r *http.Request, endpoint func(http.ResponseWriter, *http.Request)) {
	endpoint(w, r)
}
