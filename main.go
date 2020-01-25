package main

import (
	"./config"
	"./database"
	"./helpers"
	"./routes"
	"./ws"

	"log"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/handlers"
	"github.com/subosito/gotenv"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	gotenv.Load()
}

func main() {

	log.Println("Server is init...")
	conf := config.GetServerConfig()

	env := os.Getenv("ENV")

	helpers.Initialize()

	database.InitializeRedis(env)
	db := database.Initialize(env)
	defer db.Con.Close()

	go ws.Init()
	go ws.InitPersonalChannel()

	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "X-AT", "X-RT", "Content-Length", "X-Device"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	router := routes.NewRouter()

	log.Println("Server is start on " + conf.Port)
	log.Fatal(http.ListenAndServe(conf.Port, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router)))
}
