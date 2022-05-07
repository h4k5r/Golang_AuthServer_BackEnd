package main

import (
	"AuthServer/Database"
	"AuthServer/auth_functions"
	"AuthServer/constants"
	"AuthServer/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load("main/static/.env")

	if err != nil {
		println(err.Error())
		return
	}

	Database.InitDB()

	if Database.Err != nil {
		log.Fatal("Database Error")
		return
	}

	r := mux.NewRouter()

	r.HandleFunc(constants.Register, auth_functions.HandleRegister).Methods(constants.Post)
	r.HandleFunc(constants.Login, auth_functions.HandleLogin).Methods(constants.Post)
	r.HandleFunc(constants.ResetPassword, auth_functions.HandleReset).Methods(constants.Post)

	r.HandleFunc("*", utils.HandleOptions).Methods("OPTIONS")

	log.Fatal(http.ListenAndServe(":8080", r))
}
