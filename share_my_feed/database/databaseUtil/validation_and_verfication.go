package databaseUtil

import (
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func IsUniqueUsername(username string) bool {
	podcast := bson.M{}
	if err := UsersColl.FindOne(Ctx, bson.M{"username": username}).Decode(&podcast); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return true
		}
	}
	return false
}

func isAppropriateEmail(email string) bool {
	return false
}

func IsUniqueEmail(email string) bool {
	//podcast := userType{}
	//var podcast userType
	podcast := bson.M{}
	if err := UsersColl.FindOne(Ctx, bson.M{"email": email}).Decode(&podcast); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return true
		}
	}
	return false
}

func IsAppropriateName(username string) bool {
	return true
}

func isPassword(password string) (string, error) {
	// ADD - all the other needs for a secure password

	return password, nil
}

func GetPasswordHash(password string) (hash []byte, err error) {
	tempPassword, err := isPassword(password)
	if err != nil {
		return []byte{}, err
	}
	passwordBytes := []byte(tempPassword)

	hashBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	if err != nil {
		return []byte{}, err
	}

	return hashBytes, nil
}
