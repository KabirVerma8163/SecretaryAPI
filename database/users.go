package database

import (
	"LinkingAPI/database/databaseUtil"
	"encoding/json"
	"fmt"
	"github.com/kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userType struct {
	mgm.DefaultModel `bson:",inline"`
	Username         string             `json:"username" bson:"username"`
	Email            string             `json:"email" bson:"email"`
	PasswordHash     []byte             `json:"password_hash" bson:"password_hash"`
	Name             string             `json:"name" bson:"name"`
	TempUser         bool               `json:"temp_user" bson:"temp_user"`
	TempUsername     string             `json:"temp_username" bson:"temp_username"`
	UserDataID       primitive.ObjectID `json:"user_data_id" bson:"user_data_id"`
	//PasswordSalt     string             `json:"password_salt" bson:"password_salt"`
}

// MODIFY : Does this struct really go here
type password struct {
	Password string `json:"password" bson:"password"`
}

func (user *userType) initializeUser(userData []byte) (err error) {
	err = json.Unmarshal(userData, user)
	if err != nil {
		return err
	}

	if !databaseUtil.IsUniqueEmail(user.Email) {
		return fmt.Errorf("EmailTaken: Email %v has been associated with another account", user.Email)
	}

	if !databaseUtil.IsUniqueUsername(user.Username) {
		return fmt.Errorf("UsernameTaken: Username %v has been taken", user.Username)
	}

	if !databaseUtil.IsAppropriateName(user.Name) {
		return fmt.Errorf("InappropriateName: Name %v is not a valid name", user.Name)
	}

	var pass password
	err = json.Unmarshal(userData, &pass)
	if err != nil {
		return err
	}

	user.PasswordHash, err = databaseUtil.GetPasswordHash(pass.Password)
	if err != nil {
		return err
	}

	return nil
}

func (user *userType) initializeTempDiscordUser(discordData []byte) error {
	var discord discordAccount
	err := json.Unmarshal(discordData, &discord)
	if err != nil {
		return err
	}

	user.TempUser = true
	user.TempUsername = discord.DiscordUsername

	return nil
}

func (user *userType) initializeDiscordTempUser() (err error) {
	user.TempUser = true
	//user.TempUsername =
	return nil
}

// NewUser Inserts it into the database and then returns the userData
func NewUser(userData []byte) (err error) {
	var user userType

	err = user.initializeUser(userData)
	if err != nil {
		return err
	}

	userInsertResult, err := databaseUtil.UsersColl.InsertOne(databaseUtil.Ctx, user)
	if err != nil {
		return err
	}

	user.ID = userInsertResult.InsertedID.(primitive.ObjectID)

	userArray, err := json.Marshal(user)
	if err != nil {
		return err
	}

	user.UserDataID, err = newUserData(userArray, user.ID)
	if err != nil {
		return err
	}

	_, err = databaseUtil.UsersColl.UpdateOne(databaseUtil.Ctx,
		bson.M{"_id": user.ID}, bson.D{
			{"$set", bson.D{{"user_data_id", user.UserDataID}}},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// NewDiscordUser Inserts it into the database and then returns the userData
func NewDiscordUser(discordData []byte) (err error) {
	var user userType
	err = user.initializeTempDiscordUser(discordData)

	userInsertResult, err := databaseUtil.UsersColl.InsertOne(databaseUtil.Ctx, user)
	if err != nil {
		return err
	}

	user.ID = userInsertResult.InsertedID.(primitive.ObjectID)

	userData, err := newTempUserDataFromDiscord(discordData, user.ID)
	if err != nil {
		return err
	}

	_, err = databaseUtil.UsersColl.UpdateOne(databaseUtil.Ctx,
		bson.M{"_id": user.ID}, bson.D{
			{"$set", bson.D{{"user_data_id", userData.ID}}},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

//	var user userType
//	err = json.Unmarshal(userData, &user)
//	if err != nil {
//		return err
//	}
//
//	err = user.initializeUser(userData)
//	if err != nil {
//		return err
//	}
//
//	userInsertResult, err := databaseUtil.UsersColl.InsertOne(databaseUtil.Ctx, user)
//	if err != nil {
//		return err
//	}
//
//	user.ID = userInsertResult.InsertedID.(primitive.ObjectID)
//
//	userArray, err := json.Marshal(user)
//	if err != nil {
//		return err
//	}
//
//	user.UserDataId, err = newUserData(userArray, user.ID)
//	if err != nil {
//		return err
//	}
//
//	_, err = databaseUtil.UsersColl.UpdateOne(databaseUtil.Ctx,
//		bson.M{"_id": user.ID}, bson.D{
//			{"$set", bson.D{{"user_data_id", user.UserDataId}}},
//		},
//	)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
