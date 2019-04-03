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

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	collection := client.Database("nitishp1812buildingdb").Collection(collectionName)

	handler := dbHandler{collection}

	router := mux.NewRouter()
	router.HandleFunc("/", handler.showAllJSON).Methods("GET")
	router.HandleFunc("/summary/", handler.showAllSummary).Methods("GET")
	router.HandleFunc("/filter/{left}={right}", handler.equalToFilterJSON).Methods("GET")
	router.HandleFunc("/summary/filter/{left}={right}", handler.equalToFilterSummary).Methods("GET")
	router.HandleFunc("/filter/{left}>{right}", handler.greaterThanFilterJSON).Methods("GET")
	router.HandleFunc("/summary/filter/{left}>{right}", handler.greaterThanFilterSummary).Methods("GET")
	router.HandleFunc("/filter/{left}<{right}", handler.lessThanFilterJSON).Methods("GET")
	router.HandleFunc("/summary/filter/{left}<{right}", handler.lessThanFilterSummary).Methods("GET")
	router.HandleFunc("/filter/{left}>={right}", handler.greaterThanEqualToFilterJSON).Methods("GET")
	router.HandleFunc("/summary/filter/{left}>={right}", handler.greaterThanEqualToFilterSummary).Methods("GET")
	router.HandleFunc("/filter/{left}<={right}", handler.lessThanEqualToFilterJSON).Methods("GET")
	router.HandleFunc("/summary/filter/{left}<={right}", handler.lessThanEqualToFilterSummary).Methods("GET")
	router.HandleFunc("/filter/{left}!={right}", handler.notEqualToFilterJSON).Methods("GET")
	router.HandleFunc("/summary/filter/{left}!={right}", handler.notEqualToFilterSummary).Methods("GET")
	fmt.Println("The local server is now running on http://localhost:5000/")
	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Fatal(err)
	}
}
