package main

import (
	"github.com/nitishp1812/buildingloader/databaseapi"
	"github.com/nitishp1812/buildingloader/etlpipeline"
)

func main() {
	//extract the data from the API
	buildings := etlpipeline.Extract()

	//load the data into a MongoDb database
	collectionName := etlpipeline.Load(buildings)

	//Setup the API
	databaseapi.StartAPI(collectionName)
}
