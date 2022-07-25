package database

import (
	"LinkingAPI/share_my_feed/database/databaseUtil"
	"fmt"
	"github.com/kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LinksList struct {
	mgm.DefaultModel `bson:",inline"`
	//OwnerId           primitive.ObjectID              `json:"owner_id" bson:"owner_id"`
	Username          string                          `json:"owner_username" bson:"owner_username"`
	UserName          string                          `json:"user_name" bson:"owner_name"`
	UserDataId        primitive.ObjectID              `json:"user_data_id" bson:"owner_data_id"`
	AccessIds         map[primitive.ObjectID][]string `json:"access_ids" bson:"access_ids"` // Must be userDataIDs
	ListName          string                          `json:"list_name" bson:"list_name"`
	Links             []LinkType                      `json:"links" bson:"links"`
	Categories        []string                        `json:"categories" bson:"categories"`
	LinkAccessibility string                          `json:"link_accessibility" bson:"link_accessibility"`
	Description       string                          `json:"description" bson:"description"`
}

// TODO be mindful of what the user can access, like ownerdataid

func (list *LinksList) initialize(dataID primitive.ObjectID) error {
	userData, err := getUserData(dataID)
	if err != nil {
		return err
	}

	list.Username = userData.Username
	list.UserName = userData.Name
	list.UserDataId = userData.ID

	if list.AccessIds == nil {
		list.AccessIds = map[primitive.ObjectID][]string{}
	}

	// TODO: fix this thing
	if list.Links == nil {
		list.Links = []LinkType{}
	} else {
		//list.Links, err = newLinks(list.Links)
	}

	if list.Description == "" {
		list.Description = fmt.Sprintf("This is %s's list!", list.Username)
	}

	if list.Categories == nil {
		list.Categories = []string{}
	}
	//list.LinkAccessibility = "private"
	return nil
}

func getList(listID, dataID primitive.ObjectID) (list LinksList, err error) {
	listCursor := databaseUtil.LinksListsColl.FindOne(databaseUtil.Ctx, bson.M{"_id": listID})
	err = listCursor.Decode(&list)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return LinksList{}, fmt.Errorf("NoUserFound: userType with the given Id does not exist")
		}
		return LinksList{}, err
	}

	if list.LinkAccessibility != "public" {
		if dataID != list.UserDataId {
			_, ok := list.AccessIds[dataID]
			if !ok {
				return LinksList{}, fmt.Errorf("UserDoesNotHaveAccess: User does not have access to list")
			}
		}
	}
	return list, nil
}

func GetLists(dataID primitive.ObjectID) (lists []LinksList, err error) {
	userData, err := getUserData(dataID)
	for _, l := range userData.LinksList {
		list, err := getList(l, dataID)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}

	return lists, nil
}

func GetListById(listID, userDataID primitive.ObjectID) (list LinksList, err error) {
	list, err = getList(listID, userDataID)
	if err != nil {
		return LinksList{}, err
	}

	list.UserDataId = primitive.ObjectID{}
	return list, nil
}

func GetListByName(name string, dataID primitive.ObjectID) (list LinksList, err error) {
	userData, err := getUserData(dataID)

	listID, ok := userData.LinksList[name]
	if !ok {
		return LinksList{}, fmt.Errorf("list %s does not exist", name)
	}

	list, err = getList(listID, dataID)
	if err != nil {
		return LinksList{}, err
	}

	list.UserDataId = primitive.ObjectID{}
	return list, nil
}

func addListToDB(list LinksList) (linkId primitive.ObjectID, err error) {
	listInsertResult, err := databaseUtil.LinksListsColl.InsertOne(databaseUtil.Ctx, list)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return listInsertResult.InsertedID.(primitive.ObjectID), nil
}

func AddList(list LinksList, dataID primitive.ObjectID) (err error) {
	err = list.initialize(dataID)
	if err != nil {
		return err
	}

	userData, err := getUserData(dataID)
	for k, _ := range userData.LinksList {
		if k == list.ListName {
			return fmt.Errorf("list %v already exists", list.ListName)
		}
	}

	listID, err := addListToDB(list)

	_, err = databaseUtil.UserDataColl.UpdateOne(databaseUtil.Ctx,
		bson.M{"_id": list.UserDataId}, bson.D{
			{"$push", bson.D{{"links_list", listID}}},
		})
	if err != nil {
		return err
	}

	return err
}

func RemoveList(listID, dataID primitive.ObjectID) (err error) {
	list, err := getList(listID, dataID)
	if list.UserDataId != dataID {
		return fmt.Errorf("AuthenticationError: User does not own list")
	}

	_, err = databaseUtil.UserDataColl.UpdateOne(databaseUtil.Ctx,
		bson.M{"_id": list.UserDataId}, bson.D{
			{"$unset", bson.D{{fmt.Sprintf("links_list.%v", list.ListName), ""}}},
		})
	if err != nil {
		return err
	}

	_, err = databaseUtil.LinksListsColl.DeleteOne(databaseUtil.Ctx, bson.M{"_id": list.ID})
	if err != nil {
		return err
	}

	return nil
}

func RemoveListWithName(name string, dataID primitive.ObjectID) (err error) {
	data, err := getUserData(dataID)
	if err != nil {
		return err
	}

	listID, ok := data.LinksList[name]
	if !ok {
		return fmt.Errorf("AuthenticationError: User does not own list")
	}

	_, err = databaseUtil.UserDataColl.UpdateOne(databaseUtil.Ctx,
		bson.M{"_id": data}, bson.D{
			{"$unset", bson.D{{fmt.Sprintf("links_list.%v", name), ""}}},
		})
	if err != nil {
		return err
	}

	_, err = databaseUtil.LinksListsColl.DeleteOne(databaseUtil.Ctx, bson.M{"_id": listID})
	if err != nil {
		return err
	}

	return nil
}

func canUpdateList(listId primitive.ObjectID, dataID primitive.ObjectID, permsNeeded []string) (err error) {
	list, err := getList(listId, dataID)
	if err != nil {
		return err
	}

	perms := list.AccessIds[dataID]
	access := false

	for _, v := range perms {
		for _, v2 := range permsNeeded {
			if v == v2 {
				access = true
			}
		}
	}
	if !access {
		return fmt.Errorf("user does not have required permissions to edit list")
	}
	//, field string, value interface{}
	//_, err = databaseUtil.LinksListsColl.UpdateOne(databaseUtil.Ctx,
	//	bson.M{"_id": list.ID},
	//	bson.M{"$set": bson.M{field: value}},
	//)
	//if err != nil {
	//	return err
	//}

	return nil
}

func addLinkToList(listID primitive.ObjectID, userDataID primitive.ObjectID, link LinkType) (list LinksList, err error) {
	err = canUpdateList(listID, userDataID, []string{"edit", "add"})
	if err != nil {
		return LinksList{}, err
	}

	updateResult, err := databaseUtil.LinksListsColl.UpdateOne(databaseUtil.Ctx,
		bson.M{"_id": listID},
		bson.M{"$push": bson.M{"links": link}},
	)
	if err != nil {
		return LinksList{}, err
	}

	//list, err := updateResult.UnmarshalBSON()
	fmt.Println(updateResult)

	return LinksList{}, nil
}

//func AddCategoryToList(category string, listID primitive.ObjectID, userDataID primitive.ObjectID) (err error) {
//	list, err := getList(listID, userDataID)
//	if err != nil {
//		return err
//	}
//
//	if list.UserDataId != userDataID {
//		perms, ok := list.AccessIds[userDataID]
//		if !ok {
//			return fmt.Errorf("UserDoesNotHaveAccess: User does not have access to list")
//		} else {
//			if perms != "edit" && perms != "add" {
//				return fmt.Errorf("UserDoesNotHaveAccess: User does not have access to add links")
//			}
//		}
//	}
//
//	_, err = databaseUtil.LinksListsColl.UpdateOne(databaseUtil.Ctx,
//		bson.M{"_id": list.ID},
//		bson.M{"$push": bson.M{"categories": category}},
//	)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
