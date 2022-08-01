package api

import (
	"LinkingAPI/share_my_feed/database"
	"encoding/json"
	"fmt"
	guuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
)

type UnWrapper struct {
	Error string `json:"error" bson:"error"`

	DiscordID            string             `json:"discord_id"`
	DiscordUsername      string             `json:"discord_username" bson:"discord_username"`
	DiscordDiscriminator int                `json:"discord_discriminator" bson:"discord_discriminator"`
	Username             string             `json:"username"`
	UserDataID           primitive.ObjectID `json:"user_data_id" json:"owner_data_id"`

	ListName          string              `json:"list_name"`
	Links             []database.LinkType `json:"links" bson:"links"`
	ListCategories    []string            `json:"list_categories" bson:"list_categories"`
	LinkAccessibility string              `json:"link_accessibility" bson:"link_accessibility"`
	Description       string              `json:"description" bson:"description"`

	LinkID     guuid.UUID           `json:"link_id" bson:"link_id"`
	Name       string               `json:"link_name" bson:"link_name"`
	URL        string               `json:"link_url" bson:"link_url"`
	Type       string               `json:"type" bson:"type"`
	Tags       map[string]bool      `json:"link_tags" bson:"link_tags"`
	Categories []string             `json:"categories" bson:"categories"`
	Read       bool                 `json:"read" bson:"read"`
	ReadUsers  []primitive.ObjectID `json:"read_users" bson:"read_users"`
}

type command struct {
	Address         string
	AccountSpecific bool
	Testing         bool
	Function        func(http.ResponseWriter, *http.Request)
}

func (unwrapper *UnWrapper) initialize(dataBytes []byte) error {
	err := json.Unmarshal(dataBytes, &unwrapper)

	unwrapper.Error = ""
	if err != nil {
		return err
	}
	return nil
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
