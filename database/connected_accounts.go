package database

import (
	"encoding/json"
	"github.com/kamva/mgm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type connectedAccount struct {
	mgm.DefaultModel `bson:",inline"`
	Username         string             `json:"username" bson:"username"`
	UserId           primitive.ObjectID `json:"user_id" bson:"user_id"`
	UserDataId       primitive.ObjectID `json:"user_data_id" bson:"user_data_id"`
	ClientType       string             `json:"client_type" bson:"client_type"`
}

type discordAccount struct {
	connectedAccount
	DiscordUsername      string                 `json:"discord_user_name" bson:"discord_user_name"`
	DiscordDiscriminator int                    `json:"discord_tag" bson:"discord_tag"`
	DiscordId            string                 `json:"discord_id" bson:"discord_id"`
	OtherDetails         map[string]interface{} `json:"other_details" bson:"other_details"`
}

func (discord *discordAccount) initialize(discordData []byte, userID primitive.ObjectID) (err error) {
	err = json.Unmarshal(discordData, discord)
	if err != nil {
		return err
	}

	return nil
}
