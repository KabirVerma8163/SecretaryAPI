package main

import (
	"LinkingAPI/database"
	"LinkingAPI/database/databaseUtil"
	"log"
	//"LinkingAPI/db-communication"
	"encoding/json"
	"fmt"
)

func main() {

	err := databaseUtil.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer databaseUtil.Close()

	err = databaseUtil.ClearAllCollections()
	if err != nil {
		fmt.Println(err)
	}

	/* Checking insertion of accountAttributes
	 */
	//err = database2.ClearAllCollections()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//accountAttributes := map[string]interface{} {
	//	"AccountUsername" : "ToastedWaffle",
	//	"AccountDetails" : map[string]interface{} {
	//		"Username" : "ToastedWaffle",
	//		"UserTag" : 2496,
	//	},
	//
	//}
	//
	//err = database2.AddTempDiscordUser(accountAttributes)
	//if err != nil {
	//	fmt.Print(err)
	//}

	//listAttributes := map[string]interface{} {
	//	"OwnerId" : "61eca2a5cf888ce76bff8207",
	//	"ListName": "First List",
	//}
	//
	//err = database2.AddLinkList(listAttributes)
	//if err != nil {
	//	fmt.Print(err)
	//}

	//linkAttributes := map[string]interface{}{
	//	"Name": "facebook",
	//	"URL":  "facebook.com",
	//}
	//
	//linkID, err := primitive.ObjectIDFromHex("61eca2ea066ef49659cda746")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//link, err := database2.NewLink(linkAttributes)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(link)
	//
	//err = database2.AddLinkToList(link, linkID)
	//if err != nil {
	//	fmt.Print(err)
	//}
	//
	//userAttributes := map[string]interface{} {
	//	"Username" : "ToastedWaffle",
	//	"Email" : "Kabir2@gmail.com",
	//	"Name" : "Kabir Verma",
	//	"Password" : "ThisDick",
	//}
	//err = database2.AddUser(userAttributes)
	//if err != nil {
	//	fmt.Print(err)
	//}

	//userId, err := database2.GetUserIdWithUsername("ToastedWaffle")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(userId)
	////
	////userId, err = database2.GetUserIdWithEmail("Kabir2@gmail.com")
	////if err != nil {
	////	fmt.Println(err)
	////}
	////
	////fmt.Println(userId)
	//
	//listAttributes := map[string]interface{} {
	//	"ListName" : "Kabir's second list",
	//	"OwnerId" : userId,
	//	""
	//}
	//// TODO: the list is is decided by the API not the client
	//err = database2.AddLinkList(listAttributes)
	//if err != nil {
	//	fmt.Println(err)
	//}

	userData := []byte(`{
		"username": "facebook",
		"email":  "facebook@facebook.com",
		"password": "your mom",
		"name" : "facebook"
	}`)
	//
	//link, err := database.NewLink(linkData)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(link)

	err = database.NewUser(userData)
	if err != nil {
		fmt.Println(err)
	}
}

func PrettyPrint(x interface{}) {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
}
