package database

import (
	"github.com/kamva/mgm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinksList struct {
	mgm.DefaultModel  `bson:",inline"`
	OwnerId           primitive.ObjectID              `json:"owner_id" bson:"owner_id"`
	OwnerUsername     string                          `json:"owner_username" bson:"owner_username"`
	OwnerName         string                          `json:"owner_name" bson:"owner_name"`
	OwnerDataId       primitive.ObjectID              `json:"owner_data_id" bson:"owner_data_id"`
	AccessIds         map[primitive.ObjectID][]string `json:"access_ids" bson:"access_ids"` // Must be userDataIDs
	ListName          string                          `json:"list_name" bson:"list_name"`
	Links             []LinkType                      `json:"links" bson:"links"`
	Categories        []string                        `json:"categories" bson:"categories"`
	LinkAccessibility string                          `json:"link_accessibility" bson:"link_accessibility"`
}
