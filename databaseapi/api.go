package databaseapi

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//StartAPI starts running the API linked to the local MongoDB database on the 'localhost:5000/' port
func StartAPI(collectionName string) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	//check if proper connection could be established
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	collection := client.Database("nitishp1812buildingdb").Collection(collectionName)

	handler := dbHandler{collection}

	//add the routes for the API
	router := mux.NewRouter()
	router.HandleFunc("/", handler.showAllJSON).Methods("GET")
	router.HandleFunc("/summary/", handler.showAllSummary).Methods("GET")
	router.HandleFunc("/filter/{left}={right}", handler.equalToFilterJSON).Methods("GET")
	router.HandleFunc("/summary/filter/{left}={right}", handler.equalToFilterSummary).Methods("GET")
	router.HandleFunc("/filter/{left}>{right}", handler.greaterThanFilterJSON).Methods("GET")
	router.HandleFunc("/summary/filter/{left}>{right}", handler.greaterThanFilterSummary).Methods("GET")
	router.HandleFunc("/filter/{left}<{right}", handler.lessThanFilterJSON).Methods("GET")
	router.HandleFunc("/summary/filter/{left}<{right}", handler.lessThanFilterSummary).Methods("GET")
	router.HandleFunc("/filter/{left}!{right}", handler.notEqualToFilterJSON).Methods("GET")
	router.HandleFunc("/summary/filter/{left}!{right}", handler.notEqualToFilterSummary).Methods("GET")

	fmt.Println("The local server is now running on http://localhost:5000/")

	//start listening for the input
	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Fatal(err)
	}
}
