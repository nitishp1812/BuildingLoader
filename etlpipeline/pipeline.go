package etlpipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("User-Agent", "setup-test")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var buildings []APIBuilding
	if err := json.Unmarshal(body, &buildings); err != nil {
		log.Fatal(err)
	}

	return buildings
}

//Transform converts each object of the form of the data from the API to a form more suitable for storing and querying in a database
func (building *APIBuilding) Transform() *DBBuilding {
	var newBuilding DBBuilding
	var err error

	newBuilding.BaseBbl = building.BaseBbl
	newBuilding.Bin = building.Bin
	newBuilding.ConstructYear = building.ConstructYear
	newBuilding.DoittID = building.DoittID
	newBuilding.FeatCode, err = strconv.ParseInt(building.FeatCode, 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	newBuilding.Geom = building.TheGeom
	newBuilding.GeomSource = building.GeomSource
	newBuilding.GroundElev, err = strconv.ParseInt(building.GroundElev, 10, 64)

	if err != nil {
		log.Fatal(err)
	}

	newBuilding.HeightRoof, err = strconv.ParseFloat(building.HeightRoof, 64)

	if err != nil {
		log.Fatal(err)
	}

	newBuilding.Lststatype = building.Lststatype
	newBuilding.MplutoBbl = building.MplutoBbl
	newBuilding.Name = building.Name
	newBuilding.ShapeArea, err = strconv.ParseFloat(building.ShapeArea, 64)

	if err != nil {
		log.Fatal(err)
	}

	newBuilding.ShapeLen, err = strconv.ParseFloat(building.ShapeLen, 64)

	if err != nil {
		log.Fatal(err)
	}

	return &newBuilding
}

//Load loads the buildings extracted from the API to a local MongoDB collection in a database hosted at 'localhost:27017'
func Load(buildings []APIBuilding) string {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created client")

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	collectionName := setUpDB(client)

	buildingCollection := client.Database("nitishp1812buildingdb").Collection(collectionName)

	for _, building := range buildings {
		_, err := buildingCollection.InsertOne(context.Background(), *(building.Transform()))
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Data inserted into the MongoDb database 'nitishp1812buildingdb' in the collection '%s'", collectionName)

	return collectionName
}

//setUpDB sets up the collection to store the data from the API in. It allows for a max of 3 collections in the database
func setUpDB(client *mongo.Client) (intendedName string) {
	collectionsCursor, err := client.Database("nitishp1812buildingdb").ListCollections(context.Background(), bson.D{{}})
	defer collectionsCursor.Close(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	timeSignature := fmt.Sprintf("%d-%s-%d", time.Now().Year(), time.Now().Month().String(), time.Now().Day())
	intendedName = fmt.Sprintf("buildings-%s", timeSignature)

	count := 0

	for collectionsCursor.Next(context.Background()) {
		var collection mongo.Collection
		if err := collectionsCursor.Decode(&collection); err != nil {
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
