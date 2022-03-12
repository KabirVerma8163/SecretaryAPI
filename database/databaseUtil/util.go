package databaseUtil

import "go.mongodb.org/mongo-driver/bson"

func ClearAllCollections() error {
	_, err := LinksListsCollection.DeleteMany(Ctx, bson.M{})
	if err != nil {
		return err
	}
	_, err = UsersCollection.DeleteMany(Ctx, bson.M{})
	if err != nil {
		return err
	}
	_, err = UserDataCollection.DeleteMany(Ctx, bson.M{})
	if err != nil {
		return err
	}
	return nil
}
