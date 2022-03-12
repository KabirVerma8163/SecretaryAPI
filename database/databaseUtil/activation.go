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
var UserDataCollection *mongo.Collection
var UsersCollection *mongo.Collection
var LinksListsCollection *mongo.Collection

//var DiscordLinksListsColl *mongo.Collection

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
	UsersCollection = ApiDb.Collection("users")
	UserDataCollection = ApiDb.Collection("user_data")
	LinksListsCollection = ApiDb.Collection("links_lists")

	return nil
}

func Close() error {
	return Client.Disconnect(Ctx)
}
