package database

import (
	"LinkingAPI/database/databaseUtil"
	"encoding/json"
	"github.com/kamva/mgm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userDataType struct {
	mgm.DefaultModel `bson:",inline"`
	Username         string                        `json:"username" bson:"username"`
	UserId           primitive.ObjectID            `json:"user_id" bson:"user_id"`
	Name             string                        `json:"name" bson:"name"`
	LinksList        map[string]primitive.ObjectID `json:"links_list" bson:"links_list"`
	CategoriesList   []string                      `json:"categories_list" bson:"categories_list"`
	Email            string                        `json:"email" bson:"email"`
	OtherData        []map[string]interface{}      `json:"other_data" bson:"other_data"`
	//ConnectedAccounts []connectedAccount            `json:"connected_accounts" bson:"connected_accounts"`
	TempUser     bool   `json:"temp_user" bson:"temp_user"`
	TempUserName string `json:"temp_username" bson:"temp_username"`
}

func (data *userDataType) initialize(userID primitive.ObjectID) {
	data.LinksList = map[string]primitive.ObjectID{}
	data.CategoriesList = []string{}

	data.UserId = userID
}

func newUserData(userDataData []byte, userID primitive.ObjectID) (userDataId primitive.ObjectID, err error) {
	var userdata userDataType
	err = json.Unmarshal(userDataData, &userdata)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	userdata.initialize(userID)

	userDataInsertResult, err := databaseUtil.UserDataCollection.InsertOne(databaseUtil.Ctx, userdata)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return userDataInsertResult.InsertedID.(primitive.ObjectID), nil
}
