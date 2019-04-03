package databaseapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nitishp1812/buildingloader/etlpipeline"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (handler *dbHandler) equalityFilter(writer http.ResponseWriter, request *http.Request,
	comparison string) ([]etlpipeline.DBBuilding, error) {
	parameters := mux.Vars(request)
	leftParameter := parameters["left"]
	rightParameter := parameters["right"]

	isLeftDataField, leftParameterDataType := isDBField(leftParameter)
	isRightDataField, rightParameterDataType := isDBField(rightParameter)

	if leftParameterDataType == "polygon" || rightParameterDataType == "polygon" {
		fmt.Fprintln(writer, "This field cannot be compared for equality")
		return []etlpipeline.DBBuilding{}, errors.New("Invalid comparison")
	}

	filter := bson.D{}
	var err error

	if isLeftDataField && isRightDataField {
		if leftParameterDataType == rightParameterDataType {
			filter = bson.D{{
				"$expr", bson.D{{
					comparison, bson.A{"$" + leftParameter, "$" + rightParameter},
				}},
			}}
		} else {
			fmt.Fprintln(writer, "The data fields are of different types and cannot be compared")
			return []etlpipeline.DBBuilding{}, errors.New("Invalid comparison")
		}
	} else if isLeftDataField {
		filter, err = getFilter(comparison, leftParameter, rightParameter, leftParameterDataType)
		if err != nil {
			fmt.Fprintln(writer, err.Error())
			return []etlpipeline.DBBuilding{}, errors.New("Invalid type conversion attempt")
		}
	} else if isRightDataField {
		filter, err = getFilter(comparison, rightParameter, leftParameter, rightParameterDataType)
		if err != nil {
			fmt.Fprintln(writer, err.Error())
			return []etlpipeline.DBBuilding{}, errors.New("Invalid type conversion attempt")
		}
	} else {
		fmt.Fprintln(writer, "The filter is not valid. One of the operands must be the name of a field in the dataset")
		return []etlpipeline.DBBuilding{}, errors.New("No database field keyed")
	}

	buildings := handler.getFilteredBuildings(&filter)

	return buildings, nil
}

func (handler *dbHandler) floatComparisonFilter(writer http.ResponseWriter, request *http.Request,
	comparison string, invertedComparison string) ([]etlpipeline.DBBuilding, error) {
	parameters := mux.Vars(request)
	leftParameter := parameters["left"]
	rightParameter := parameters["right"]

	isLeftDataField, leftParameterDataType := isDBField(leftParameter)
	isRightDataField, rightParameterDataType := isDBField(rightParameter)

	if isLeftDataField && leftParameterDataType != "float" {
		fmt.Fprintln(writer, "The given parameter does not support this type of comparison")
		return []etlpipeline.DBBuilding{}, errors.New("Invalid comparison")
	}

	if isRightDataField && rightParameterDataType != "float" {
		fmt.Fprintln(writer, "The given parameter does not support this type of comparison")
		return []etlpipeline.DBBuilding{}, errors.New("Invalid comparison")
	}

	filter := bson.D{}
	var err error

	if isLeftDataField && isRightDataField {
		filter = bson.D{{
			"$expr", bson.D{{
				comparison, bson.A{"$" + leftParameter, "$" + rightParameter},
			}},
		}}
	} else if isLeftDataField {
		filter, err = getFilter(comparison, leftParameter, rightParameter, leftParameterDataType)
		if err != nil {
			fmt.Fprintln(writer, err.Error())
			return []etlpipeline.DBBuilding{}, errors.New("Invalid type conversion attempt")
		}
	} else if isRightDataField {
		filter, err = getFilter(invertedComparison, rightParameter, leftParameter, rightParameterDataType)
		if err != nil {
			fmt.Fprintln(writer, err.Error())
			return []etlpipeline.DBBuilding{}, errors.New("Invalid type conversion attempt")
		}
	} else {
		fmt.Fprintln(writer, "The filter is not valid. One of the operands must be the name of a field in the dataset")
		return []etlpipeline.DBBuilding{}, errors.New("No database field keyed")
	}

	buildings := handler.getFilteredBuildings(&filter)
	return buildings, nil
}

func (handler *dbHandler) getFilteredBuildings(filter *bson.D) []etlpipeline.DBBuilding {
	var buildings []etlpipeline.DBBuilding

	cursor, err := handler.collection.Find(context.Background(), *filter, options.Find())
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

	return buildings
}

func getFilter(operation string, parameter string, valueString string, dataType string) (filter bson.D, err error) {
	if dataType == "float" {
		value, err := strconv.ParseFloat(valueString, 64)
		if err != nil {
			return nil, errors.New("The given value could not be parsed to a decimal")
		}
		filter = bson.D{{
			parameter, bson.D{{
				operation, value,
			}},
		}}
	} else {
		filter = bson.D{{
			parameter, bson.D{{
				operation, valueString,
			}},
		}}
	}
	return
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
		return true, "float"
	case parameter == "feat_code":
		return true, "string"
	case parameter == "height_roof":
		return true, "float"
	case parameter == "doitt_id":
		return true, "string"
	case parameter == "lststatus":
		return true, "string"
	case parameter == "lstmoddate":
		return true, "float"
	case parameter == "construct_year":
		return true, "float"
	case parameter == "bin":
		return true, "string"
	case parameter == "geom":
		return true, "polygon"
	default:
		return false, ""
	}
}
