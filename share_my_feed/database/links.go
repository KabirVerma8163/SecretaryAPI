package database

import (
	"LinkingAPI/share_my_feed/database/databaseUtil"
	"encoding/json"
	guuid "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinkType struct {
	LinkID     guuid.UUID           `json:"link_id" bson:"link_id"`
	Name       string               `json:"link_name" bson:"link_name"`
	URL        string               `json:"link_url" bson:"link_url"`
	Type       string               `json:"type" bson:"type"`
	Tags       map[string]bool      `json:"link_tags" bson:"link_tags"`
	Categories []string             `json:"categories" bson:"categories"`
	Read       bool                 `json:"read" bson:"read"`
	ReadUsers  []primitive.ObjectID `json:"read_users" bson:"read_users"`
	//LinkComments     []map[string], probably []LinkComments             `json:"link_comments" bson:"link_comments"`
}

func (link *LinkType) initialize() (err error) {
	link.LinkID = guuid.New()

	link.URL, err = databaseUtil.CheckUrl(link.URL)
	if err != nil {
		return err
	}

	if link.Name == "" {
		link.Name = link.URL
	}

	// TODO: Make sure you make this work properly
	if link.Tags == nil {
		link.Tags = map[string]bool{}
	}

	if link.Categories == nil {
		link.Categories = []string{}
	}

	link.Read = false

	link.ReadUsers = []primitive.ObjectID{}

	return nil
}

func newLink(linkData []byte) (link LinkType, err error) {
	err = json.Unmarshal(linkData, &link)
	if err != nil {
		return LinkType{}, err
	}

	err = link.initialize()
	if err != nil {
		return LinkType{}, err
	}

	return link, nil
}

func newLinks(linksData []byte) (links []LinkType, err error) {
	err = json.Unmarshal(linksData, &links)
	if err != nil {
		return []LinkType{}, nil
	}

	for _, link := range links {
		err = link.initialize()
		if err != nil {
			return []LinkType{}, err
		}
	}
	return links, nil
}

//func AddLink(linkData []byte, listID primitive.ObjectID, userDataID primitive.ObjectID) (err error) {
//
//	list, err = addLinkToList(listID, userDataID, linkData)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

//func AddLinkAsDiscordUser(linkData []byte, listID primitive.ObjectID, discordUsername string) (err error) {
//	userDataID, err := getUserDataIDWithDiscordID(discordUsername)
//	if err != nil {
//		return err
//	}
//
//	err = addLinkToList(listID, userDataID, linkData)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

//
//	/*
//	URL must be tested
//
//	LinkID must be changed
//	read must be false
//	read users must be empty
//	*/
//
//	return link, err
