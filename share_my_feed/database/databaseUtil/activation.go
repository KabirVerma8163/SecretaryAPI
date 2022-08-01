package databaseUtil

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
)

var Client *mongo.Client
var Ctx context.Context
var ApiDb *mongo.Database
var UserDataColl *mongo.Collection
var UsersColl *mongo.Collection
var LinksListsColl *mongo.Collection
var DiscordConnectedAccountsColl *mongo.Collection

func Start() error {
	var err error
	//URI := os.Getenv("TEST_MGM_URL")
	//dbName := os.Getenv("TEST_DB")

	err = godotenv.Load(".env")
	if err != nil {
		return err
	}

	URI := os.Getenv("MONGO_URI")
	Client, err = mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		return err
	}

	Ctx = context.Background()

	err = Client.Connect(Ctx)
	if err != nil {
		return err
	}

	if err := Client.Ping(context.Background(), readpref.Primary()); err != nil {
		// Can't connect to Mongo server
		return err
	}
	fmt.Println("Connection Secured")

	ApiDb = Client.Database("SecretaryAPI")
	UsersColl = ApiDb.Collection("users")
	UserDataColl = ApiDb.Collection("user_data")
	DiscordConnectedAccountsColl = ApiDb.Collection("discord_connected_accounts")

	LinksListsColl = ApiDb.Collection("links_lists")
	return nil
}

func Close() error {
	return Client.Disconnect(Ctx)
}
