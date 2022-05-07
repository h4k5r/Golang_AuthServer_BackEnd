package auth_functions

import (
	"AuthServer/Database"
	"AuthServer/models"
	"AuthServer/utils"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func HandleRegister(writer http.ResponseWriter, request *http.Request) {
	//Database.Database.Collection("users").FindOne()
	credentials := struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_Password"`
	}{}

	err := json.NewDecoder(request.Body).Decode(&credentials)

	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	result := Database.Database.Collection("users").FindOne(context.TODO(), models.CreateUser(credentials.Email, credentials.Password))

	if result.Err() == nil {
		writer.Write([]byte("User Exists"))
		return
	}
	doc := bson.D{
		{"email", credentials.Email},
		{"password", credentials.Password},
	}
	//result, err := Database.Database.Collection("users").InsertOne(context.TODO(), doc)

	insertResult, insertErr := Database.Database.Collection("users").InsertOne(context.TODO(), doc)
	if insertErr != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	println(insertResult.InsertedID)

	writer.Header().Set("Content-Type", "application/json")
	utils.SetCORS(writer)
	println()

}

func HandleLogin(writer http.ResponseWriter, request *http.Request) {
	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(request.Body).Decode(&credentials)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	println(credentials.Email)
	println(credentials.Password)

	writer.Header().Set("Content-Type", "application/json")
	utils.SetCORS(writer)
	writer.Write([]byte("Error marshalling json"))
}

func HandleReset(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	utils.SetCORS(writer)
	credentials := struct {
		Email string `json:"email"`
	}{}
	err := json.NewDecoder(request.Body).Decode(&credentials)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	println(credentials.Email)
	writer.Write([]byte("Error marshalling json"))
}
