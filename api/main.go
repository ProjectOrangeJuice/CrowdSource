package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var conn *mongo.Database
var debug bool

func main() {
	//load the env file
	err := godotenv.Load()
	failOnError(err, "Error getting env vars")

	debug, err = strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		debug = false
	}

	//Create a database connection
	conn, err = configDB(context.Background())
	failOnError(err, "Connecting to database failed")

	router := mux.NewRouter()

	router.HandleFunc("/product/{barcode}", getProduct).Methods("GET")
	router.HandleFunc("/product", changeProduct).Methods("POST")
	// router.HandleFunc("/vote/{barcode}", voteOnProduct).Methods("POST")
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
}

//Auth0
type customClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}

func getUsername(r *http.Request) string {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	log.Printf("Length %v", len(authHeaderParts))

	//when auth0 is turned off, all users are "test"
	if len(authHeaderParts) < 2 && debug {
		log.Printf("Token not found. Giving test username")
		return "test"
	}
	tokenString := authHeaderParts[1]
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, nil)
	if err != nil && debug {
		log.Printf("Token not found. Giving test username (error) %e", err)
		return "test"
	} else if err != nil {
		failOnError(err, "Token not found")
	}
	claims, _ := token.Claims.(*customClaims)
	fmt.Printf("(user) %s ", claims.Subject)
	return claims.Subject

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
