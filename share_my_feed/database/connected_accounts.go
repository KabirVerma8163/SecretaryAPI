package database

import (
	"LinkingAPI/share_my_feed/database/databaseUtil"
	"encoding/json"
	"fmt"
	"github.com/kamva/mgm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConnectedAccount struct {
	mgm.DefaultModel `bson:",inline"`
	Username         string             `json:"username" bson:"username"`
	UserID           primitive.ObjectID `json:"user_id" bson:"user_id"`
	UserDataID       primitive.ObjectID `json:"user_data_id" bson:"user_data_id"`
	ClientType       string             `json:"client_type" bson:"client_type"`
}

type discordAccount struct {
	ConnectedAccount     `bson:",inline"`
	DiscordUsername      string                 `json:"discord_username" bson:"discord_username"`
	DiscordDiscriminator int                    `json:"discord_discriminator" bson:"discord_discriminator"`
	DiscordID            string                 `json:"discord_id" bson:"discord_id"`
	OtherDetails         map[string]interface{} `json:"other_details" bson:"other_details"`
}

// api_redundant -> newTempUser -> newUserData -> newDiscordAccount

func (discord *discordAccount) initialize(userData UserDataType) (err error) {
	//discord.UserID = userData.UserID
	//discord.UserDataId = userData.ID
	discord.ClientType = "discord"

	discord.Username = userData.Username
	discord.UserID = userData.UserID
	discord.UserDataID = userData.ID
	//discord.ID, err = primitive.ObjectIDFromHex(discord.DiscordID)
	//if err != nil {
	//	return err
	//}

	return nil
}

func newDiscordConnectedAccount(discordData []byte, userData UserDataType) (discordID primitive.ObjectID, err error) {
	var discord discordAccount
	err = json.Unmarshal(discordData, &discord)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	err = discord.initialize(userData)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	accountInsertResult, err := databaseUtil.DiscordConnectedAccountsColl.InsertOne(databaseUtil.Ctx, discord)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	discord.ID = accountInsertResult.InsertedID.(primitive.ObjectID)

	discord.OtherDetails = map[string]interface{}{}

	return discord.ID, err
}

func getDiscordDataIDFromDiscordID(discordID string) (discordAccount, error) {
	accountCursor := databaseUtil.DiscordConnectedAccountsColl.FindOne(databaseUtil.Ctx, bson.M{"discord_id": discordID})
	var account discordAccount
	err := accountCursor.Decode(&account)
	if err != nil {
		return discordAccount{}, err
	}

	return account, nil
}

func GetUserDataIDWithDiscordID(discordID string) (dataID primitive.ObjectID, err error) {
	userDataCursor := databaseUtil.DiscordConnectedAccountsColl.FindOne(databaseUtil.Ctx, bson.M{"discord_id": discordID})
	var discordAcc discordAccount
	err = userDataCursor.Decode(&discordAcc)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return primitive.ObjectID{}, fmt.Errorf("ServerError: UserDataType for given user does not exist")
		}
		return primitive.ObjectID{}, err
	}
	return discordAcc.UserDataID, err
}
