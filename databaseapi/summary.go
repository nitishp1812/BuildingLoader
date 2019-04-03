package databaseapi

import (
	"fmt"
	"sort"
	"sync"

	"github.com/nitishp1812/buildingloader/etlpipeline"
)

//get all the summary output for the given array of database objects
func getSummary(buildings []etlpipeline.DBBuilding) string {
	//if there are no items after the filter
	if len(buildings) == 0 {
		return "There were no items that matched the filter"
	}

	var group sync.WaitGroup

	mainChannel := make(chan string, 8)

	//get the values of the numnber and string parameters in list form from the building objects
	numberParameters := getNumberParameters(buildings)
	stringParameters := getStringParameters(buildings)

	numberParameterNames := []string{"Roof Height", "Shape Length", "Shape Area", "Ground Elevation", "Year Constructed"}
	stringParameterNames := []string{"Geom Source", "Feature Code", "Last Status"}

	//range through each number parameter and start a goroutine to get the output for it
	for index, parameter := range numberParameters {
		group.Add(1)
		go getNumberSummary(parameter, mainChannel, &group, numberParameterNames[index])
	}

	//range through each string parameter and start a goroutine to get the output for it
	for index, parameter := range stringParameters {
		group.Add(1)
		go getStringSummary(parameter, mainChannel, &group, stringParameterNames[index])
	}

	//wait till all of the methods have generated their output
	group.Wait()

	//close the channel to indicate no more values will be put into it
	close(mainChannel)

	output := ""

	//iterate through each output and add it to get a combined output
	for value := range mainChannel {
		output = output + value
	}

	return output
}

//load the summary output into the channel for the given slice of numerical values
func getNumberSummary(values []float64, channel chan<- string, group *sync.WaitGroup, name string) {
	defer group.Done()

	numberSummary := ""

	//concurrently calculate the mean and median for the values
	var wg sync.WaitGroup
	wg.Add(2)

	meanChannel := make(chan float64, 1)
	medianChannel := make(chan float64, 1)

	go getMean(values, meanChannel, &wg)
	go getMedian(values, medianChannel, &wg)

	wg.Wait()

	//close the channels and load the mean and median
	close(meanChannel)
	close(medianChannel)

	mean := <-meanChannel
	median := <-medianChannel

	//generate output based on the values of the mean and median and the data being calculated
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

	//add this summary to the channel which was passed
	channel <- numberSummary
}

//load the summary output into the channel for the given slice of string values
func getStringSummary(values []string, channel chan<- string, group *sync.WaitGroup, name string) {
	defer group.Done()

	stringSummary := ""

	//get the count for each typ eof occurrence from the given values
	arguments := getCounts(values)

	max := 0
	maxType := ""

	//get the value with the max number of occurrences from the given values
	for value, count := range arguments {
		if count > max {
			max = count
			maxType = value
		}
	}

	//generate the output based on the calculated value
	stringSummary = fmt.Sprintf("The parameter '%s' has the value '%s' in %d occurrences out of %d\n\n\n",
		name, maxType, max, len(values))

	//add this summary to the channel which was passed
	channel <- stringSummary
}

//method to calculate the number of occurrences of each string as a map in the given slice of values
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

//method to get the individual numerical parameters from the object
func getNumberParameters(buildings []etlpipeline.DBBuilding) [][]float64 {
	numberParams := [][]float64{}

	//initialize the slices as empty slices
	heights := []float64{}
	shapeLengths := []float64{}
	shapeAreas := []float64{}
	groundElevations := []float64{}
	constructYears := []float64{}

	//add each value for that field to its respective array
	for _, building := range buildings {
		heights = append(heights, building.HeightRoof)
		shapeLengths = append(shapeLengths, building.ShapeLen)
		shapeAreas = append(shapeAreas, building.ShapeArea)
		groundElevations = append(groundElevations, building.GroundElev)
		constructYears = append(constructYears, building.ConstructYear)
	}

	//add the individual arrays to a 2d array for easy access
	numberParams = append(numberParams, heights, shapeLengths, shapeAreas, groundElevations, constructYears)

	return numberParams
}

//method to get the individual string parameters from the object
func getStringParameters(buildings []etlpipeline.DBBuilding) [][]string {
	stringParams := [][]string{}

	//initialize the slices as empty slices
	geomSources := []string{}
	featureCodes := []string{}
	lastStatusTypes := []string{}

	//add each value for that field to its respective array
	for _, building := range buildings {
		geomSources = append(geomSources, building.GeomSource)
		featureCodes = append(featureCodes, building.FeatCode)
		lastStatusTypes = append(lastStatusTypes, building.Lststatype)
	}

	//add the individual arrays to a 2d array for easy access
	stringParams = append(stringParams, geomSources, featureCodes, lastStatusTypes)

	return stringParams
}

//method to get the mean for the given set of values and pass it to the channel
func getMean(values []float64, channel chan<- float64, wg *sync.WaitGroup) {
	defer wg.Done()

	var average float64
	average = 0
	for _, value := range values {
		average += value
	}

	average /= float64(len(values))

	channel <- average
}

//method to get the median for the given set of values and pass it to the channel
func getMedian(values []float64, channel chan<- float64, wg *sync.WaitGroup) {
	defer wg.Done()

	//sort the array and then access the middle element(s) to get the mean
	sort.Float64s(values)

	var median float64

	if (len(values) % 2) == 0 {
		median = (values[(len(values)/2)] + values[(len(values)/2)+1]) / 2
	} else {
		median = values[(len(values)+1)/2]
	}

	channel <- median
}
