package database

import (
	"LinkingAPI/database/databaseUtil"
	"encoding/json"
	"fmt"
	"github.com/kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinksList struct {
	mgm.DefaultModel `bson:",inline"`
	//OwnerId           primitive.ObjectID              `json:"owner_id" bson:"owner_id"`
	OwnerUsername     string                        `json:"owner_username" bson:"owner_username"`
	OwnerName         string                        `json:"owner_name" bson:"owner_name"`
	OwnerDataId       primitive.ObjectID            `json:"owner_data_id" bson:"owner_data_id"`
	AccessIds         map[primitive.ObjectID]string `json:"access_ids" bson:"access_ids"` // Must be userDataIDs
	ListName          string                        `json:"list_name" bson:"list_name"`
	Links             []LinkType                    `json:"links" bson:"links"`
	Categories        []string                      `json:"categories" bson:"categories"`
	LinkAccessibility string                        `json:"link_accessibility" bson:"link_accessibility"`
}

// TODO be mindful of what the user can access, like ownerdataid

func (list *LinksList) initialize(username string) error {
	userData, err := getUserDataWithUsername(username)
	if err != nil {
		return err
	}

	list.OwnerUsername = username
	list.OwnerName = userData.Name
	list.OwnerDataId = userData.ID

	if list.AccessIds == nil {
		list.AccessIds = map[primitive.ObjectID]string{}
	}

	// TODO: fix this thing
	if list.Links == nil {
		list.Links = []LinkType{}
	} else {
		//list.Links, err = newLinks(list.Links)
	}

	if list.Categories == nil {
		list.Categories = []string{}
	}
	//list.LinkAccessibility = "private"
	return nil
}

func AddListAsUser(listData []byte, username string) (err error) {
	var list LinksList
	err = json.Unmarshal(listData, &list)
	if err != nil {
		return err
	}

	err = list.initialize(username)
	if err != nil {
		return err
	}

	listID, err := addList(list)

	_, err = databaseUtil.UserDataColl.UpdateOne(databaseUtil.Ctx,
		bson.M{"username": username}, bson.D{
			{"$push", bson.D{{"links_list", listID}}},
		})
	if err != nil {
		return err
	}

	return err
}

func AddListAsDiscordUser(listData []byte, userID string) (err error) {
	var list LinksList
	err = json.Unmarshal(listData, &list)
	if err != nil {
		return err
	}

	discord, err := getDiscordDataIDFromDiscordID(userID)

	err = list.initialize(discord.Username)
	if err != nil {
		return err
	}

	listID, err := addList(list)

	_, err = databaseUtil.UserDataColl.UpdateOne(databaseUtil.Ctx,
		bson.M{"_id": discord.UserDataID},
		bson.M{"$set": bson.M{fmt.Sprintf("links_list.%v", list.ListName): listID}},
	)
	if err != nil {
		return err
	}

	return err
}

func addList(list LinksList) (linkId primitive.ObjectID, err error) {
	listInsertResult, err := databaseUtil.LinksListsColl.InsertOne(databaseUtil.Ctx, list)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return listInsertResult.InsertedID.(primitive.ObjectID), nil
}

func getList(listID primitive.ObjectID, userDataID primitive.ObjectID) (list LinksList, err error) {
	listCursor := databaseUtil.LinksListsColl.FindOne(databaseUtil.Ctx, bson.M{"_id": listID})
	err = listCursor.Decode(&list)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return LinksList{}, fmt.Errorf("NoUserFound: userType with the given Id does not exist")
		}
		return LinksList{}, err
	}

	if list.LinkAccessibility != "public" {
		if userDataID != list.OwnerDataId {
			_, ok := list.AccessIds[userDataID]
			if !ok {
				return LinksList{}, fmt.Errorf("UserDoesNotHaveAccess: User does not have access to list")
			}
		}
	}
	return list, nil
}

func GetList(listID primitive.ObjectID, userDataID primitive.ObjectID) (list LinksList, err error) {
	list, err = getList(listID, userDataID)
	if err != nil {
		return LinksList{}, err
	}

	list.OwnerDataId = primitive.ObjectID{}

	return list, nil
}

func GetListAsDiscordUser(listID primitive.ObjectID, discordID string) (list LinksList, err error) {
	userDataID, err := getUserDataIDWithDiscordID(discordID)

	list, err = GetList(listID, userDataID)
	if err != nil {
		return LinksList{}, err
	}

	return list, nil
}

func addLinkToList(listID primitive.ObjectID, userDataID primitive.ObjectID, linkData []byte) (err error) {
	link, err := newLink(linkData)
	if err != nil {
		return err
	}

	list, err := GetList(listID, userDataID)
	if err != nil {
		return err
	}

	if list.OwnerDataId != userDataID {
		perms, ok := list.AccessIds[userDataID]
		if !ok {
			return fmt.Errorf("UserDoesNotHaveAccess: User does not have access to list")
		} else {
			if perms != "edit" && perms != "add" {
				return fmt.Errorf("UserDoesNotHaveAccess: User does not have access to add links")
			}
		}
	}

	_, err = databaseUtil.LinksListsColl.UpdateOne(databaseUtil.Ctx,
		bson.M{"_id": list.ID},
		bson.M{"$push": bson.M{"links": link}},
	)
	if err != nil {
		return err
	}

	return nil
}

func AddCategoryToList(category string, listID primitive.ObjectID, userDataID primitive.ObjectID) (err error) {
	list, err := GetList(listID, userDataID)
	if err != nil {
		return err
	}

	if list.OwnerDataId != userDataID {
		perms, ok := list.AccessIds[userDataID]
		if !ok {
			return fmt.Errorf("UserDoesNotHaveAccess: User does not have access to list")
		} else {
			if perms != "edit" && perms != "add" {
				return fmt.Errorf("UserDoesNotHaveAccess: User does not have access to add links")
			}
		}
	}

	_, err = databaseUtil.LinksListsColl.UpdateOne(databaseUtil.Ctx,
		bson.M{"_id": list.ID},
		bson.M{"$push": bson.M{"categories": category}},
	)
	if err != nil {
		return err
	}

	return nil
}

func AddCategoryToListAsDiscordUser(category string, listID primitive.ObjectID, discordID string) (err error) {
	userDataID, err := getUserDataIDWithDiscordID(discordID)

	err = AddCategoryToList(category, listID, userDataID)
	if err != nil {
		return err
	}

	return nil
}

//func GetListByName(listName string, userDataID primitive.ObjectID) (list LinksList, err error) {
//	listCursor := databaseUtil.LinksListsColl.FindOne(databaseUtil.Ctx, bson.M{"list_name": listName})
//	err = listCursor.Decode(&list)
//	if err != nil {
//		if err.Error() == "mongo: no documents in result" {
//			return LinksList{}, fmt.Errorf("NoUserFound: userType with the given Id does not exist")
//		}
//		return LinksList{}, err
//	}
//
//	if list.LinkAccessibility != "public" {
//		if userDataID != list.OwnerDataId {
//			_, ok := list.AccessIds[userDataID]
//			if !ok {
//				return LinksList{}, fmt.Errorf("UserDoesNotHaveAccess: User does not have access to list")
//			}
//		}
//	}
//	return list, nil
//}

//func getListByName(listName string) (list LinksList, err error) {
//	listCursor := databaseUtil.LinksListsColl.FindOne(databaseUtil.Ctx, bson.M{"list_name": listName})
//	err = listCursor.Decode(&list)
//	if err != nil {
//		if err.Error() == "mongo: no documents in result" {
//			return LinksList{}, fmt.Errorf("NoUserFound: userType with the given Id does not exist")
//		}
//		return LinksList{}, err
//	}
//	return list, nil
//
//}
