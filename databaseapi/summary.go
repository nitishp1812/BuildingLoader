package databaseapi

import (
	"fmt"
	"sync"

	"github.com/nitishp1812/buildingloader/etlpipeline"
)

func getSummary(buildings []etlpipeline.DBBuilding) string {
	if len(buildings) == 0 {
		return "There were no items that matched the filter"
	}
	var group sync.WaitGroup

	mainChannel := make(chan string, 8)

	numberParameters := getNumberParameters(buildings)
	stringParameters := getStringParameters(buildings)

	numberParameterNames := []string{"Roof Height", "Shape Length", "Shape Area", "Ground Elevation", "Year constructed"}
	stringParameterNames := []string{"Geom Source", "Feature Code", "Last Status"}

	for index, parameter := range numberParameters {
		group.Add(1)
		go getNumberSummary(parameter, mainChannel, &group, numberParameterNames[index])
	}

	for index, parameter := range stringParameters {
		group.Add(1)
		go getStringSummary(parameter, mainChannel, &group, stringParameterNames[index])
	}

	group.Wait()

	close(mainChannel)

	output := ""

	for value := range mainChannel {
		output = output + value
	}

	return output
}

func getNumberSummary(values []float64, channel chan string, group *sync.WaitGroup, name string) {
	defer group.Done()

	numberSummary := ""
	channel <- numberSummary
}

func getStringSummary(values []string, channel chan string, group *sync.WaitGroup, name string) {
	defer group.Done()

	stringSummary := ""
	arguments := getCounts(values)

	max := 0
	maxType := ""

	for value, count := range arguments {
		if count > max {
			max = count
			maxType = value
		}
	}

	stringSummary = fmt.Sprintf("The parameter '%s' has the value '%s' in %d occurrences out of %d\n",
		name, maxType, max, len(values))

	channel <- stringSummary
}

func getCounts(values []string) map[string]int {
	countMap := map[string]int{}

	for _, value := range values {
		if _, present := countMap[value]; present {
			countMap[value]++
		} else {
			countMap[value] = 1
		}
	}

	return countMap
}

func getNumberParameters(buildings []etlpipeline.DBBuilding) [][]float64 {
	numberParams := [][]float64{}

	heights := []float64{}
	shapeLengths := []float64{}
	shapeAreas := []float64{}
	groundElevations := []float64{}
	constructYears := []float64{}

	for _, building := range buildings {
		heights = append(heights, building.HeightRoof)
		shapeLengths = append(shapeLengths, building.ShapeLen)
		shapeAreas = append(shapeAreas, building.ShapeArea)
		groundElevations = append(groundElevations, building.GroundElev)
		constructYears = append(constructYears, building.ConstructYear)
	}

	numberParams = append(numberParams, heights, shapeLengths, shapeAreas, groundElevations, constructYears)

	return numberParams
}

func getStringParameters(buildings []etlpipeline.DBBuilding) [][]string {
	stringParams := [][]string{}

	geomSources := []string{}
	featureCodes := []string{}
	lastStatusTypes := []string{}

	for _, building := range buildings {
		geomSources = append(geomSources, building.GeomSource)
		featureCodes = append(featureCodes, building.FeatCode)
		lastStatusTypes = append(lastStatusTypes, building.Lststatype)
	}

	stringParams = append(stringParams, geomSources, featureCodes, lastStatusTypes)

	return stringParams
}

func getMean(buildings []etlpipeline.DBBuilding) {

}
