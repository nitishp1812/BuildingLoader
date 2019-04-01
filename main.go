package main

import (
	"fmt"

	"github.com/nitishp1812/buildingloader/databaseapi"
	"github.com/nitishp1812/buildingloader/etlpipeline"
)

func main() {
	fmt.Println("jrfbj4rh")
	buildings := etlpipeline.Extract()
	fmt.Println("Extracted")
	collectionName := etlpipeline.Load(buildings)
	fmt.Println("Saved")

	databaseapi.StartAPI(collectionName)
}
