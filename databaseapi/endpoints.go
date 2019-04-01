package databaseapi

import (
	"context"
	"log"
	"net/http"

	"github.com/nitishp1812/buildingloader/etlpipeline"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbHandler struct {
	collection *mongo.Collection
}

func (handler *dbHandler) showAll(writer http.ResponseWriter, request *http.Request) {
	var buildings []etlpipeline.DBBuilding

	cursor, err := handler.collection.Find(context.Background(), bson.M{}, options.Find())
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var building etlpipeline.DBBuilding
		if err := cursor.Decode(&building); err != nil {
			log.Fatal(err)
		}

		buildings = append(buildings, building)
	}

	generateOutput(writer, buildings)
}

func (handler *dbHandler) nameEqualTo(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	name := params["name"]
	filter := bson.M{"name": bson.M{"$eq": name}}

	var buildings []etlpipeline.DBBuilding

	cursor, err := handler.collection.Find(context.Background(), filter, options.Find())
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var building etlpipeline.DBBuilding
		if err := cursor.Decode(&building); err != nil {
			log.Fatal(err)
		}

		buildings = append(buildings, building)
	}

	generateOutput(writer, buildings)
}

func (handler *dbHandler) nameNotEqualTo(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	name := params["name"]
	filter := bson.M{"groundelev": bson.M{"$ne": name}}

	var buildings []etlpipeline.DBBuilding

	cursor, err := handler.collection.Find(context.Background(), filter, options.Find())
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var building etlpipeline.DBBuilding
		if err := cursor.Decode(&building); err != nil {
			log.Fatal(err)
		}
		buildings = append(buildings, building)
	}

	generateOutput(writer, buildings)
}

func (handler *dbHandler) groundElevGreaterThan(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	groundElev := params["groundelev"]
	filter := bson.M{"groundelev": bson.M{"$gt": groundElev}}

	var buildings []etlpipeline.DBBuilding

	cursor, err := handler.collection.Find(context.Background(), filter, options.Find())
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var building etlpipeline.DBBuilding
		if err := cursor.Decode(&building); err != nil {
			log.Fatal(err)
		}

		buildings = append(buildings, building)
	}
}

func (handler *dbHandler) groundElevLessThan(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	groundElev := params["groundelev"]
	filter := bson.M{"groundelev": bson.M{"$lt": groundElev}}

	var buildings []etlpipeline.DBBuilding

	cursor, err := handler.collection.Find(context.Background(), filter, options.Find())
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var building etlpipeline.DBBuilding
		if err := cursor.Decode(&building); err != nil {
			log.Fatal(err)
		}

		buildings = append(buildings, building)
	}
}

func (handler *dbHandler) groundElevEqualTo(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	groundElev := params["groundelev"]
	filter := bson.M{"groundelev": bson.M{"$eq": groundElev}}

	var buildings []etlpipeline.DBBuilding

	cursor, err := handler.collection.Find(context.Background(), filter, options.Find())
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var building etlpipeline.DBBuilding
		if err := cursor.Decode(&building); err != nil {
			log.Fatal(err)
		}

		buildings = append(buildings, building)
	}
}

func (handler *dbHandler) groundElevNotEqulaTo(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	groundElev := params["groundelev"]
	filter := bson.M{"groundelev": bson.M{"$ne": groundElev}}

	var buildings []etlpipeline.DBBuilding

	cursor, err := handler.collection.Find(context.Background(), filter, options.Find())
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var building etlpipeline.DBBuilding
		if err := cursor.Decode(&building); err != nil {
			log.Fatal(err)
		}

		buildings = append(buildings, building)
	}
}
