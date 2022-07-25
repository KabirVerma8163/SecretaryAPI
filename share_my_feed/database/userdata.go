package database

import (
	"LinkingAPI/share_my_feed/database/databaseUtil"
	"encoding/json"
	"fmt"
	"github.com/kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDataType struct {
	mgm.DefaultModel  `bson:",inline"`
	Username          string                        `json:"username" bson:"username"`
	UserID            primitive.ObjectID            `json:"user_id" bson:"user_id"`
	Name              string                        `json:"name" bson:"name"`
	LinksList         map[string]primitive.ObjectID `json:"links_list" bson:"links_list"`
	CategoriesList    []string                      `json:"categories_list" bson:"categories_list"`
	Email             string                        `json:"email" bson:"email"`
	OtherData         []map[string]interface{}      `json:"other_data" bson:"other_data"`
	ConnectedAccounts map[string]primitive.ObjectID `json:"connected_accounts" bson:"connected_accounts"`
	TempUser          bool                          `json:"temp_user" bson:"temp_user"`
	TempUsername      string                        `json:"temp_username" bson:"temp_username"`
}

func (data *UserDataType) initialize(userID primitive.ObjectID) {
	data.LinksList = map[string]primitive.ObjectID{}
	data.CategoriesList = []string{}

	data.UserID = userID

}

func (data *UserDataType) initializeDiscordTempUser(discord discordAccount, tempUserID primitive.ObjectID) error {
	data.initialize(tempUserID)

	data.TempUser = true
	data.TempUsername = discord.DiscordUsername
	data.UserID = tempUserID

	data.ConnectedAccounts = map[string]primitive.ObjectID{}
	data.ConnectedAccounts["discord"] = discord.ID

	return nil
}

func newUserData(userDataData []byte, userID primitive.ObjectID) (userDataID primitive.ObjectID, err error) {
	var userData UserDataType
	err = json.Unmarshal(userDataData, &userData)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	userData.initialize(userID)

	userDataInsertResult, err := databaseUtil.UserDataColl.InsertOne(databaseUtil.Ctx, userData)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return userDataInsertResult.InsertedID.(primitive.ObjectID), nil
}

func newTempUserDataFromDiscord(discordData []byte, tempUserID primitive.ObjectID) (userData UserDataType, err error) {
	var discord discordAccount
	err = json.Unmarshal(discordData, &discord)
	if err != nil {
		return userData, err
	}

	err = userData.initializeDiscordTempUser(discord, tempUserID)
	if err != nil {
		return UserDataType{}, err
	}

	userDataInsertResult, err := databaseUtil.UserDataColl.InsertOne(databaseUtil.Ctx, userData)
	if err != nil {
		return UserDataType{}, err
	}

	userData.ID = userDataInsertResult.InsertedID.(primitive.ObjectID)

	discordID, err := newDiscordConnectedAccount(discordData, userData)
	if err != nil {
		return UserDataType{}, err
	}

	_, err = databaseUtil.UserDataColl.UpdateOne(databaseUtil.Ctx,
		bson.M{"_id": userData.ID}, bson.D{
			{"$set", bson.D{{"connected_accounts.discord", discordID}}},
		},
	)
	return userData, nil
}

func GetUserDataWithUserID(userID primitive.ObjectID) (data UserDataType, err error) {
	userDataCursor := databaseUtil.UserDataColl.FindOne(databaseUtil.Ctx, bson.M{"user_id": userID})
	err = userDataCursor.Decode(&data)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return UserDataType{}, fmt.Errorf("ServerError: UserDataType for given user does not exist")
		}
		return UserDataType{}, err
	}
	return data, err
}

func GetUserDataIDWithUsername(username string) (userDataID primitive.ObjectID, err error) {
	userDataCursor := databaseUtil.UserDataColl.FindOne(databaseUtil.Ctx, bson.M{"username": username})
	var userData UserDataType
	err = userDataCursor.Decode(&userData)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return primitive.ObjectID{}, fmt.Errorf("ServerError: UserDataType for given user does not exist")
		}
		return primitive.ObjectID{}, err
	}

	return userData.ID, nil
}

func GetUserDataWithUsername(username string) (data UserDataType, err error) {
	userDataCursor := databaseUtil.UserDataColl.FindOne(databaseUtil.Ctx, bson.M{"username": username})
	err = userDataCursor.Decode(&data)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return UserDataType{}, fmt.Errorf("ServerError: UserDataType for given user does not exist")
		}
		return UserDataType{}, err
	}
	return data, err
}

func getUserData(userDataID primitive.ObjectID) (data UserDataType, err error) {
	userDataCursor := databaseUtil.UserDataColl.FindOne(databaseUtil.Ctx, bson.M{"_id": userDataID})
	err = userDataCursor.Decode(&data)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return UserDataType{}, fmt.Errorf("ServerError: UserDataType for given user does not exist")
		}
		return UserDataType{}, err
	}
	return data, err
}
