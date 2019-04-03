package databaseapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type dbHandler struct {
	collection *mongo.Collection
}

func (handler *dbHandler) showAllJSON(writer http.ResponseWriter, request *http.Request) {
	buildings := handler.getFilteredBuildings(&bson.D{{}})

	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(buildings)
}

func (handler *dbHandler) showAllSummary(writer http.ResponseWriter, request *http.Request) {
	buildings := handler.getFilteredBuildings(&bson.D{{}})

	output := getSummary(buildings)
	fmt.Fprintln(writer, output)
}

func (handler *dbHandler) equalToFilterJSON(writer http.ResponseWriter, request *http.Request) {
	buildings, err := handler.equalityFilter(writer, request, "$eq")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(buildings)
}

func (handler *dbHandler) equalToFilterSummary(writer http.ResponseWriter, request *http.Request) {
	buildings, err := handler.equalityFilter(writer, request, "$eq")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	output := getSummary(buildings)
	fmt.Fprintln(writer, output)
}

func (handler *dbHandler) notEqualToFilterJSON(writer http.ResponseWriter, request *http.Request) {
	buildings, err := handler.equalityFilter(writer, request, "$ne")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(buildings)

}

func (handler *dbHandler) notEqualToFilterSummary(writer http.ResponseWriter, request *http.Request) {
	buildings, err := handler.equalityFilter(writer, request, "$ne")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	output := getSummary(buildings)
	fmt.Fprintln(writer, output)

}

func (handler *dbHandler) greaterThanFilterJSON(writer http.ResponseWriter, request *http.Request) {
	buildings, err := handler.floatComparisonFilter(writer, request, "$gt", "$lt")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(buildings)

}

func (handler *dbHandler) greaterThanFilterSummary(writer http.ResponseWriter, request *http.Request) {
	buildings, err := handler.floatComparisonFilter(writer, request, "$gt", "$lt")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	output := getSummary(buildings)
	fmt.Fprintln(writer, output)
}

func (handler *dbHandler) lessThanFilterJSON(writer http.ResponseWriter, request *http.Request) {
	buildings, err := handler.floatComparisonFilter(writer, request, "$lt", "$gt")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(buildings)
}

func (handler *dbHandler) lessThanFilterSummary(writer http.ResponseWriter, request *http.Request) {
	buildings, err := handler.floatComparisonFilter(writer, request, "$lt", "$gt")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	output := getSummary(buildings)
	fmt.Fprintln(writer, output)
}

func (handler *dbHandler) greaterThanEqualToFilterJSON(writer http.ResponseWriter, request *http.Request) {
	buildings, err := handler.floatComparisonFilter(writer, request, "$gte", "$lte")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(buildings)
}

func (handler *dbHandler) greaterThanEqualToFilterSummary(writer http.ResponseWriter, request *http.Request) {
	buildings, err := handler.floatComparisonFilter(writer, request, "$gte", "$lte")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	output := getSummary(buildings)
	fmt.Fprintln(writer, output)
}

func (handler *dbHandler) lessThanEqualToFilterJSON(writer http.ResponseWriter, request *http.Request) {
	buildings, err := handler.floatComparisonFilter(writer, request, "$lte", "$gte")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	writer.Header().Set("Content-type", "application/json")
	json.NewEncoder(writer).Encode(buildings)
}

func (handler *dbHandler) lessThanEqualToFilterSummary(writer http.ResponseWriter, request *http.Request) {
	buildings, err := handler.floatComparisonFilter(writer, request, "$lte", "$gte")
	if err != nil {
		fmt.Fprintln(writer, err.Error())
	}

	output := getSummary(buildings)
	fmt.Fprintln(writer, output)
}
