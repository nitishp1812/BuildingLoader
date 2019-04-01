package main

import (
	"github.com/nitishp1812/buildingloader/databaseapi"
	"github.com/nitishp1812/buildingloader/etlpipeline"
)

func main() {
	buildings := etlpipeline.Extract()
	collectionName := etlpipeline.Load(buildings)
	databaseapi.StartAPI(collectionName)
}
