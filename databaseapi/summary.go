package databaseapi

import (
	"fmt"
	"sort"
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

	numberParameterNames := []string{"Roof Height", "Shape Length", "Shape Area", "Ground Elevation", "Year Constructed"}
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

func getNumberSummary(values []float64, channel chan<- string, group *sync.WaitGroup, name string) {
	defer group.Done()

	numberSummary := ""

	var wg sync.WaitGroup
	wg.Add(2)

	meanChannel := make(chan float64, 1)
	medianChannel := make(chan float64, 1)

	go getMean(values, meanChannel, &wg)
	go getMedian(values, medianChannel, &wg)

	wg.Wait()

	close(meanChannel)
	close(medianChannel)

	mean := <-meanChannel
	median := <-medianChannel

	if name == "Year Constructed" {
		numberSummary = fmt.Sprintf("The mean of the values in the field '%s' is %0.0f.\nThe median is %0.0f.\n",
			name, mean, median)
	} else {
		numberSummary = fmt.Sprintf("The mean of the values in the field '%s' is %0.3f.\nThe median is %0.3f.\n",
			name, mean, median)
	}

	if mean > median {
		numberSummary = numberSummary + "Since the mean is greater than the median, we can infer that " +
			"the values are skewed to the left\n\n\n"
	} else if mean < median {
		numberSummary = numberSummary + "Since the median is greater than the mean, we can infer that " +
			"the values are skewed to the right\n\n\n"
	} else {
		numberSummary = numberSummary + "Since the mean is equal to the medina, we can infer that " +
			"the values are normally distributed\n\n\n"
	}

	channel <- numberSummary
}

func getStringSummary(values []string, channel chan<- string, group *sync.WaitGroup, name string) {
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

	stringSummary = fmt.Sprintf("The parameter '%s' has the value '%s' in %d occurrences out of %d\n\n\n",
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

func getMean(values []float64, channel chan float64, wg *sync.WaitGroup) {
	defer wg.Done()

	var average float64
	average = 0
	for _, value := range values {
		average += value
	}

	average /= float64(len(values))

	channel <- average
}

func getMedian(values []float64, channel chan float64, wg *sync.WaitGroup) {
	defer wg.Done()

	sort.Float64s(values)

	var median float64

	if (len(values) % 2) == 0 {
		median = (values[(len(values)/2)] + values[(len(values)/2)+1]) / 2
	} else {
		median = values[(len(values)+1)/2]
	}

	channel <- median
}
