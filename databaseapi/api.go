package databaseapi

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/nitishp1812/buildingloader/etlpipeline"

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
	router.HandleFunc("/", handler.showAll).Methods("GET")
	router.HandleFunc("/filter/{param1}={param2}", handler.showAll).Methods("GET")
	router.HandleFunc("/filter/{param1}>{param2}", handler.showAll).Methods("GET")
	router.HandleFunc("/filter/{param1}<{param2}", handler.lessThanFilter)
	router.HandleFunc("/filter/{param1}>={param2}", handler.showAll).Methods("GET")
	router.HandleFunc("/filter/{param1}<={param2}", handler.showAll).Methods("GET")
	router.HandleFunc("/filter/{param1}!={param2}", handler.showAll).Methods("GET")
	fmt.Println("The local server is now running on http://localhost:5000/")
	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Fatal(err)
	}
}

func generateOutput(writer http.ResponseWriter, buildings []etlpipeline.DBBuilding) {
	fmt.Fprintln(writer, buildings)
}
