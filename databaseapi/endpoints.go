package databaseapi

import (
	"context"
	"fmt"
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

func (handler *dbHandler) lessThanFilter(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	param1 := params["param1"]
	param2 := params["param2"]

	isFirstParam, firstType := isDBField(param1)
	isSecondParam, secondType := isDBField(param2)

	if isFirstParam {
		if firstType != "int" && firstType != "float" {
			fmt.Fprintln(writer, "The given parameter does not support this type of comparison")
			return
		}
	}

	if isSecondParam {
		if secondType != "int" && secondType != "float" {
			fmt.Fprintln(writer, "The given parameter does not support this type of comparison")
			return
		}
	}
}

func isDBField(parameter string) (bool, string) {
	switch {
	case parameter == "base_bbl":
		return true, "string"
	case parameter == "mpluto_bbl":
		return true, "string"
	case parameter == "shape_len":
		return true, "float"
	case parameter == "shape_area":
		return true, "float"
	case parameter == "geom_source":
		return true, "string"
	case parameter == "ground_elev":
		return true, "int"
	case parameter == "feat_code":
		return true, "string"
	case parameter == "height_roof":
		return true, "float"
	case parameter == "doitt_id":
		return true, "string"
	case parameter == "lststatus":
		return true, "string"
	case parameter == "lstmoddate":
		return true, "int"
	case parameter == "construct_year":
		return true, "int"
	case parameter == "bin":
		return true, "string"
	case parameter == "geom":
		return true, "polygon"
	default:
		return false, ""
	}
}
