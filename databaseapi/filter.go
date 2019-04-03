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

//method to generate the filtered list on the basis of its arguments and return the list
func (handler *dbHandler) equalityFilter(writer http.ResponseWriter, request *http.Request,
	comparison string) ([]etlpipeline.DBBuilding, error) {
	//get the arguments passed by the user
	parameters := mux.Vars(request)
	leftParameter := parameters["left"]
	rightParameter := parameters["right"]

	//check whether the parameters are valid for the comparison
	isLeftDataField, leftParameterDataType := isDBField(leftParameter)
	isRightDataField, rightParameterDataType := isDBField(rightParameter)

	if leftParameterDataType == "polygon" || rightParameterDataType == "polygon" {
		fmt.Fprintln(writer, "This field cannot be compared for equality")
		return []etlpipeline.DBBuilding{}, errors.New("Invalid comparison")
	}

	filter := bson.D{}
	var err error

	//generate a filter based on whether both the parameters are database fields, or one of the two is
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
		//throw an error if no parameter is a database field
		fmt.Fprintln(writer, "The filter is not valid. One of the operands must be the name of a field in the dataset")
		return []etlpipeline.DBBuilding{}, errors.New("No database field keyed")
	}

	//generate the filtered list on the basis of the filter
	buildings := handler.getFilteredBuildings(&filter)

	return buildings, nil
}

//method to generate the filtered list on the basis of its arguments and return the list
func (handler *dbHandler) floatComparisonFilter(writer http.ResponseWriter, request *http.Request,
	comparison string, invertedComparison string) ([]etlpipeline.DBBuilding, error) {
	//get the arguments passed by the user
	parameters := mux.Vars(request)
	leftParameter := parameters["left"]
	rightParameter := parameters["right"]

	//check whether the parameters are valid for the comparison
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

	//generate a filter based on whether both the parameters are database fields, or one of the two is
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
		//throw an error if no parameter is a database field
		fmt.Fprintln(writer, "The filter is not valid. One of the operands must be the name of a field in the dataset")
		return []etlpipeline.DBBuilding{}, errors.New("No database field keyed")
	}

	//generate the filtered list on the basis of the filter
	buildings := handler.getFilteredBuildings(&filter)

	return buildings, nil
}

//method to filter the objects on the basis of the passed filter and the passed collection
func (handler *dbHandler) getFilteredBuildings(filter *bson.D) []etlpipeline.DBBuilding {
	var buildings []etlpipeline.DBBuilding

	//get a cursor pointing to the documents in the collection which pass the filter
	cursor, err := handler.collection.Find(context.Background(), *filter, options.Find())
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	//iterate through every document in the filtered collection and add it to the filtered slice
	for cursor.Next(context.Background()) {
		var building etlpipeline.DBBuilding
		if err := cursor.Decode(&building); err != nil {
			log.Fatal(err)
		}

		buildings = append(buildings, building)
	}

	return buildings
}

//method to get the filter object based on the given operation, the parameter name and the value to be operated with
func getFilter(operation string, parameter string, valueString string, dataType string) (filter bson.D, err error) {
	if dataType == "float" {
		//raise an error if the value cannot be converted to a float
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

//method to determine whether the field exists in the database and if it does, get its data type
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
	case parameter == "construct_year":
		return true, "float"
	case parameter == "bin":
		return true, "string"
	default:
		return false, ""
	}
}
