package database

import (
	"encoding/json"
	guuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinkType struct {
	LinkID     guuid.UUID           `json:"link_id" bson:"link_id"`
	Name       string               `json:"link_name" bson:"link_name"`
	URL        string               `json:"link_url" bson:"link_url"`
	Type       string               `json:"type" bson:"type"`
	//Tags       map[string]bool      `json:"link_tags" bson:"link_tags"`
	Categories []string             `json:"categories" bson:"categories"`
	Read       bool                 `json:"read" bson:"read"`
	ReadUsers  []primitive.ObjectID `json:"read_users" bson:"read_users"`
	//LinkComments     []map[string], probably []LinkComments             `json:"link_comments" bson:"link_comments"`
}

func NewLink (linkData []byte) (link LinkType, err error) {
	err = json.Unmarshal(linkData, &link)
	if err != nil {
		return LinkType{}, err
	}

	/*
	URL must be tested

	LinkID must be changed
	read must be false
	read users must be empty
	*/

	return link, err
}