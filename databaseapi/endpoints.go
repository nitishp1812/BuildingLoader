package databaseapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//structure to store the pointer to the collection being accessed
type dbHandler struct {
	collection *mongo.Collection
}

//endpoint to generate all objects in JSON format
func (handler *dbHandler) showAllJSON(writer http.ResponseWriter, request *http.Request) {
	//generate the filtered buildings by passing appropriate comparison parameters
	buildings := handler.getFilteredBuildings(&bson.D{{}})

	//write the JSON output of the filtered list to the screen
	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(buildings)
}

//endpoint to generate the summary of all objects
func (handler *dbHandler) showAllSummary(writer http.ResponseWriter, request *http.Request) {
	//generate the filtered buildings by passing appropriate comparison parameters
	buildings := handler.getFilteredBuildings(&bson.D{{}})

	//get the output summary for the filtered list and print it on screen
	output := getSummary(buildings)
	fmt.Fprintln(writer, output)
}

//endpoint to generate the JSON output for all objects with the given 2 parameters equal
func (handler *dbHandler) equalToFilterJSON(writer http.ResponseWriter, request *http.Request) {
	//generate the filtered buildings by passing appropriate comparison parameters
	buildings, err := handler.equalityFilter(writer, request, "$eq")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	//write the JSON output of the filtered list to the screen
	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(buildings)
}

//endpoint to generate the summary output for all objects with the given 2 parameters equal
func (handler *dbHandler) equalToFilterSummary(writer http.ResponseWriter, request *http.Request) {
	//generate the filtered buildings by passing appropriate comparison parameters
	buildings, err := handler.equalityFilter(writer, request, "$eq")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	//get the output summary for the filtered list and print it on screen
	output := getSummary(buildings)
	fmt.Fprintln(writer, output)
}

//endpoint to generate the JSON output for all objects with the given 2 parameters unequal
func (handler *dbHandler) notEqualToFilterJSON(writer http.ResponseWriter, request *http.Request) {
	//generate the filtered buildings by passing appropriate comparison parameters
	buildings, err := handler.equalityFilter(writer, request, "$ne")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	//write the JSON output of the filtered list to the screen
	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(buildings)
}

//endpoint to generate the summary output for all objects with the given 2 parameters unequal
func (handler *dbHandler) notEqualToFilterSummary(writer http.ResponseWriter, request *http.Request) {
	//generate the filtered buildings by passing appropriate comparison parameters
	buildings, err := handler.equalityFilter(writer, request, "$ne")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	//get the output summary for the filtered list and print it on screen
	output := getSummary(buildings)
	fmt.Fprintln(writer, output)

}

//endpoint to generate the JSON output for all objects with the given parameter greater than the other
func (handler *dbHandler) greaterThanFilterJSON(writer http.ResponseWriter, request *http.Request) {
	//generate the filtered buildings by passing appropriate comparison parameters
	buildings, err := handler.floatComparisonFilter(writer, request, "$gt", "$lt")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	//write the JSON output of the filtered list to the screen
	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(buildings)
}

//endpoint to generate the summary output for all objects with the given parameter greater than the other
func (handler *dbHandler) greaterThanFilterSummary(writer http.ResponseWriter, request *http.Request) {
	//generate the filtered buildings by passing appropriate comparison parameters
	buildings, err := handler.floatComparisonFilter(writer, request, "$gt", "$lt")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	//get the output summary for the filtered list and print it on screen
	output := getSummary(buildings)
	fmt.Fprintln(writer, output)
}

//endpoint to generate the JSON output for all objects with the given parameter lesser than the other
func (handler *dbHandler) lessThanFilterJSON(writer http.ResponseWriter, request *http.Request) {
	//generate the filtered buildings by passing appropriate comparison parameters
	buildings, err := handler.floatComparisonFilter(writer, request, "$lt", "$gt")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	//write the JSON output of the filtered list to the screen
	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(buildings)
}

//endpoint to generate the summary output for all objects with the given parameter lesser than the other
func (handler *dbHandler) lessThanFilterSummary(writer http.ResponseWriter, request *http.Request) {
	//generate the filtered buildings by passing appropriate comparison parameters
	buildings, err := handler.floatComparisonFilter(writer, request, "$lt", "$gt")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	//get the output summary for the filtered list and print it on screen
	output := getSummary(buildings)
	fmt.Fprintln(writer, output)
}
