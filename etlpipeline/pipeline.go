package etlpipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Extract returns a slice of structs which are used to represent the data extracted from the New York City Building footprints dataset.
func Extract() []APIBuilding {
	buildingURL := "https://data.cityofnewyork.us/resource/k8ez-gyqp.json"
	client := http.Client{
		Timeout: time.Second * 20,
	}

	request, err := http.NewRequest(http.MethodGet, buildingURL, nil)
	checkErr(err)

	request.Header.Set("User-Agent", "setup-test")
	response, err := client.Do(request)
	checkErr(err)

	body, err := ioutil.ReadAll(response.Body)
	checkErr(err)

	var buildings []APIBuilding
	if err := json.Unmarshal(body, &buildings); err != nil {
		log.Fatal(err)
	}

	return buildings
}

//Transform converts each object of the form of the data from the API to a form more suitable for storing and querying in a database
func (building *APIBuilding) Transform() *DBBuilding {
	return &DBBuilding{}
}

//Load loads the buildings extracted from the API to a local MongoDB collection in a database hosted at 'localhost:27017'
func Load(buildings []APIBuilding) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	checkErr(err)

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected")

	collectionName := setUpDB(client)

	buildingCollection := client.Database("nitishp1812buildingdb").Collection(collectionName)

	for _, building := range buildings {
		_, err := buildingCollection.InsertOne(context.Background(), *(building.Transform()))
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Data inserted into the MongoDb database 'nitishp1812buildingdb' in the collection '%s'", collectionName)
}

//setUpDB sets up the collection to store the data from the API in. It allows for a max of 3 collections in the database
func setUpDB(client *mongo.Client) (intendedName string) {
	var filter interface{}

	collectionsCursor, err := client.Database("nitishp1812buildingdb").ListCollections(context.Background(), filter)
	checkErr(err)

	timeSignature := fmt.Sprintf("%d-%s-%d", time.Now().Year(), time.Now().Month().String(), time.Now().Day())
	intendedName = fmt.Sprintf("buildings-%s", timeSignature)

	count := 0

	for collectionsCursor.Next(context.Background()) {
		var collection *mongo.Collection
		if err := collectionsCursor.Decode(collection); err != nil {
			log.Fatal(err)
		}

		if count == 2 {
			collection.Drop(context.Background())
			return
		}

		if collection.Name() == intendedName {
			collection.Drop(context.Background())
			return
		}

	}

	return
}

//checkErr checks whether the error is nil and prints the error and stops the program if it is
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
