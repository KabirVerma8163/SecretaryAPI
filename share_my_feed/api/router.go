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

						dataBytes, err := ioutil.ReadAll(r.Body)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}

						var unwrapper UnWrapper
						err = unwrapper.initialize(dataBytes)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							return
						}

						if r.URL.Path == "/users/New-discord" {
							_, err := database.GetUserDataIDWithDiscordID(unwrapper.DiscordID)
							if err == nil {
								unwrapper.Error = "RequestError: User already has an account"
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
								//fmt.Println(string(body))
								r.Body = io.NopCloser(bytes.NewReader(data))
								w.WriteHeader(http.StatusNotFound)
								w.Header().Set("Content-Type", "application/json")
								_, err = w.Write(data)
								if err != nil {
									w.WriteHeader(http.StatusBadRequest)
								}
								return
							}
						}

						userDataID, err := database.GetUserDataIDWithDiscordID(unwrapper.DiscordID)
						if err != nil {
							if err.Error() == "ServerError: UserDataType for given user does not exist" && r.URL.Path != "/users/New-discord" {
								unwrapper.Error = "ServerError: UserDataType for given user does not exist"
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
								//fmt.Println(string(body))
								r.Body = io.NopCloser(bytes.NewReader(data))
								w.WriteHeader(http.StatusNotFound)
								w.Header().Set("Content-Type", "application/json")
								_, err = w.Write(data)
								if err != nil {
									w.WriteHeader(http.StatusBadRequest)
								}
								return
							} else if r.URL.Path != "/users/New-discord" {
								w.WriteHeader(http.StatusNotFound)
								return
							}

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
