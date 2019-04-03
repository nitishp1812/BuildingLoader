package etlpipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//method to webscrape the URl for the data
func getURL() string {
	website := "https://data.cityofnewyork.us/Housing-Development/Building-Footprints/nqwf-w8eh"
	response, err := http.Get(website)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Get the response body as a string
	byteData, err := ioutil.ReadAll(response.Body)
	pageContent := string(byteData)

	// Find the start index for the data link
	urlStartIndex := strings.Index(pageContent, "\"https://data.cityofnewyork.us/resource/")
	if urlStartIndex == -1 {
		log.Fatal("No such element found")
	}

	//increase by 1 to ignore the starting quotes
	urlStartIndex++

	//the length of the string without the quotes is 53
	urlEndIndex := urlStartIndex + 53

	//get url by getting the byte elements between the two indexes
	url := []byte(pageContent[urlStartIndex:urlEndIndex])

	return string(url)
}

//Extract returns a slice of structs which are used to represent the data extracted from the New York City Building footprints dataset.
func Extract() []APIBuilding {
	buildingURL := getURL()

	//make the client to be used to call the API
	client := http.Client{
		Timeout: time.Second * 20,
	}

	//make the request object to get the data from the URL
	request, err := http.NewRequest(http.MethodGet, buildingURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("User-Agent", "setup-test")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	//get the body (data) from teh response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	//unpack the data into a slice of structs
	var buildings []APIBuilding
	if err := json.Unmarshal(body, &buildings); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully extracted the data from the API")

	return buildings
}

//Transform converts each object of the form of the data from the API to a form more suitable for storing and querying in a database
func (building *APIBuilding) Transform() bson.D {
	constructYear, err := strconv.ParseInt(building.ConstructYear, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	groundElev, err := strconv.ParseInt(building.GroundElev, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	heightRoof, err := strconv.ParseFloat(building.HeightRoof, 64)
	if err != nil {
		log.Fatal(err)
	}

	shapeArea, err := strconv.ParseFloat(building.ShapeArea, 64)
	if err != nil {
		log.Fatal(err)
	}

	shapeLen, err := strconv.ParseFloat(building.ShapeLen, 64)
	if err != nil {
		log.Fatal(err)
	}

	modTime, err := time.Parse(time.RFC3339Nano, building.Lstmoddate)
	if err != nil {
		log.Fatal(err)
	}

	//structure of the data in the MongoDB database
	document := bson.D{
		{"base_bbl", building.BaseBbl},
		{"bin", building.Bin},
		{"construct_year", constructYear},
		{"doitt_id", building.DoittID},
		{"feat_code", building.FeatCode},
		{"geom", bson.D{
			{"type", building.TheGeom.Type},
			{"coordinates", building.TheGeom.Coordinates},
		}},
		{"geom_source", building.GeomSource},
		{"ground_elev", groundElev},
		{"height_roof", heightRoof},
		{"last_mod_date", modTime},
		{"last_status_type", building.Lststatype},
		{"mpluto_bbl", building.MplutoBbl},
		{"shape_area", shapeArea},
		{"shape_len", shapeLen},
	}

	return document
}

//Load loads the buildings extracted from the API to a local MongoDB collection in a database hosted at 'localhost:27017'
func Load(buildings []APIBuilding) string {
	//set up a connection to a mongodb local server
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	//set up the database and get the name of the collection to store the data in
	collectionName := setUpDB(client)

	buildingCollection := client.Database("nitishp1812buildingdb").Collection(collectionName)

	//iterate through each building object and store in the proper format in the database
	for _, building := range buildings {
		_, err := buildingCollection.InsertOne(context.Background(), building.Transform())
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Data inserted into the MongoDb database 'nitishp1812buildingdb' in the collection '%s' hosted at 'mongodb://localhost:27017'\n", collectionName)

	return collectionName
}

//setUpDB sets up the collection to store the data from the API in. It allows for a max of 3 collections in the database
func setUpDB(client *mongo.Client) (intendedName string) {
	//drop any existing databases with the same name
	if err := client.Database("nitishp1812buildingdb").Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	daySignature := fmt.Sprintf("%d-%s-%d", time.Now().Year(), time.Now().Month().String(), time.Now().Day())
	intendedName = fmt.Sprintf("buildings-%s", daySignature)

	return
}
