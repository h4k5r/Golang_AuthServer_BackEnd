package auth_functions

import (
	"AuthServer/Database"
	"AuthServer/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"os"
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
	doc := bson.D{
		{"email", credentials.Email},
	}
	result := Database.Database.Collection("users").FindOne(context.TODO(), doc)
	//result, _ := Database.Database.Collection("users").Find(context.TODO(), bson.)

	if result.Err() == nil {
		writer.Write([]byte("User Exists"))
		return
	}
	hashedPassword, hashErr := utils.HashPassword(credentials.Password)

	if hashErr != nil {
		writer.Write([]byte("Hashing Error"))
		return
	}

	doc = bson.D{
		{"email", credentials.Email},
		{"password", hashedPassword},
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

	doc := bson.D{
		{"email", credentials.Email},
	}

	result := Database.Database.Collection("users").FindOne(context.TODO(), doc)

	if result.Err() != nil {
		writer.Write([]byte(result.Err().Error()))
		return
	}
	var user bson.D
	err = result.Decode(&user)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	userMap := user.Map()
	value, _ := userMap["password"]

	password := value.(string)
	if !utils.CheckPasswordHash(credentials.Password, password) {
		writer.Write([]byte("Password Wrong"))
		return
	}
	email, _ := userMap["password"]

	jwt, err := utils.GetJWT(email.(string))
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	utils.SetCORS(writer)
	//writer.Write([]byte("Success"))
	err = json.NewEncoder(writer).Encode(struct {
		Token string `token:"email"`
	}{Token: jwt})
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}
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
	doc := bson.D{
		{"email", credentials.Email},
	}

	result := Database.Database.Collection("users").FindOne(context.TODO(), doc)

	if result.Err() != nil {
		writer.Write([]byte(result.Err().Error()))
		return
	}
	var user bson.D
	err = result.Decode(&user)
	if err != nil {
		writer.Write([]byte("if the email found the reset link would be sent"))
		return
	}
	userMap := user.Map()
	value, _ := userMap["email"]

	email := value.(string)
	println(email)
	from := mail.NewEmail("Potato Lord", "rahul16086@gmail.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail("Example User", email)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	writer.Write([]byte("if the email found the reset link would be sent"))

}
