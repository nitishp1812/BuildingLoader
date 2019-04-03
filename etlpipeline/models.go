package etlpipeline

import "time"

//APIBuilding represents the format in which the data is stored in the API result
type APIBuilding struct {
	GeomSource    string       "json:\"geomsource\""
	MplutoBbl     string       "json:\"mpluto_bbl\""
	BaseBbl       string       "json:\"base_bbl\""
	ShapeLen      string       "json:\"shape_len\""
	ShapeArea     string       "json:\"shape_area\""
	GroundElev    string       "json:\"groundelev\""
	FeatCode      string       "json:\"feat_code\""
	HeightRoof    string       "json:\"heightroof\""
	DoittID       string       "json:\"doitt_id\""
	Lststatype    string       "json:\"lststatype\""
	Lstmoddate    string       "json:\"lstmoddate\""
	Name          string       "json:\"name\""
	ConstructYear string       "json:\"cnstrct_yr\""
	Bin           string       "json:\"bin\""
	TheGeom       Multipolygon "json:\"the_geom\""
}

//DBBuilding represents the format which is used to store the data in the MongoDb database
type DBBuilding struct {
	GeomSource    string       "bson:\"geom_source\"      json:\"geom_source\""
	MplutoBbl     string       "bson:\"mpluto_bbl\"       json:\"mpluto_bbl\""
	BaseBbl       string       "bson:\"base_bbl\"         json:\"base_bbl\""
	ShapeLen      float64      "bson:\"shape_len\"        json:\"shape_len\""
	ShapeArea     float64      "bson:\"shape_area\"       json:\"shape_area\""
	GroundElev    float64      "bson:\"ground_elev\"      json:\"ground_elev\""
	FeatCode      string       "bson:\"feat_code\"        json:\"feat_code\""
	HeightRoof    float64      "bson:\"height_roof\"      json:\"height_roof\""
	DoittID       string       "bson:\"doitt_id\"         json:\"doitt_id\""
	Lststatype    string       "bson:\"last_status_type\" json:\"lststatype\""
	Lstmoddate    time.Time    "bson:\"last_mod_date\"    json:\"lstmoddate\""
	ConstructYear float64      "bson:\"construct_year\"   json:\"construct_year\""
	Bin           string       "bson:\"bin\"              json:\"bin\""
	Geom          Multipolygon "bson:\"geom\"             json:\"geom\""
}

//Multipolygon is the structure whixh represents the GeoJSON multipolygon structure
type Multipolygon struct {
	Type        string          "json:\"type\""
	Coordinates [][][][]float64 "json:\"coordinates\""
}
