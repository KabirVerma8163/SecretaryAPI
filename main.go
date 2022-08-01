package main

import (
	"LinkingAPI/share_my_feed/api"
	"LinkingAPI/share_my_feed/database"
	"LinkingAPI/share_my_feed/database/databaseUtil"
	"fmt"
	"log"
)

func main() {

	err := databaseUtil.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer databaseUtil.Close()

	//err = databaseUtil.ClearAllCollections()
	//if err != nil {
	//	fmt.Println(err)
	//}

	//userData := []byte(`{
	//	"username": "facebook",
	//	"email":  "facebook@facebook.com",
	//	"password": "your mom",
	//	"name" : "facebook"
	//}`)
	//
	//discordData := []byte(`{
	//	"discord_username": "ToastedWaffle",
	//	"discord_discriminator":  1234,
	//	"discord_id": "545286684007464973"
	//}`)
	//
	//listData := []byte(`{
	//	"list_name": "list one",
	//	"categories": ["science", "technology"]
	//}`)
	//
	//err = database.NewUser(userData)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//err = database.NewDiscordUser(discordData)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//err = database.AddListAsDiscordUser(listData, "545286684007464973")
	//if err != nil {
	//	fmt.Println(err)
	//}
	database.DiscordCacheInit()
	err = api.Start()
	if err != nil {
		fmt.Println(err)
	}
	//link := []byte(`{
	//"link_name": "fred",
	//"link_url" : "lunar.cf"
	//}`)
	//
	//discordID := "545286684007464973"
	//
	//listID, err := primitive.ObjectIDFromHex("6243d4aff2dbc8a477d1fa2c")
	//if err != nil {
	//	fmt.Println(err)
	//}

	//err = database.AddLinkAsDiscordUser(link, listID, discordID)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//list, err := database.GetListAsDiscordUser(listID, discordID)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//databaseUtil.PrettyPrint(list)
}
