package api

import (
	"LinkingAPI/share_my_feed/database"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
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

var allCommands map[string]command

func Start() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	DiscordPassword = os.Getenv("DISCORD_PASSWORD")

	allCommands = map[string]command{}

	ListHandler := listHandler{}
	ListHandler.initialize()
	for k, v := range ListHandler.ListsFunctions {
		http.Handle("/"+k, UniversalWrapper(listWrapper, v))
	}

	UserHandler := userHandler{}
	UserHandler.initialize()
	for k, v := range UserHandler.UsersFunctions {
		http.Handle("/"+k, UniversalWrapper(userWrapper, v))
	}

	http.HandleFunc("/test", func(writer http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Body)
	})

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

						var unwrapper UnWrapper
						err := unwrapper.initialize(r)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}

						command, ok := allCommands[r.URL.Path]
						if !ok {
							w.WriteHeader(http.StatusNotFound)
							return
						}

						userDataID, errUserData := database.GetUserDataIDWithDiscordID(unwrapper.DiscordID)
						if command.Address == "users/New-discord" {
							if errUserData == nil {
								unwrapper.Error = "RequestError: User already has an account"
								data, err := dataMarshall(unwrapper, r)
								if err != nil {
									w.WriteHeader(http.StatusBadRequest)
									return
								}

								r.Body = io.NopCloser(bytes.NewReader(data))
								w.WriteHeader(http.StatusNotFound)
								w.Header().Set("Content-Type", "application/json")
								_, err = w.Write(data)
								if err != nil {
									w.WriteHeader(http.StatusBadRequest)
									return
								}
								return
							}
						} else if command.AccountSpecific {
							if errUserData != nil {
								if err.Error() == "ServerError: UserDataType for given user does not exist" {
									unwrapper.Error = "ServerError: UserDataType for given user does not exist"
									data, err := dataMarshall(unwrapper, r)
									if err != nil {
										w.WriteHeader(http.StatusBadRequest)
										return
									}

									r.Body = io.NopCloser(bytes.NewReader(data))
									w.WriteHeader(http.StatusNotFound)
									w.Header().Set("Content-Type", "application/json")

									_, err = w.Write(data)
									if err != nil {
										w.WriteHeader(http.StatusBadRequest)
										return
									}
									return
								} else {
									w.WriteHeader(http.StatusNotFound)
									return
								}
							}
							unwrapper.UserDataID = userDataID
						} else {

						}

						data, err := dataMarshall(unwrapper, r)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}
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

func dataMarshall(unwrapper UnWrapper, r *http.Request) ([]byte, error) {
	data, err := json.Marshal(unwrapper)
	if err != nil {
		return []byte{}, err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []byte{}, err
	}

	body = append(body, data...)
	//fmt.Println(string(body))
	return data, nil
}
